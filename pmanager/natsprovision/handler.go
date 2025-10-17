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

package natsprovision

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/metaform/connector-fabric-manager/common/dmodel"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/nats-io/nats.go/jetstream"
)

type natsDeploymentHandler struct {
	natsclient.RetriableMessageProcessor[dmodel.DeploymentManifest]
}

func newNatsDeploymentHandler(
	client natsclient.MsgClient,
	monitor monitor.LogMonitor) *natsDeploymentHandler {
	return &natsDeploymentHandler{
		RetriableMessageProcessor: natsclient.RetriableMessageProcessor[dmodel.DeploymentManifest]{
			Client:     client,
			Monitor:    monitor,
			Processing: atomic.Bool{},
			Dispatcher: func(ctx context.Context, payload dmodel.DeploymentManifest) error {
				fmt.Println("Received manifest: %v", payload)
				return nil
			},
		},
	}
}

func (n *natsDeploymentHandler) Init(ctx context.Context, consumer jetstream.Consumer) error {
	go func() {
		err := n.ProcessLoop(ctx, consumer)
		if err != nil {
			n.Monitor.Warnf("Error processing message: %v", err)
		}
	}()
	return nil
}
