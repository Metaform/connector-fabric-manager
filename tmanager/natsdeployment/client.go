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
	"fmt"
	"time"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/dmodel"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/nats-io/nats.go/jetstream"
)

type natsDeploymentClient struct {
	client     *natsclient.NatsClient
	dispatcher api.DeploymentCallbackDispatcher
	monitor    monitor.LogMonitor
}

func (n natsDeploymentClient) Init(ctx context.Context, consumer jetstream.Consumer) error {
	go func() {
		err := n.processLoop(ctx, consumer)
		if err != nil {
			n.monitor.Warnf("Error processing message: %v", err)
		}
	}()
	return nil
}

func (n natsDeploymentClient) Deploy(ctx context.Context, manifest dmodel.DeploymentManifest) error {
	_, err := n.client.JetStream.Publish(ctx, "subject", nil)
	if err != nil {
		return err
	}
	return nil
}

// processLoop handles the main loop for consuming and processing messages from a JetStream consumer.
// It runs continuously until the provided context is canceled or an error occurs.
// Returns an error if message fetching or processing fails.
func (n natsDeploymentClient) processLoop(ctx context.Context, consumer jetstream.Consumer) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			messageBatch, err := consumer.Fetch(1, jetstream.FetchMaxWait(time.Second))
			if err != nil {
				return err
			}

			for message := range messageBatch.Messages() {
				if err = n.processMessage(ctx, message); err != nil {
					n.monitor.Warnf("Error processing deployment message: %v", err)
				}
			}
		}
	}
}

func (n natsDeploymentClient) processMessage(ctx context.Context, message jetstream.Msg) error {
	var dResponse api.DeploymentResponse
	if err := json.Unmarshal(message.Data(), &dResponse); err != nil {
		err := n.ackMessage(dResponse.ID, message)
		if err != nil {
			n.monitor.Warnf("Failed to ACK message %s: %v", dResponse.ID, err)
		}
		return fmt.Errorf("failed to unmarshal deployment response message: %w", err)
	}

	n.monitor.Debugf("Received deployment response %s for %s", dResponse.ID, dResponse.ManifestID)
	resultErr := n.dispatcher.Dispatch(ctx, dResponse)
	if resultErr == nil {
		return n.ackMessage(dResponse.ID, message)
	}

	switch {
	case model.IsRecoverable(resultErr):
		if err := message.Nak(); err != nil {
			return fmt.Errorf("retriable failure when dispatching deployment response message and NAK response %s (errors: %w, %v)",
				dResponse.ID, resultErr, err)
		}
		return fmt.Errorf("retriable failure when dispatching deployment response %s: %w", dResponse.ID, resultErr)
	default:
		// All other errors are fatal
		if err := message.Ack(); err != nil {
			return fmt.Errorf("fatal failure when dispatching deployment response %s (errors: %w, %v)",
				dResponse.ID, resultErr, err)
		}
		return fmt.Errorf("fatal failure when dispatching deployment response %s: %w", dResponse.ID, resultErr)
	}
}

func (n natsDeploymentClient) ackMessage(id string, message jetstream.Msg) error {
	if err := message.Ack(); err != nil {
		return fmt.Errorf("failed to ACK activity message %s: %w", id, err)
	}
	return nil
}
