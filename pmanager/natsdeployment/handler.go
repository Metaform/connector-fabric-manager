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
	"sync/atomic"

	"github.com/metaform/connector-fabric-manager/common/dmodel"
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/nats-io/nats.go/jetstream"
)

type natsDeploymentHandler struct {
	natsclient.RetriableMessageProcessor[dmodel.DeploymentManifest]
}

func newNatsDeploymentHandler(
	client natsclient.MsgClient,
	provisionManager api.ProvisionManager,
	monitor monitor.LogMonitor) *natsDeploymentHandler {
	return &natsDeploymentHandler{
		RetriableMessageProcessor: natsclient.RetriableMessageProcessor[dmodel.DeploymentManifest]{
			Client:     client,
			Monitor:    monitor,
			Processing: atomic.Bool{},
			Dispatcher: func(ctx context.Context, payload dmodel.DeploymentManifest) error {
				fmt.Println("Received manifest: %v", payload)
				_, err := provisionManager.Start(context.Background(), &payload)
				if err != nil {
					fmt.Println("Error, %v", err)
					switch {
					case model.IsClientError(err):
						// return error response
					case model.IsRecoverable(err):
						// return natsclient.NakError(, err)
					case model.IsFatal(err):
						// return error response
					default:
						// return error response
					}
				}
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
