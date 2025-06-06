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
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// The image used for testing
const natsImage = "nats:2.10-alpine"

type natsTestContainer struct {
	container testcontainers.Container
	uri       string
	client    *natsClient
}

func setupNatsContainer(ctx context.Context, bucket string) (*natsTestContainer, error) {
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

	natsClient, err := newNatsClient(uri, bucket)
	if err != nil {
		return nil, err
	}
	//return natsClient, nil

	return &natsTestContainer{
		container: container,
		uri:       uri,
		client:    natsClient,
	}, nil
}

func teardownNatsContainer(ctx context.Context, nt *natsTestContainer) error {
	if nt.client != nil {
		nt.client.Close()
	}
	if nt.container != nil {
		return nt.container.Terminate(ctx)
	}
	return nil
}
