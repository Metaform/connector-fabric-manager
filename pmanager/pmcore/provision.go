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

func (p provisionManager) Start(ctx context.Context, manifest *api.DeploymentManifest) error {
	// TODO implement validation
	// TODO perform de-duplication
	definition := api.OrchestrationDefinition{}

	deploymentID := ""

	orchestration, err := api.InstantiateOrchestration(deploymentID, definition, manifest.Payload)
	if err != nil {
		return fmt.Errorf("error instantiating orchestration for deployment %s: %w", deploymentID, err)
	}
	err = p.orchestrator.ExecuteOrchestration(ctx, orchestration)
	if err != nil {
		return fmt.Errorf("error executing orchestration %s for deployment %s: %w", orchestration.ID, deploymentID, err)
	}
	return nil
}

func (p provisionManager) Cancel(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
