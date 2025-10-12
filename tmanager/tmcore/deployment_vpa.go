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

package tmcore

import (
	"context"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/dmodel"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/metaform/connector-fabric-manager/tmanager/tmstore"
)

type participantDeployer struct {
	participantGenerator participantGenerator
}

func (t participantDeployer) Deploy(
	_ context.Context,
	identifier string,
	properties map[string]any) error {

	// TODO perform property validation against a custom schema

	// TODO get cells and dataspace profiles from store
	cells := make([]api.Cell, 0)
	dProfiles := make([]api.DataspaceProfile, 0)

	participantProfile, err := t.participantGenerator.Generate(identifier, properties, cells, dProfiles)
	if err != nil {
		return err
	}
	dManifest := dmodel.DeploymentManifest{
		ID:             uuid.New().String(),
		DeploymentType: dmodel.VpaDeploymentType,
		Payload:        make(map[string]any),
	}

	vpaManifests := make([]dmodel.VPAManifest, 0, len(participantProfile.VPAs))
	for _, vpa := range participantProfile.VPAs {
		vpaManifest := dmodel.VPAManifest{
			ID:         vpa.ID,
			VPAType:    vpa.Type,
			Cell:       vpa.Cell.ID,
			Properties: vpa.DeploymentProperties,
		}
		vpaManifests = append(vpaManifests, vpaManifest)
	}

	dManifest.Payload[dmodel.VpaPayloadType] = vpaManifests
	// TODO finish
	return nil
}

type vpaDeploymentCallbackHandler struct {
	TenantStore tmstore.TenantStore
}

func (h vpaDeploymentCallbackHandler) handle(_ context.Context, response api.DeploymentResponse) error {
	if !response.Success {
		// TODO move to error state
		return nil
	}
	return nil
}
