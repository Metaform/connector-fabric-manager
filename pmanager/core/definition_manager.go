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
	"strings"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

type definitionManager struct {
	trxContext store.TransactionContext
	store      api.DefinitionStore
}

func (d definitionManager) CreateOrchestrationDefinition(ctx context.Context, definition *api.OrchestrationDefinition) (*api.OrchestrationDefinition, error) {
	return store.Trx[api.OrchestrationDefinition](d.trxContext).AndReturn(ctx, func(ctx context.Context) (*api.OrchestrationDefinition, error) {
		var missingErrors []error

		// Verify that all referenced activities exist
		for _, activity := range definition.Activities {
			exists, err := d.store.ExistsActivityDefinition(ctx, activity.Type)
			if err != nil {
				return nil, err
			}
			if !exists {
				missingErrors = append(missingErrors, types.NewClientError("activity type '%s' not found", activity.Type))
			}
		}

		if len(missingErrors) > 0 {
			return nil, errors.Join(missingErrors...)
		}

		persisted, err := d.store.StoreOrchestrationDefinition(ctx, definition)
		if err != nil {
			return nil, err
		}
		return persisted, nil
	})
}

func (d definitionManager) DeleteOrchestrationDefinition(
	// TODO this method should check outstanding orchestrations when the orchestration index is implemented
	ctx context.Context,
	orchestrationType model.OrchestrationType) error {

	return d.trxContext.Execute(ctx, func(ctx context.Context) error {
		exists, err := d.store.ExistsOrchestrationDefinition(ctx, orchestrationType)
		if err != nil {
			return types.NewRecoverableWrappedError(err, "failed to check orchestration definition for type %s", orchestrationType)
		}
		if !exists {
			return types.ErrNotFound
		}

		deleted, err := d.store.DeleteOrchestrationDefinition(ctx, orchestrationType)
		if err != nil {
			return types.NewRecoverableWrappedError(err, "failed to delete orchestration definition for type %s", orchestrationType)
		}
		if !deleted {
			return types.NewClientError("unable to delete orchestration definition type %s", orchestrationType)
		}
		return nil
	})
}

func (d definitionManager) CreateActivityDefinition(ctx context.Context, definition *api.ActivityDefinition) (*api.ActivityDefinition, error) {
	return store.Trx[api.ActivityDefinition](d.trxContext).AndReturn(ctx, func(ctx context.Context) (*api.ActivityDefinition, error) {
		definition, err := d.store.StoreActivityDefinition(ctx, definition)
		if err != nil {
			return nil, err
		}
		return definition, nil
	})
}

func (d definitionManager) DeleteActivityDefinition(ctx context.Context, atype api.ActivityType) error {
	return d.trxContext.Execute(ctx, func(ctx context.Context) error {
		exists, err := d.store.ExistsActivityDefinition(ctx, atype)
		if err != nil {
			return types.NewRecoverableWrappedError(err, "failed to check activity definition for type %s", atype)
		}
		if !exists {
			return types.ErrNotFound
		}
		referenced, err := d.store.ActivityDefinitionReferences(ctx, atype)

		if err != nil {
			return types.NewRecoverableWrappedError(err, "failed to check activity definition references for type %s", atype)
		}
		if len(referenced) > 0 {
			return types.NewClientError("activity type '%s' is referenced by an orchestration definition: %s", atype, strings.Join(referenced, ", "))
		}

		deleted, err := d.store.DeleteActivityDefinition(ctx, atype)
		if err != nil {
			return types.NewRecoverableWrappedError(err, "failed to check activity definition references for type %s", atype)
		}
		if !deleted {
			return types.NewClientError("unable to delete activity definition type %s", atype)
		}
		return nil
	})

}
