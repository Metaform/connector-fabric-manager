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

package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/metaform/connector-fabric-manager/assembly/routing"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

type response struct {
	Message string `json:"message"`
}

type HandlerServiceAssembly struct {
	system.DefaultServiceAssembly
}

func (h *HandlerServiceAssembly) Name() string {
	return "Provision Manager Handlers"
}

func (h *HandlerServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{routing.RouterKey, api.ProvisionManagerKey, api.DefinitionStoreKey}
}

func (h *HandlerServiceAssembly) Init(context *system.InitContext) error {
	router := context.Registry.Resolve(routing.RouterKey).(chi.Router)
	router.Use(middleware.Recoverer)

	provisionManager := context.Registry.Resolve(api.ProvisionManagerKey).(api.ProvisionManager)
	definitionManager := context.Registry.Resolve(api.DefinitionManagerKey).(api.DefinitionManager)
	handler := NewHandler(provisionManager, definitionManager, context.LogMonitor)

	router.Get("/health", handler.health)
	router.Post("/deployment", handler.deployment)
	router.Post("/activity-definition", handler.activityDefinition)
	router.Post("/deployment-definition", handler.deploymentDefinition)

	return nil
}
