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
	"context"
	"fmt"
	"iter"
	"sync"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/query"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

var defaultMatcher = &query.DefaultFieldMatcher{}

// MemoryDefinitionStore is a thread-safe in-memory store for orchestration and activity definitions.
type MemoryDefinitionStore struct {
	mutex                    sync.RWMutex
	orchestrationDefinitions map[string]*api.OrchestrationDefinition
	activityDefinitions      map[string]*api.ActivityDefinition
}

// NewDefinitionStore creates a new thread-safe in-memory definition store
func NewDefinitionStore() *MemoryDefinitionStore {
	return &MemoryDefinitionStore{
		orchestrationDefinitions: make(map[string]*api.OrchestrationDefinition),
		activityDefinitions:      make(map[string]*api.ActivityDefinition),
	}
}

func (d *MemoryDefinitionStore) FindOrchestrationDefinition(_ context.Context, orchestrationType model.OrchestrationType) (*api.OrchestrationDefinition, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	definition, exists := d.orchestrationDefinitions[orchestrationType.String()]
	if !exists {
		return nil, store.ErrNotFound
	}

	// Return a copy to prevent external modifications
	definitionCopy := *definition
	return &definitionCopy, nil
}

func (d *MemoryDefinitionStore) FindOrchestrationDefinitionsByPredicate(_ context.Context, predicate query.Predicate) iter.Seq2[api.OrchestrationDefinition, error] {
	return findDefinitionsByPredicate(d, predicate, d.orchestrationDefinitions)
}

func (d *MemoryDefinitionStore) ExistsOrchestrationDefinition(_ context.Context, orchestrationType model.OrchestrationType) (bool, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	_, exists := d.orchestrationDefinitions[orchestrationType.String()]
	return exists, nil
}

func (d *MemoryDefinitionStore) GetOrchestrationDefinitionCount(_ context.Context, predicate query.Predicate) (int, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return countByPredicate(d.orchestrationDefinitions, predicate), nil
}

func (d *MemoryDefinitionStore) FindActivityDefinition(_ context.Context, activityType api.ActivityType) (*api.ActivityDefinition, error) {
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

func (d *MemoryDefinitionStore) FindActivityDefinitionsByPredicate(_ context.Context, predicate query.Predicate) iter.Seq2[api.ActivityDefinition, error] {
	return findDefinitionsByPredicate(d, predicate, d.activityDefinitions)
}

func (d *MemoryDefinitionStore) ExistsActivityDefinition(_ context.Context, activityType api.ActivityType) (bool, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	_, exists := d.activityDefinitions[activityType.String()]
	return exists, nil
}

func (d *MemoryDefinitionStore) GetActivityDefinitionCount(_ context.Context, predicate query.Predicate) (int, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return countByPredicate(d.activityDefinitions, predicate), nil
}

func (d *MemoryDefinitionStore) StoreOrchestrationDefinition(_ context.Context, definition *api.OrchestrationDefinition) (*api.OrchestrationDefinition, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.orchestrationDefinitions[definition.Type.String()] != nil {
		return nil, types.ErrConflict
	}

	// Store a copy to prevent external modifications
	definitionCopy := *definition
	d.orchestrationDefinitions[definitionCopy.Type.String()] = &definitionCopy
	return definition, nil
}

func (d *MemoryDefinitionStore) StoreActivityDefinition(_ context.Context, definition *api.ActivityDefinition) (*api.ActivityDefinition, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.activityDefinitions[definition.Type.String()] != nil {
		return nil, types.ErrConflict
	}
	// Store a copy to prevent external modifications
	definitionCopy := *definition
	d.activityDefinitions[definitionCopy.Type.String()] = &definitionCopy
	return definition, nil
}

func (d *MemoryDefinitionStore) DeleteOrchestrationDefinition(_ context.Context, orchestrationType model.OrchestrationType) (bool, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, exists := d.orchestrationDefinitions[orchestrationType.String()]
	if exists {
		delete(d.orchestrationDefinitions, orchestrationType.String())
	}
	return exists, nil
}

func (d *MemoryDefinitionStore) DeleteActivityDefinition(_ context.Context, activityType api.ActivityType) (bool, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, exists := d.activityDefinitions[activityType.String()]
	if exists {
		delete(d.activityDefinitions, activityType.String())
	}
	return exists, nil
}

func (d *MemoryDefinitionStore) ListOrchestrationDefinitions(_ context.Context, offset, limit int) ([]*api.OrchestrationDefinition, bool, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return listDefinitions[api.OrchestrationDefinition](d.orchestrationDefinitions, offset, limit)
}

func (d *MemoryDefinitionStore) ListActivityDefinitions(_ context.Context, offset, limit int) ([]*api.ActivityDefinition, bool, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return listDefinitions[api.ActivityDefinition](d.activityDefinitions, offset, limit)
}

// Clear removes all stored definitions
func (d *MemoryDefinitionStore) Clear() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.orchestrationDefinitions = make(map[string]*api.OrchestrationDefinition)
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

// findDefinitionsByPredicate is a generic helper that filters definitions by predicate
func findDefinitionsByPredicate[T any](d *MemoryDefinitionStore, predicate query.Predicate, definitionMap map[string]*T) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		d.mutex.RLock()
		defer d.mutex.RUnlock()

		for _, definition := range definitionMap {
			if predicate.Matches(definition, defaultMatcher) {
				// Return a copy to prevent external modifications
				definitionCopy := *definition
				if !yield(definitionCopy, nil) {
					return
				}
			}
		}
	}
}

// countByPredicate is a generic helper that counts definitions matching a predicate
func countByPredicate[T any](definitionMap map[string]*T, predicate query.Predicate) int {
	count := 0
	for _, definition := range definitionMap {
		if predicate.Matches(definition, defaultMatcher) {
			count++
		}
	}
	return count
}
