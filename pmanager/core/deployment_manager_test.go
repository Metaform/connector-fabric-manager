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

package core

import (
	"context"
	"errors"
	"testing"

	"github.com/metaform/connector-fabric-manager/common/model"
	cstore "github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/memorystore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefinitionManager_CreateDeploymentDefinition_Success(t *testing.T) {
	// Given
	store := memorystore.NewDefinitionStore()
	manager := definitionManager{
		trxContext: cstore.NoOpTransactionContext{},
		store:      store,
	}
	ctx := context.Background()

	// Create required activity definition first
	activityDef := &api.ActivityDefinition{
		Type:         "test-activity",
		Description:  "Test activity",
		InputSchema:  map[string]any{"type": "object"},
		OutputSchema: map[string]any{"type": "object"},
	}
	_, err := store.StoreActivityDefinition(activityDef)
	require.NoError(t, err, "Failed to store activity definition")

	// Create deployment definition that references the activity
	deploymentDef := &api.DeploymentDefinition{
		Type:   model.DeploymentType("test-deployment"),
		Active: true,
		Schema: map[string]any{"type": "object"},
		Activities: []api.Activity{
			{
				ID:        "activity-1",
				Type:      "test-activity",
				Inputs:    []api.MappingEntry{{Source: "input1", Target: "target1"}},
				DependsOn: []string{},
			},
		},
	}

	result, err := manager.CreateDeploymentDefinition(ctx, deploymentDef)

	require.NoError(t, err, "CreateDeploymentDefinition should succeed")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, deploymentDef.Type, result.Type, "Deployment type should match")
	assert.Equal(t, deploymentDef.Active, result.Active, "Active flag should match")
	assert.Equal(t, len(deploymentDef.Activities), len(result.Activities), "Activities count should match")
}

func TestDefinitionManager_CreateDeploymentDefinition_MissingActivityDefinition(t *testing.T) {
	// Given
	store := memorystore.NewDefinitionStore()
	manager := definitionManager{
		trxContext: cstore.NoOpTransactionContext{},
		store:      store,
	}
	ctx := context.Background()

	// Create deployment definition that references a non-existent activity
	deploymentDef := &api.DeploymentDefinition{
		Type:   model.DeploymentType("test-deployment"),
		Active: true,
		Schema: map[string]any{"type": "object"},
		Activities: []api.Activity{
			{
				ID:        "activity-1",
				Type:      "non-existent-activity", // This activity definition doesn't exist
				Inputs:    []api.MappingEntry{{Source: "input1", Target: "target1"}},
				DependsOn: []string{},
			},
		},
	}

	result, err := manager.CreateDeploymentDefinition(ctx, deploymentDef)

	require.Error(t, err, "CreateDeploymentDefinition should fail when activity definition is missing")
	assert.Nil(t, result, "Result should be nil on error")

	// Verify the error is a client error about the missing activity
	var clientErr types.ClientError
	require.True(t, errors.As(err, &clientErr), "Error should be a ClientError")
	assert.Contains(t, err.Error(), "activity type 'non-existent-activity' not found",
		"Error message should mention the missing activity type")
}

func TestDefinitionManager_CreateDeploymentDefinition_MultipleMissingActivityDefinitions(t *testing.T) {
	store := memorystore.NewDefinitionStore()
	manager := definitionManager{
		trxContext: cstore.NoOpTransactionContext{},
		store:      store,
	}
	ctx := context.Background()

	// Create one valid activity definition
	activityDef := &api.ActivityDefinition{
		Type:         "valid-activity",
		Description:  "Valid activity",
		InputSchema:  map[string]any{"type": "object"},
		OutputSchema: map[string]any{"type": "object"},
	}
	_, err := store.StoreActivityDefinition(activityDef)
	require.NoError(t, err, "Failed to store valid activity definition")

	// Create deployment definition with mix of valid and invalid activity references
	deploymentDef := &api.DeploymentDefinition{
		Type:   model.DeploymentType("mixed-deployment"),
		Active: true,
		Schema: map[string]any{"type": "object"},
		Activities: []api.Activity{
			{
				ID:        "activity-1",
				Type:      "valid-activity", // This exists
				Inputs:    []api.MappingEntry{{Source: "input1", Target: "target1"}},
				DependsOn: []string{},
			},
			{
				ID:        "activity-2",
				Type:      "missing-activity-1", // This doesn't exist
				Inputs:    []api.MappingEntry{{Source: "input2", Target: "target2"}},
				DependsOn: []string{"activity-1"},
			},
			{
				ID:        "activity-3",
				Type:      "missing-activity-2", // This also doesn't exist
				Inputs:    []api.MappingEntry{{Source: "input3", Target: "target3"}},
				DependsOn: []string{"activity-1"},
			},
		},
	}

	result, err := manager.CreateDeploymentDefinition(ctx, deploymentDef)

	require.Error(t, err, "CreateDeploymentDefinition should fail when multiple activity definitions are missing")
	assert.Nil(t, result, "Result should be nil on error")

	// Verify the error mentions both missing activities
	assert.Contains(t, err.Error(), "missing-activity-1", "Error should mention first missing activity")
	assert.Contains(t, err.Error(), "missing-activity-2", "Error should mention second missing activity")
}

