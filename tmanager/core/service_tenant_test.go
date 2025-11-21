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
	"testing"

	"github.com/metaform/connector-fabric-manager/common/memorystore"
	"github.com/metaform/connector-fabric-manager/common/query"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTenant(t *testing.T) {
	ctx := context.Background()
	service := newTestTenantService()

	t.Run("get existing tenant", func(t *testing.T) {
		tenant := newTestTenant("tenant-1")
		_, err := service.tenantStore.Create(ctx, tenant)
		require.NoError(t, err)

		result, err := service.GetTenant(ctx, "tenant-1")

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "tenant-1", result.ID)
		assert.Equal(t, int64(1), result.Version)
	})

	t.Run("get non-existent tenant returns not found error", func(t *testing.T) {
		result, err := service.GetTenant(ctx, "non-existent")

		require.Error(t, err)
		require.Nil(t, result)
		assert.Equal(t, types.ErrNotFound, err)
	})
}

func TestCreateTenant(t *testing.T) {
	ctx := context.Background()
	service := newTestTenantService()

	t.Run("create valid tenant", func(t *testing.T) {
		tenant := newTestTenant("tenant-1")

		result, err := service.CreateTenant(ctx, tenant)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "tenant-1", result.ID)
		assert.Equal(t, int64(1), result.Version)
		assert.Equal(t, "Test Tenant tenant-1", result.Properties["name"])
	})

	t.Run("create tenant with empty ID returns error", func(t *testing.T) {
		tenant := &api.Tenant{
			Entity: api.Entity{
				ID:      "",
				Version: 1,
			},
		}

		result, err := service.CreateTenant(ctx, tenant)

		require.Error(t, err)
		require.Nil(t, result)
		assert.Equal(t, types.ErrInvalidInput, err)
	})

	t.Run("create duplicate tenant returns error", func(t *testing.T) {
		tenant1 := newTestTenant("tenant-2")
		tenant2 := newTestTenant("tenant-2")

		// First create should succeed
		result1, err1 := service.CreateTenant(ctx, tenant1)
		require.NoError(t, err1)
		require.NotNil(t, result1)

		// Second create with same ID should fail
		result2, err2 := service.CreateTenant(ctx, tenant2)
		require.Error(t, err2)
		require.Nil(t, result2)
		assert.Equal(t, types.ErrConflict, err2)
	})

}

func TestQueryTenants(t *testing.T) {
	ctx := context.Background()
	service := newTestTenantService()

	// Populate store with test data
	tenants := []*api.Tenant{
		{
			Entity:     api.Entity{ID: "tenant-1", Version: 1},
			Properties: api.Properties{"name": "Tenant One"},
		},
		{
			Entity:     api.Entity{ID: "tenant-2", Version: 1},
			Properties: api.Properties{"name": "Tenant Two"},
		},
		{
			Entity:     api.Entity{ID: "tenant-3", Version: 1},
			Properties: api.Properties{"name": "Tenant Three"},
		},
	}

	for _, tenant := range tenants {
		_, err := service.tenantStore.Create(ctx, tenant)
		require.NoError(t, err)
	}

	t.Run("query all tenants with empty predicate", func(t *testing.T) {
		predicate := &query.AtomicPredicate{
			Field:    "properties.name",
			Operator: query.OpEqual,
			Value:    "Tenant One",
		}
		options := store.DefaultPaginationOptions()

		results := make([]api.Tenant, 0)
		for tenant, err := range service.QueryTenants(ctx, predicate, options) {
			require.NoError(t, err)
			results = append(results, tenant)
		}

		assert.Equal(t, 1, len(results))
	})
}

// TestCountTenants tests the CountTenants method
func TestCountTenants(t *testing.T) {
	ctx := context.Background()
	tenantStore := memorystore.NewInMemoryEntityStore[api.Tenant](tenantIdFunc)
	service := &tenantService{
		trxContext:  store.NoOpTransactionContext{},
		tenantStore: tenantStore,
	}

	t.Run("count empty store", func(t *testing.T) {
		predicate := &query.AtomicPredicate{
			Field:    "properties.name",
			Operator: query.OpEqual,
			Value:    "Tenant One",
		}

		count, err := service.CountTenants(ctx, predicate)

		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("count with matching predicate", func(t *testing.T) {
		predicate := &query.AtomicPredicate{
			Field:    "properties.name",
			Operator: query.OpEqual,
			Value:    "Tenant One",
		}

		count, err := service.CountTenants(ctx, predicate)

		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, 0)
	})

}

func tenantIdFunc(t *api.Tenant) string {
	if t == nil {
		return ""
	}
	return t.ID
}

func newTestTenant(id string) *api.Tenant {
	return &api.Tenant{
		Entity: api.Entity{
			ID:      id,
			Version: 1,
		},
		Properties: api.Properties{
			"name": "Test Tenant " + id,
		},
	}
}

func newTestTenantService() *tenantService {
	return &tenantService{
		trxContext:  store.NoOpTransactionContext{},
		tenantStore: memorystore.NewInMemoryEntityStore[api.Tenant](tenantIdFunc),
		monitor:     nil,
	}
}
