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
	"fmt"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

type provisionManager struct {
	orchestrator api.DeploymentOrchestrator
	logMonitor   monitor.LogMonitor
}

func (p provisionManager) Start(ctx context.Context, manifest *api.DeploymentManifest) (*api.Orchestration, error) {

	// TODO return definition
	definition := api.OrchestrationDefinition{}

	// TODO implement validation

	deploymentID := manifest.ID

	// perform de-duplication
	orchestration, err := p.orchestrator.GetOrchestration(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("error checking for orchestration %s: %w", deploymentID, err)
	}

	if orchestration != nil {
		// Already exists, return its representation
		return orchestration, nil
	}

	// Does not exist, create the orchestration
	orchestration, err = api.InstantiateOrchestration(manifest.ID, definition, manifest.Payload)
	if err != nil {
		return nil, fmt.Errorf("error instantiating orchestration for deployment %s: %w", deploymentID, err)
	}
	err = p.orchestrator.ExecuteOrchestration(ctx, orchestration)
	if err != nil {
		return nil, fmt.Errorf("error executing orchestration %s for deployment %s: %w", orchestration.ID, deploymentID, err)
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
