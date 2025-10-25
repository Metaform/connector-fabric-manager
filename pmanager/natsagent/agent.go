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

package natsagent

import (
	"context"
	"fmt"
	"time"

	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsorchestration"
)

const (
	timeout = 10 * time.Second
)

// AgentServiceAssembly provides common functionality for NATS-based agents
type AgentServiceAssembly struct {
	agentName    string
	activityType string
	uri          string
	bucket       string
	streamName   string
	newProcessor func(monitor system.LogMonitor) api.ActivityProcessor

	system.DefaultServiceAssembly

	natsClient *natsclient.NatsClient
	cancel     context.CancelFunc
}

func (a *AgentServiceAssembly) Name() string {
	return a.agentName
}

func (a *AgentServiceAssembly) Start(startCtx *system.StartContext) error {
	var err error
	a.natsClient, err = natsclient.NewNatsClient(a.uri, a.bucket)
	if err != nil {
		return fmt.Errorf("failed to create NATS client: %w", err)
	}

	if err = a.setupConsumer(a.natsClient); err != nil {
		return fmt.Errorf("failed to create setup agent consumer: %w", err)
	}

	executor := &natsorchestration.NatsActivityExecutor{
		Client:            natsclient.NewMsgClient(a.natsClient),
		StreamName:        a.streamName,
		ActivityType:      a.activityType,
		ActivityProcessor: a.newProcessor(startCtx.LogMonitor),
		Monitor:           startCtx.LogMonitor,
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	return executor.Execute(ctx)
}

func (a *AgentServiceAssembly) Shutdown() error {
	if a.cancel != nil {
		a.cancel()
	}

	if a.natsClient != nil {
		a.natsClient.Connection.Close()
	}
	return nil
}

func (a *AgentServiceAssembly) setupConsumer(natsClient *natsclient.NatsClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stream, err := natsclient.SetupStream(ctx, natsClient, a.streamName)

	if err != nil {
		return fmt.Errorf("error setting up agent stream: %w", err)
	}

	_, err = natsclient.SetupConsumer(ctx, stream, a.activityType)

	if err != nil {
		return fmt.Errorf("error setting up agent consumer: %w", err)
	}

	return nil
}
