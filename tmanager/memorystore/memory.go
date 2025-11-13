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
	"iter"
	"sync"

	"github.com/metaform/connector-fabric-manager/common/query"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

func NewInMemoryEntityStore[T any](idFunc func(*T) string) *InMemoryEntityStore[T] {
	store := &InMemoryEntityStore[T]{
		cache:   make(map[string]T),
		idFunc:  idFunc,
		matcher: &query.DefaultFieldMatcher{},
	}
	return store
}

type InMemoryEntityStore[T any] struct {
	cache   map[string]T
	mu      sync.RWMutex
	idFunc  func(*T) string
	matcher query.FieldMatcher
}

func (s *InMemoryEntityStore[T]) FindById(_ context.Context, id string) (*T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entity, exists := s.cache[id]
	if !exists {
		return nil, types.ErrNotFound
	}

	return &entity, nil
}

func (s *InMemoryEntityStore[T]) Exists(_ context.Context, id string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.cache[id]
	return exists, nil
}

func (s *InMemoryEntityStore[T]) Create(_ context.Context, entity *T) (*T, error) {
	if s.idFunc(entity) == "" {
		return nil, types.ErrInvalidInput
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.cache[s.idFunc(entity)]; exists {
		return nil, types.ErrConflict
	}

	s.cache[s.idFunc(entity)] = *entity
	return entity, nil
}

func (s *InMemoryEntityStore[T]) Update(_ context.Context, entity *T) error {
	if entity == nil {
		return types.ErrInvalidInput
	}
	if s.idFunc(entity) == "" {
		return types.ErrInvalidInput
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.cache[s.idFunc(entity)]; !exists {
		return types.ErrNotFound
	}

	s.cache[s.idFunc(entity)] = *entity
	return nil
}

func (s *InMemoryEntityStore[T]) Delete(_ context.Context, id string) error {
	if id == "" {
		return types.ErrInvalidInput
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.cache[id]; !exists {
		return types.ErrNotFound
	}

	delete(s.cache, id)
	return nil
}

func (s *InMemoryEntityStore[T]) GetAll(ctx context.Context) iter.Seq2[T, error] {
	return s.GetAllPaginated(ctx, api.DefaultPaginationOptions())
}

func (s *InMemoryEntityStore[T]) GetAllPaginated(ctx context.Context, opts api.PaginationOptions) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		// Convert map to slice for consistent ordering and pagination
		entities := make([]T, 0, len(s.cache))
		for _, entity := range s.cache {
			entities = append(entities, entity)
		}

		// Apply offset
		start := opts.Offset
		if start < 0 {
			start = 0
		}
		if start >= len(entities) {
			return // No items to return
		}

		// Apply limit
		end := len(entities)
		if opts.Limit > 0 {
			requestedEnd := start + opts.Limit
			if requestedEnd < end {
				end = requestedEnd
			}
		}

		// Yield entities within the paginated range
		for i := start; i < end; i++ {
			// Check if context is cancelled
			select {
			case <-ctx.Done():
				var zero T
				yield(zero, ctx.Err())
				return
			default:
			}

			// Yield the entity with nil error
			if !yield(entities[i], nil) {
				return // Consumer stopped iteration
			}
		}
	}
}

func (s *InMemoryEntityStore[T]) FindByPredicate(ctx context.Context, predicate query.Predicate) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		for _, entity := range s.cache {
			if predicate.Matches(entity, s.matcher) {
				if !yield(entity, nil) {
					return
				}
			}
		}
	}
}

// FindByPredicatePaginated returns entities matching the predicate with pagination applied
func (s *InMemoryEntityStore[T]) FindByPredicatePaginated(
	ctx context.Context,
	predicate query.Predicate,
	opts api.PaginationOptions) iter.Seq2[T, error] {

	return func(yield func(T, error) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		// Filter entities matching the predicate into a slice
		var filtered []T
		for _, entity := range s.cache {
			if predicate.Matches(entity, s.matcher) {
				filtered = append(filtered, entity)
			}
		}

		// Apply offset
		start := opts.Offset
		if start < 0 {
			start = 0
		}
		if start >= len(filtered) {
			return // No items to return
		}

		// Apply limit
		end := len(filtered)
		if opts.Limit > 0 {
			requestedEnd := start + opts.Limit
			if requestedEnd < end {
				end = requestedEnd
			}
		}

		// Yield entities within the paginated range
		for i := start; i < end; i++ {
			// Check if context is canceled
			select {
			case <-ctx.Done():
				var zero T
				yield(zero, ctx.Err())
				return
			default:
			}

			// Yield the entity with nil error
			if !yield(filtered[i], nil) {
				return // Consumer stopped iteration
			}
		}
	}
}

// FindFirstByPredicate returns the first entity matching the predicate or types.ErrNotFound if none found
func (s *InMemoryEntityStore[T]) FindFirstByPredicate(ctx context.Context, predicate query.Predicate) (*T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, entity := range s.cache {
		if predicate.Matches(entity, s.matcher) {
			return &entity, nil
		}
	}
	return nil, types.ErrNotFound
}

func (s *InMemoryEntityStore[T]) CountByPredicate(ctx context.Context, predicate query.Predicate) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, entity := range s.cache {
		if predicate.Matches(entity, s.matcher) {
			count++
		}
	}
	return count, nil
}

