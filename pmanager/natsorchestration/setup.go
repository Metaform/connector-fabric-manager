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
	"github.com/metaform/connector-fabric-manager/pmanager/natsclient"
	"github.com/nats-io/nats.go/jetstream"
	"strings"
)

// SetupStream configures a JetStream stream configured for activity messages.
func SetupStream(ctx context.Context, client *natsclient.NatsClient, streamName string) (jetstream.Stream, error) {
	cfg := jetstream.StreamConfig{
		Name:      streamName,
		Retention: jetstream.WorkQueuePolicy,
		Subjects:  []string{ActivitySubjectPrefix + ".*"},
	}

	return client.JetStream.CreateOrUpdateStream(ctx, cfg)
}

// SetupConsumer creates or updates a NATS JetStream consumer for an activity processor.
func SetupConsumer(ctx context.Context, stream jetstream.Stream, subject string) (jetstream.Consumer, error) {
	sanitizedSubject := strings.ReplaceAll(subject, ".", "-") // convert to `-` because NATs uses dot-notation to denote subject hierarchies
	return stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       sanitizedSubject,
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: "event." + sanitizedSubject,
	})
}
