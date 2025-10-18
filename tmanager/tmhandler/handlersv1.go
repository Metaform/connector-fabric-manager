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
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

type TMHandler struct {
	participantDeployer api.ParticipantDeployer
	logMonitor          monitor.LogMonitor
}

func NewHandler(participantDeployer api.ParticipantDeployer, logMonitor monitor.LogMonitor) *TMHandler {
	return &TMHandler{
		participantDeployer: participantDeployer,
		logMonitor:          logMonitor,
	}
}

func (h *TMHandler) deployParticipant(w http.ResponseWriter, req *http.Request) {
	// Only allow POST requests
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO externalize and pass as a parameter
	identifier := chi.URLParam(req, "id")
	if identifier == "" {
		http.Error(w, "Missing identifier parameter", http.StatusBadRequest)
		return
	}

	// Read the request body
	_, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	err = h.participantDeployer.Deploy(req.Context(), identifier, make(api.VpaPropMap), make(map[string]interface{}))
	if err != nil {
		// TODO distinguish
		http.Error(w, "Failed to deploy participant", http.StatusInternalServerError)
		return
	}
}
