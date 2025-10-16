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

//go:generate mockery

package natsclient

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

// MsgClient is an interface for interacting with NATS. This interface is used to allow for mocking in unit tests that
// verify correct behavior in response to error conditions (i.e., negative tests).
type MsgClient interface {
	Update(ctx context.Context, key string, value []byte, version uint64) (uint64, error)
	Stream(ctx context.Context, streamName string) (jetstream.Stream, error)
	Get(ctx context.Context, key string) (jetstream.KeyValueEntry, error)
	Publish(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error)
}

func NewMsgClient(nc *NatsClient) MsgClient {
	return natsClientAdapter{Client: nc}
}

// Wraps the NatsClient to satisfy the MsgClient interface.
type natsClientAdapter struct {
	Client *NatsClient
}

func (a natsClientAdapter) Update(ctx context.Context, key string, value []byte, version uint64) (uint64, error) {
	return a.Client.KVStore.Update(ctx, key, value, version)
}

func (a natsClientAdapter) Stream(ctx context.Context, streamName string) (jetstream.Stream, error) {
	return a.Client.JetStream.Stream(ctx, streamName)
}

func (a natsClientAdapter) Get(ctx context.Context, key string) (jetstream.KeyValueEntry, error) {
	return a.Client.KVStore.Get(ctx, key)
}

func (a natsClientAdapter) Publish(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) (*jetstream.PubAck, error) {
	return a.Client.JetStream.Publish(ctx, subject, payload, opts...)
}
