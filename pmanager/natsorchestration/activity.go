//  Copyright (c) 2025 Metaform Systems, Inc
//
//  This program and the accompanying materials are made available under the
//  terms of the Apache License, Version 2.0 which is available at
//  https://www.apache.org/licenses/LICENSE-2.0
//
//  SPDX-License-Identifier: Apache-2.0
//
//  Contributors:
//       Metaform Systems, Inc. - initial API and implementation
//

package natsorchestration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/nats-io/nats.go/jetstream"
	"time"
)

const (
	streamName      = "cfm-activity"
	durableConsumer = "cfm-durable-activity"
)

type NatsActivityExecutor struct {
	id                string
	client            MsgClient
	activityName      string
	activityProcessor api.ActivityProcessor
	monitor           monitor.LogMonitor
}

// Execute starts a goroutine to process messages from the activity queue.
func (e *NatsActivityExecutor) Execute(ctx context.Context) error {
	stream, err := e.client.Stream(ctx, streamName)
	if err != nil {
		return fmt.Errorf("error opening stream: %w", err)
	}

	consumer, err := stream.Consumer(ctx, durableConsumer)
	if err != nil {
		return fmt.Errorf("error connecting to consumer: %w", err)
	}

	go func() {
		err := e.processLoop(ctx, consumer)
		if err != nil {
			e.monitor.Debugf("Error processing message: %v", err)
		}
	}()
	return nil
}

// processLoop handles the main loop for consuming and processing messages from a JetStream consumer.
// It runs continuously until the provided context is canceled or an error occurs.
// Returns an error if message fetching or processing fails.
func (e *NatsActivityExecutor) processLoop(ctx context.Context, consumer jetstream.Consumer) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			messageBatch, err := consumer.Fetch(1, jetstream.FetchMaxWait(time.Second))
			if err != nil {
				return err // TODO handle a better way
			}

			for message := range messageBatch.Messages() {
				if err = e.processMessage(ctx, message); err != nil {
					e.monitor.Debugf("Error processing message: %v", err)
				}
			}
		}
	}
}

// processMessage processes a single message from the JetStream consumer by delegating to its ActivityProcessor. When
// processing is complete, the orchestration state is updated, messages for the next activities are enqueued if the
// orchestration can proceed, and the original message is acknowledged.
//
// Returns an error if message processing fails.
func (e *NatsActivityExecutor) processMessage(ctx context.Context, message jetstream.Msg) error {
	var oMessage api.ActivityMessage
	if err := json.Unmarshal(message.Data(), &oMessage); err != nil {
		// TODO pass to DLQ
		return fmt.Errorf("failed to unmarshal orchestration message: %w", err)
	}

	orchestration, revision, err := ReadOrchestration(ctx, oMessage.OrchestrationID, e.client)
	if err != nil {
		return fmt.Errorf("failed to read orchestration data: %w", err)
	}

	activityContext := newActivityContext(ctx, orchestration.ID, oMessage.Activity)
	e.monitor.Debugf("Received message: %s\n", e.id)
	result := e.activityProcessor.Process(activityContext)

	switch result.Result {
	case api.ActivityResultRetryError:
		return e.handleRetryError(message, oMessage.OrchestrationID, result.Error)

	case api.ActivityResultFatalError:
		return e.handleFatalError(ctx, message, orchestration, revision, oMessage.OrchestrationID, result.Error)

	case api.ActivityResultWait:
		return e.ackMessage(message, oMessage.OrchestrationID)

	case api.ActivityResultSchedule:
		if err := message.NakWithDelay(result.WaitMillis); err != nil {
			return fmt.Errorf("failed to reschedule schedule activity %s: %w", oMessage.OrchestrationID, err)
		}
		return nil
	}

	return e.processOrchestration(ctx, message, oMessage)
}

