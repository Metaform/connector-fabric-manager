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
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/metaform/connector-fabric-manager/tmanager/model/v1alpha1"
)

type TMHandler struct {
	handler.HttpHandler
	participantDeployer api.ParticipantProfileDeployer
	cellDeployer        api.CellDeployer
	dataspaceDeployer   api.DataspaceProfileDeployer
}

func NewHandler(
	participantDeployer api.ParticipantProfileDeployer,
	cellDeployer api.CellDeployer,
	dataspaceDeployer api.DataspaceProfileDeployer,
	monitor system.LogMonitor) *TMHandler {
	return &TMHandler{
		HttpHandler: handler.HttpHandler{
			Monitor: monitor,
		},
		participantDeployer: participantDeployer,
		cellDeployer:        cellDeployer,
		dataspaceDeployer:   dataspaceDeployer,
	}
}

func (h *TMHandler) createDeployParticipant(w http.ResponseWriter, req *http.Request) {
	if h.InvalidMethod(w, req, http.MethodPost) {
		return
	}

	var newDeployment v1alpha1.NewParticipantProfileDeployment
	if !h.ReadPayload(w, req, &newDeployment) {
		return
	}

	// TODO support specific cell selection
	profile, err := h.participantDeployer.DeployProfile(
		req.Context(),
		newDeployment.Identifier,
		*api.ToVPAMap(newDeployment.VPAProperties),
		newDeployment.Properties)
	if err != nil {
		h.HandleError(w, err)
	}

	response := v1alpha1.ToParticipantProfile(profile)
	h.ResponseAccepted(w, response)
}

func (h *TMHandler) getParticipantProfile(w http.ResponseWriter, req *http.Request, id string) {
	if h.InvalidMethod(w, req, http.MethodGet) {
		return
	}

	profile, err := h.participantDeployer.GetProfile(req.Context(), id)
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

	recordedCell, err := h.cellDeployer.RecordExternalDeployment(req.Context(), *cell)
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

	profile, err := h.dataspaceDeployer.CreateProfile(req.Context(), newProfile.Artifacts, newProfile.Properties)
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

	err := h.dataspaceDeployer.DeployProfile(req.Context(), newDeployment.ProfileID, newDeployment.CellID)
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
