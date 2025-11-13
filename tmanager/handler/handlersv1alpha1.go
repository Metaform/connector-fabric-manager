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
	"net/http"

	"github.com/metaform/connector-fabric-manager/common/handler"
	"github.com/metaform/connector-fabric-manager/common/query"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/metaform/connector-fabric-manager/tmanager/model/v1alpha1"
)

type TMHandler struct {
	handler.HttpHandler
	tenantService      api.TenantService
	participantService api.ParticipantProfileService
	cellService        api.CellService
	dataspaceService   api.DataspaceProfileService
}

func NewHandler(
	tenantService api.TenantService,
	participantService api.ParticipantProfileService,
	cellService api.CellService,
	dataspaceService api.DataspaceProfileService,
	monitor system.LogMonitor) *TMHandler {
	return &TMHandler{
		HttpHandler: handler.HttpHandler{
			Monitor: monitor,
		},
		tenantService:      tenantService,
		participantService: participantService,
		cellService:        cellService,
		dataspaceService:   dataspaceService,
	}
}

func (h *TMHandler) deployParticipant(w http.ResponseWriter, req *http.Request, tenantID string) {
	if h.InvalidMethod(w, req, http.MethodPost) {
		return
	}

	var newDeployment v1alpha1.NewParticipantProfileDeployment
	if !h.ReadPayload(w, req, &newDeployment) {
		return
	}

	properties := newDeployment.Properties
	if properties == nil {
		properties = make(map[string]any)
	}
	// TODO support specific cell selection
	profile, err := h.participantService.DeployProfile(
		req.Context(),
		tenantID,
		newDeployment.Identifier,
		*api.ToVPAMap(newDeployment.VPAProperties),
		properties)

	if err != nil {
		h.HandleError(w, err)
	}

	response := v1alpha1.ToParticipantProfile(profile)
	h.ResponseAccepted(w, response)
}

func (h *TMHandler) disposeParticipant(w http.ResponseWriter, req *http.Request, tenantID string, participantID string) {
	if h.InvalidMethod(w, req, http.MethodDelete) {
		return
	}

	err := h.participantService.DisposeProfile(req.Context(), tenantID, participantID)
	if err != nil {
		h.HandleError(w, err)
	}

	h.Accepted(w)
}

func (h *TMHandler) createTenant(w http.ResponseWriter, req *http.Request) {
	if h.InvalidMethod(w, req, http.MethodPost) {
		return
	}

	var newTenant v1alpha1.NewTenant
	if !h.ReadPayload(w, req, &newTenant) {
		return
	}

	tenant, err := h.tenantService.CreateTenant(req.Context(), v1alpha1.NewAPITenant(&newTenant))
	if err != nil {
		h.HandleError(w, err)
		return
	}

	response := v1alpha1.ToTenant(tenant)
	h.ResponseOK(w, response)
}

func (h *TMHandler) queryTenant(w http.ResponseWriter, req *http.Request) {
	if h.InvalidMethod(w, req, http.MethodPost) {
		return
	}
	var tenantQuery v1alpha1.Query
	if !h.ReadPayload(w, req, &tenantQuery) {
		return
	}
	predicate, err := query.ParsePredicate(tenantQuery.Predicate)
	if err != nil {
		h.WriteError(w, fmt.Sprintf("Client error: %v", err), http.StatusBadRequest)
		return
	}
	h.OK(w)
	_, err = w.Write([]byte("["))
	if err != nil {
		h.Monitor.Infow("Error writing response: %v", err)
		return
	}
	first := true
	for tenant, err := range h.tenantService.QueryTenants(req.Context(), predicate, api.PaginationOptions{
		Offset: 0,
		Limit:  10000,
	}) {
		if err != nil {
			h.Monitor.Infow("Error streaming tenant: %v", err)
			break
		}

		if !first {
			_, err = w.Write([]byte(","))
			if err != nil {
				h.Monitor.Infow("Error writing response: %v", err)
				return
			}
		}
		first = false

		response := v1alpha1.ToTenant(&tenant)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			h.Monitor.Infow("Error encoding tenant response: %v", err)
			break
		}

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	_, err = w.Write([]byte("]"))
	if err != nil {
		h.Monitor.Infow("Error writing response: %v", err)
		return
	}

}

func (h *TMHandler) getParticipantProfile(w http.ResponseWriter, req *http.Request, tenantID string, participantID string) {
	if h.InvalidMethod(w, req, http.MethodGet) {
		return
	}

	profile, err := h.participantService.GetProfile(req.Context(), tenantID, participantID)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	response := v1alpha1.ToParticipantProfile(profile)
	h.ResponseOK(w, response)
}

func (h *TMHandler) createCell(w http.ResponseWriter, req *http.Request) {
	if h.InvalidMethod(w, req, http.MethodPost) {
		return
	}

	var newCell v1alpha1.NewCell
	if !h.ReadPayload(w, req, &newCell) {
		return
	}

	cell := v1alpha1.NewAPICell(newCell)

	recordedCell, err := h.cellService.RecordExternalDeployment(req.Context(), cell)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	response := v1alpha1.ToCell(*recordedCell)
	h.ResponseOK(w, response)
}

func (h *TMHandler) createDataspaceProfile(w http.ResponseWriter, req *http.Request) {
	if h.InvalidMethod(w, req, http.MethodPost) {
		return
	}

	var newProfile v1alpha1.NewDataspaceProfile
	if !h.ReadPayload(w, req, &newProfile) {
		return
	}

	profile, err := h.dataspaceService.CreateProfile(req.Context(), newProfile.Artifacts, newProfile.Properties)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	response := v1alpha1.ToDataspaceProfile(profile)
	h.ResponseOK(w, response)
}

func (h *TMHandler) deployDataspaceProfile(w http.ResponseWriter, req *http.Request) {
	if h.InvalidMethod(w, req, http.MethodPost) {
		return
	}

	var newDeployment v1alpha1.NewDataspaceProfileDeployment
	if !h.ReadPayload(w, req, &newDeployment) {
		return
	}

	err := h.dataspaceService.DeployProfile(req.Context(), newDeployment.ProfileID, newDeployment.CellID)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	h.Accepted(w)
}

func (h *TMHandler) health(w http.ResponseWriter, _ *http.Request) {
	response := response{Message: "OK"}
	h.ResponseOK(w, response)
}
