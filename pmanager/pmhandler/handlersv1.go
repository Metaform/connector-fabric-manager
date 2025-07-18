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

package pmhandler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"io"
	"net/http"
)

type PMHandler struct {
	provisionManager api.ProvisionManager
	definitionStore  api.DefinitionStore
	logMonitor       monitor.LogMonitor
}

func NewHandler(provisionManager api.ProvisionManager, definitionStore api.DefinitionStore, logMonitor monitor.LogMonitor) *PMHandler {
	return &PMHandler{
		provisionManager: provisionManager,
		definitionStore:  definitionStore,
		logMonitor:       logMonitor,
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

	var definition api.ActivityDefinition
	if err := json.Unmarshal(body, &definition); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	h.definitionStore.StoreActivityDefinition(&definition)
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

	var definition api.DeploymentDefinition
	if err := json.Unmarshal(body, &definition); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	h.definitionStore.StoreDeploymentDefinition(&definition)
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
	var manifest api.DeploymentManifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	orchestration, err := h.provisionManager.Start(req.Context(), &manifest)
	if err != nil {
		switch {
		case model.IsClientError(err):
			http.Error(w, fmt.Sprintf("Invalid deployment: %s", err.Error()), http.StatusBadRequest)
			return
		case model.IsRecoverable(err):
			id := uuid.New().String()
			h.logMonitor.Infof("Recoverable error encountered during deployment [%s]: %w ", id, err)
			http.Error(w, fmt.Sprintf("Recoverable error encountered during deployment [%s], id"), http.StatusServiceUnavailable)
			return
		case model.IsFatal(err):
			id := uuid.New().String()
			h.logMonitor.Infof("Fatal error encountered during deployment [%s]: %w ", id, err)
			http.Error(w, fmt.Sprintf("Fatal error encountered during deployment [%s], id"), http.StatusInternalServerError)
			return
		default:
			http.Error(w, "Failed to initiate orchestration", http.StatusInternalServerError)
			return
		}
	}
	// Return success response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orchestration); err != nil {
		h.logMonitor.Infow("Error encoding response: %v", err)
	}
}
