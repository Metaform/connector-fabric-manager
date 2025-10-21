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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/model/v1alpha1"
)

type PMHandler struct {
	provisionManager api.ProvisionManager
	definitionStore  api.DefinitionStore
	monitor          system.LogMonitor
}

func NewHandler(provisionManager api.ProvisionManager, definitionStore api.DefinitionStore, monitor system.LogMonitor) *PMHandler {
	return &PMHandler{
		provisionManager: provisionManager,
		definitionStore:  definitionStore,
		monitor:          monitor,
	}
}

func (h *PMHandler) health(w http.ResponseWriter, _ *http.Request) {
	response := response{Message: "OK"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PMHandler) activityDefinition(w http.ResponseWriter, req *http.Request) {
	// Only allow POST requests
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	var definition v1alpha1.ActivityDefinition
	if err := json.Unmarshal(body, &definition); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	h.definitionStore.StoreActivityDefinition(v1alpha1.ToAPIActivityDefinition(&definition))
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
}

func (h *PMHandler) deploymentDefinition(w http.ResponseWriter, req *http.Request) {
	// Only allow POST requests
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	var definition v1alpha1.DeploymentDefinition
	if err := json.Unmarshal(body, &definition); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	h.definitionStore.StoreDeploymentDefinition(v1alpha1.ToAPIDeploymentDefinition(&definition))
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
}

func (h *PMHandler) deployment(w http.ResponseWriter, req *http.Request) {
	// Only allow POST requests
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	// Deserialize the DeploymentManifest from JSON
	var manifest model.DeploymentManifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	orchestration, err := h.provisionManager.Start(req.Context(), &manifest)
	if err != nil {
		h.handleError(w, err)
		return
	}
	// Return success response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orchestration); err != nil {
		h.monitor.Infow("Error encoding response: %v", err)
	}
}

func (h *PMHandler) handleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *types.BadRequestError:
		http.Error(w, fmt.Sprintf("Bad request: %s", e.Message), http.StatusBadRequest)
	case *types.SystemError:
		id := uuid.New().String()
		h.monitor.Infow("Internal Error [%s]: %v", id, err)
		http.Error(w, fmt.Sprintf("Internal server error occurred [%s]", id), http.StatusInternalServerError)
	case types.FatalError:
		http.Error(w, "A fatal error occurred", http.StatusInternalServerError)
	default:
		http.Error(w, fmt.Sprintf("Operation failed: %s", err.Error()), http.StatusInternalServerError)
	}
}
