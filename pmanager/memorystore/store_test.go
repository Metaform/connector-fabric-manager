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
	"testing"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefinitionStore(t *testing.T) {
	definitionStore := NewDefinitionStore()

	assert.NotNil(t, definitionStore)
	assert.NotNil(t, definitionStore.orchestrationDefinitions)
	assert.NotNil(t, definitionStore.activityDefinitions)
	assert.Equal(t, 0, len(definitionStore.orchestrationDefinitions))
	assert.Equal(t, 0, len(definitionStore.activityDefinitions))
}

func TestDefinitionStore_OrchestrationDefinition_StoreAndFind(t *testing.T) {
	definitionStore := NewDefinitionStore()

	var oType model.OrchestrationType = "test-orchestration-1"
	definition := &api.OrchestrationDefinition{
		Type: oType,
	}

	definitionStore.StoreOrchestrationDefinition(definition)

	result, err := definitionStore.FindOrchestrationDefinition(oType)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, definition.Type, result.Type)

	// Verify it's a copy (different memory address)
	assert.NotSame(t, definition, result)
}

func TestDefinitionStore_OrchestrationDefinition_FindNotFound(t *testing.T) {
	definitionStore := NewDefinitionStore()

	result, err := definitionStore.FindOrchestrationDefinition("non-existent")

	assert.Error(t, err)
	assert.Equal(t, store.ErrNotFound, err)
	assert.Nil(t, result)
}

func TestDefinitionStore_ActivityDefinition_StoreAndFind(t *testing.T) {
	definitionStore := NewDefinitionStore()

	var activityType api.ActivityType = "test-activity-1"
	definition := &api.ActivityDefinition{
		Type:        activityType,
		Description: "Test activity",
	}

	definitionStore.StoreActivityDefinition(definition)

	result, err := definitionStore.FindActivityDefinition(activityType)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, definition.Type, result.Type)
	assert.Equal(t, definition.Description, result.Description)

	// Verify it's a copy (different memory address)
	assert.NotSame(t, definition, result)
}

func TestDefinitionStore_ActivityDefinition_FindNotFound(t *testing.T) {
	definitionStore := NewDefinitionStore()

	result, err := definitionStore.FindActivityDefinition("non-existent")

	assert.Error(t, err)
	assert.Equal(t, store.ErrNotFound, err)
	assert.Nil(t, result)
}

func TestDefinitionStore_OrchestrationDefinition_Delete(t *testing.T) {
	definitionStore := NewDefinitionStore()

	var oType model.OrchestrationType = "test-orchestration-1"
	definition := &api.OrchestrationDefinition{Type: oType}
	definitionStore.StoreOrchestrationDefinition(definition)

	_, err := definitionStore.FindOrchestrationDefinition(oType)
	assert.NoError(t, err)

	deleted, err := definitionStore.DeleteOrchestrationDefinition(oType)
	assert.Nil(t, err)
	assert.True(t, deleted)

	_, err = definitionStore.FindOrchestrationDefinition(oType)
	assert.Equal(t, store.ErrNotFound, err)

	deleted, err = definitionStore.DeleteOrchestrationDefinition(oType)
	assert.Nil(t, err)
	assert.False(t, deleted)
}

func TestDefinitionStore_ActivityDefinition_Delete(t *testing.T) {
	definitionStore := NewDefinitionStore()

	var activityType api.ActivityType = "test-activity-1"
	definition := &api.ActivityDefinition{Type: activityType}
	definitionStore.StoreActivityDefinition(definition)

	_, err := definitionStore.FindActivityDefinition(activityType)
	assert.NoError(t, err)

	deleted, err := definitionStore.DeleteActivityDefinition(activityType)
	assert.Nil(t, err)
	assert.True(t, deleted)

	_, err = definitionStore.FindActivityDefinition(activityType)
	assert.Equal(t, store.ErrNotFound, err)

	deleted, err = definitionStore.DeleteActivityDefinition(activityType)
	assert.Nil(t, err)
	assert.False(t, deleted)
}

func TestDefinitionStore_DataIsolation(t *testing.T) {
	definitionStore := NewDefinitionStore()

	var originalType model.OrchestrationType = "original-type"
	originalDef := &api.OrchestrationDefinition{
		Type: originalType,
	}
	definitionStore.StoreOrchestrationDefinition(originalDef)

	originalDef.Type = "modified-type"

	retrievedDef, err := definitionStore.FindOrchestrationDefinition(originalType)
	require.NoError(t, err)

	assert.Equal(t, originalType, retrievedDef.Type)
	assert.NotEqual(t, originalDef.Type, retrievedDef.Type)

	retrievedDef.Type = "retrieved-modified"

	retrievedDef2, err := definitionStore.FindOrchestrationDefinition(originalType)
	require.NoError(t, err)
	assert.Equal(t, originalType, retrievedDef2.Type)
}

func TestDefinitionStore_StoreOverwrite(t *testing.T) {
	definitionStore := NewDefinitionStore()

	var dType model.OrchestrationType = "test-orchestration"

	// Store first definition
	definition1 := &api.OrchestrationDefinition{
		Type: dType,
	}
	definitionStore.StoreOrchestrationDefinition(definition1)

	// Store second definition with same ID (overwrite)
	definition2 := &api.OrchestrationDefinition{
		Type: dType,
	}
	definitionStore.StoreOrchestrationDefinition(definition2)

	// Verify the second definition is stored
	result, err := definitionStore.FindOrchestrationDefinition(dType)
	require.NoError(t, err)
	assert.Equal(t, dType, result.Type)

	// Verify only one definition exists
	defintions, _, err := definitionStore.ListOrchestrationDefinitions(0, 1000)
	assert.Equal(t, 1, len(defintions))
}

