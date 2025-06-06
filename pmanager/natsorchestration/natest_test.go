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
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewJetStreamPubSub(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	nt, err := setupNatsContainer(ctx, "test-bucket")
	require.NoError(t, err)
	//goland:noinspection GoUnhandledErrorResult
	defer teardownNatsContainer(ctx, nt)

	cfg := jetstream.StreamConfig{
		Name:      "test-activity",
		Retention: jetstream.WorkQueuePolicy,
		Subjects:  []string{"event.*"},
	}

	stream, err := nt.client.jetStream.CreateOrUpdateStream(ctx, cfg)
	require.NoError(t, err)

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "test-activity",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	require.NoError(t, err)

	_, err = nt.client.jetStream.Publish(ctx, "event.foo", []byte("Test message"))
	require.NoError(t, err)

	msg, err := consumer.Next(jetstream.FetchMaxWait(100 * time.Millisecond))

	require.NoError(t, err)
	require.NotNil(t, msg)

	cancel()
}
