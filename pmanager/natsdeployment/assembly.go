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

package natsdeployment

import (
	"context"
	"fmt"

	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

const (
	setupStreamKey = "setupStream"
)

type natsDeploymentServiceAssembly struct {
	streamName        string
	natsClient        *natsclient.NatsClient
	deploymentHandler *natsDeploymentHandler
	system.DefaultServiceAssembly
	processCancel context.CancelFunc
}

func NewDeploymentServiceAssembly(streamName string) system.ServiceAssembly {
	return &natsDeploymentServiceAssembly{
		streamName: streamName,
	}
}

func (a *natsDeploymentServiceAssembly) Name() string {
	return "NATs Deployment"
}

func (a *natsDeploymentServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{api.ProvisionManagerKey, natsclient.NatsClientKey}
}

func (a *natsDeploymentServiceAssembly) Init(ctx *system.InitContext) error {

	a.natsClient = ctx.Registry.Resolve(natsclient.NatsClientKey).(*natsclient.NatsClient)

	natsContext := context.Background()
	defer natsContext.Done()

	provisionManager := ctx.Registry.Resolve(api.ProvisionManagerKey).(api.ProvisionManager)
	client := natsclient.NewMsgClient(a.natsClient)
	a.deploymentHandler = newNatsDeploymentHandler(client, provisionManager, ctx.LogMonitor)

	return nil
}

func (a *natsDeploymentServiceAssembly) Start(_ *system.StartContext) error {
	var ctx context.Context
	natsContext := context.Background()
	defer natsContext.Done()

	stream, err := natsclient.SetupStream(natsContext, a.natsClient, a.streamName)
	if err != nil {
		return fmt.Errorf("error initializing NATS stream: %w", err)
	}

	consumer, err := natsclient.SetupConsumer(natsContext, stream, natsclient.CFMDeployment)
	if err != nil {
		return fmt.Errorf("error initializing NATS deployment manifest consumer: %w", err)
	}

	ctx, a.processCancel = context.WithCancel(context.Background())
	return a.deploymentHandler.Init(ctx, consumer)
}

func (a *natsDeploymentServiceAssembly) Shutdown() error {
	if a.processCancel != nil {
		a.processCancel()
	}
	if a.natsClient != nil {
		a.natsClient.Connection.Close()
	}
	return nil
}
