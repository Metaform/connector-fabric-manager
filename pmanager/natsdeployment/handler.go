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

	"github.com/google/uuid"
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
			Dispatcher: func(ctx context.Context, manifest dmodel.DeploymentManifest) error {
				_, err := provisionManager.Start(ctx, &manifest)
				if err != nil {
					switch {
					case model.IsRecoverable(err):
						// Return error to NAK the message and retry
						return err
					default:
						// return error response
						m := &dmodel.DeploymentResponse{
							ID:             uuid.New().String(),
							Success:        false,
							ErrorDetail:    err.Error(),
							ManifestID:     manifest.ID,
							DeploymentType: manifest.DeploymentType,
							Properties:     make(map[string]any),
						}
						ser, err := json.Marshal(m)
						if err != nil {
							return model.NewRecoverableError("failed to marshal response: %s", err.Error())
						}
						_, err = client.Publish(ctx, natsclient.CFMDeploymentResponseSubject, ser)
						if err != nil {
							return model.NewRecoverableError("failed to publish response: %s", err.Error())
						}

						return nil // ack message back
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
