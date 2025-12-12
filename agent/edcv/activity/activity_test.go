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

package activity

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/metaform/connector-fabric-manager/common/mocks"
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/common/token"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ConfigOptions func(*Config)

func WithTokenProvider(provider token.TokenProvider) ConfigOptions {
	return func(config *Config) {
		config.TokenProvider = provider
	}
}

func validConfig(opts ...ConfigOptions) *Config {
	c := Config{
		VaultClient:        NewMockVaultClient("client-123", "123"),
		HTTPClient:         &http.Client{},
		Monitor:            system.NoopMonitor{},
		IdentityHubBaseURL: "http://identity.example.com:8765/foo/bar",
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

func TestEDCVActivityProcessor_Process_WithValidData(t *testing.T) {
	tokenProvider := mocks.NewMockTokenProvider(t)
	tokenProvider.On("GetToken").Return("someToken", nil)
	processor := NewProcessor(validConfig(WithTokenProvider(tokenProvider)))

	ctx := context.Background()
	processingData := map[string]any{
		model.ParticipantIdentifier: "participant-abc",
		"clientID.vaultAccess":      "client-123",
		"clientID.apiAccess":        "client-456",
	}
	outputData := make(map[string]any)

	activity := api.Activity{
		ID:   "test-activity",
		Type: "edcv",
	}

	activityContext := api.NewActivityContext(ctx, "orch-123", activity, processingData, outputData)

	result := processor.Process(activityContext)

	assert.Equal(t, api.ActivityResultType(api.ActivityResultComplete), result.Result)
	assert.NoError(t, result.Error)
}

func TestEDCVActivityProcessor_Process_WithPublicURL(t *testing.T) {
	tokenProvider := mocks.NewMockTokenProvider(t)
	tokenProvider.On("GetToken").Return("someToken", nil)
	processor := NewProcessor(validConfig(WithTokenProvider(tokenProvider)))

	ctx := context.Background()
	processingData := map[string]any{
		model.ParticipantIdentifier: "participant-abc",
		"clientID.vaultAccess":      "client-123",
		"clientID.apiAccess":        "client-456",
		"publicURL":                 "http://secure.example.com:1234/fizz/buzz",
	}
	outputData := make(map[string]any)

	activity := api.Activity{
		ID:   "test-activity",
		Type: "edcv",
	}

	activityContext := api.NewActivityContext(ctx, "orch-123", activity, processingData, outputData)

	result := processor.Process(activityContext)

	assert.Equal(t, api.ActivityResultType(api.ActivityResultComplete), result.Result)
	assert.NoError(t, result.Error)
}

func TestEDCVActivityProcessor_Process_MissingParticipantID(t *testing.T) {

	processor := NewProcessor(validConfig(WithTokenProvider(mocks.NewMockTokenProvider(t))))
	ctx := context.Background()
	processingData := map[string]any{
		"clientID.vaultAccess": "client-123",
		"clientID.apiAccess":   "client-456",
		// Missing model.ParticipantIdentifier
	}
	outputData := make(map[string]any)

	activity := api.Activity{
		ID:            "activity-1",
		Type:          "edcv",
		Discriminator: api.DeployDiscriminator,
	}

	activityContext := api.NewActivityContext(ctx, "orchestration-1", activity, processingData, outputData)

	result := processor.Process(activityContext)

	require.NotNil(t, result)
	assert.Equal(t, api.ActivityResultType(api.ActivityResultFatalError), result.Result)
	assert.NotNil(t, result.Error)
	assert.Contains(t, result.Error.Error(), "error processing EDC-V activity")
}

func TestEDCVActivityProcessor_Process_MissingClientID(t *testing.T) {
	processor := NewProcessor(validConfig(WithTokenProvider(mocks.NewMockTokenProvider(t))))
	ctx := context.Background()
	processingData := map[string]any{
		model.ParticipantIdentifier: "participant-123",
		// Missing "clientID"
	}
	outputData := make(map[string]any)

	activity := api.Activity{
		ID:            "activity-2",
		Type:          "edcv",
		Discriminator: api.DeployDiscriminator,
	}

	activityContext := api.NewActivityContext(ctx, "orchestration-2", activity, processingData, outputData)

	result := processor.Process(activityContext)

	require.NotNil(t, result)
	assert.Equal(t, api.ActivityResultType(api.ActivityResultFatalError), result.Result)
	assert.NotNil(t, result.Error)
	assert.Contains(t, result.Error.Error(), "error processing EDC-V activity")
}

func TestEDCVActivityProcessor_Process_EmptyProcessingData(t *testing.T) {
	processor := NewProcessor(validConfig(WithTokenProvider(mocks.NewMockTokenProvider(t))))
	ctx := context.Background()
	processingData := make(map[string]any)
	outputData := make(map[string]any)

	activity := api.Activity{
		ID:   "activity-3",
		Type: "edcv",
	}

	activityContext := api.NewActivityContext(ctx, "orchestration-3", activity, processingData, outputData)

	result := processor.Process(activityContext)

	require.NotNil(t, result)
	assert.Equal(t, api.ActivityResultType(api.ActivityResultFatalError), result.Result)
	assert.NotNil(t, result.Error)
}

func TestEDCVActivityProcessor_Process_InvalidDataTypes(t *testing.T) {
	processor := NewProcessor(validConfig(WithTokenProvider(mocks.NewMockTokenProvider(t))))
	ctx := context.Background()
	processingData := map[string]any{
		model.ParticipantIdentifier: 123, // Should be string
		"clientID":                  456, // Should be string
	}
	outputData := make(map[string]any)

	activity := api.Activity{
		ID:   "activity-4",
		Type: "edcv",
	}

	activityContext := api.NewActivityContext(ctx, "orchestration-4", activity, processingData, outputData)

	result := processor.Process(activityContext)

	require.NotNil(t, result)
	assert.Equal(t, api.ActivityResultType(api.ActivityResultFatalError), result.Result)
	assert.NotNil(t, result.Error)
}

func TestEDCVActivityProcessor_Process_OrchestrationIDInError(t *testing.T) {
	processor := NewProcessor(validConfig(WithTokenProvider(mocks.NewMockTokenProvider(t))))

	ctx := context.Background()
	processingData := map[string]any{
		model.ParticipantIdentifier: "participant-123",
		// Missing clientID
	}
	outputData := make(map[string]any)

	activity := api.Activity{
		ID:   "activity-5",
		Type: "edcv",
	}

	orchestrationID := "test-orch-12345"
	activityContext := api.NewActivityContext(ctx, orchestrationID, activity, processingData, outputData)

	result := processor.Process(activityContext)

	require.NotNil(t, result)
	assert.Equal(t, api.ActivityResultType(api.ActivityResultFatalError), result.Result)
	require.NotNil(t, result.Error)
	assert.Contains(t, result.Error.Error(), orchestrationID)
}

func TestEDCVActivityProcessor_Process_MultipleUnknownFields(t *testing.T) {
	tokenProvider := mocks.NewMockTokenProvider(t)
	tokenProvider.On("GetToken").Return("someToken", nil)
	processor := NewProcessor(validConfig(WithTokenProvider(tokenProvider)))

	ctx := context.Background()
	processingData := map[string]any{
		model.ParticipantIdentifier: "participant-multi",
		"clientID.vaultAccess":      "client-123",
		"clientID.apiAccess":        "client-456",
		"field1":                    "value1",
		"field2":                    "value2",
		"field3":                    "value3",
	}
	outputData := make(map[string]any)

	activity := api.Activity{
		ID:   "activity-multi",
		Type: "edcv",
	}

	activityContext := api.NewActivityContext(ctx, "orch-multi", activity, processingData, outputData)

	result := processor.Process(activityContext)

	require.NotNil(t, result)
	assert.Equal(t, api.ActivityResultType(api.ActivityResultComplete), result.Result)
	assert.Nil(t, result.Error)
}

func TestEDCVActivityProcessor_Process_MissingVaultEntry(t *testing.T) {

	tokenProvider := mocks.NewMockTokenProvider(t)
	tokenProvider.On("GetToken").Return("someToken", nil)

	invalidConfig := validConfig(WithTokenProvider(tokenProvider))
	invalidConfig.VaultClient = NewMockVaultClient()
	processor := NewProcessor(invalidConfig)

	ctx := context.Background()
	processingData := map[string]any{
		model.ParticipantIdentifier: "participant-multi",
		"clientID.vaultAccess":      "client-123",
		"clientID.apiAccess":        "client-456",
	}
	outputData := make(map[string]any)

	activity := api.Activity{
		ID:   "activity-multi",
		Type: "edcv",
	}

	activityContext := api.NewActivityContext(ctx, "orch-multi", activity, processingData, outputData)

	result := processor.Process(activityContext)

	assert.NotNil(t, result.Error)
	assert.Equal(t, api.ActivityResultType(api.ActivityResultFatalError), result.Result)
}

func TestEDCVActivityProcessor_Process_TokenFailure(t *testing.T) {
	tokenProvider := mocks.NewMockTokenProvider(t)
	tokenProvider.On("GetToken").Return("", fmt.Errorf("some error"))
	processor := NewProcessor(validConfig(WithTokenProvider(tokenProvider)))

	ctx := context.Background()
	processingData := map[string]any{
		model.ParticipantIdentifier: "participant-abc",
		"clientID.vaultAccess":      "client-123",
		"clientID.apiAccess":        "client-456",
	}
	outputData := make(map[string]any)

	activity := api.Activity{
		ID:   "test-activity",
		Type: "edcv",
	}

	activityContext := api.NewActivityContext(ctx, "orch-123", activity, processingData, outputData)

	result := processor.Process(activityContext)

	assert.Equal(t, api.ActivityResultType(api.ActivityResultFatalError), result.Result)
	assert.Error(t, result.Error, "some error")
}

type MockVaultClient struct {
	cache map[string]string
}

func NewMockVaultClient(secrets ...string) MockVaultClient {
	cache := make(map[string]string)
	for i := 0; i < len(secrets); i += 2 {
		if i+1 < len(secrets) {
			cache[secrets[i]] = secrets[i+1]
		}
	}
	return MockVaultClient{
		cache: cache,
	}
}

func (m MockVaultClient) ResolveSecret(ctx context.Context, path string) (string, error) {
	if value, ok := m.cache[path]; ok {
		return value, nil
	}
	return "", types.ErrNotFound
}

func (m MockVaultClient) StoreSecret(ctx context.Context, path string, value string) error {
	m.cache[path] = value
	return nil
}

func (m MockVaultClient) DeleteSecret(ctx context.Context, path string) error {
	delete(m.cache, path)
	return nil
}

func (m MockVaultClient) Close() error {
	return nil
}

func (m MockVaultClient) Health(ctx context.Context) error {
	return nil
}
