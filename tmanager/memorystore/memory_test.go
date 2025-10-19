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
	"testing"

	"github.com/metaform/connector-fabric-manager/common/collection"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testEntity is a simple test entity for testing the InMemoryEntityStore
type testEntity struct {
	ID    string
	Value string
}

func testIdFunc(e *testEntity) string {
	return e.ID
}

func TestNewInMemoryEntityStore(t *testing.T) {
	store := NewInMemoryEntityStore[testEntity](testIdFunc)

	require.NotNil(t, store)
	require.NotNil(t, store.cache)
	require.NotNil(t, store.idFunc)
	assert.Equal(t, 0, len(store.cache))
}

func TestInMemoryEntityStore_Create(t *testing.T) {
	store := NewInMemoryEntityStore[testEntity](testIdFunc)
	ctx := context.Background()

	t.Run("successful create", func(t *testing.T) {
		entity := &testEntity{ID: "test-1", Value: "value1"}

		result, err := store.Create(ctx, entity)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "test-1", result.ID)
		assert.Equal(t, "value1", result.Value)
		assert.Equal(t, 1, len(store.cache))
	})

	t.Run("create with empty ID should fail", func(t *testing.T) {
		entity := &testEntity{ID: "", Value: "value1"}

		result, err := store.Create(ctx, entity)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Equal(t, types.ErrInvalidInput, err)
	})

	t.Run("create duplicate should fail", func(t *testing.T) {
		entity := &testEntity{ID: "test-2", Value: "value1"}
		duplicate := &testEntity{ID: "test-2", Value: "value2"}

		// First create should succeed
		result1, err1 := store.Create(ctx, entity)
		require.NoError(t, err1)
		require.NotNil(t, result1)

		// Second create with same ID should fail
		result2, err2 := store.Create(ctx, duplicate)
		require.Error(t, err2)
		require.Nil(t, result2)
		assert.Equal(t, types.ErrConflict, err2)
	})
}

