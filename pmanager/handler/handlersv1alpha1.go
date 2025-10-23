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

	"github.com/metaform/connector-fabric-manager/common/handler"
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/model/v1alpha1"
)

type PMHandler struct {
	handler.HttpHandler
	provisionManager  api.ProvisionManager
	definitionManager api.DefinitionManager
}

func NewHandler(provisionManager api.ProvisionManager, definitionManager api.DefinitionManager, monitor system.LogMonitor) *PMHandler {
	return &PMHandler{
		HttpHandler: handler.HttpHandler{
			Monitor: monitor,
		},
		provisionManager:  provisionManager,
		definitionManager: definitionManager,
	}
}

func (h *PMHandler) activityDefinition(w http.ResponseWriter, req *http.Request) {
	if h.InvalidMethod(w, req, http.MethodPost) {
		return
	}

	var definition v1alpha1.ActivityDefinition
	if !h.ReadPayload(w, req, &definition) {
		return
	}

	_, err := h.definitionManager.CreateActivityDefinition(req.Context(), v1alpha1.ToAPIActivityDefinition(&definition))
	if err != nil {
		h.HandleError(w, err)
		return
	}

	h.Created(w)
}

func (h *PMHandler) deploymentDefinition(w http.ResponseWriter, req *http.Request) {
	if h.InvalidMethod(w, req, http.MethodPost) {
		return
	}

	var definition v1alpha1.DeploymentDefinition
	if !h.ReadPayload(w, req, &definition) {
		return
	}

	_, err := h.definitionManager.CreateDeploymentDefinition(req.Context(), v1alpha1.ToAPIDeploymentDefinition(&definition))
	if err != nil {
		h.HandleError(w, err)
		return
	}

	h.Created(w)
}

func (h *PMHandler) deployment(w http.ResponseWriter, req *http.Request) {
	if h.InvalidMethod(w, req, http.MethodPost) {
		return
	}

	var manifest model.DeploymentManifest
	if !h.ReadPayload(w, req, &manifest) {
		return
	}

	orchestration, err := h.provisionManager.Start(req.Context(), &manifest)
	if err != nil {
		h.HandleError(w, err)
		return
	}
	h.ResponseAccepted(w, orchestration)
}

func (h *PMHandler) health(w http.ResponseWriter, _ *http.Request) {
	response := response{Message: "OK"}
	h.ResponseOK(w, response)
}
