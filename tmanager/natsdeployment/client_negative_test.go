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
	"errors"
	"testing"
	"time"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNatsDeploymentClient_ProcessMessage_Errors(t *testing.T) {
	tests := []struct {
		name            string
		messageData     []byte
		setupDispatcher func(*mockDeploymentDispatcher)
		setupMessage    func(*mockJetStreamMsg)
		expectedError   string
	}{
		{
			name:        "unmarshal error - invalid JSON",
			messageData: []byte(`invalid json`),
			setupDispatcher: func(d *mockDeploymentDispatcher) {
				// Dispatcher should not be called for invalid JSON
			},
			setupMessage: func(m *mockJetStreamMsg) {
				m.On("Data").Return([]byte(`invalid json`))
				m.On("Ack").Return(nil)
			},
			expectedError: "failed to unmarshal ",
		},
		{
			name:        "dispatcher returns recoverable error",
			messageData: []byte(`{"id":"test-id","success":true,"manifestID":"manifest-1"}`),
			setupDispatcher: func(d *mockDeploymentDispatcher) {
				d.On("Dispatch", mock.Anything, mock.MatchedBy(func(r model.DeploymentResponse) bool {
					return r.ID == "test-id"
				})).Return(types.NewRecoverableError("temporary failure"))
			},
			setupMessage: func(m *mockJetStreamMsg) {
				m.On("Data").Return([]byte(`{"id":"test-id","success":true,"manifestID":"manifest-1"}`))
				m.On("Nak").Return(nil)
			},
			expectedError: "retriable failure when dispatching ",
		},
		{
			name:        "dispatcher returns recoverable error and NAK fails",
			messageData: []byte(`{"id":"test-id-2","success":true,"manifestID":"manifest-2"}`),
			setupDispatcher: func(d *mockDeploymentDispatcher) {
				d.On("Dispatch", mock.Anything, mock.MatchedBy(func(r model.DeploymentResponse) bool {
					return r.ID == "test-id-2"
				})).Return(types.NewRecoverableError("temporary failure"))
			},
			setupMessage: func(m *mockJetStreamMsg) {
				m.On("Data").Return([]byte(`{"id":"test-id-2","success":true,"manifestID":"manifest-2"}`))
				m.On("Nak").Return(errors.New("NAK failed"))
			},
			expectedError: "retriable failure when dispatching ",
		},
		{
			name:        "dispatcher returns fatal error",
			messageData: []byte(`{"id":"test-id-3","success":false,"manifestID":"manifest-3"}`),
			setupDispatcher: func(d *mockDeploymentDispatcher) {
				d.On("Dispatch", mock.Anything, mock.MatchedBy(func(r model.DeploymentResponse) bool {
					return r.ID == "test-id-3"
				})).Return(types.NewFatalError("permanent failure"))
			},
			setupMessage: func(m *mockJetStreamMsg) {
				m.On("Data").Return([]byte(`{"id":"test-id-3","success":false,"manifestID":"manifest-3"}`))
				m.On("Ack").Return(nil)
			},
			expectedError: "fatal failure when dispatching ",
		},
		{
			name:        "dispatcher returns fatal error and ACK fails",
			messageData: []byte(`{"id":"test-id-4","success":false,"manifestID":"manifest-4"}`),
			setupDispatcher: func(d *mockDeploymentDispatcher) {
				d.On("Dispatch", mock.Anything, mock.MatchedBy(func(r model.DeploymentResponse) bool {
					return r.ID == "test-id-4"
				})).Return(types.NewFatalError("permanent failure"))
			},
			setupMessage: func(m *mockJetStreamMsg) {
				m.On("Data").Return([]byte(`{"id":"test-id-4","success":false,"manifestID":"manifest-4"}`))
				m.On("Ack").Return(errors.New("ACK failed"))
			},
			expectedError: "fatal failure when dispatching ",
		},
		{
			name:        "ACK message error after successful dispatch",
			messageData: []byte(`{"id":"test-id-5","success":true,"manifestID":"manifest-5"}`),
			setupDispatcher: func(d *mockDeploymentDispatcher) {
				d.On("Dispatch", mock.Anything, mock.MatchedBy(func(r model.DeploymentResponse) bool {
					return r.ID == "test-id-5"
				})).Return(nil)
			},
			setupMessage: func(m *mockJetStreamMsg) {
				m.On("Data").Return([]byte(`{"id":"test-id-5","success":true,"manifestID":"manifest-5"}`))
				m.On("Ack").Return(errors.New("ACK failed"))
			},
			expectedError: "failed to ACK ",
		},
		{
			name:        "ACK message error for invalid JSON",
			messageData: []byte(`invalid json`),
			setupDispatcher: func(d *mockDeploymentDispatcher) {
				// Dispatcher should not be called for invalid JSON
			},
			setupMessage: func(m *mockJetStreamMsg) {
				m.On("Data").Return([]byte(`invalid json`))
				m.On("Ack").Return(errors.New("ACK failed"))
			},
			expectedError: "failed to unmarshal ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDispatcher := &mockDeploymentDispatcher{}
			tt.setupDispatcher(mockDispatcher)

			mockMessage := &mockJetStreamMsg{}
			tt.setupMessage(mockMessage)

			client := newNatsDeploymentClient(nil, mockDispatcher, system.NoopMonitor{})

			err := client.ProcessMessage(context.Background(), mockMessage)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)

			mockDispatcher.AssertExpectations(t)
			mockMessage.AssertExpectations(t)
		})
	}
}

