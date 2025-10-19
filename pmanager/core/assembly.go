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
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

type PMCoreServiceAssembly struct {
	system.DefaultServiceAssembly
}

func (m PMCoreServiceAssembly) Name() string {
	return "Provision Manager Core"
}

func (m PMCoreServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{api.ProvisionManagerKey}
}

func (m PMCoreServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{api.DefinitionStoreKey, api.DeploymentOrchestratorKey}
}

func (m PMCoreServiceAssembly) Init(context *system.InitContext) error {
	context.Registry.Register(api.ProvisionManagerKey, provisionManager{
		orchestrator: context.Registry.Resolve(api.DeploymentOrchestratorKey).(api.DeploymentOrchestrator),
		store:        context.Registry.Resolve(api.DefinitionStoreKey).(api.DefinitionStore),
		monitor:      context.LogMonitor,
	})
	return nil
}
