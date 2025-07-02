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

package natstestutil

import (
	"context"
	"fmt"
	"github.com/metaform/connector-fabric-manager/pmanager/natsclient"
	"github.com/metaform/connector-fabric-manager/pmanager/natsorchestration"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

// The image used for testing
const natsImage = "nats:2.10-alpine"

type NatsTestContainer struct {
	Container testcontainers.Container
	Uri       string
	Client    *natsclient.NatsClient
}

func SetupNatsContainer(ctx context.Context, bucket string) (*NatsTestContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        natsImage,
		ExposedPorts: []string{"0:4222/tcp"},
		WaitingFor:   wait.ForLog("Server is ready"),
		Cmd: []string{
			"-js", // Enable JetStream
			"-DV", // Debug and trace
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "4222")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("nats://%s:%s", hostIP, mappedPort.Port())

	natsClient, err := natsclient.NewNatsClient(uri, bucket)
	if err != nil {
		return nil, err
	}
	//return NatsClient, nil

	return &NatsTestContainer{
		Container: container,
		Uri:       uri,
		Client:    natsClient,
	}, nil
}

func TeardownNatsContainer(ctx context.Context, nt *NatsTestContainer) {
	if nt.Client != nil {
		nt.Client.Close()
	}
	if nt.Container != nil {
		err := nt.Container.Terminate(ctx)
		if err != nil {
			fmt.Println("Error terminating container: ", err)
		}
	}
}

func SetupStream(t *testing.T, ctx context.Context, client *natsclient.NatsClient, streamName string) jetstream.Stream {
	stream, err := natsorchestration.SetupStream(ctx, client, streamName)
	require.NoError(t, err)
	return stream
}

func SetupConsumer(t *testing.T, ctx context.Context, stream jetstream.Stream, subject string) jetstream.Consumer {
	consumer, err := natsorchestration.SetupConsumer(ctx, stream, subject)
	require.NoError(t, err)
	return consumer
}
