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

package natsclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/nats-io/nats.go/jetstream"
)

// RetriableMessageProcessor delegates to a dispatcher to process messages from a JetStream consumer and retries on failure.
type RetriableMessageProcessor[T any] struct {
	Client     MsgClient
	Dispatcher func(ctx context.Context, payload T) error
	Monitor    system.LogMonitor
	Processing atomic.Bool
}

// ProcessLoop handles the main loop for consuming and processing messages from a JetStream consumer.
// It runs continuously until the provided context is canceled or an error occurs.
// Returns an error if message fetching or processing fails.
func (n *RetriableMessageProcessor[T]) ProcessLoop(ctx context.Context, consumer jetstream.Consumer) error {
	n.Processing.Store(true)
	for {
		select {
		case <-ctx.Done():
			n.Processing.Store(false)
			return ctx.Err()
		default:
			messageBatch, err := consumer.Fetch(1, jetstream.FetchMaxWait(time.Second))
			if err != nil {
				return err
			}

			for message := range messageBatch.Messages() {
				if err = n.ProcessMessage(ctx, message); err != nil {
					n.Monitor.Warnf("Error processing received message: %v", err)
				}
			}
		}
	}
}

func (n *RetriableMessageProcessor[T]) ProcessMessage(ctx context.Context, message jetstream.Msg) error {
	var payload T
	if err := json.Unmarshal(message.Data(), &payload); err != nil {
		err2 := AckMessage(message)
		if err2 != nil {
			n.Monitor.Warnf("Failed to ACK message %s: %v", err2)
		}
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	resultErr := n.Dispatcher(ctx, payload)
	if resultErr == nil {
		return AckMessage(message)
	}

	switch {
	case types.IsRecoverable(resultErr):
		if err := message.Nak(); err != nil {
			return fmt.Errorf("retriable failure when dispatching message and NAK response (errors: %w, %v)", resultErr, err)
		}
		return fmt.Errorf("retriable failure when dispatching message: %w", resultErr)
	default:
		// All other errors are fatal
		if err := message.Ack(); err != nil {
			return fmt.Errorf("fatal failure when dispatching message (errors: %w, %v)", resultErr, err)
		}
		return fmt.Errorf("fatal failure when dispatching message: %w", resultErr)
	}
}

func AckMessage(message jetstream.Msg) error {
	if err := message.Ack(); err != nil {
		return fmt.Errorf("failed to ACK message: %w", err)
	}
	return nil
}

func NakError(message jetstream.Msg, err error) error {
	err2 := message.Nak() // Attempt redelivery
	if err2 != nil {
		err = errors.Join(err, err2)
	}
	return err
}
