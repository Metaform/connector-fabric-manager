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
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"time"
)

const (
	defaultDuration = 20 * time.Second
	defaultPings    = 5
	forever         = -1
)

type NatsClient struct {
	Connection *nats.Conn
	JetStream  jetstream.JetStream
	KVStore    jetstream.KeyValue
}

// NewNatsClient creates and returns a new NatsClient instance connected to the specified URL and bucket with given options.
// If options are not provided, default Connection settings are used for the NATS Client configuration.
// Returns an error if the Connection to NATS or JetStream initialization fails.
func NewNatsClient(url string, bucket string, options ...nats.Option) (*NatsClient, error) {
	if options == nil || len(options) == 0 {
		options = []nats.Option{nats.PingInterval(defaultDuration),
			nats.MaxPingsOutstanding(defaultPings),
			nats.ReconnectWait(time.Second),
			nats.RetryOnFailedConnect(true),
			nats.MaxReconnects(forever)}
	}
	connection, err := nats.Connect(url, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	jetStream, err := jetstream.New(connection)
	if err != nil {
		connection.Close()
		return nil, fmt.Errorf("failed to create jetstream context: %w", err)
	}

	kvManager, err := jetStream.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{Bucket: bucket})
	if err != nil {
		connection.Close()
		return nil, fmt.Errorf("failed to create jetstream key value manager: %w", err)
	}

	return &NatsClient{
		Connection: connection,
		JetStream:  jetStream,
		KVStore:    kvManager,
	}, nil
}

// Close closes the NATS Connection.
func (nc *NatsClient) Close() {
	if nc.Connection != nil {
		nc.Connection.Close()
	}
}
