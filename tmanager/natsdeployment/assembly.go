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
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

type natsDeploymentServiceAssembly struct {
	uri              string
	bucket           string
	streamName       string
	natsClient       *natsclient.NatsClient
	deploymentClient *natsDeploymentClient
	processCancel    context.CancelFunc

	system.DefaultServiceAssembly
}

func NewNatsDeploymentServiceAssembly(uri string, bucket string, streamName string) system.ServiceAssembly {
	return &natsDeploymentServiceAssembly{
		uri:        uri,
		bucket:     bucket,
		streamName: streamName,
	}
}

func (a *natsDeploymentServiceAssembly) Name() string {
	return "NATs Deployment Client"
}

func (a *natsDeploymentServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{api.DeploymentClientKey}
}

func (d *natsDeploymentServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{}
}

func (a *natsDeploymentServiceAssembly) Init(ctx *system.InitContext) error {
	natsClient, err := natsclient.NewNatsClient(a.uri, a.bucket)
	if err != nil {
		return err
	}

	a.natsClient = natsClient

	dispatcher := newDeploymentCallbackService()
	ctx.Registry.Register(api.DeploymentHandlerRegistryKey, dispatcher)

	client := natsclient.NewMsgClient(natsClient)
	a.deploymentClient = newNatsDeploymentClient(client, dispatcher, ctx.LogMonitor)
	ctx.Registry.Register(api.DeploymentClientKey, a.deploymentClient)

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

	consumer, err := natsclient.SetupConsumer(natsContext, stream, natsclient.CFMDeploymentResponse)
	if err != nil {
		return fmt.Errorf("error initializing NATS deployment consumer: %w", err)
	}

	ctx, a.processCancel = context.WithCancel(context.Background())

	return a.deploymentClient.Init(ctx, consumer)
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
