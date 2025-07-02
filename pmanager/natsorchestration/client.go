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

package natsorchestration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsclient"
	"github.com/nats-io/nats.go/jetstream"
	"strings"
)

const ActivitySubjectPrefix = "event"

// MsgClient is an interface for interacting with NATS. This interface is used to allow for mocking in unit tests that
// verify correct behavior in response to error conditions (i.e., negative tests).
type MsgClient interface {
	Update(ctx context.Context, key string, value []byte, version uint64) (uint64, error)
	Stream(ctx context.Context, streamName string) (jetstream.Stream, error)
	Get(ctx context.Context, key string) (jetstream.KeyValueEntry, error)
	Publish(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error)
}

// Wraps the NatsClient to satisfy the MsgClient interface.
type NatsClientAdapter struct {
	Client *natsclient.NatsClient
}

func (a NatsClientAdapter) Update(ctx context.Context, key string, value []byte, version uint64) (uint64, error) {
	return a.Client.KVStore.Update(ctx, key, value, version)
}

func (a NatsClientAdapter) Stream(ctx context.Context, streamName string) (jetstream.Stream, error) {
	return a.Client.JetStream.Stream(ctx, streamName)
}

func (a NatsClientAdapter) Get(ctx context.Context, key string) (jetstream.KeyValueEntry, error) {
	return a.Client.KVStore.Get(ctx, key)
}

func (a NatsClientAdapter) Publish(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	return a.Client.JetStream.Publish(ctx, subject, payload, opts...)
}

// EnqueueActivityMessages enqueues the given activities for processing.
//
// Messages are sent to a named durable queue corresponding to the activity type. For example, messages for the
// 'test-activity' type will be routed to the 'event.test-activity' queue.
func EnqueueActivityMessages(ctx context.Context, orchestrationID string, activities []api.Activity, client MsgClient) error {
	for _, activity := range activities {
		// route to queue
		payload, err := json.Marshal(api.ActivityMessage{
			OrchestrationID: orchestrationID,
			Activity:        activity,
		})
		if err != nil {
			return fmt.Errorf("error marshalling activity payload: %w", err)
		}

		// Strip out periods since they denote a subject hierarchy for NATS
		subject := ActivitySubjectPrefix + "." + strings.ReplaceAll(activity.Type, ".", "-")
		_, err = client.Publish(ctx, subject, payload)
		if err != nil {
			return fmt.Errorf("error publishing to stream: %w", err)
		}
	}
	return nil
}

// ReadOrchestration reads the orchestration state from the KV store.
func ReadOrchestration(ctx context.Context, orchestrationID string, client MsgClient) (api.Orchestration, uint64, error) {
	oEntry, err := client.Get(ctx, orchestrationID)
	if err != nil {
		return api.Orchestration{}, 0, fmt.Errorf("failed to get orchestration state %s: %w", orchestrationID, err)
	}

	var orchestration api.Orchestration
	if err = json.Unmarshal(oEntry.Value(), &orchestration); err != nil {
		return api.Orchestration{}, 0, fmt.Errorf("failed to unmarshal orchestration %s: %w", orchestrationID, err)
	}

	return orchestration, oEntry.Revision(), nil
}

// UpdateOrchestration updates the orchestration state in the KV store using optimistic concurrency by comparing the
// last known revision.
func UpdateOrchestration(
	ctx context.Context,
	orchestration api.Orchestration,
	revision uint64,
	client MsgClient,
	updateFn func(*api.Orchestration)) (api.Orchestration, uint64, error) {
	for {
		updateFn(&orchestration)
		// TODO break after number of retries using exponential backoff
		serialized, err := json.Marshal(orchestration)
		if err != nil {
			return api.Orchestration{}, 0, fmt.Errorf("failed to marshal orchestration %s: %w", orchestration.ID, err)
		}
		_, err = client.Update(ctx, orchestration.ID, serialized, revision)
		if err == nil {
			break
		}
		orchestration, revision, err = ReadOrchestration(ctx, orchestration.ID, client)
		if err != nil {
			return api.Orchestration{}, 0, fmt.Errorf("failed to read orchestration data for update: %w", err)
		}
	}
	return orchestration, revision, nil
}
