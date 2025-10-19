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

package tmcore

import (
	"github.com/metaform/connector-fabric-manager/common/dmodel"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/metaform/connector-fabric-manager/tmanager/tmstore"
)

type TMCoreServiceAssembly struct {
	system.DefaultServiceAssembly
	vpaGenerator participantGenerator
}

func (a *TMCoreServiceAssembly) Name() string {
	return "Tenant Manager Core"
}

func (a *TMCoreServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{api.DeploymentClientKey, store.TransactionContextKey, tmstore.CellStoreKey, tmstore.DataspaceProfileStoreKey}
}

func (a *TMCoreServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{api.ParticipantDeployerKey}
}

func (a *TMCoreServiceAssembly) Init(context *system.InitContext) error {
	a.vpaGenerator = participantGenerator{
		CellSelector: defaultVPASelector, // Register the default selector, which may be overridden
	}

	trxContext := context.Registry.Resolve(store.TransactionContextKey).(store.TransactionContext)
	deploymentClient := context.Registry.Resolve(api.DeploymentClientKey).(api.DeploymentClient)
	cellStore := context.Registry.Resolve(tmstore.CellStoreKey).(tmstore.EntityStore[api.Cell])
	dProfileStore := context.Registry.Resolve(tmstore.DataspaceProfileStoreKey).(tmstore.EntityStore[api.DataspaceProfile])

	participantDeployer := participantDeployer{
		participantGenerator: a.vpaGenerator,
		deploymentClient:     deploymentClient,
		trxContext:           trxContext,
		cellStore:            cellStore,
		dProfileStore:        dProfileStore,
	}
	context.Registry.Register(api.ParticipantDeployerKey, participantDeployer)

	registry := context.Registry.Resolve(api.DeploymentHandlerRegistryKey).(api.DeploymentHandlerRegistry)
	handler := vpaDeploymentCallbackHandler{}
	registry.RegisterDeploymentHandler(dmodel.VpaDeploymentType, handler.handle)

	return nil
}

func (a *TMCoreServiceAssembly) Prepare(context *system.InitContext) error {
	selector, found := context.Registry.ResolveOptional(api.CellSelectorKey)
	if found {
		// Override the default selector with a custom implementation
		a.vpaGenerator = selector.(participantGenerator)
	}
	return nil
}
