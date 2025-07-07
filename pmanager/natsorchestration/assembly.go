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
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsclient"
)

type natsOrchestratorServiceAssembly struct {
	uri        string
	bucket     string
	streamName string
	natsClient *natsclient.NatsClient
	system.DefaultServiceAssembly
}

func NewOrchestratorServiceAssembly(uri string, bucket string, streamName string) system.ServiceAssembly {
	return &natsOrchestratorServiceAssembly{
		uri:        uri,
		bucket:     bucket,
		streamName: streamName,
	}
}

func (a *natsOrchestratorServiceAssembly) Name() string {
	return "NATs Deployment OrchestratorKey"
}

func (a *natsOrchestratorServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{api.DeploymentOrchestratorKey}
}

func (a *natsOrchestratorServiceAssembly) Init(context *system.InitContext) error {
	natsClient, err := natsclient.NewNatsClient(a.uri, a.bucket)
	if err != nil {
		return err
	}

	a.natsClient = natsClient

	context.Registry.Register(api.DeploymentOrchestratorKey, &NatsDeploymentOrchestrator{
		Client:  NatsClientAdapter{Client: natsClient},
		Monitor: context.LogMonitor,
	})
	return nil
}

func (a *natsOrchestratorServiceAssembly) Shutdown() error {
	if a.natsClient != nil {
		a.natsClient.Connection.Close()
	}
	return nil
}
