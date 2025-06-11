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

	orchestration, oRevision, err := e.readOrchestration(ctx, oMessage.OrchestrationID)
	if err != nil {
		return fmt.Errorf("failed to read orchestration data: %w", err)
	}

	activityContext := newActivityContext(ctx, orchestration.ID, oMessage.Activity)
	e.monitor.Debugf("Received message: %s\n", e.id)
	result := e.activityProcessor.Process(activityContext)
	if result.Result == api.ActivityResultError {
		err := message.Nak()
		if err != nil {
			return fmt.Errorf("failed to execute activity message and NAK response %s (errors: %w, %v)",
				oMessage.OrchestrationID, result.Error, err)
		}
		return fmt.Errorf("failed to execute activity %s: %w", oMessage.OrchestrationID, result.Error)
	} else if result.Result == api.ActivityResultWait {
		err := message.Ack()
		if err != nil {
			return fmt.Errorf("failed to ACK activity message %s: %w", oMessage.OrchestrationID, err)
		}
		return err
	} else if result.Result == api.ActivityResultSchedule {
		err := message.NakWithDelay(result.WaitMillis)
		if err != nil {
			return fmt.Errorf("failed to reschedule schedule activity %s: %w", oMessage.OrchestrationID, err)
		}
		return nil
	}

	// Completed activity execution, re-read the orchestration and update it
	orchestration, oRevision, err = e.readOrchestration(ctx, oMessage.OrchestrationID)
	if err != nil {
		return fmt.Errorf("failed to read orchestration data for update: %w", err)
	}
	orchestration.Completed[oMessage.Activity.ID] = struct{}{}
	orchestration, oRevision, err = e.saveOrchestration(ctx, orchestration, oMessage.Activity.ID, oMessage.OrchestrationID, oRevision)
	if err != nil {
		return err
	}

	canProceed, err := orchestration.CanProceedToNextStep(oMessage.Activity.ID)
	if err != nil {
		return fmt.Errorf("failed to proceed with orchestration %s: %w", oMessage.OrchestrationID, err)
	}

	if !canProceed {
		return nil
	}
	next := orchestration.GetNextStepActivities(oMessage.Activity.ID)
	if len(next) == 0 {
		e.monitor.Debugf("Finished orchestration: %s", oMessage.OrchestrationID)
		orchestration.State = api.OrchestrationStateStateCompleted
		orchestration, oRevision, err = e.saveOrchestration(ctx, orchestration, oMessage.Activity.ID, oMessage.OrchestrationID, oRevision)
		if err != nil {
			return err
		}
		return nil
	}
	err = EnqueueActivityMessages(ctx, orchestration.ID, next, e.client)
	if err != nil {
		return fmt.Errorf("failed to enqueue next orchestration activities %s: %w", oMessage.OrchestrationID, err)
	}
	e.monitor.Debugf("Finished activity: %s", oMessage.Activity.Type)

	if err = message.Ack(); err != nil {
		return fmt.Errorf("failed to ACK activity message %s: %w", oMessage.OrchestrationID, err)
	}
	return nil
}

// readOrchestration reads the orchestration state from the KV store.
func (e *NatsActivityExecutor) readOrchestration(ctx context.Context, orchestrationID string) (api.Orchestration, uint64, error) {
	oEntry, err := e.client.Get(ctx, orchestrationID)
	if err != nil {
		return api.Orchestration{}, 0, fmt.Errorf("failed to get orchestration state %s: %w", orchestrationID, err)
	}

	var orchestration api.Orchestration
	if err = json.Unmarshal(oEntry.Value(), &orchestration); err != nil {
		// TODO pass to DLQ
		return api.Orchestration{}, 0, fmt.Errorf("failed to unmarshal orchestration %s: %w", orchestrationID, err)
	}

	oRevision := oEntry.Revision()
	return orchestration, oRevision, nil
}

// saveOrchestration updates the orchestration state in the KV store using optimistic concurrency by comparing the last known revision.
func (e *NatsActivityExecutor) saveOrchestration(
	ctx context.Context,
	orchestration api.Orchestration,
	completedActivityID string,
	orchestrationID string,
	revision uint64) (api.Orchestration, uint64, error) {
	for {
		// TODO break after number of retries using exponential backoff
		serialized, err := json.Marshal(orchestration)
		if err != nil {
			return api.Orchestration{}, 0, fmt.Errorf("failed to marshal orchestration %s: %w", orchestrationID, err)
		}
		_, err = e.client.Update(ctx, orchestrationID, serialized, revision)
		if err == nil {
			break
		}
		orchestration, revision, err = e.readOrchestration(ctx, orchestrationID)
		if err != nil {
			return api.Orchestration{}, 0, fmt.Errorf("failed to read orchestration data for update: %w", err)
		}
		orchestration.Completed[completedActivityID] = struct{}{}
	}
	return orchestration, revision, nil
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
