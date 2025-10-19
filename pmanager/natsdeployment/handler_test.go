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
	"errors"
	"testing"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/common/natsclient/mocks"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNatsDeploymentHandler_Dispatcher_SuccessfulDispatch(t *testing.T) {
	mockClient := mocks.NewMockMsgClient(t)
	mockProvisionManager := &MockProvisionManager{}

	handler := newNatsDeploymentHandler(mockClient, mockProvisionManager, system.NoopMonitor{})

	ctx := context.Background()
	manifest := model.DeploymentManifest{
		ID:             "test-manifest-id",
		DeploymentType: model.VpaDeploymentType,
		Payload:        map[string]any{"key": "value"},
	}

	orchestration := &api.Orchestration{
		ID:    "orchestration-123",
		State: api.OrchestrationStateRunning,
	}

	// Mock expectations - successful provision manager call
	mockProvisionManager.On("Start", ctx, &manifest).Return(orchestration, nil)

	// No publish call should be made on success
	err := handler.RetriableMessageProcessor.Dispatcher(ctx, manifest)

	// Should return no error (message gets ACK'd)
	require.NoError(t, err)

	mockProvisionManager.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestNatsDeploymentHandler_Dispatcher_ValidResponseStructure(t *testing.T) {
	mockClient := mocks.NewMockMsgClient(t)
	mockProvisionManager := &MockProvisionManager{}

	handler := newNatsDeploymentHandler(mockClient, mockProvisionManager, system.NoopMonitor{})

	ctx := context.Background()
	manifest := model.DeploymentManifest{
		ID:             "test-manifest-id",
		DeploymentType: model.VpaDeploymentType,
		Payload:        map[string]any{"key": "value"},
	}

	nonRecoverableErr := errors.New("deployment failed")
	pubAck := &jetstream.PubAck{}

	// Mock expectations
	mockProvisionManager.On("Start", ctx, &manifest).Return(nil, nonRecoverableErr)

	var capturedPayload []byte
	mockClient.EXPECT().Publish(ctx, natsclient.CFMDeploymentResponseSubject, mock.AnythingOfType("[]uint8")).
		Run(func(ctx context.Context, subject string, payload []byte, opts ...jetstream.PublishOpt) {
			capturedPayload = payload
		}).Return(pubAck, nil)

	err := handler.RetriableMessageProcessor.Dispatcher(ctx, manifest)

	require.NoError(t, err)

	var response model.DeploymentResponse
	err = json.Unmarshal(capturedPayload, &response)
	require.NoError(t, err)

	assert.False(t, response.Success)
	assert.Equal(t, "deployment failed", response.ErrorDetail)
	assert.Equal(t, "test-manifest-id", response.ManifestID)
	assert.Equal(t, model.VpaDeploymentType, response.DeploymentType)
	assert.NotEmpty(t, response.ID) // Should have a generated UUID
	assert.NotNil(t, response.Properties)
	assert.Equal(t, 0, len(response.Properties)) // Should be empty map

	mockProvisionManager.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestNatsDeploymentHandler_Dispatcher_ErrorTypes(t *testing.T) {
	testCases := []struct {
		name          string
		error         error
		shouldPublish bool
	}{
		{"RecoverableError", types.NewRecoverableError("network timeout"), false},
		{"ClientError", types.NewClientError("invalid manifest"), true},
		{"FatalError", types.NewFatalError("fatal processing error"), true},
		{"StandardError", errors.New("standard error"), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockClient := mocks.NewMockMsgClient(t)
			mockProvisionManager := &MockProvisionManager{}

			handler := newNatsDeploymentHandler(mockClient, mockProvisionManager, system.NoopMonitor{})

			ctx := context.Background()
			manifest := model.DeploymentManifest{
				ID:             "test-manifest-id",
				DeploymentType: model.VpaDeploymentType,
				Payload:        map[string]any{"key": "value"},
			}

			// Mock expectations
			mockProvisionManager.On("Start", ctx, &manifest).Return(nil, tc.error)

			if tc.shouldPublish {
				pubAck := &jetstream.PubAck{}
				mockClient.EXPECT().Publish(ctx, natsclient.CFMDeploymentResponseSubject, mock.Anything).Return(pubAck, nil)
			}

			err := handler.RetriableMessageProcessor.Dispatcher(ctx, manifest)

			if tc.shouldPublish {
				assert.NoError(t, err, "Non-recoverable errors should result in ACK")
			} else {
				assert.Error(t, err, "Recoverable errors should be returned")
				assert.True(t, types.IsRecoverable(err))
			}

			mockProvisionManager.AssertExpectations(t)
			mockClient.AssertExpectations(t)
		})
	}
}

type MockProvisionManager struct {
	mock.Mock
}

func (m *MockProvisionManager) Start(ctx context.Context, manifest *model.DeploymentManifest) (*api.Orchestration, error) {
	args := m.Called(ctx, manifest)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.Orchestration), args.Error(1)
}

func (m *MockProvisionManager) Cancel(ctx context.Context, deploymentID string) error {
	args := m.Called(ctx, deploymentID)
	return args.Error(0)
}

func (m *MockProvisionManager) GetOrchestration(ctx context.Context, deploymentID string) (*api.Orchestration, error) {
	args := m.Called(ctx, deploymentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*api.Orchestration), args.Error(1)
}
