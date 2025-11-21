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

func TestDefinitionManager_CreateOrchestrationDefinition_Success(t *testing.T) {
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
	_, err := store.StoreActivityDefinition(context.Background(), activityDef)
	require.NoError(t, err, "Failed to store activity definition")

	// Create an orchestration definition that references the activity
	orchestrationDef := &api.OrchestrationDefinition{
		Type:   model.OrchestrationType("test-orchestration"),
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

	result, err := manager.CreateOrchestrationDefinition(ctx, orchestrationDef)

	require.NoError(t, err, "CreateOrchestrationDefinition should succeed")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, orchestrationDef.Type, result.Type, "Orchestration type should match")
	assert.Equal(t, orchestrationDef.Active, result.Active, "Active flag should match")
	assert.Equal(t, len(orchestrationDef.Activities), len(result.Activities), "Activities count should match")
}

func TestDefinitionManager_CreateOrchestrationDefinition_MissingActivityDefinition(t *testing.T) {
	// Given
	store := memorystore.NewDefinitionStore()
	manager := definitionManager{
		trxContext: cstore.NoOpTransactionContext{},
		store:      store,
	}
	ctx := context.Background()

	// Create an orchestration definition that references a non-existent activity
	orchestrationDef := &api.OrchestrationDefinition{
		Type:   model.OrchestrationType("test-orchestration"),
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

	result, err := manager.CreateOrchestrationDefinition(ctx, orchestrationDef)

	require.Error(t, err, "CreateOrchestrationDefinition should fail when activity definition is missing")
	assert.Nil(t, result, "Result should be nil on error")

	// Verify the error is a client error about the missing activity
	var clientErr types.ClientError
	require.True(t, errors.As(err, &clientErr), "Error should be a ClientError")
	assert.Contains(t, err.Error(), "activity type 'non-existent-activity' not found",
		"Error message should mention the missing activity type")
}

func TestDefinitionManager_CreateOrchestrationDefinition_MultipleMissingActivityDefinitions(t *testing.T) {
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
	_, err := store.StoreActivityDefinition(context.Background(), activityDef)
	require.NoError(t, err, "Failed to store valid activity definition")

	// Create orchestration definition with mix of valid and invalid activity references
	orchestrationDef := &api.OrchestrationDefinition{
		Type:   model.OrchestrationType("mixed-orchestration"),
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

	result, err := manager.CreateOrchestrationDefinition(ctx, orchestrationDef)

	require.Error(t, err, "CreateOrchestrationDefinition should fail when multiple activity definitions are missing")
	assert.Nil(t, result, "Result should be nil on error")

	// Verify the error mentions both missing activities
	assert.Contains(t, err.Error(), "missing-activity-1", "Error should mention first missing activity")
	assert.Contains(t, err.Error(), "missing-activity-2", "Error should mention second missing activity")
}

func TestDefinitionManager_CreateOrchestrationDefinition_EmptyActivities(t *testing.T) {
	store := memorystore.NewDefinitionStore()
	manager := definitionManager{
		trxContext: cstore.NoOpTransactionContext{},
		store:      store,
	}
	ctx := context.Background()

	// Create an orchestration definition with no activities
	orchestrationDef := &api.OrchestrationDefinition{
		Type:       model.OrchestrationType("empty-orchestration"),
		Active:     true,
		Schema:     map[string]any{"type": "object"},
		Activities: []api.Activity{}, // Empty activities slice
	}

	result, err := manager.CreateOrchestrationDefinition(ctx, orchestrationDef)

	require.NoError(t, err, "CreateOrchestrationDefinition should succeed with empty activities")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Equal(t, 0, len(result.Activities), "Should have 0 activities")
}

func TestDefinitionManager_CreateOrchestrationDefinition_StoreError(t *testing.T) {
	store := memorystore.NewDefinitionStore()
	manager := definitionManager{
		trxContext: cstore.NoOpTransactionContext{},
		store:      store,
	}
	ctx := context.Background()

	// Create and store the orchestration definition first
	orchestrationDef := &api.OrchestrationDefinition{
		Type:       model.OrchestrationType("duplicate-orchestration"),
		Active:     true,
		Schema:     map[string]any{"type": "object"},
		Activities: []api.Activity{},
	}

	_, err := store.StoreOrchestrationDefinition(context.Background(), orchestrationDef)
	require.NoError(t, err, "First store should succeed")

	// Attempt to store the same orchestration definition again
	result, err := manager.CreateOrchestrationDefinition(ctx, orchestrationDef)

	require.Error(t, err, "CreateOrchestrationDefinition should fail on duplicate")
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

	_, err := store.StoreActivityDefinition(context.Background(), activityDef)
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
		Description:  "Send application",
		InputSchema:  map[string]any{"type": "object"},
		OutputSchema: map[string]any{"type": "object"},
	}

	result1, err := manager.CreateActivityDefinition(ctx, activityDef1)
	require.NoError(t, err, "Should create first activity definition")
	assert.Equal(t, activityDef1.Type, result1.Type, "First activity type should match")

	result2, err := manager.CreateActivityDefinition(ctx, activityDef2)
	require.NoError(t, err, "Should create second activity definition")
	assert.Equal(t, activityDef2.Type, result2.Type, "Second activity type should match")

	// Create an orchestration definition that uses both activities
	orchestrationDef := &api.OrchestrationDefinition{
		Type:   model.OrchestrationType("full-orchestration"),
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
				ID:        "orchestration-step",
				Type:      "deploy",
				Inputs:    []api.MappingEntry{{Source: "artifact", Target: "deployment_artifact"}},
				DependsOn: []string{"prepare-step"},
			},
		},
	}

	result3, err := manager.CreateOrchestrationDefinition(ctx, orchestrationDef)

	require.NoError(t, err, "Should create orchestration definition")
	assert.NotNil(t, result3, "Result should not be nil")
	assert.Equal(t, orchestrationDef.Type, result3.Type, "Orchestration type should match")
	assert.Equal(t, 2, len(result3.Activities), "Should have 2 activities")

	// Verify the activities are correctly referenced
	activityTypes := make(map[api.ActivityType]bool)
	for _, activity := range result3.Activities {
		activityTypes[activity.Type] = true
	}
	assert.True(t, activityTypes["prepare"], "Should contain prepare activity")
	assert.True(t, activityTypes["deploy"], "Should contain deploy activity")

	// Verify stored definitions can be retrieved
	retrievedOrchestration, err := store.FindOrchestrationDefinition(context.Background(), orchestrationDef.Type)
	require.NoError(t, err, "Should retrieve orchestration definition")
	assert.Equal(t, orchestrationDef.Type, retrievedOrchestration.Type, "Retrieved deployment type should match")

	retrievedActivity1, err := store.FindActivityDefinition(ctx, "prepare")
	require.NoError(t, err, "Should retrieve first activity definition")
	assert.Equal(t, api.ActivityType("prepare"), retrievedActivity1.Type, "Retrieved activity type should match")

	retrievedActivity2, err := store.FindActivityDefinition(ctx, "deploy")
	require.NoError(t, err, "Should retrieve second activity definition")
	assert.Equal(t, api.ActivityType("deploy"), retrievedActivity2.Type, "Retrieved activity type should match")
}
