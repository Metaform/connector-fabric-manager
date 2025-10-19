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

package tmhandler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/metaform/connector-fabric-manager/assembly/routing"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/metaform/connector-fabric-manager/tmanager/tmstore"
)

type response struct {
	Message string `json:"message"`
}

type HandlerServiceAssembly struct {
	system.DefaultServiceAssembly
}

func (h *HandlerServiceAssembly) Name() string {
	return "Tenant Manager Handlers"
}

func (h *HandlerServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{}
}

func (h *HandlerServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{api.ParticipantDeployerKey, tmstore.CellStoreKey, tmstore.DataspaceProfileStoreKey, routing.RouterKey}
}

func (h *HandlerServiceAssembly) Init(context *system.InitContext) error {
	router := context.Registry.Resolve(routing.RouterKey).(chi.Router)
	router.Use(middleware.Recoverer)

	deployer := context.Registry.Resolve(api.ParticipantDeployerKey).(api.ParticipantDeployer)
	cellStore := context.Registry.Resolve(tmstore.CellStoreKey).(tmstore.EntityStore[api.Cell])
	dProfileStore := context.Registry.Resolve(tmstore.DataspaceProfileStoreKey).(tmstore.EntityStore[api.DataspaceProfile])

	handler := NewHandler(deployer, cellStore, dProfileStore, context.LogMonitor)

	router.Post("/participant/{id}", handler.deployParticipant)

	return nil
}
