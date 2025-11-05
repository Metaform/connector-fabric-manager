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
	"errors"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go/jetstream"
)

const CFMSubjectPrefix = "event"
const CFMOrchestration = "cfm-orchestration"
const CFMOrchestrationSubject = CFMSubjectPrefix + "." + CFMOrchestration
const CFMOrchestrationResponse = "cfm-orchestration-response"
const CFMOrchestrationResponseSubject = CFMSubjectPrefix + "." + CFMOrchestrationResponse

// SetupStream configures a JetStream stream used for component messaging. If the stream does not exist, it is created.
func SetupStream(ctx context.Context, client *NatsClient, streamName string) (jetstream.Stream, error) {
	stream, err := client.JetStream.Stream(ctx, streamName)
	if err == nil {
		return stream, nil
	}

	// If stream doesn't exist, create it
	if errors.Is(err, jetstream.ErrStreamNotFound) {
		cfg := jetstream.StreamConfig{
			Name:      streamName,
			Retention: jetstream.WorkQueuePolicy,
			Subjects:  []string{CFMSubjectPrefix + ".*"},
		}
		return client.JetStream.CreateOrUpdateStream(ctx, cfg)
	}

	return nil, fmt.Errorf("unable to access NATS stream: %w", err)
}

// SetupConsumer creates or updates a NATS JetStream consumer for an activity processor.
func SetupConsumer(ctx context.Context, stream jetstream.Stream, subject string) (jetstream.Consumer, error) {
	sanitizedSubject := strings.ReplaceAll(subject, ".", "-") // convert to `-` because NATs uses dot-notation to denote subject hierarchies
	return stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       sanitizedSubject,
		AckPolicy:     jetstream.AckExplicitPolicy,
		FilterSubject: CFMSubjectPrefix + "." + sanitizedSubject,
	})
}
