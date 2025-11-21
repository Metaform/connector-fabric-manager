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

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

type provisionManager struct {
	orchestrator api.Orchestrator
	store        api.DefinitionStore
	monitor      system.LogMonitor
}

func (p provisionManager) Start(ctx context.Context, manifest *model.OrchestrationManifest) (*api.Orchestration, error) {

	manifestID := manifest.ID

	definition, err := p.store.FindOrchestrationDefinition(ctx, manifest.OrchestrationType)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// Not found is a client error
			return nil, types.NewClientError("orchestration type '%s' not found", manifest.OrchestrationType)
		}
		return nil, types.NewFatalWrappedError(err, "unable to find orchestration definition for manifest %s", manifestID)
	}

	// Validate required fields
	if manifest.ID == "" {
		return nil, types.NewClientError("Missing required field: id")
	}

	if manifest.OrchestrationType == "" {
		return nil, types.NewClientError("Missing required field: orchestrationType")
	}

	// perform de-duplication
	orchestration, err := p.orchestrator.GetOrchestration(ctx, manifestID)
	if err != nil {
		return nil, types.NewFatalWrappedError(err, "error performing de-duplication for %s", manifestID)
	}

	if orchestration != nil {
		// Already exists, return its representation
		return orchestration, nil
	}

	// Does not exist, create the orchestration
	orchestration, err = api.InstantiateOrchestration(manifest.ID, manifest.CorrelationID, manifest.OrchestrationType, definition.Activities, manifest.Payload)
	if err != nil {
		return nil, types.NewFatalWrappedError(err, "error instantiating orchestration for %s", manifestID)
	}
	err = p.orchestrator.Execute(ctx, orchestration)
	if err != nil {
		return nil, types.NewFatalWrappedError(err, "error executing orchestration %s for %s", orchestration.ID, manifestID)
	}
	return orchestration, nil
}

func (p provisionManager) Cancel(ctx context.Context, orchestrationID string) error {
	//TODO implement me
	panic("implement me")
}

func (p provisionManager) GetOrchestration(ctx context.Context, orchestrationID string) (*api.Orchestration, error) {
	return p.orchestrator.GetOrchestration(ctx, orchestrationID)
}

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

func (d definitionManager) CreateActivityDefinition(ctx context.Context, definition *api.ActivityDefinition) (*api.ActivityDefinition, error) {
	return store.Trx[api.ActivityDefinition](d.trxContext).AndReturn(ctx, func(ctx context.Context) (*api.ActivityDefinition, error) {
		definition, err := d.store.StoreActivityDefinition(ctx, definition)
		if err != nil {
			return nil, err
		}
		return definition, nil
	})
}
