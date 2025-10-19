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

package tmstore

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

// InMemoryTManagerStore is an in-memory implementation of the TManagerStore interface.
type InMemoryTManagerStore struct {
	dProfileStorage   storage[api.DataspaceProfile]
	cellStorage       storage[api.Cell]
	deploymentStorage storage[api.DeploymentRecord]
}

func NewInMemoryTManagerStore(seed bool) *InMemoryTManagerStore {
	store := &InMemoryTManagerStore{
		dProfileStorage:   newStorage(func(p *api.DataspaceProfile) string { return p.ID }),
		cellStorage:       newStorage(func(c *api.Cell) string { return c.ID }),
		deploymentStorage: newStorage(func(r *api.DeploymentRecord) string { return r.ID }),
	}
	if seed {
		cells, profiles := seedData()
		for _, cell := range cells {
			store.cellStorage.Create(cell)
		}
		for _, profile := range profiles {
			store.dProfileStorage.Create(profile)
		}
	}
	return store
}

func (s *InMemoryTManagerStore) GetCells() ([]api.Cell, error) {
	return s.cellStorage.GetAll(), nil
}

func (s *InMemoryTManagerStore) GetDataspaceProfiles() ([]api.DataspaceProfile, error) {
	return s.dProfileStorage.GetAll(), nil
}

func (s *InMemoryTManagerStore) FindDeployment(id string) (*api.DeploymentRecord, error) {
	record := s.deploymentStorage.FindById(id)
	if record == nil {
		return nil, model.ErrNotFound
	}
	return record, nil
}

func (s *InMemoryTManagerStore) DeploymentExists(id string) (bool, error) {
	record := s.deploymentStorage.FindById(id)
	if record == nil {
		return false, nil
	}
	return true, nil
}

func (s *InMemoryTManagerStore) CreateDeployment(record api.DeploymentRecord) (*api.DeploymentRecord, error) {
	record.ID = uuid.New().String()
	err := s.deploymentStorage.Create(record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *InMemoryTManagerStore) UpdateDeployment(record api.DeploymentRecord) error {
	return s.deploymentStorage.Save(&record)
}

type storage[T any] struct {
	cache  map[string]T
	idFunc func(*T) string
	mu     sync.RWMutex
}

func newStorage[T any](idFunc func(*T) string) storage[T] {
	return storage[T]{
		cache:  make(map[string]T),
		idFunc: idFunc,
		mu:     sync.RWMutex{},
	}
}

func (s *storage[T]) FindById(id string) *T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entity, exists := s.cache[id]
	if !exists {
		return nil
	}

	return &entity
}

func (s *storage[T]) Create(entity T) error {
	if s.idFunc(&entity) == "" {
		return model.ErrInvalidInput
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.cache[s.idFunc(&entity)]; exists {
		return model.ErrConflict
	}

	s.cache[s.idFunc(&entity)] = entity
	return nil
}

func (s *storage[T]) Save(entity *T) error {
	if entity == nil {
		return model.ErrInvalidInput
	}
	if s.idFunc(entity) == "" {
		return model.ErrInvalidInput
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.cache[s.idFunc(entity)]; !exists {
		return model.ErrNotFound
	}

	s.cache[s.idFunc(entity)] = *entity
	return nil
}

func (s *storage[T]) Delete(ctx context.Context, id string) error {
	if id == "" {
		return model.ErrInvalidInput
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.cache[id]; !exists {
		return model.ErrNotFound
	}

	delete(s.cache, id)
	return nil
}

func (s *storage[T]) GetAll() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]T, 0, len(s.cache))
	for _, entity := range s.cache {
		result = append(result, entity)
	}
	return result
}

// seedData temporary function to initialize and return sample cells and dataspace profiles for use in deployment workflows.
func seedData() ([]api.Cell, []api.DataspaceProfile) {
	cells := []api.Cell{
		{
			DeployableEntity: api.DeployableEntity{
				Entity: api.Entity{
					ID:      "cell-001",
					Version: 1,
				},
				State:          api.DeploymentStateActive,
				StateTimestamp: time.Now(),
			},
			Properties: api.Properties{
				"region": "us-east-1",
				"type":   "kubernetes",
			},
		},
	}

	dProfiles := []api.DataspaceProfile{
		{
			Entity: api.Entity{
				ID:      "dataspace-profile-001",
				Version: 1,
			},
			Artifacts: []string{"connector-runtime", "policy-engine"},
			Deployments: []api.DataspaceDeployment{
				{
					DeployableEntity: api.DeployableEntity{
						Entity: api.Entity{
							ID:      "deployment-001",
							Version: 1,
						},
						State:          api.DeploymentStateActive,
						StateTimestamp: time.Now(),
					},
					Cell:       cells[0], // Reference to the first cell
					Properties: api.Properties{},
				},
			},
			Properties: api.Properties{},
		},
	}
	return cells, dProfiles
}
