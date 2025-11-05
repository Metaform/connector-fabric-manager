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
	cstore "github.com/metaform/connector-fabric-manager/common/store"
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
	return []system.ServiceType{api.ProvisionManagerKey, api.DefinitionManagerKey}
}

func (m PMCoreServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{api.DefinitionStoreKey, api.OrchestratorKey, cstore.TransactionContextKey}
}

func (m PMCoreServiceAssembly) Init(context *system.InitContext) error {
	store := context.Registry.Resolve(api.DefinitionStoreKey).(api.DefinitionStore)
	context.Registry.Register(api.ProvisionManagerKey, provisionManager{
		orchestrator: context.Registry.Resolve(api.OrchestratorKey).(api.Orchestrator),
		store:        store,
		monitor:      context.LogMonitor,
	})

	context.Registry.Register(api.DefinitionManagerKey, definitionManager{
		trxContext: context.Registry.Resolve(cstore.TransactionContextKey).(cstore.TransactionContext),
		store:      store,
	})
	return nil
}
