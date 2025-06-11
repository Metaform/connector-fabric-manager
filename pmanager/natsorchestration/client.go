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
	"github.com/nats-io/nats.go/jetstream"
	"strings"
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
		subject := "event." + strings.ReplaceAll(activity.Type, ".", "-")
		_, err = client.Publish(ctx, subject, payload)
		if err != nil {
			return fmt.Errorf("error publishing to stream: %w", err)
		}
	}
	return nil
}
