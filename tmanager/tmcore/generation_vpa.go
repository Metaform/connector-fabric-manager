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
	"time"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

// participantGenerator generates participant profiles and VPAs that can be deployed to cells.
type participantGenerator struct {
	CellSelector api.CellSelector
}

func (g participantGenerator) Generate(
	identifier string,
	vpaProperties api.VpaPropMap,
	properties map[string]any,
	cells []api.Cell,
	dProfiles []api.DataspaceProfile) (*api.ParticipantProfile, error) {

	cell, err := g.CellSelector(model.VpaDeploymentType, cells, dProfiles)
	if err != nil {
		return nil, err
	}

	connector := g.generateVPA(model.ConnectorType, vpaProperties, cell)
	cService := g.generateVPA(model.CredentialServiceType, vpaProperties, cell)
	dPlane := g.generateVPA(model.DataPlaneType, vpaProperties, cell)
	vpas := []api.VirtualParticipantAgent{connector, cService, dPlane}

	pProfile := &api.ParticipantProfile{
		Entity: api.Entity{
			ID:      uuid.New().String(),
			Version: 0,
		},
		Identifier:        identifier,
		DataSpaceProfiles: dProfiles,
		VPAs:              vpas,
		Properties:        properties,
	}
	return pProfile, nil
}

// generateVPA creates a VPA targeted at given cell.
func (g participantGenerator) generateVPA(
	vpaType model.VPAType,
	vpaProperties api.VpaPropMap,
	cell *api.Cell) api.VirtualParticipantAgent {

	vpa := api.VirtualParticipantAgent{
		DeployableEntity: api.DeployableEntity{
			Entity: api.Entity{
				ID:      uuid.New().String(),
				Version: 0,
			},
			State:          api.DeploymentStateActive,
			StateTimestamp: time.Now().UTC(),
		},
		Type:       vpaType,
		Cell:       *cell,
		Properties: make(api.Properties),
	}

	// Decompose the properties and add them to the VPA
	props, found := vpaProperties[vpaType]
	if found {
		vpa.Properties = make(api.Properties)
		for k, v := range props {
			vpa.Properties[k] = v
		}
	}

	return vpa
}
