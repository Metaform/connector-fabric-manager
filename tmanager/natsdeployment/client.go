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
	"encoding/json"
	"sync/atomic"

	"github.com/metaform/connector-fabric-manager/common/dmodel"
	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/nats-io/nats.go/jetstream"
)

type natsDeploymentClient struct {
	natsclient.RetriableMessageProcessor[dmodel.DeploymentResponse]
}

func newNatsDeploymentClient(
	client natsclient.MsgClient,
	dispatcher deploymentCallbackDispatcher,
	monitor system.LogMonitor) *natsDeploymentClient {
	return &natsDeploymentClient{
		RetriableMessageProcessor: natsclient.RetriableMessageProcessor[dmodel.DeploymentResponse]{
			Client:     client,
			Monitor:    monitor,
			Processing: atomic.Bool{},
			Dispatcher: func(ctx context.Context, payload dmodel.DeploymentResponse) error {
				return dispatcher.Dispatch(ctx, payload)
			},
		},
	}
}

func (n *natsDeploymentClient) Init(ctx context.Context, consumer jetstream.Consumer) error {
	go func() {
		err := n.ProcessLoop(ctx, consumer)
		if err != nil {
			n.Monitor.Warnf("Error Processing message: %v", err)
		}
	}()
	return nil
}

func (n *natsDeploymentClient) Deploy(ctx context.Context, manifest dmodel.DeploymentManifest) error {
	serialized, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	_, err = n.Client.Publish(ctx, natsclient.CFMDeploymentSubject, serialized)
	if err != nil {
		return err
	}
	return nil
}
