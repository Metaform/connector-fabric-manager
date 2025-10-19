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
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/common/type"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/metaform/connector-fabric-manager/tmanager/tmstore"
)

type TMHandler struct {
	participantDeployer api.ParticipantDeployer
	cellStore           tmstore.EntityStore[api.Cell]
	dProfileStore       tmstore.EntityStore[api.DataspaceProfile]
	monitor             system.LogMonitor
}

func NewHandler(
	participantDeployer api.ParticipantDeployer,
	cellStore tmstore.EntityStore[api.Cell],
	dProfileStore tmstore.EntityStore[api.DataspaceProfile],
	monitor system.LogMonitor) *TMHandler {
	return &TMHandler{
		participantDeployer: participantDeployer,
		cellStore:           cellStore,
		dProfileStore:       dProfileStore,
		monitor:             monitor,
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
		switch e := err.(type) {
		case *_type.BadRequestError:
			http.Error(w, fmt.Sprintf("Bad request: %s", e.Message), http.StatusBadRequest)
		case *_type.SystemError:
			id := uuid.New().String()
			h.monitor.Infow("Internal Error [%s]: %v", id, err)
			http.Error(w, fmt.Sprintf("Internal server error occurred during participant deployment [%s]", id), http.StatusInternalServerError)
		case _type.FatalError:
			http.Error(w, "A fatal error occurred", http.StatusInternalServerError)
		default:
			http.Error(w, fmt.Sprintf("Failed to deploy participant: %s", err.Error()), http.StatusInternalServerError)
		}

	}
}