func TestDefinitionStore_ListOrchestrationDefinitions_WithPagination(t *testing.T) {
	definitionStore := NewDefinitionStore()

	// Test empty store
	definitions, hasMore, err := definitionStore.ListOrchestrationDefinitions(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(definitions))
	assert.False(t, hasMore)

	// Add test data
	def1 := &api.OrchestrationDefinition{Type: "type1"}
	def2 := &api.OrchestrationDefinition{Type: "type2"}
	def3 := &api.OrchestrationDefinition{Type: "type3"}
	def4 := &api.OrchestrationDefinition{Type: "type4"}
	def5 := &api.OrchestrationDefinition{Type: "type5"}

	definitionStore.StoreOrchestrationDefinition(def1)
	definitionStore.StoreOrchestrationDefinition(def2)
	definitionStore.StoreOrchestrationDefinition(def3)
	definitionStore.StoreOrchestrationDefinition(def4)
	definitionStore.StoreOrchestrationDefinition(def5)

	// Test first page
	definitions, hasMore, err = definitionStore.ListOrchestrationDefinitions(0, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(definitions))
	assert.True(t, hasMore)

	// Test second page
	definitions, hasMore, err = definitionStore.ListOrchestrationDefinitions(2, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(definitions))
	assert.True(t, hasMore)

	// Test last page
	definitions, hasMore, err = definitionStore.ListOrchestrationDefinitions(4, 2)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(definitions))
	assert.False(t, hasMore)

	// Test get all
	definitions, hasMore, err = definitionStore.ListOrchestrationDefinitions(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(definitions))
	assert.False(t, hasMore)

	// Test offset beyond total
	definitions, hasMore, err = definitionStore.ListOrchestrationDefinitions(10, 2)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(definitions))
	assert.False(t, hasMore)

	// Test partial last page
	definitions, hasMore, err = definitionStore.ListOrchestrationDefinitions(3, 5)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(definitions))
	assert.False(t, hasMore)
}

func TestDefinitionStore_ListOrchestrationDefinitions_ValidationErrors(t *testing.T) {
	definitionStore := NewDefinitionStore()

	// Test negative offset
	_, _, err := definitionStore.ListOrchestrationDefinitions(-1, 10)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "offset cannot be negative")

	// Test zero limit
	_, _, err = definitionStore.ListOrchestrationDefinitions(0, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit must be positive")

	// Test negative limit
	_, _, err = definitionStore.ListOrchestrationDefinitions(0, -1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit must be positive")
}

func TestDefinitionStore_ListOrchestrationDefinitions_DataIsolation(t *testing.T) {
	definitionStore := NewDefinitionStore()

	var originalType model.OrchestrationType = "original"
	originalDef := &api.OrchestrationDefinition{Type: originalType}
	definitionStore.StoreOrchestrationDefinition(originalDef)

	definitions, hasMore, err := definitionStore.ListOrchestrationDefinitions(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(definitions))
	assert.False(t, hasMore)

	// Modify the original definition
	originalDef.Type = "modified"

	// Verify returned definition is not affected
	assert.Equal(t, originalType, definitions[0].Type)

	// Modify the returned definition
	definitions[0].Type = "returned-modified"

	// Verify stored definition is not affected
	storedDef, err := definitionStore.FindOrchestrationDefinition(originalType)
	assert.NoError(t, err)
	assert.Equal(t, originalType, storedDef.Type)
}

func TestDefinitionStore_ListActivityDefinitions_WithPagination(t *testing.T) {
	definitionStore := NewDefinitionStore()

	// Test empty store
	definitions, hasMore, err := definitionStore.ListActivityDefinitions(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(definitions))
	assert.False(t, hasMore)

	// Add test data
	def1 := &api.ActivityDefinition{Type: "type1", Description: "desc1"}
	def2 := &api.ActivityDefinition{Type: "type2", Description: "desc2"}
	def3 := &api.ActivityDefinition{Type: "type3", Description: "desc3"}
	def4 := &api.ActivityDefinition{Type: "type4", Description: "desc4"}
	def5 := &api.ActivityDefinition{Type: "type5", Description: "desc5"}

	definitionStore.StoreActivityDefinition(def1)
	definitionStore.StoreActivityDefinition(def2)
	definitionStore.StoreActivityDefinition(def3)
	definitionStore.StoreActivityDefinition(def4)
	definitionStore.StoreActivityDefinition(def5)

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
	definitionStore := NewDefinitionStore()

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
	definitionStore := NewDefinitionStore()

	var originalType api.ActivityType = "original"
	originalDef := &api.ActivityDefinition{Type: originalType, Description: "desc1"}
	definitionStore.StoreActivityDefinition(originalDef)

	definitions, hasMore, err := definitionStore.ListActivityDefinitions(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(definitions))
	assert.False(t, hasMore)

	// Modify the original definition
	originalDef.Type = "modified"

	// Verify returned definition is not affected
	assert.Equal(t, originalType, definitions[0].Type)

	// Modify the returned definition
	definitions[0].Type = "returned-modified"

	// Verify stored definition is not affected
	storedDef, err := definitionStore.FindActivityDefinition(originalType)
	assert.NoError(t, err)
	assert.Equal(t, originalType, storedDef.Type)
}