func TestInMemoryEntityStore_FindById(t *testing.T) {
	store := NewInMemoryEntityStore[testEntity](testIdFunc)
	ctx := context.Background()

	t.Run("find existing entity", func(t *testing.T) {
		entity := &testEntity{ID: "test-1", Value: "value1"}
		_, err := store.Create(ctx, entity)
		require.NoError(t, err)

		result, err := store.FindById(ctx, "test-1")

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "test-1", result.ID)
		assert.Equal(t, "value1", result.Value)
	})

	t.Run("find non-existing entity", func(t *testing.T) {
		result, err := store.FindById(ctx, "non-existing")

		require.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestInMemoryEntityStore_Exists(t *testing.T) {
	store := NewInMemoryEntityStore[testEntity](testIdFunc)
	ctx := context.Background()

	t.Run("entity exists", func(t *testing.T) {
		entity := &testEntity{ID: "test-1", Value: "value1"}
		_, err := store.Create(ctx, entity)
		require.NoError(t, err)

		exists, err := store.Exists(ctx, "test-1")

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("entity does not exist", func(t *testing.T) {
		exists, err := store.Exists(ctx, "non-existing")

		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestInMemoryEntityStore_Update(t *testing.T) {
	store := NewInMemoryEntityStore[testEntity](testIdFunc)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		entity := &testEntity{ID: "test-1", Value: "value1"}
		_, err := store.Create(ctx, entity)
		require.NoError(t, err)

		updatedEntity := &testEntity{ID: "test-1", Value: "updated-value"}
		err = store.Update(ctx, updatedEntity)

		require.NoError(t, err)

		// Verify the update
		result, err := store.FindById(ctx, "test-1")
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "updated-value", result.Value)
	})

	t.Run("update nil entity should fail", func(t *testing.T) {
		err := store.Update(ctx, nil)

		require.Error(t, err)
		assert.Equal(t, types.ErrInvalidInput, err)
	})

	t.Run("update entity with empty ID should fail", func(t *testing.T) {
		entity := &testEntity{ID: "", Value: "value1"}

		err := store.Update(ctx, entity)

		require.Error(t, err)
		assert.Equal(t, types.ErrInvalidInput, err)
	})

	t.Run("update non-existing entity should fail", func(t *testing.T) {
		entity := &testEntity{ID: "non-existing", Value: "value1"}

		err := store.Update(ctx, entity)

		require.Error(t, err)
		assert.Equal(t, types.ErrNotFound, err)
	})
}

func TestInMemoryEntityStore_Delete(t *testing.T) {
	store := NewInMemoryEntityStore[testEntity](testIdFunc)
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		entity := &testEntity{ID: "test-1", Value: "value1"}
		_, err := store.Create(ctx, entity)
		require.NoError(t, err)

		err = store.Delete(ctx, "test-1")

		require.NoError(t, err)

		// Verify deletion
		result, err := store.FindById(ctx, "test-1")
		require.Error(t, err)
		assert.Nil(t, result)

		exists, err := store.Exists(ctx, "test-1")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("delete with empty ID should fail", func(t *testing.T) {
		err := store.Delete(ctx, "")

		require.Error(t, err)
		assert.Equal(t, types.ErrInvalidInput, err)
	})

	t.Run("delete non-existing entity should fail", func(t *testing.T) {
		err := store.Delete(ctx, "non-existing")

		require.Error(t, err)
		assert.Equal(t, types.ErrNotFound, err)
	})
}

func TestInMemoryEntityStore_GetAll(t *testing.T) {
	store := NewInMemoryEntityStore[testEntity](testIdFunc)
	ctx := context.Background()

	t.Run("get all from empty store", func(t *testing.T) {
		entities, err := collection.CollectAll(store.GetAll(ctx))

		require.NoError(t, err)
		assert.Equal(t, 0, len(entities))
	})

	t.Run("get all entities", func(t *testing.T) {
		// Create test entities
		entities := []*testEntity{
			{ID: "test-1", Value: "value1"},
			{ID: "test-2", Value: "value2"},
			{ID: "test-3", Value: "value3"},
		}

		for _, entity := range entities {
			_, err := store.Create(ctx, entity)
			require.NoError(t, err)
		}

		result, err := collection.CollectAll(store.GetAll(ctx))

		require.NoError(t, err)
		assert.Equal(t, 3, len(result))

		// Verify all entities are present
		ids := make(map[string]bool)
		for _, entity := range result {
			ids[entity.ID] = true
		}
		assert.True(t, ids["test-1"])
		assert.True(t, ids["test-2"])
		assert.True(t, ids["test-3"])
	})
}

func TestInMemoryEntityStore_GetAllPaginated(t *testing.T) {
	store := NewInMemoryEntityStore[testEntity](testIdFunc)
	ctx := context.Background()

	// Create test entities
	entities := []*testEntity{
		{ID: "test-1", Value: "value1"},
		{ID: "test-2", Value: "value2"},
		{ID: "test-3", Value: "value3"},
		{ID: "test-4", Value: "value4"},
		{ID: "test-5", Value: "value5"},
	}

	for _, entity := range entities {
		_, err := store.Create(ctx, entity)
		require.NoError(t, err)
	}

	t.Run("pagination with limit", func(t *testing.T) {
		opts := api.PaginationOptions{Offset: 0, Limit: 3}
		result, err := collection.CollectAll(store.GetAllPaginated(ctx, opts))

		require.NoError(t, err)
		assert.Equal(t, 3, len(result))
	})

	t.Run("pagination with offset", func(t *testing.T) {
		opts := api.PaginationOptions{Offset: 2, Limit: 2}
		result, err := collection.CollectAll(store.GetAllPaginated(ctx, opts))

		require.NoError(t, err)
		assert.Equal(t, 2, len(result))
	})

	t.Run("pagination with offset beyond range", func(t *testing.T) {
		opts := api.PaginationOptions{Offset: 10, Limit: 2}
		result, err := collection.CollectAll(store.GetAllPaginated(ctx, opts))

		require.NoError(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("pagination with negative offset", func(t *testing.T) {
		opts := api.PaginationOptions{Offset: -1, Limit: 2}
		result, err := collection.CollectAll(store.GetAllPaginated(ctx, opts))

		require.NoError(t, err)
		assert.Equal(t, 2, len(result))
	})

	t.Run("pagination with no limit", func(t *testing.T) {
		opts := api.PaginationOptions{Offset: 0, Limit: 0}
		result, err := collection.CollectAll(store.GetAllPaginated(ctx, opts))

		require.NoError(t, err)
		assert.Equal(t, 5, len(result))
	})
}

func TestInMemoryEntityStore_ConcurrentAccess(t *testing.T) {
	store := NewInMemoryEntityStore[testEntity](testIdFunc)
	ctx := context.Background()

	t.Run("concurrent create and read", func(t *testing.T) {
		done := make(chan bool, 2)

		// Goroutine 1: Create entities
		go func() {
			for i := 0; i < 10; i++ {
				entity := &testEntity{ID: fmt.Sprintf("test-%d", i), Value: fmt.Sprintf("value%d", i)}
				store.Create(ctx, entity)
			}
			done <- true
		}()

		// Goroutine 2: Read entities
		go func() {
			for i := 0; i < 10; i++ {
				store.FindById(ctx, fmt.Sprintf("test-%d", i))
			}
			done <- true
		}()

		// Wait for both goroutines
		<-done
		<-done

		// Verify some entities were created
		entities, err := collection.CollectAll(store.GetAll(ctx))
		require.NoError(t, err)
		assert.True(t, len(entities) > 0)
	})
}

func TestInMemoryEntityStore_ContextCancellation(t *testing.T) {
	store := NewInMemoryEntityStore[testEntity](testIdFunc)

	// Create some test data
	for i := 0; i < 5; i++ {
		entity := &testEntity{ID: fmt.Sprintf("test-%d", i), Value: fmt.Sprintf("value%d", i)}
		_, err := store.Create(context.Background(), entity)
		require.NoError(t, err)
	}

	t.Run("cancelled context in GetAllPaginated", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		var results []testEntity
		var lastErr error

		for entity, err := range store.GetAllPaginated(ctx, api.DefaultPaginationOptions()) {
			if err != nil {
				lastErr = err
				break
			}
			results = append(results, entity)
		}

		// Should get a context cancellation error
		require.Error(t, lastErr)
		assert.Equal(t, context.Canceled, lastErr)
	})
}

func TestInMemoryEntityStore_CopyIsolation(t *testing.T) {
	store := NewInMemoryEntityStore[testEntity](testIdFunc)
	ctx := context.Background()

	entity := &testEntity{ID: "test-1", Value: "original-value"}
	_, err := store.Create(ctx, entity)
	require.NoError(t, err)

	// Retrieve the entity (first retrieval)
	retrieved1, err := store.FindById(ctx, "test-1")
	require.NoError(t, err)
	require.NotNil(t, retrieved1)
	assert.Equal(t, "original-value", retrieved1.Value)

	// Modify the retrieved copy
	retrieved1.Value = "modified-value"

	// Retrieve the entity again (second retrieval)
	retrieved2, err := store.FindById(ctx, "test-1")
	require.NoError(t, err)
	require.NotNil(t, retrieved2)

	// The second retrieval should still have the original value
	// The modification to retrieved1 should not be visible in retrieved2
	assert.Equal(t, "original-value", retrieved2.Value)
	assert.NotEqual(t, retrieved1.Value, retrieved2.Value)
}