func TestNatsDeploymentClient_ProcessLoop_Errors(t *testing.T) {
	tests := []struct {
		name          string
		setupConsumer func(*mockConsumer)
		expectedError string
	}{
		{
			name: "consumer fetch error",
			setupConsumer: func(c *mockConsumer) {
				c.On("Fetch", 1, mock.Anything).Return(nil, errors.New("fetch failed"))
			},
			expectedError: "fetch failed",
		},
		{
			name: "consumer connection lost",
			setupConsumer: func(c *mockConsumer) {
				c.On("Fetch", 1, mock.Anything).Return(nil, &jetstream.APIError{ErrorCode: jetstream.JSErrCodeStreamNotFound})
			},
			// NATS stream not found error code 10059
			expectedError: "err_code=10059",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConsumer := &mockConsumer{}
			tt.setupConsumer(mockConsumer)

			mockDispatcher := &mockDeploymentDispatcher{}
			client := newNatsDeploymentClient(nil, mockDispatcher, system.NoopMonitor{})

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err := client.ProcessLoop(ctx, mockConsumer)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)

			mockConsumer.AssertExpectations(t)
		})
	}
}

// Mock implementations

type mockDeploymentDispatcher struct {
	mock.Mock
}

func (m *mockDeploymentDispatcher) Dispatch(ctx context.Context, response model.DeploymentResponse) error {
	args := m.Called(ctx, response)
	return args.Error(0)
}

type mockJetStreamMsg struct {
	mock.Mock
}

func (m *mockJetStreamMsg) Headers() nats.Header {
	panic("not implemented")
}

func (m *mockJetStreamMsg) Metadata() (*jetstream.MsgMetadata, error) {
	panic("not implemented")
}

func (m *mockJetStreamMsg) DoubleAck(_ context.Context) error {
	panic("not implemented")
}

func (m *mockJetStreamMsg) NakWithDelay(_ time.Duration) error {
	panic("not implemented")
}

func (m *mockJetStreamMsg) TermWithReason(_ string) error {
	panic("not implemented")
}

func (m *mockJetStreamMsg) Data() []byte {
	args := m.Called()
	return args.Get(0).([]byte)
}

func (m *mockJetStreamMsg) Ack() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockJetStreamMsg) Nak() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockJetStreamMsg) InProgress() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockJetStreamMsg) Term() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockJetStreamMsg) Reply() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockJetStreamMsg) Subject() string {
	args := m.Called()
	return args.String(0)
}

type mockConsumer struct {
	mock.Mock
}

func (m *mockConsumer) FetchBytes(_ int, _ ...jetstream.FetchOpt) (jetstream.MessageBatch, error) {
	panic("not implemented")
}

func (m *mockConsumer) FetchNoWait(_ int) (jetstream.MessageBatch, error) {
	panic("not implemented")
}

func (m *mockConsumer) Consume(handler jetstream.MessageHandler, opts ...jetstream.PullConsumeOpt) (jetstream.ConsumeContext, error) {
	args := m.Called(handler, opts)
	return args.Get(0).(jetstream.ConsumeContext), args.Error(1)
}

func (m *mockConsumer) Messages(opts ...jetstream.PullMessagesOpt) (jetstream.MessagesContext, error) {
	args := m.Called(opts)
	return args.Get(0).(jetstream.MessagesContext), args.Error(1)
}

func (m *mockConsumer) Fetch(batch int, opts ...jetstream.FetchOpt) (jetstream.MessageBatch, error) {
	args := m.Called(batch, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(jetstream.MessageBatch), args.Error(1)
}

func (m *mockConsumer) Next(opts ...jetstream.FetchOpt) (jetstream.Msg, error) {
	args := m.Called(opts)
	return args.Get(0).(jetstream.Msg), args.Error(1)
}

func (m *mockConsumer) Info(ctx context.Context) (*jetstream.ConsumerInfo, error) {
	args := m.Called(ctx)
	return args.Get(0).(*jetstream.ConsumerInfo), args.Error(1)
}

func (m *mockConsumer) CachedInfo() *jetstream.ConsumerInfo {
	args := m.Called()
	return args.Get(0).(*jetstream.ConsumerInfo)
}
