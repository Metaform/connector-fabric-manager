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

package common

import (
	"context"
	"fmt"

	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsorchestration"
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
}

func (b *AgentServiceAssembly) Name() string {
	return b.agentName
}

func (b *AgentServiceAssembly) Start(startCtx *system.StartContext) error {
	natsClient, err := natsclient.NewNatsClient(b.uri, b.bucket)
	if err != nil {
		return err
	}

	if err = b.setupConsumer(natsClient); err != nil {
		return err
	}

	executor := &natsorchestration.NatsActivityExecutor{
		Client:            natsclient.NewMsgClient(natsClient),
		StreamName:        b.streamName,
		ActivityType:      b.activityType,
		ActivityProcessor: b.newProcessor(startCtx.LogMonitor),
		Monitor:           system.NoopMonitor{},
	}

	ctx := context.Background()
	return executor.Execute(ctx)
}

func (b *AgentServiceAssembly) setupConsumer(natsClient *natsclient.NatsClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stream, err := natsclient.SetupStream(ctx, natsClient, b.streamName)

	if err != nil {
		return fmt.Errorf("error setting up agent stream: %w", err)
	}

	_, err = natsclient.SetupConsumer(ctx, stream, b.activityType)

	if err != nil {
		return fmt.Errorf("error setting up agent consumer: %w", err)
	}

	return nil
}