func (e *NatsActivityExecutor) processOrchestration(ctx context.Context, message jetstream.Msg, oMessage api.ActivityMessage) error {
	// Re-read the orchestration and update it
	orchestration, revision, err := ReadOrchestration(ctx, oMessage.OrchestrationID, e.client)
	if err != nil {
		return fmt.Errorf("failed to read orchestration data for update: %w", err)
	}

	orchestration, revision, err = CompleteOrchestrationActivity(ctx, orchestration, oMessage.Activity.ID, revision, e.client)
	if err != nil {
		return err
	}

	// Return if orchestration is in the error state
	if orchestration.State == api.OrchestrationStateErrored {
		return e.ackMessage(message, oMessage.OrchestrationID)
	}

	canProceed, err := orchestration.CanProceedToNextStep(oMessage.Activity.ID)
	if err != nil {
		return fmt.Errorf("failed to proceed with orchestration %s: %w", oMessage.OrchestrationID, err)
	}

	if !canProceed {
		return e.ackMessage(message, oMessage.OrchestrationID)
	}

	next := orchestration.GetNextStepActivities(oMessage.Activity.ID)
	if len(next) == 0 {
		// No more steps, mark as completed
		return e.handleOrchestrationCompletion(ctx, message, orchestration, revision, oMessage.OrchestrationID)
	}

	e.monitor.Debugf("Finished activity: %s", oMessage.Activity.Type)

	// Enqueue next activities
	if err := EnqueueActivityMessages(ctx, orchestration.ID, next, e.client); err != nil {
		return fmt.Errorf("failed to enqueue next orchestration activities %s: %w", oMessage.OrchestrationID, err)
	}

	return e.ackMessage(message, oMessage.OrchestrationID)
}

func (e *NatsActivityExecutor) handleOrchestrationCompletion(ctx context.Context, message jetstream.Msg, orchestration api.Orchestration, revision uint64, orchestrationID string) error {
	if _, _, err := MarkOrchestrationCompleted(ctx, orchestration, revision, e.client); err != nil {
		e.monitor.Infof("Failed to mark orchestration %s as completed: %v", orchestrationID, err)
	}

	e.monitor.Debugf("Finished orchestration: %s", orchestrationID)
	return e.ackMessage(message, orchestrationID)
}

func (e *NatsActivityExecutor) handleRetryError(message jetstream.Msg, orchestrationID string, resultErr error) error {
	// Nak to redeliver message
	if err := message.Nak(); err != nil {
		return fmt.Errorf("failed to execute activity message and NAK response %s (errors: %w, %v)",
			orchestrationID, resultErr, err)
	}
	return fmt.Errorf("failed to execute activity %s: %w", orchestrationID, resultErr)
}

func (e *NatsActivityExecutor) handleFatalError(ctx context.Context, message jetstream.Msg, orchestration api.Orchestration, revision uint64, orchestrationID string, resultErr error) error {
	// Update the orchestration before acking back. If the update fails, just log it to ensure the ack is sent to avoid message re-delivery
	if _, _, err := MarkOrchestrationErrored(ctx, orchestration, revision, e.client); err != nil {
		e.monitor.Infof("Failed to mark orchestration %s as fatal: %v", orchestrationID, err)
	}

	if err := message.Ack(); err != nil {
		return fmt.Errorf("fatal failure while executing activity %s (errors: %w, %v)",
			orchestrationID, resultErr, err)
	}
	return fmt.Errorf("fatal failure while executing activity %s: %w", orchestrationID, resultErr)
}

func (e *NatsActivityExecutor) ackMessage(message jetstream.Msg, orchestrationID string) error {
	if err := message.Ack(); err != nil {
		return fmt.Errorf("failed to ACK activity message %s: %w", orchestrationID, err)
	}
	return nil
}

type defaultActivityContext struct {
	activity api.Activity
	oID      string
	context  context.Context
	data     map[string]any
}

func newActivityContext(ctx context.Context, oID string, activity api.Activity) api.ActivityContext {
	return defaultActivityContext{
		activity: activity,
		oID:      oID,
		context:  ctx,
		data:     make(map[string]any),
	}
}

// Context returns the current request context
func (d defaultActivityContext) Context() context.Context {
	return d.context
}

// ID returns the ID of the current active
func (d defaultActivityContext) ID() string {
	return d.activity.ID
}

// OID returns the ID of the current orchestration
func (d defaultActivityContext) OID() string {
	return d.oID
}

func (d defaultActivityContext) SetValue(key string, value any) {
	//TODO implement me
	panic("implement me")
}

func (d defaultActivityContext) Value(key string) any {
	//TODO implement me
	panic("implement me")
}
