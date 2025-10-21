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

	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

const (
	setupStreamKey = "setupStream"
)

type natsOrchestratorServiceAssembly struct {
	uri        string
	bucket     string
	streamName string
	natsClient *natsclient.NatsClient
	system.DefaultServiceAssembly
	processCancel context.CancelFunc
}

func NewOrchestratorServiceAssembly(uri string, bucket string, streamName string) system.ServiceAssembly {
	return &natsOrchestratorServiceAssembly{
		uri:        uri,
		bucket:     bucket,
		streamName: streamName,
	}
}

func (a *natsOrchestratorServiceAssembly) Name() string {
	return "NATs Provisioning"
}

func (a *natsOrchestratorServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{api.DeploymentOrchestratorKey, natsclient.NatsClientKey}
}

func (a *natsOrchestratorServiceAssembly) Init(ctx *system.InitContext) error {
	natsClient, err := natsclient.NewNatsClient(a.uri, a.bucket)
	if err != nil {
		return err
	}

	a.natsClient = natsClient
	ctx.Registry.Register(natsclient.NatsClientKey, natsClient)

	natsContext := context.Background()
	defer natsContext.Done()

	setupStream := true
	if ctx.Config.IsSet(setupStreamKey) {
		setupStream = ctx.Config.GetBool(setupStreamKey)
	}

	if setupStream {
		_, err = natsclient.SetupStream(natsContext, natsClient, a.streamName)
		if err != nil {
			return fmt.Errorf("error initializing NATS stream: %w", err)
		}
	}

	ctx.Registry.Register(api.DeploymentOrchestratorKey, NewNatsDeploymentOrchestrator(natsclient.NewMsgClient(natsClient), ctx.LogMonitor))

	return nil
}

func (a *natsOrchestratorServiceAssembly) Shutdown() error {
	if a.processCancel != nil {
		a.processCancel()
	}
	if a.natsClient != nil {
		a.natsClient.Connection.Close()
	}
	return nil
}
