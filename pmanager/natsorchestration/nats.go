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
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"time"
)

const (
	defaultDuration = 20 * time.Second
	defaultPings    = 5
	forever         = -1
)

type natsClient struct {
	connection *nats.Conn
	jetStream  jetstream.JetStream
	kvStore    jetstream.KeyValue
}

// newNatsClient creates and returns a new natsClient instance connected to the specified URL and bucket with given options.
// If options are not provided, default connection settings are used for the NATS client configuration.
// Returns an error if the connection to NATS or JetStream initialization fails.
func newNatsClient(url string, bucket string, options ...nats.Option) (*natsClient, error) {
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

	return &natsClient{
		connection: connection,
		jetStream:  jetStream,
		kvStore:    kvManager,
	}, nil
}

// Close closes the NATS connection.
func (nc *natsClient) Close() {
	if nc.connection != nil {
		nc.connection.Close()
	}
}
