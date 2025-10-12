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
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/dmodel"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

type TMCoreServiceAssembly struct {
	system.DefaultServiceAssembly
	vpaGenerator participantGenerator
}

func (a *TMCoreServiceAssembly) Name() string {
	return "Tenant Manager Core"
}

func (a *TMCoreServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{}
}

func (a *TMCoreServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{api.DeploymentCallbackDispatcherKey, api.DeploymentHandlerRegistryKey, api.ParticipantDeployerKey}
}

func (a *TMCoreServiceAssembly) Init(context *system.InitContext) error {
	a.vpaGenerator = participantGenerator{
		CellSelector: defaultVPASelector, // Register the default selector, which may be overridden
	}

	participantDeployer := participantDeployer{
		participantGenerator: a.vpaGenerator,
	}
	context.Registry.Register(api.ParticipantDeployerKey, participantDeployer)

	callbackService := deploymentCallbackService{handlers: make(map[string]api.DeploymentCallbackHandler)}
	handler := vpaDeploymentCallbackHandler{}
	callbackService.RegisterDeploymentHandler(dmodel.VpaDeploymentType, handler.handle)

	context.Registry.Register(api.DeploymentCallbackDispatcherKey, callbackService)
	context.Registry.Register(api.DeploymentHandlerRegistryKey, callbackService)
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
