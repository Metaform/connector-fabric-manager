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

package pmcore

import (
	"context"
	"errors"
	"testing"

	"github.com/metaform/connector-fabric-manager/common/dmodel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/api/mocks"
	"github.com/metaform/connector-fabric-manager/pmanager/memorystore"
)

func TestProvisionManager_Start(t *testing.T) {
	tests := []struct {
		name           string
		manifest       *dmodel.DeploymentManifest
		setupStore     func(store api.DefinitionStore)
		setupOrch      func(orch *mocks.DeploymentOrchestrator)
		expectedError  string
		expectedResult *api.Orchestration
	}{
		{
			name: "successful deployment with new orchestration",
			manifest: &dmodel.DeploymentManifest{
				ID:             "test-deployment-1",
				DeploymentType: "test-type",
				Payload:        map[string]interface{}{"key": "value"},
			},
			setupStore: func(store api.DefinitionStore) {
				definition := &api.DeploymentDefinition{
					Type:       "test-type",
					ApiVersion: "1.0",
					Versions: []api.Version{
						{
							Version: "1.0.0",
							Active:  true,
							Activities: []api.Activity{
								{
									ID:   "activity1",
									Type: "test-activity",
								},
							},
						},
					},
				}
				store.StoreDeploymentDefinition(definition)
			},
			setupOrch: func(orch *mocks.DeploymentOrchestrator) {
				orch.EXPECT().GetOrchestration(mock.Anything, "test-deployment-1").Return(nil, nil)
				orch.EXPECT().ExecuteOrchestration(mock.Anything, mock.AnythingOfType("*api.Orchestration")).Return(nil)
			},
			expectedResult: &api.Orchestration{
				ID: "test-deployment-1",
			},
		},
		{
			name: "deduplication - successful deployment with existing orchestration",
			manifest: &dmodel.DeploymentManifest{
				ID:             "test-deployment-2",
				DeploymentType: "test-type",
				Payload:        map[string]interface{}{"key": "value"},
			},
			setupStore: func(store api.DefinitionStore) {
				definition := &api.DeploymentDefinition{
					Type:       "test-type",
					ApiVersion: "1.0",
					Versions: []api.Version{
						{
							Version: "1.0.0",
							Active:  true,
							Activities: []api.Activity{
								{
									ID:   "activity1",
									Type: "test-activity",
								},
							},
						},
					},
				}
				store.StoreDeploymentDefinition(definition)
			},
			setupOrch: func(orch *mocks.DeploymentOrchestrator) {
				existingOrch := &api.Orchestration{
					ID: "test-deployment-2",
				}
				orch.EXPECT().GetOrchestration(mock.Anything, "test-deployment-2").Return(existingOrch, nil)
			},
			expectedResult: &api.Orchestration{
				ID: "test-deployment-2",
			},
		},
		{
			name: "deployment definition not found",
			manifest: &dmodel.DeploymentManifest{
				ID:             "test-deployment-3",
				DeploymentType: "non-existent-type",
				Payload:        map[string]interface{}{"key": "value"},
			},
			setupStore: func(store api.DefinitionStore) {
				// Don't store any definitions
			},
			setupOrch: func(orch *mocks.DeploymentOrchestrator) {
				// No orchestrator calls expected
			},
			expectedError: "deployment type 'non-existent-type' not found",
		},
		{
			name: "deployment definition has no active version",
			manifest: &dmodel.DeploymentManifest{
				ID:             "test-deployment-4",
				DeploymentType: "test-type-inactive",
				Payload:        map[string]interface{}{"key": "value"},
			},
			setupStore: func(store api.DefinitionStore) {
				definition := &api.DeploymentDefinition{
					Type:       "test-type-inactive",
					ApiVersion: "1.0",
					Versions: []api.Version{
						{
							Version: "1.0.0",
							Active:  false,
							Activities: []api.Activity{
								{
									ID:   "activity1",
									Type: "test-activity",
								},
							},
						},
					},
				}
				store.StoreDeploymentDefinition(definition)
			},
			setupOrch: func(orch *mocks.DeploymentOrchestrator) {
				// No orchestrator calls expected
			},
			expectedError: "error deploying test-deployment-4",
		},
		{
			name: "orchestrator get orchestration error",
			manifest: &dmodel.DeploymentManifest{
				ID:             "test-deployment-5",
				DeploymentType: "test-type",
				Payload:        map[string]interface{}{"key": "value"},
			},
			setupStore: func(store api.DefinitionStore) {
				definition := &api.DeploymentDefinition{
					Type:       "test-type",
					ApiVersion: "1.0",
					Versions: []api.Version{
						{
							Version: "1.0.0",
							Active:  true,
							Activities: []api.Activity{
								{
									ID:   "activity1",
									Type: "test-activity",
								},
							},
						},
					},
				}
				store.StoreDeploymentDefinition(definition)
			},
			setupOrch: func(orch *mocks.DeploymentOrchestrator) {
				orch.EXPECT().GetOrchestration(mock.Anything, "test-deployment-5").Return(nil, errors.New("orchestrator error"))
			},
			expectedError: "error checking for orchestration test-deployment-5: orchestrator error",
		},
		{
			name: "orchestrator execute orchestration error",
			manifest: &dmodel.DeploymentManifest{
				ID:             "test-deployment-6",
				DeploymentType: "test-type",
				Payload:        map[string]interface{}{"key": "value"},
			},
			setupStore: func(store api.DefinitionStore) {
				definition := &api.DeploymentDefinition{
					Type:       "test-type",
					ApiVersion: "1.0",
					Versions: []api.Version{
						{
							Version: "1.0.0",
							Active:  true,
							Activities: []api.Activity{
								{
									ID:   "activity1",
									Type: "test-activity",
								},
							},
						},
					},
				}
				store.StoreDeploymentDefinition(definition)
			},
			setupOrch: func(orch *mocks.DeploymentOrchestrator) {
				orch.EXPECT().GetOrchestration(mock.Anything, "test-deployment-6").Return(nil, nil)
				orch.EXPECT().ExecuteOrchestration(mock.Anything, mock.AnythingOfType("*api.Orchestration")).Return(errors.New("execution error"))
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
			mockOrch := mocks.NewDeploymentOrchestrator(t)
			tt.setupOrch(mockOrch)

			// Create provision manager
			pm := &provisionManager{
				orchestrator: mockOrch,
				store:        store,
				logMonitor:   &monitor.NoopMonitor{},
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

// Helper function to create a test deployment definition
func createTestDeploymentDefinition(deploymentType string, active bool) *api.DeploymentDefinition {
	return &api.DeploymentDefinition{
		Type:       dmodel.DeploymentType(deploymentType),
		ApiVersion: "1.0",
		Versions: []api.Version{
			{
				Version: "1.0.0",
				Active:  active,
				Activities: []api.Activity{
					{
						ID:   "activity1",
						Type: "test-activity",
					},
				},
			},
		},
	}
}

// Test helper to verify orchestration instantiation
func TestProvisionManager_Start_OrchestrationInstantiation(t *testing.T) {
	// Setup memory store with test definition
	store := memorystore.NewDefinitionStore()
	definition := createTestDeploymentDefinition("test-type", true)
	store.StoreDeploymentDefinition(definition)

	// Setup mock orchestrator
	mockOrch := mocks.NewDeploymentOrchestrator(t)
	mockOrch.EXPECT().GetOrchestration(mock.Anything, "test-deployment").Return(nil, nil)
	mockOrch.EXPECT().ExecuteOrchestration(mock.Anything, mock.MatchedBy(func(orch *api.Orchestration) bool {
		// Verify orchestration properties
		return orch.ID == "test-deployment"
	})).Return(nil)

	// Create provision manager
	pm := &provisionManager{
		orchestrator: mockOrch,
		store:        store,
		logMonitor:   &monitor.NoopMonitor{},
	}

	// Create test manifest
	manifest := &dmodel.DeploymentManifest{
		ID:             "test-deployment",
		DeploymentType: "test-type",
		Payload:        map[string]interface{}{"key": "value"},
	}

	// Execute test
	ctx := context.Background()
	result, err := pm.Start(ctx, manifest)

	// Assert results
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "test-deployment", result.ID)
}
