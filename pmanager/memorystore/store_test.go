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

package memorystore

import (
	"github.com/metaform/connector-fabric-manager/common/store"
	"testing"

	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefinitionStore(t *testing.T) {
	definitionStore := newDefinitionStore()

	assert.NotNil(t, definitionStore)
	assert.NotNil(t, definitionStore.deploymentDefinitions)
	assert.NotNil(t, definitionStore.activityDefinitions)
	assert.Equal(t, 0, len(definitionStore.deploymentDefinitions))
	assert.Equal(t, 0, len(definitionStore.activityDefinitions))
}

func TestDefinitionStore_DeploymentDefinition_StoreAndFind(t *testing.T) {
	definitionStore := newDefinitionStore()

	deploymentID := "test-deployment-1"
	definition := &api.DeploymentDefinition{
		Type:       "test-type",
		ApiVersion: "v1",
	}

	definitionStore.StoreDeploymentDefinition(deploymentID, definition)

	result, err := definitionStore.FindDeploymentDefinition(deploymentID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, definition.Type, result.Type)
	assert.Equal(t, definition.ApiVersion, result.ApiVersion)

	// Verify it's a copy (different memory address)
	assert.NotSame(t, definition, result)
}

func TestDefinitionStore_DeploymentDefinition_FindNotFound(t *testing.T) {
	definitionStore := newDefinitionStore()

	result, err := definitionStore.FindDeploymentDefinition("non-existent")

	assert.Error(t, err)
	assert.Equal(t, store.ErrNotFound, err)
	assert.Nil(t, result)
}

func TestDefinitionStore_ActivityDefinition_StoreAndFind(t *testing.T) {
	definitionStore := newDefinitionStore()

	activityID := "test-activity-1"
	definition := &api.ActivityDefinition{
		Type:        "test-activity-type",
		Provider:    "test-provider",
		Description: "Test activity",
	}

	definitionStore.StoreActivityDefinition(activityID, definition)

	result, err := definitionStore.FindActivityDefinition(activityID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, definition.Type, result.Type)
	assert.Equal(t, definition.Provider, result.Provider)
	assert.Equal(t, definition.Description, result.Description)

	// Verify it's a copy (different memory address)
	assert.NotSame(t, definition, result)
}

func TestDefinitionStore_ActivityDefinition_FindNotFound(t *testing.T) {
	definitionStore := newDefinitionStore()

	result, err := definitionStore.FindActivityDefinition("non-existent")

	assert.Error(t, err)
	assert.Equal(t, store.ErrNotFound, err)
	assert.Nil(t, result)
}

func TestDefinitionStore_DeploymentDefinition_Delete(t *testing.T) {
	definitionStore := newDefinitionStore()

	deploymentID := "test-deployment-1"
	definition := &api.DeploymentDefinition{Type: "test-type"}
	definitionStore.StoreDeploymentDefinition(deploymentID, definition)

	_, err := definitionStore.FindDeploymentDefinition(deploymentID)
	assert.NoError(t, err)

	deleted := definitionStore.DeleteDeploymentDefinition(deploymentID)
	assert.True(t, deleted)

	_, err = definitionStore.FindDeploymentDefinition(deploymentID)
	assert.Equal(t, store.ErrNotFound, err)

	deleted = definitionStore.DeleteDeploymentDefinition(deploymentID)
	assert.False(t, deleted)
}

func TestDefinitionStore_ActivityDefinition_Delete(t *testing.T) {
	definitionStore := newDefinitionStore()

	activityID := "test-activity-1"
	definition := &api.ActivityDefinition{Type: "test-type"}
	definitionStore.StoreActivityDefinition(activityID, definition)

	_, err := definitionStore.FindActivityDefinition(activityID)
	assert.NoError(t, err)

	deleted := definitionStore.DeleteActivityDefinition(activityID)
	assert.True(t, deleted)

	_, err = definitionStore.FindActivityDefinition(activityID)
	assert.Equal(t, store.ErrNotFound, err)

	deleted = definitionStore.DeleteActivityDefinition(activityID)
	assert.False(t, deleted)
}

func TestDefinitionStore_DataIsolation(t *testing.T) {
	definitionStore := newDefinitionStore()

	deploymentID := "test-deployment"
	originalDef := &api.DeploymentDefinition{
		Type:       "original-type",
		ApiVersion: "v1",
	}
	definitionStore.StoreDeploymentDefinition(deploymentID, originalDef)

	originalDef.Type = "modified-type"

	retrievedDef, err := definitionStore.FindDeploymentDefinition(deploymentID)
	require.NoError(t, err)

	assert.Equal(t, "original-type", retrievedDef.Type)
	assert.NotEqual(t, originalDef.Type, retrievedDef.Type)

	retrievedDef.Type = "retrieved-modified"

	retrievedDef2, err := definitionStore.FindDeploymentDefinition(deploymentID)
	require.NoError(t, err)
	assert.Equal(t, "original-type", retrievedDef2.Type)
}

func TestDefinitionStore_StoreOverwrite(t *testing.T) {
	definitionStore := newDefinitionStore()

	deploymentID := "test-deployment"

	// Store first definition
	definition1 := &api.DeploymentDefinition{
		Type:       "type1",
		ApiVersion: "v1",
	}
	definitionStore.StoreDeploymentDefinition(deploymentID, definition1)

	// Store second definition with same ID (overwrite)
	definition2 := &api.DeploymentDefinition{
		Type:       "type2",
		ApiVersion: "v2",
	}
	definitionStore.StoreDeploymentDefinition(deploymentID, definition2)

	// Verify the second definition is stored
	result, err := definitionStore.FindDeploymentDefinition(deploymentID)
	require.NoError(t, err)
	assert.Equal(t, "type2", result.Type)
	assert.Equal(t, "v2", result.ApiVersion)

	// Verify only one definition exists
	defintions, _, err := definitionStore.ListDeploymentDefinitions(0, 1000)
	assert.Equal(t, 1, len(defintions))
}

func TestDefinitionStore_ListDeploymentDefinitions_WithPagination(t *testing.T) {
	definitionStore := newDefinitionStore()

	// Test empty store
	definitions, hasMore, err := definitionStore.ListDeploymentDefinitions(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(definitions))
	assert.False(t, hasMore)

	// Add test data
	def1 := &api.DeploymentDefinition{Type: "type1", ApiVersion: "v1"}
	def2 := &api.DeploymentDefinition{Type: "type2", ApiVersion: "v2"}
	def3 := &api.DeploymentDefinition{Type: "type3", ApiVersion: "v3"}
	def4 := &api.DeploymentDefinition{Type: "type4", ApiVersion: "v4"}
	def5 := &api.DeploymentDefinition{Type: "type5", ApiVersion: "v5"}

	definitionStore.StoreDeploymentDefinition("deploy-1", def1)
	definitionStore.StoreDeploymentDefinition("deploy-2", def2)
	definitionStore.StoreDeploymentDefinition("deploy-3", def3)
	definitionStore.StoreDeploymentDefinition("deploy-4", def4)
	definitionStore.StoreDeploymentDefinition("deploy-5", def5)

	// Test first page
	definitions, hasMore, err = definitionStore.ListDeploymentDefinitions(0, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(definitions))
	assert.True(t, hasMore)

	// Test second page
	definitions, hasMore, err = definitionStore.ListDeploymentDefinitions(2, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(definitions))
	assert.True(t, hasMore)

	// Test last page
	definitions, hasMore, err = definitionStore.ListDeploymentDefinitions(4, 2)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(definitions))
	assert.False(t, hasMore)

	// Test get all
	definitions, hasMore, err = definitionStore.ListDeploymentDefinitions(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(definitions))
	assert.False(t, hasMore)

	// Test offset beyond total
	definitions, hasMore, err = definitionStore.ListDeploymentDefinitions(10, 2)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(definitions))
	assert.False(t, hasMore)

	// Test partial last page
	definitions, hasMore, err = definitionStore.ListDeploymentDefinitions(3, 5)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(definitions))
	assert.False(t, hasMore)
}

