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

//go:build test
// +build test

package natsorchestration_test_test

import (
	"context"
	"encoding/json"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsorchestration"
	"github.com/metaform/connector-fabric-manager/pmanager/natsorchestration/mocks"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	activityType = "test.activity"
	subject      = "event.test-activity"
)

func TestExecuteOrchestration_Errors(t *testing.T) {
	tests := []struct {
		name          string
		orchestration api.Orchestration
		setupMock     func(*mocks.MsgClient)
		expectedError string
		expectedCalls int
		parallelCheck bool
	}{
		{
			name: "orchestration already exists",
			orchestration: api.Orchestration{
				ID: "test-3",
				Steps: []api.OrchestrationStep{
					{
						Activities: []api.Activity{{ID: "A1", Type: activityType}},
					},
				},
			},
			setupMock: func(m *mocks.MsgClient) {
				m.EXPECT().Update(mock.Anything, "test-3", mock.Anything, uint64(0)).
					Return(uint64(0), &jetstream.APIError{ErrorCode: jetstream.JSErrCodeStreamWrongLastSequence})
			},
			expectedCalls: 0,
		},
		{
			name: "update error",
			orchestration: api.Orchestration{
				ID: "test-4",
				Steps: []api.OrchestrationStep{
					{
						Activities: []api.Activity{{ID: "A1", Type: activityType}},
					},
				},
			},
			setupMock: func(m *mocks.MsgClient) {
				m.EXPECT().Update(mock.Anything, "test-4", mock.Anything, uint64(0)).
					Return(uint64(0), assert.AnError)
			},
			expectedError: "error storing orchestration",
			expectedCalls: 0,
		},
		{
			name: "publish error",
			orchestration: api.Orchestration{
				ID: "test-5",
				Steps: []api.OrchestrationStep{
					{
						Activities: []api.Activity{{ID: "A1", Type: activityType}},
					},
				},
			},
			setupMock: func(m *mocks.MsgClient) {
				m.EXPECT().Update(mock.Anything, "test-5", mock.Anything, uint64(0)).
					Return(uint64(1), nil)
				m.EXPECT().Publish(mock.Anything, subject, mock.Anything).
					Return(nil, assert.AnError)
			},
			expectedError: "error publishing to stream",
			expectedCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMsgClient(t)
			tt.setupMock(mockClient)

			orchestrator := &natsorchestration.NatsDeploymentOrchestrator{
				Client:  mockClient,
				Monitor: monitor.NoopMonitor{},
			}

			err := orchestrator.ExecuteOrchestration(context.Background(), &tt.orchestration)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			// Verify the published messages contain correct data
			if tt.expectedCalls > 0 {
				calls := mockClient.Calls
				publishCalls := 0
				for _, call := range calls {
					if call.Method == "Publish" {
						publishCalls++
						payload := call.Arguments[2].([]byte)
						var msg api.ActivityMessage
						err := json.Unmarshal(payload, &msg)
						require.NoError(t, err)
						assert.Equal(t, tt.orchestration.ID, msg.OrchestrationID)
					}
				}
				assert.Equal(t, tt.expectedCalls, publishCalls)
			}
		})
	}
}

func TestEnqueueMessages_Errors(t *testing.T) {
	tests := []struct {
		name       string
		activities []api.Activity
		setupMock  func(*mocks.MsgClient)
		wantErr    bool
	}{
		{
			name: "publish error",
			activities: []api.Activity{
				{ID: "A1", Type: activityType},
			},
			setupMock: func(m *mocks.MsgClient) {
				m.EXPECT().Publish(mock.Anything, subject, mock.Anything).
					Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMsgClient(t)
			tt.setupMock(mockClient)

			err := natsorchestration.EnqueueActivityMessages(context.Background(), "test-oid", tt.activities, mockClient)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