func TestDefinitionManager_CreateDeploymentDefinition_EmptyActivities(t *testing.T) {
	// Given
	store := memorystore.NewDefinitionStore()
	manager := definitionManager{
		trxContext: cstore.NoOpTransactionContext{},
		store:      store,
	}
	ctx := context.Background()

	// Create deployment definition with no activities
	deploymentDef := &api.DeploymentDefinition{
		Type:       model.DeploymentType("empty-deployment"),
		Active:     true,
		Schema:     map[string]any{"type": "object"},
		Activities: []api.Activity{}, // Empty activities slice
	}

	result, err := manager.CreateDeploymentDefinition(ctx, deploymentDef)

	require.NoError(t, err, "CreateDeploymentDefinition should succeed with empty activities")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, 0, len(result.Activities), "Should have 0 activities")
}

func TestDefinitionManager_CreateDeploymentDefinition_StoreError(t *testing.T) {
	store := memorystore.NewDefinitionStore()
	manager := definitionManager{
		trxContext: cstore.NoOpTransactionContext{},
		store:      store,
	}
	ctx := context.Background()

	// Create and store deployment definition first
	deploymentDef := &api.DeploymentDefinition{
		Type:       model.DeploymentType("duplicate-deployment"),
		Active:     true,
		Schema:     map[string]any{"type": "object"},
		Activities: []api.Activity{},
	}

	_, err := store.StoreDeploymentDefinition(deploymentDef)
	require.NoError(t, err, "First store should succeed")

	// Attempt to store the same deployment definition again
	result, err := manager.CreateDeploymentDefinition(ctx, deploymentDef)

	require.Error(t, err, "CreateDeploymentDefinition should fail on duplicate")
	assert.Nil(t, result, "Result should be nil on error")

	// Verify the error is a conflict error
	require.True(t, errors.Is(err, types.ErrConflict), "Error should be a ConflictError")
}

func TestDefinitionManager_CreateActivityDefinition_Success(t *testing.T) {
	store := memorystore.NewDefinitionStore()
	manager := definitionManager{
		trxContext: cstore.NoOpTransactionContext{},
		store:      store,
	}
	ctx := context.Background()

	activityDef := &api.ActivityDefinition{
		Type:         "test-activity",
		Description:  "Test activity definition",
		InputSchema:  map[string]any{"type": "object", "properties": map[string]any{"input": map[string]any{"type": "string"}}},
		OutputSchema: map[string]any{"type": "object", "properties": map[string]any{"output": map[string]any{"type": "string"}}},
	}

	result, err := manager.CreateActivityDefinition(ctx, activityDef)

	require.NoError(t, err, "CreateActivityDefinition should succeed")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, activityDef.Type, result.Type, "Activity type should match")
	assert.Equal(t, activityDef.Description, result.Description, "Description should match")
	assert.Equal(t, activityDef.InputSchema, result.InputSchema, "InputSchema should match")
	assert.Equal(t, activityDef.OutputSchema, result.OutputSchema, "OutputSchema should match")
}

