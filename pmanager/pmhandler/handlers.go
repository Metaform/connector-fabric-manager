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
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"io"
	"net/http"
)

type PMHandler struct {
	provisionManager api.ProvisionManager
	logMonitor       monitor.LogMonitor
}

func NewHandler(provisionManager api.ProvisionManager, logMonitor monitor.LogMonitor) *PMHandler {
	return &PMHandler{
		provisionManager: provisionManager,
		logMonitor:       logMonitor,
	}
}

func (h *PMHandler) health(w http.ResponseWriter, _ *http.Request) {
	response := response{Message: "OK"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	// Validate required fields
	if manifest.ID == "" {
		http.Error(w, "Missing required field: id", http.StatusBadRequest)
		return
	}

	if manifest.DeploymentType == "" {
		http.Error(w, "Missing required field: deploymentType", http.StatusBadRequest)
		return
	}

	orchestration, err := h.provisionManager.Start(req.Context(), &manifest)
	if err != nil {
		http.Error(w, "Failed to initiate orchestration", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orchestration); err != nil {
		h.logMonitor.Infow("Error encoding response: %v", err)
	}
}
