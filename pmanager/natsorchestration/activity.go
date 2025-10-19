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
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/dmodel"
	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsActivityExecutor struct {
	Client            natsclient.MsgClient
	StreamName        string
	ActivityType      string
	ActivityProcessor api.ActivityProcessor
	Monitor           system.LogMonitor
}

// Execute starts a goroutine to process messages from the activity queue.
func (e *NatsActivityExecutor) Execute(ctx context.Context) error {
	stream, err := e.Client.Stream(ctx, e.StreamName)
	if err != nil {
		return fmt.Errorf("error opening stream: %w", err)
	}

	consumerName := strings.ReplaceAll(e.ActivityType, ".", "-")
	consumer, err := stream.Consumer(ctx, consumerName)
	if err != nil {
		return fmt.Errorf("error connecting to consumer %s: %w", consumerName, err)
	}

	go func() {
		err := e.processLoop(ctx, consumer)
		if err != nil {
			e.Monitor.Warnf("Error processing message: %v", err)
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
				return err
			}

			for message := range messageBatch.Messages() {
				if err = e.processMessage(ctx, message); err != nil {
					e.Monitor.Warnf("Error processing message: %v", err)
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
		err := natsclient.AckMessage(message)
		if err != nil {
			e.Monitor.Warnf("Failed to ACK message: %v", err)
		}
		return fmt.Errorf("failed to unmarshal orchestration message: %w", err)
	}

	orchestration, revision, err := readOrchestration(ctx, oMessage.OrchestrationID, e.Client)
	if err != nil {
		return fmt.Errorf("failed to read orchestration data: %w", err)
	}

	activityContext := newActivityContext(ctx, orchestration.ID, oMessage.Activity)
	e.Monitor.Debugf("Received activity message %s for orchestration %s", oMessage.Activity.ID, oMessage.OrchestrationID)
	result := e.ActivityProcessor.Process(activityContext)

	switch result.Result {
	case api.ActivityResultRetryError:
		return e.handleRetryError(activityContext, orchestration, revision, message, result.Error)

	case api.ActivityResultFatalError:
		return e.handleFatalError(activityContext, orchestration, revision, result.Error, message)

	case api.ActivityResultWait:
		e.persistState(activityContext, orchestration, revision)
		return natsclient.AckMessage(message)

	case api.ActivityResultSchedule:
		e.persistState(activityContext, orchestration, revision)
		if err := message.NakWithDelay(result.WaitMillis); err != nil {
			return fmt.Errorf("failed to reschedule schedule activity %s: %w", oMessage.OrchestrationID, err)
		}
		return nil
	}

	return e.processOnActivityCompletion(activityContext, orchestration, revision, message, oMessage)
}

func (e *NatsActivityExecutor) persistState(activityContext api.ActivityContext, orchestration api.Orchestration, revision uint64) {
	if _, _, err := updateOrchestration(activityContext.Context(), orchestration, revision, e.Client, func(o *api.Orchestration) {
		for key, value := range activityContext.Values() {
			orchestration.ProcessingData[key] = value
		}
	}); err != nil {
		e.Monitor.Warnf("Failed to persist orchestration state for %s: %v", orchestration.ID, err)
	}
}

func (e *NatsActivityExecutor) processOnActivityCompletion(
	activityContext api.ActivityContext,
	orchestration api.Orchestration,
	revision uint64,
	message jetstream.Msg,
	oMessage api.ActivityMessage) error {

	// The orchestration state must be saved and re-read to determine if activities completed after the last read and the orchestration is complete.
	orchestration, revision, err := updateOrchestration(activityContext.Context(), orchestration, revision, e.Client, func(o *api.Orchestration) {
		for key, value := range activityContext.Values() {
			orchestration.ProcessingData[key] = value
		}
		o.Completed[oMessage.Activity.ID] = struct{}{} // Mark current activity as completed
	})
	if err != nil {
		err = natsclient.NakError(message, err)
		return err
	}

	// Return if orchestration is in the error state since processing should stop
	if orchestration.State == api.OrchestrationStateErrored {
		return natsclient.AckMessage(message)
	}

	// Check if all parallel activities have completed and the orchestration can continue to the next step
	canProceed, err := orchestration.CanProceedToNextStep(oMessage.Activity.ID)
	if err != nil {
		err = natsclient.NakError(message, err)
		return fmt.Errorf("failed to proceed with orchestration %s: %w", oMessage.OrchestrationID, err)
	}

	if !canProceed {
		// Waiting for parallel activities to complete
		return natsclient.AckMessage(message)
	}

	next := orchestration.GetNextStepActivities(oMessage.Activity.ID)
	if len(next) == 0 {
		// No more steps, mark as completed
		return e.handleOrchestrationCompletion(activityContext, orchestration, revision, message)
	}

	// Enqueue next activities
	if err := enqueueActivityMessages(activityContext.Context(), orchestration.ID, next, e.Client); err != nil {
		// Failed redeliver the message
		err = natsclient.NakError(message, err)
		return fmt.Errorf("failed to enqueue next orchestration activities %s: %w", oMessage.OrchestrationID, err)
	}

	return natsclient.AckMessage(message)
}

func (e *NatsActivityExecutor) handleOrchestrationCompletion(
	activityContext api.ActivityContext,
	orchestration api.Orchestration,
	revision uint64,
	message jetstream.Msg) error {
	// Mark as completed
	_, _, err := updateOrchestration(activityContext.Context(), orchestration, revision, e.Client, func(o *api.Orchestration) {
		o.State = api.OrchestrationStateCompleted
	})
	if err != nil {
		// Error marking, redeliver the message
		err = natsclient.NakError(message, err)
		return fmt.Errorf("failed to mark orchestration %s as completed: %v", orchestration.ID, err)
	}

	err = e.publishResponse(activityContext, orchestration)
	if err != nil {
		return err
	}

	return natsclient.AckMessage(message)
}

func (e *NatsActivityExecutor) publishResponse(activityContext api.ActivityContext, orchestration api.Orchestration) error {
	dr := &dmodel.DeploymentResponse{
		ID:             uuid.New().String(),
		ManifestID:     orchestration.ID,
		Success:        true,
		DeploymentType: orchestration.DeploymentType,
		Properties:     make(map[string]any),
	}
	ser, err := json.Marshal(dr)
	if err != nil {
		return fmt.Errorf("failed to marshal deployment response: %w", err)
	}
	_, err = e.Client.Publish(activityContext.Context(), natsclient.CFMDeploymentResponseSubject, ser)
	return err
}

// handleRetryError handles retriable errors by persisting the orchestration state and re-delivering the message using a Nak.
func (e *NatsActivityExecutor) handleRetryError(
	activityContext api.ActivityContext,
	orchestration api.Orchestration,
	revision uint64,
	message jetstream.Msg,
	resultErr error) error {

	e.persistState(activityContext, orchestration, revision)
	// Nak to redeliver message
	if err := message.Nak(); err != nil {
		return fmt.Errorf("retriable failure when executing activity message and NAK response %s (errors: %w, %v)",
			orchestration.ID, resultErr, err)
	}
	return fmt.Errorf("retriable failure when executing activity %s: %w", orchestration.ID, resultErr)
}

// handleFatalError handles unrecoverable errors by updating the orchestration state to "Errored" and acknowledging the message.
// It ensures acknowledgments are sent to avoid message re-delivery, even if the state update fails.
// Returns an error with specific details about the fatal failure.
func (e *NatsActivityExecutor) handleFatalError(
	activityContext api.ActivityContext,
	orchestration api.Orchestration,
	revision uint64,
	resultErr error,
	message jetstream.Msg) error {
	// Update the orchestration before acking back. If the update fails, just log it to ensure the ack is sent to avoid message re-delivery
	if _, _, err := updateOrchestration(activityContext.Context(), orchestration, revision, e.Client, func(o *api.Orchestration) {
		for key, value := range activityContext.Values() {
			orchestration.ProcessingData[key] = value
		}
		o.State = api.OrchestrationStateErrored
	}); err != nil {
		e.Monitor.Warnf("Failed to mark orchestration %s as fatal: %v", orchestration.ID, err)
	}

	if err := message.Ack(); err != nil {
		return fmt.Errorf("fatal failure while executing activity %s (errors: %w, %v)",
			orchestration.ID, resultErr, err)
	}
	return fmt.Errorf("fatal failure while executing activity %s: %w", orchestration.ID, resultErr)
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
	d.data[key] = value
}

func (d defaultActivityContext) Value(key string) (any, bool) {
	value, ok := d.data[key]
	return value, ok
}

func (d defaultActivityContext) Values() map[string]any {
	return d.data
}
