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

package core

import (
	"context"
	"errors"
	"testing"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/api/mocks"
	"github.com/metaform/connector-fabric-manager/pmanager/memorystore"
)

func TestProvisionManager_Start(t *testing.T) {
	tests := []struct {
		name           string
		manifest       *model.OrchestrationManifest
		setupStore     func(store api.DefinitionStore)
		setupOrch      func(orch *mocks.MockOrchestrator)
		expectedError  string
		expectedResult *api.Orchestration
	}{
		{
			name: "successful deployment with new orchestration",
			manifest: &model.OrchestrationManifest{
				ID:                "test-orchestration-1",
				OrchestrationType: "test-type",
				Payload:           map[string]any{"key": "value"},
			},
			setupStore: func(store api.DefinitionStore) {
				definition := &api.OrchestrationDefinition{
					Type: "test-type",
					Activities: []api.Activity{
						{
							ID:   "activity1",
							Type: "test-activity",
						},
					},
				}
				store.StoreOrchestrationDefinition(definition)
			},
			setupOrch: func(orch *mocks.MockOrchestrator) {
				orch.EXPECT().GetOrchestration(mock.Anything, "test-orchestration-1").Return(nil, nil)
				orch.EXPECT().Execute(mock.Anything, mock.AnythingOfType("*api.Orchestration")).Return(nil)
			},
			expectedResult: &api.Orchestration{
				ID: "test-orchestration-1",
			},
		},
		{
			name: "deduplication - successful deployment with existing orchestration",
			manifest: &model.OrchestrationManifest{
				ID:                "test-orchestration-2",
				OrchestrationType: "test-type",
				Payload:           map[string]any{"key": "value"},
			},
			setupStore: func(store api.DefinitionStore) {
				definition := &api.OrchestrationDefinition{
					Type: "test-type",
					Activities: []api.Activity{
						{
							ID:   "activity1",
							Type: "test-activity",
						},
					},
				}
				store.StoreOrchestrationDefinition(definition)
			},
			setupOrch: func(orch *mocks.MockOrchestrator) {
				existingOrch := &api.Orchestration{
					ID: "test-orchestration-2",
				}
				orch.EXPECT().GetOrchestration(mock.Anything, "test-orchestration-2").Return(existingOrch, nil)
			},
			expectedResult: &api.Orchestration{
				ID: "test-orchestration-2",
			},
		},
		{
			name: "orchestration definition not found",
			manifest: &model.OrchestrationManifest{
				ID:                "test-orchestration-3",
				OrchestrationType: "non-existent-type",
				Payload:           map[string]any{"key": "value"},
			},
			setupStore: func(store api.DefinitionStore) {
				// Don't store any definitions
			},
			setupOrch: func(orch *mocks.MockOrchestrator) {
				// No orchestrator calls expected
			},
			expectedError: "orchestration type 'non-existent-type' not found",
		},
		{
			name: "orchestrator get orchestration error",
			manifest: &model.OrchestrationManifest{
				ID:                "test-orchestration-5",
				OrchestrationType: "test-type",
				Payload:           map[string]any{"key": "value"},
			},
			setupStore: func(store api.DefinitionStore) {
				definition := &api.OrchestrationDefinition{
					Type: "test-type",
					Activities: []api.Activity{
						{
							ID:   "activity1",
							Type: "test-activity",
						},
					},
				}
				store.StoreOrchestrationDefinition(definition)
			},
			setupOrch: func(orch *mocks.MockOrchestrator) {
				orch.EXPECT().GetOrchestration(mock.Anything, "test-orchestration-5").Return(nil, errors.New("orchestrator error"))
			},
			expectedError: "error checking for orchestration test-orchestration-5: orchestrator error",
		},
		{
			name: "orchestrator execute orchestration error",
			manifest: &model.OrchestrationManifest{
				ID:                "test-orchestration-6",
				OrchestrationType: "test-type",
				Payload:           map[string]any{"key": "value"},
			},
			setupStore: func(store api.DefinitionStore) {
				definition := &api.OrchestrationDefinition{
					Type: "test-type",
					Activities: []api.Activity{
						{
							ID:   "activity1",
							Type: "test-activity",
						},
					},
				}
				store.StoreOrchestrationDefinition(definition)
			},
			setupOrch: func(orch *mocks.MockOrchestrator) {
				orch.EXPECT().GetOrchestration(mock.Anything, "test-orchestration-6").Return(nil, nil)
				orch.EXPECT().Execute(mock.Anything, mock.AnythingOfType("*api.Orchestration")).Return(errors.New("execution error"))
			},
			expectedError: "error executing orchestration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup memory store
			store := memorystore.NewDefinitionStore()
			tt.setupStore(store)

			// Setup mock orchestrator
			mockOrch := mocks.NewMockOrchestrator(t)
			tt.setupOrch(mockOrch)

			// Create provision manager
			pm := &provisionManager{
				orchestrator: mockOrch,
				store:        store,
				monitor:      &system.NoopMonitor{},
			}

			// Execute test
			ctx := context.Background()
			result, err := pm.Start(ctx, tt.manifest)

			// Assert results
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.ID, result.ID)
				//if tt.expectedResult.Status != "" {
				//	assert.Equal(t, tt.expectedResult.Status, result.Status)
				//}
			}
		})
	}
}

// Helper function to create a test orchestration definition
func createTestOrchestrationDefinition(orchestrationType string, active bool) *api.OrchestrationDefinition {
	return &api.OrchestrationDefinition{
		Type: model.OrchestrationType(orchestrationType),
		Activities: []api.Activity{
			{
				ID:   "activity1",
				Type: "test-activity",
			},
		},
	}
}

// Test helper to verify orchestration instantiation
func TestProvisionManager_Start_OrchestrationInstantiation(t *testing.T) {
	// Setup memory store with test definition
	store := memorystore.NewDefinitionStore()
	definition := createTestOrchestrationDefinition("test-type", true)
	store.StoreOrchestrationDefinition(definition)

	// Setup mock orchestrator
	mockOrch := mocks.NewMockOrchestrator(t)
	mockOrch.EXPECT().GetOrchestration(mock.Anything, "test-deployment").Return(nil, nil)
	mockOrch.EXPECT().Execute(mock.Anything, mock.MatchedBy(func(orch *api.Orchestration) bool {
		// Verify orchestration properties
		return orch.ID == "test-deployment"
	})).Return(nil)

	// Create provision manager
	pm := &provisionManager{
		orchestrator: mockOrch,
		store:        store,
		monitor:      &system.NoopMonitor{},
	}

	// Create test manifest
	manifest := &model.OrchestrationManifest{
		ID:                "test-deployment",
		OrchestrationType: "test-type",
		Payload:           map[string]any{"key": "value"},
	}

	// Execute test
	ctx := context.Background()
	result, err := pm.Start(ctx, manifest)

	// Assert results
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "test-deployment", result.ID)
}