func TestDefinitionStore_ListDeploymentDefinitions_ValidationErrors(t *testing.T) {
	definitionStore := newDefinitionStore()

	// Test negative offset
	_, _, err := definitionStore.ListDeploymentDefinitions(-1, 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "offset cannot be negative")

	// Test zero limit
	_, _, err = definitionStore.ListDeploymentDefinitions(0, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit must be positive")

	// Test negative limit
	_, _, err = definitionStore.ListDeploymentDefinitions(0, -1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit must be positive")
}

func TestDefinitionStore_ListDeploymentDefinitions_DataIsolation(t *testing.T) {
	definitionStore := newDefinitionStore()

	originalDef := &api.DeploymentDefinition{Type: "original", ApiVersion: "v1"}
	definitionStore.StoreDeploymentDefinition("deploy-1", originalDef)

	definitions, hasMore, err := definitionStore.ListDeploymentDefinitions(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(definitions))
	assert.False(t, hasMore)

	// Modify the original definition
	originalDef.Type = "modified"

	// Verify returned definition is not affected
	assert.Equal(t, "original", definitions[0].Type)

	// Modify the returned definition
	definitions[0].Type = "returned-modified"

	// Verify stored definition is not affected
	storedDef, err := definitionStore.FindDeploymentDefinition("deploy-1")
	assert.NoError(t, err)
	assert.Equal(t, "original", storedDef.Type)
}

func TestDefinitionStore_ListActivityDefinitions_WithPagination(t *testing.T) {
	definitionStore := newDefinitionStore()

	// Test empty store
	definitions, hasMore, err := definitionStore.ListActivityDefinitions(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(definitions))
	assert.False(t, hasMore)

	// Add test data
	def1 := &api.ActivityDefinition{Type: "type1", Provider: "provider1", Description: "desc1"}
	def2 := &api.ActivityDefinition{Type: "type2", Provider: "provider2", Description: "desc2"}
	def3 := &api.ActivityDefinition{Type: "type3", Provider: "provider3", Description: "desc3"}
	def4 := &api.ActivityDefinition{Type: "type4", Provider: "provider4", Description: "desc4"}
	def5 := &api.ActivityDefinition{Type: "type5", Provider: "provider5", Description: "desc5"}

	definitionStore.StoreActivityDefinition("activity-1", def1)
	definitionStore.StoreActivityDefinition("activity-2", def2)
	definitionStore.StoreActivityDefinition("activity-3", def3)
	definitionStore.StoreActivityDefinition("activity-4", def4)
	definitionStore.StoreActivityDefinition("activity-5", def5)

	// Test first page
	definitions, hasMore, err = definitionStore.ListActivityDefinitions(0, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(definitions))
	assert.True(t, hasMore)

	// Test second page
	definitions, hasMore, err = definitionStore.ListActivityDefinitions(2, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(definitions))
	assert.True(t, hasMore)

	// Test last page
	definitions, hasMore, err = definitionStore.ListActivityDefinitions(4, 2)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(definitions))
	assert.False(t, hasMore)

	// Test get all
	definitions, hasMore, err = definitionStore.ListActivityDefinitions(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(definitions))
	assert.False(t, hasMore)

	// Test offset beyond total
	definitions, hasMore, err = definitionStore.ListActivityDefinitions(10, 2)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(definitions))
	assert.False(t, hasMore)

	// Test partial last page
	definitions, hasMore, err = definitionStore.ListActivityDefinitions(3, 5)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(definitions))
	assert.False(t, hasMore)
}

