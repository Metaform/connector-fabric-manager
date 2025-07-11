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

package natsorchestration_test_test

import (
	"context"
	"github.com/metaform/connector-fabric-manager/pmanager/natsorchestration"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewJetStreamPubSub(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	nt, err := natsorchestration.SetupNatsContainer(ctx, "test-bucket")
	require.NoError(t, err)

	defer natsorchestration.TeardownNatsContainer(ctx, nt)

	stream := natsorchestration.SetupTestStream(t, ctx, nt.Client, "test-activity")

	consumer := natsorchestration.SetupTestConsumer(t, ctx, stream, "foo")

	_, err = nt.Client.JetStream.Publish(ctx, "event.foo", []byte("Test message"))
	require.NoError(t, err)

	msg, err := consumer.Next(jetstream.FetchMaxWait(100 * time.Millisecond))

	require.NoError(t, err)
	require.NotNil(t, msg)

	cancel()
}
