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

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/metaform/connector-fabric-manager/tmanager/model/v1alpha1"
)

const json_content = "application/json"

type TMHandler struct {
	participantDeployer api.ParticipantProfileDeployer
	cellDeployer        api.CellDeployer
	dataspaceDeployer   api.DataspaceProfileDeployer
	monitor             system.LogMonitor
}

func NewHandler(
	participantDeployer api.ParticipantProfileDeployer,
	cellDeployer api.CellDeployer,
	dataspaceDeployer api.DataspaceProfileDeployer,
	monitor system.LogMonitor) *TMHandler {
	return &TMHandler{
		participantDeployer: participantDeployer,
		cellDeployer:        cellDeployer,
		dataspaceDeployer:   dataspaceDeployer,
		monitor:             monitor,
	}
}

// FIXME should take a type with properties
func (h *TMHandler) createParticipant(w http.ResponseWriter, req *http.Request) {
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

	var profileDeployment v1alpha1.NewParticipantProfileDeployment
	if err := json.Unmarshal(body, &profileDeployment); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// TODO support specific cell selection
	result, err := h.participantDeployer.DeployProfile(req.Context(), profileDeployment.Identifier, make(api.VpaPropMap), make(map[string]interface{}))
	if err != nil {
		handleError(w, err, h)
	}
	output := v1alpha1.ToParticipantProfile(result)
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", json_content)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		h.monitor.Infow("Failed to serialize profile response: %v", err)
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
	}
}

func (h *TMHandler) getParticipantProfile(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO externalize and pass as a parameter
	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "Missing identifier parameter", http.StatusBadRequest)
		return
	}

	result, err := h.participantDeployer.GetProfile(req.Context(), id)

	if err != nil {
		handleError(w, err, h)
		return
	}

	output := v1alpha1.ToParticipantProfile(result)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", json_content)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		h.monitor.Infow("Failed to serialize cell response: %v", err)
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
	}
}

func (h *TMHandler) createCell(w http.ResponseWriter, req *http.Request) {
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

	var newCell v1alpha1.NewCell
	if err := json.Unmarshal(body, &newCell); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// TODO NewCell validation
	cell := v1alpha1.NewAPICell(newCell)

	result, err := h.cellDeployer.RecordExternalDeployment(req.Context(), *cell)

	if err != nil {
		handleError(w, err, h)
		return
	}
	mCell := v1alpha1.ToCell(*result)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", json_content)
	if err := json.NewEncoder(w).Encode(mCell); err != nil {
		h.monitor.Infow("Failed to serialize cell response: %v", err)
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
	}
}

func (h *TMHandler) createDataspaceProfile(w http.ResponseWriter, req *http.Request) {
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

	var newProfile v1alpha1.NewDataspaceProfile
	if err := json.Unmarshal(body, &newProfile); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// TODO validation
	result, err := h.dataspaceDeployer.CreateProfile(req.Context(), newProfile.Artifacts, newProfile.Properties)

	if err != nil {
		handleError(w, err, h)
		return
	}

	// move to transformer
	mProfile := v1alpha1.DataspaceProfile{
		Entity: v1alpha1.Entity{
			ID:      result.ID,
			Version: result.Version,
		},
		Artifacts:   result.Artifacts,
		Deployments: []v1alpha1.DataspaceDeployment{
			// TODO finish
		},
		Properties: result.Properties,
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", json_content)
	if err := json.NewEncoder(w).Encode(mProfile); err != nil {
		h.monitor.Infow("Failed to serialize cell response: %v", err)
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
	}
}

func (h *TMHandler) deployDataspaceProfile(w http.ResponseWriter, req *http.Request) {
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

	var deployment v1alpha1.NewDataspaceProfileDeployment
	if err := json.Unmarshal(body, &deployment); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	err = h.dataspaceDeployer.DeployProfile(req.Context(), deployment.ProfileID, deployment.CellID)
	if err != nil {
		handleError(w, err, h)
		return
	}
}

func handleError(w http.ResponseWriter, err error, h *TMHandler) {
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