func TestDefinitionStore_ListActivityDefinitions_ValidationErrors(t *testing.T) {
	definitionStore := newDefinitionStore()

	// Test negative offset
	_, _, err := definitionStore.ListActivityDefinitions(-1, 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "offset cannot be negative")

	// Test zero limit
	_, _, err = definitionStore.ListActivityDefinitions(0, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit must be positive")

	// Test negative limit
	_, _, err = definitionStore.ListActivityDefinitions(0, -1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit must be positive")
}

func TestDefinitionStore_ListActivityDefinitions_DataIsolation(t *testing.T) {
	definitionStore := newDefinitionStore()

	originalDef := &api.ActivityDefinition{Type: "original", Provider: "provider1", Description: "desc1"}
	definitionStore.StoreActivityDefinition("activity-1", originalDef)

	definitions, hasMore, err := definitionStore.ListActivityDefinitions(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(definitions))
	assert.False(t, hasMore)

	// Modify the original definition
	originalDef.Type = "modified"

	// Verify returned definition is not affected
	assert.Equal(t, "original", definitions[0].Type)

	// Modify the returned definition
	definitions[0].Type = "returned-modified"

	// Verify stored definition is not affected
	storedDef, err := definitionStore.FindActivityDefinition("activity-1")
	assert.NoError(t, err)
	assert.Equal(t, "original", storedDef.Type)
}
