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

package pmcore

import (
	"context"
	"errors"
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

type provisionManager struct {
	orchestrator api.DeploymentOrchestrator
	store        api.DefinitionStore
	logMonitor   monitor.LogMonitor
}

func (p provisionManager) Start(ctx context.Context, manifest *api.DeploymentManifest) (*api.Orchestration, error) {

	deploymentID := manifest.ID

	definition, err := p.store.FindDeploymentDefinition(manifest.DeploymentType)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// Not found is a client error
			return nil, model.NewClientError("deployment type '%s' not found", manifest.DeploymentType)
		}
		return nil, model.NewFatalWrappedError(err, "no deployment definition for deployment %s", deploymentID)
	}

	activeVersion, err := definition.GetActiveVersion()
	if err != nil {
		return nil, model.NewFatalError("error deploying %s: unable to get active version", deploymentID)
	}

	// Validate required fields
	if manifest.ID == "" {
		return nil, model.NewClientError("Missing required field: id")
	}

	if manifest.DeploymentType == "" {
		return nil, model.NewClientError("Missing required field: deploymentType")
	}

	// perform de-duplication
	orchestration, err := p.orchestrator.GetOrchestration(ctx, deploymentID)
	if err != nil {
		return nil, model.NewFatalWrappedError(err, "error checking for orchestration %s", deploymentID)
	}

	if orchestration != nil {
		// Already exists, return its representation
		return orchestration, nil
	}

	// Does not exist, create the orchestration
	orchestration, err = api.InstantiateOrchestration(manifest.ID, activeVersion.Activities, manifest.Payload)
	if err != nil {
		return nil, model.NewFatalWrappedError(err, "error instantiating orchestration for deployment %s", deploymentID)
	}
	err = p.orchestrator.ExecuteOrchestration(ctx, orchestration)
	if err != nil {
		return nil, model.NewFatalWrappedError(err, "error executing orchestration %s for deployment %s", orchestration.ID, deploymentID)
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
