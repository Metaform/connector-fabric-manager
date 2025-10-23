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
	orchestrator api.DeploymentOrchestrator
	store        api.DefinitionStore
	monitor      system.LogMonitor
}

func (p provisionManager) Start(ctx context.Context, manifest *model.DeploymentManifest) (*api.Orchestration, error) {

	deploymentID := manifest.ID

	definition, err := p.store.FindDeploymentDefinition(manifest.DeploymentType)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// Not found is a client error
			return nil, types.NewClientError("deployment type '%s' not found", manifest.DeploymentType)
		}
		return nil, types.NewFatalWrappedError(err, "unable to find deployment definition for deployment %s", deploymentID)
	}

	// Validate required fields
	if manifest.ID == "" {
		return nil, types.NewClientError("Missing required field: id")
	}

	if manifest.DeploymentType == "" {
		return nil, types.NewClientError("Missing required field: deploymentType")
	}

	// perform de-duplication
	orchestration, err := p.orchestrator.GetOrchestration(ctx, deploymentID)
	if err != nil {
		return nil, types.NewFatalWrappedError(err, "error checking for orchestration %s", deploymentID)
	}

	if orchestration != nil {
		// Already exists, return its representation
		return orchestration, nil
	}

	// Does not exist, create the orchestration
	orchestration, err = api.InstantiateOrchestration(manifest.ID, manifest.CorrelationID, manifest.DeploymentType, definition.Activities, manifest.Payload)
	if err != nil {
		return nil, types.NewFatalWrappedError(err, "error instantiating orchestration for deployment %s", deploymentID)
	}
	err = p.orchestrator.ExecuteOrchestration(ctx, orchestration)
	if err != nil {
		return nil, types.NewFatalWrappedError(err, "error executing orchestration %s for deployment %s", orchestration.ID, deploymentID)
	}
	return orchestration, nil
}

func (p provisionManager) Cancel(ctx context.Context, deploymentID string) error {
	//TODO implement me
	panic("implement me")
}

func (p provisionManager) GetOrchestration(ctx context.Context, deploymentID string) (*api.Orchestration, error) {
	return p.orchestrator.GetOrchestration(ctx, deploymentID)
}

type definitionManager struct {
	trxContext store.TransactionContext
	store      api.DefinitionStore
}

func (d definitionManager) CreateDeploymentDefinition(ctx context.Context, definition *api.DeploymentDefinition) (*api.DeploymentDefinition, error) {
	return store.Trx[api.DeploymentDefinition](d.trxContext).AndReturn(ctx, func(ctx context.Context) (*api.DeploymentDefinition, error) {
		var missingErrors []error

		// Verify that all referenced activities exist
		for _, activity := range definition.Activities {
			exists, err := d.store.ExistsActivityDefinition(activity.Type)
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

		definition, err := d.store.StoreDeploymentDefinition(definition)
		if err != nil {
			return nil, err
		}
		return definition, nil
	})
}

func (d definitionManager) CreateActivityDefinition(ctx context.Context, definition *api.ActivityDefinition) (*api.ActivityDefinition, error) {
	return store.Trx[api.ActivityDefinition](d.trxContext).AndReturn(ctx, func(ctx context.Context) (*api.ActivityDefinition, error) {
		definition, err := d.store.StoreActivityDefinition(definition)
		if err != nil {
			return nil, err
		}
		return definition, nil
	})
}
