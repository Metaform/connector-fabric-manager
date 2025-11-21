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
	"iter"

	"github.com/metaform/connector-fabric-manager/common/query"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

type tenantService struct {
	trxContext  store.TransactionContext
	tenantStore store.EntityStore[api.Tenant]
	monitor     system.LogMonitor
}

func (t tenantService) GetTenant(ctx context.Context, tenantID string) (*api.Tenant, error) {
	return store.Trx[api.Tenant](t.trxContext).AndReturn(ctx, func(ctx context.Context) (*api.Tenant, error) {
		return t.tenantStore.FindById(ctx, tenantID)
	})
}

func (t tenantService) CreateTenant(ctx context.Context, tenant *api.Tenant) (*api.Tenant, error) {
	return store.Trx[api.Tenant](t.trxContext).AndReturn(ctx, func(ctx context.Context) (*api.Tenant, error) {
		return t.tenantStore.Create(ctx, tenant)
	})
}

func (t tenantService) DeleteTenant(ctx context.Context, tenantID string) error {
	//TODO implement me
	panic("implement me")
}

func (t tenantService) QueryTenants(ctx context.Context, predicate query.Predicate, options store.PaginationOptions) iter.Seq2[api.Tenant, error] {
	return func(yield func(api.Tenant, error) bool) {
		err := t.trxContext.Execute(ctx, func(ctx context.Context) error {
			for tenant, err := range t.tenantStore.FindByPredicatePaginated(ctx, predicate, options) {
				if !yield(tenant, err) {
					return context.Canceled
				}
			}
			return nil
		})
		if err != nil {
			yield(api.Tenant{}, err)
		}
	}
}

func (t tenantService) CountTenants(ctx context.Context, predicate query.Predicate) (int, error) {
	var count int
	err := t.trxContext.Execute(ctx, func(ctx context.Context) error {
		c, err := t.tenantStore.CountByPredicate(ctx, predicate)
		count = c
		return err
	})
	return count, err
}
