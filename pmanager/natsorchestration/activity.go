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

//go:generate mockery --name msgClient --filename msg_client_mock.go --with-expecter --outpkg mocks --dir . --output ./mocks

// Package natsorchestration provides a NATS-based deployment orchestrator.
package natsorchestration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/nats-io/nats.go/jetstream"
	"strings"
	"time"
)

const (
	streamName = "cfm-activity"
)

// MsgClient is an interface for interacting with NATS. This interface is used to allow for mocking in unit tests that
// verify correct behavior in response to error conditions (i.e., negative tests).
type MsgClient interface {
	Update(ctx context.Context, key string, value []byte, version uint64) (uint64, error)
	Stream(ctx context.Context, streamName string) (jetstream.Stream, error)
	Get(ctx context.Context, key string) (jetstream.KeyValueEntry, error)
	Publish(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error)
}

// Wraps the natsClient to satisfy the MsgClient interface.
type natsClientAdapter struct {
	client *natsClient
}

func (a natsClientAdapter) Update(ctx context.Context, key string, value []byte, version uint64) (uint64, error) {
	return a.client.kvStore.Update(ctx, key, value, version)
}

func (a natsClientAdapter) Stream(ctx context.Context, streamName string) (jetstream.Stream, error) {
	return a.client.jetStream.Stream(ctx, streamName)
}

func (a natsClientAdapter) Get(ctx context.Context, key string) (jetstream.KeyValueEntry, error) {
	return a.client.kvStore.Get(ctx, key)
}

func (a natsClientAdapter) Publish(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	return a.client.jetStream.Publish(ctx, subject, payload, opts...)
}

// NatsDeploymentOrchestrator is responsible for executing an orchestration using NATS for reliable messaging. For each
// activity, a message is published to a durable queue based on the activity type. Activity messages are then dequeued
// and reliably processed by a NatsActivityExecutor that handles the activity type.
type NatsDeploymentOrchestrator struct {
	client  MsgClient
	monitor monitor.LogMonitor
}

// ExecuteOrchestration asynchronously executes the given orchestration by dispatching messages to durable activity
// queues, where they can be dequeued and reliably processed by NatsActivityExecutors.
//
// A Jetstream KV entry is used to maintain durable state and is updated as the orchestration progresses. This
// state is passed to the executors, which access and update it.

func (o *NatsDeploymentOrchestrator) ExecuteOrchestration(ctx context.Context, orchestration api.Orchestration) error {
	// TODO validate orchestration - this should include a check to see if there are no steps or steps with no activities

	serializedOrchestration, err := json.Marshal(orchestration)
	if err != nil {
		return fmt.Errorf("error marshalling orchestration: %w", err)
	}

	// Use update to check if the orchestration already exists
	_, err = o.client.Update(ctx, orchestration.ID, serializedOrchestration, 0)
	if err != nil {
		var jsErr *jetstream.APIError
		if errors.As(err, &jsErr) {
			if jsErr.APIError().ErrorCode == jetstream.JSErrCodeStreamWrongLastSequence {
				// Orchestration already exists, return
				return nil
			}
		}
		return fmt.Errorf("error storing orchestration: %w", err)
	}

	activities, parallel := getInitialActivities(orchestration)
	err = enqueueActivityMessages(ctx, orchestration.ID, activities, parallel, o.client)
	if err != nil {
		return err
	}
	return nil

}

// Enqueues the given activities for processing.
//
// Messages are sent to a named durable queue corresponding to the activity type. For example, messages for the
// 'test-activity' type will be routed to the 'event.test-activity' queue.
func enqueueActivityMessages(ctx context.Context, orchestrationID string, activities []api.Activity, parallel bool, client MsgClient) error {
	for _, activity := range activities {
		// route to queue
		payload, err := json.Marshal(activityMessage{
			OrchestrationID: orchestrationID,
			Activity:        activity,
			Parallel:        parallel,
		})
		if err != nil {
			return fmt.Errorf("error marshalling activity payload: %w", err)
		}

		// Strip out periods since they denote a subject hierarchy for NATS
		subject := "event." + strings.ReplaceAll(activity.Type, ".", "-")
		_, err = client.Publish(ctx, subject, payload)
		if err != nil {
			return fmt.Errorf("error publishing to stream: %w", err)
		}
	}
	return nil
}

// Returns the initial activities for the given orchestration.
//
// If the orchestration's first step is parallel, all contained activities are returned. Otherwise, the first activity is
// returned. If the orchestration has no activities, an empty list is returned.
func getInitialActivities(orchestration api.Orchestration) ([]api.Activity, bool) {
	for _, step := range orchestration.Steps {
		if step.Parallel {
			if len(step.Activities) > 0 {
				return step.Activities, true
			}
		} else {
			if len(step.Activities) > 0 {
				return step.Activities[0:1], false
			}
		}
	}
	return []api.Activity{}, false
}

// Message sent to the activity queue.
type activityMessage struct {
	OrchestrationID string       `json:"orchestrationID"`
	Activity        api.Activity `json:"activity"`
	Parallel        bool         `json:"parallel"`
}

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

	consumer, err := stream.Consumer(ctx, "cfm-durable-activity")
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
	var oMessage activityMessage
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
	advance, err := e.activityProcessor.Process(activityContext)
	if err != nil {
		return fmt.Errorf("failed to process activity %s: %w", oMessage.OrchestrationID, err)
	}
	if !advance {
		return nil
	}

	//re-read the orchestration and update it
	orchestration, oRevision, err = e.readOrchestration(ctx, oMessage.OrchestrationID)
	if err != nil {
		return fmt.Errorf("failed to read orchestration data for update: %w", err)
	}
	orchestration.Completed[oMessage.Activity.ID] = struct{}{}
	orchestration, oRevision, err = e.saveOrchestration(ctx, orchestration, oMessage.Activity.ID, oMessage.OrchestrationID, oRevision)
	if err != nil {
		return err
	}

	if oMessage.Parallel {
		canProceed, err := orchestration.CanProceedToNextActivity(oMessage.Activity.ID, func(completed []string) bool {
			for _, id := range completed {
				if _, exists := orchestration.Completed[id]; !exists {
					return false
				}
			}
			return true
		})
		if err != nil {
			return fmt.Errorf("failed to proceed with orchestration %s: %w", oMessage.OrchestrationID, err)
		}

		if !canProceed {
			return nil
		}
	}
	next, parallel := orchestration.GetNextActivities(oMessage.Activity.ID)
	if len(next) == 0 {
		e.monitor.Debugf("Finished orchestration: %s", oMessage.OrchestrationID)
		orchestration.State = api.OrchestrationStateStateCompleted
		orchestration, oRevision, err = e.saveOrchestration(ctx, orchestration, oMessage.Activity.ID, oMessage.OrchestrationID, oRevision)
		if err != nil {
			return err
		}
		return nil
	}
	err = enqueueActivityMessages(ctx, orchestration.ID, next, parallel, e.client)
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
