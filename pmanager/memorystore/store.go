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
	"fmt"
	"sync"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

// MemoryDefinitionStore is a thread-safe in-memory store for deployment and activity definitions.
type MemoryDefinitionStore struct {
	mutex                 sync.RWMutex
	deploymentDefinitions map[string]*api.DeploymentDefinition
	activityDefinitions   map[string]*api.ActivityDefinition
}

// NewDefinitionStore creates a new thread-safe in-memory definition store
func NewDefinitionStore() *MemoryDefinitionStore {
	return &MemoryDefinitionStore{
		deploymentDefinitions: make(map[string]*api.DeploymentDefinition),
		activityDefinitions:   make(map[string]*api.ActivityDefinition),
	}
}

func (d *MemoryDefinitionStore) FindDeploymentDefinition(deploymentType model.DeploymentType) (*api.DeploymentDefinition, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	definition, exists := d.deploymentDefinitions[deploymentType.String()]
	if !exists {
		return nil, store.ErrNotFound
	}

	// Return a copy to prevent external modifications
	definitionCopy := *definition
	return &definitionCopy, nil
}

func (d *MemoryDefinitionStore) FindActivityDefinition(activityType api.ActivityType) (*api.ActivityDefinition, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	definition, exists := d.activityDefinitions[activityType.String()]
	if !exists {
		return nil, store.ErrNotFound
	}

	// Return a copy to prevent external modifications
	definitionCopy := *definition
	return &definitionCopy, nil
}

func (d *MemoryDefinitionStore) StoreDeploymentDefinition(definition *api.DeploymentDefinition) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Store a copy to prevent external modifications
	definitionCopy := *definition
	d.deploymentDefinitions[definitionCopy.Type.String()] = &definitionCopy
}

func (d *MemoryDefinitionStore) StoreActivityDefinition(definition *api.ActivityDefinition) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Store a copy to prevent external modifications
	definitionCopy := *definition
	d.activityDefinitions[definitionCopy.Type.String()] = &definitionCopy
}

func (d *MemoryDefinitionStore) DeleteDeploymentDefinition(deploymentType model.DeploymentType) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, exists := d.deploymentDefinitions[deploymentType.String()]
	if exists {
		delete(d.deploymentDefinitions, deploymentType.String())
	}
	return exists
}

func (d *MemoryDefinitionStore) DeleteActivityDefinition(activityType api.ActivityType) bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, exists := d.activityDefinitions[activityType.String()]
	if exists {
		delete(d.activityDefinitions, activityType.String())
	}
	return exists
}

func (d *MemoryDefinitionStore) ListDeploymentDefinitions(offset, limit int) ([]*api.DeploymentDefinition, bool, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return listDefinitions[api.DeploymentDefinition](d.deploymentDefinitions, offset, limit)
}

func (d *MemoryDefinitionStore) ListActivityDefinitions(offset, limit int) ([]*api.ActivityDefinition, bool, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return listDefinitions[api.ActivityDefinition](d.activityDefinitions, offset, limit)
}

// Clear removes all stored definitions
func (d *MemoryDefinitionStore) Clear() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.deploymentDefinitions = make(map[string]*api.DeploymentDefinition)
	d.activityDefinitions = make(map[string]*api.ActivityDefinition)
}

// listDefinitions lists definitions with pagination
func listDefinitions[T any](definitionMap map[string]*T, offset, limit int) ([]*T, bool, error) {
	if offset < 0 {
		return nil, false, fmt.Errorf("offset cannot be negative")
	}
	if limit <= 0 {
		return nil, false, fmt.Errorf("limit must be positive")
	}

	// Get all definitions
	allDefinitions := make([]*T, 0, len(definitionMap))
	for _, definition := range definitionMap {
		// Return a copy to prevent external modifications
		definitionCopy := *definition
		allDefinitions = append(allDefinitions, &definitionCopy)
	}

	total := len(allDefinitions)

	// Check overflow
	if offset >= total {
		return []*T{}, false, nil
	}

	// Calculate end index
	end := offset + limit
	if end > total {
		end = total
	}

	hasMore := end < total
	return allDefinitions[offset:end], hasMore, nil
}