func TestDefinitionManager_CreateActivityDefinition_Duplicate(t *testing.T) {
	store := memorystore.NewDefinitionStore()
	manager := definitionManager{
		trxContext: cstore.NoOpTransactionContext{},
		store:      store,
	}
	ctx := context.Background()

	activityDef := &api.ActivityDefinition{
		Type:         "duplicate-activity",
		Description:  "Test activity definition",
		InputSchema:  map[string]any{"type": "object"},
		OutputSchema: map[string]any{"type": "object"},
	}

	_, err := store.StoreActivityDefinition(activityDef)
	require.NoError(t, err, "First store should succeed")

	// Try to create the same activity definition again
	result, err := manager.CreateActivityDefinition(ctx, activityDef)

	require.Error(t, err, "CreateActivityDefinition should fail on duplicate")
	assert.Nil(t, result, "Result should be nil on error")

	require.True(t, errors.Is(err, types.ErrConflict), "Error should be a ConflictError")
}

func TestDefinitionManager_Integration_CompleteWorkflow(t *testing.T) {
	store := memorystore.NewDefinitionStore()
	manager := definitionManager{
		trxContext: cstore.NoOpTransactionContext{},
		store:      store,
	}
	ctx := context.Background()

	// Create activity definitions first
	activityDef1 := &api.ActivityDefinition{
		Type:         "prepare",
		Description:  "Prepare resources",
		InputSchema:  map[string]any{"type": "object"},
		OutputSchema: map[string]any{"type": "object"},
	}

	activityDef2 := &api.ActivityDefinition{
		Type:         "deploy",
		Description:  "Deploy application",
		InputSchema:  map[string]any{"type": "object"},
		OutputSchema: map[string]any{"type": "object"},
	}

	result1, err := manager.CreateActivityDefinition(ctx, activityDef1)
	require.NoError(t, err, "Should create first activity definition")
	assert.Equal(t, activityDef1.Type, result1.Type, "First activity type should match")

	result2, err := manager.CreateActivityDefinition(ctx, activityDef2)
	require.NoError(t, err, "Should create second activity definition")
	assert.Equal(t, activityDef2.Type, result2.Type, "Second activity type should match")

	// Create deployment definition that uses both activities
	deploymentDef := &api.DeploymentDefinition{
		Type:   model.DeploymentType("full-deployment"),
		Active: true,
		Schema: map[string]any{"type": "object"},
		Activities: []api.Activity{
			{
				ID:        "prepare-step",
				Type:      "prepare",
				Inputs:    []api.MappingEntry{{Source: "config", Target: "configuration"}},
				DependsOn: []string{},
			},
			{
				ID:        "deploy-step",
				Type:      "deploy",
				Inputs:    []api.MappingEntry{{Source: "artifact", Target: "deployment_artifact"}},
				DependsOn: []string{"prepare-step"},
			},
		},
	}

	result3, err := manager.CreateDeploymentDefinition(ctx, deploymentDef)

	require.NoError(t, err, "Should create deployment definition")
	assert.NotNil(t, result3, "Result should not be nil")
	assert.Equal(t, deploymentDef.Type, result3.Type, "Deployment type should match")
	assert.Equal(t, 2, len(result3.Activities), "Should have 2 activities")

	// Verify the activities are correctly referenced
	activityTypes := make(map[api.ActivityType]bool)
	for _, activity := range result3.Activities {
		activityTypes[activity.Type] = true
	}
	assert.True(t, activityTypes["prepare"], "Should contain prepare activity")
	assert.True(t, activityTypes["deploy"], "Should contain deploy activity")

	// Verify stored definitions can be retrieved
	retrievedDeployment, err := store.FindDeploymentDefinition(deploymentDef.Type)
	require.NoError(t, err, "Should retrieve deployment definition")
	assert.Equal(t, deploymentDef.Type, retrievedDeployment.Type, "Retrieved deployment type should match")

	retrievedActivity1, err := store.FindActivityDefinition("prepare")
	require.NoError(t, err, "Should retrieve first activity definition")
	assert.Equal(t, api.ActivityType("prepare"), retrievedActivity1.Type, "Retrieved activity type should match")

	retrievedActivity2, err := store.FindActivityDefinition("deploy")
	require.NoError(t, err, "Should retrieve second activity definition")
	assert.Equal(t, api.ActivityType("deploy"), retrievedActivity2.Type, "Retrieved activity type should match")
}
