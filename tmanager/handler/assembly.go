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
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/metaform/connector-fabric-manager/assembly/routing"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
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
	return []system.ServiceType{
		api.ParticipantProfileDeployerKey,
		api.CellDeployerKey,
		api.DataspaceProfileDeployerKey,
		routing.RouterKey}
}

func (h *HandlerServiceAssembly) Init(context *system.InitContext) error {
	router := context.Registry.Resolve(routing.RouterKey).(chi.Router)
	router.Use(middleware.Recoverer)

	participantDeployer := context.Registry.Resolve(api.ParticipantProfileDeployerKey).(api.ParticipantProfileDeployer)
	cellDeployer := context.Registry.Resolve(api.CellDeployerKey).(api.CellDeployer)
	dataspaceDeployer := context.Registry.Resolve(api.DataspaceProfileDeployerKey).(api.DataspaceProfileDeployer)

	handler := NewHandler(participantDeployer, cellDeployer, dataspaceDeployer, context.LogMonitor)

	router.Get("/participants/{id}", func(w http.ResponseWriter, req *http.Request) {
		id, found := handler.ExtractPathVariable(w, req, "id")
		if !found {
			return
		}
		handler.getParticipantProfile(w, req, id)
	})

	router.Post("/participants", handler.createDeployParticipant)
	router.Post("/cells", handler.createCell)
	router.Post("/dataspace-profiles", handler.createDataspaceProfile)
	router.Post("/dataspace-profiles/{id}/deployments", handler.deployDataspaceProfile)

	return nil
}
