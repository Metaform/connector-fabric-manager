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
	"github.com/metaform/connector-fabric-manager/dmodel"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

// participantGenerator generates participant profiles and VPAs that can be deployed to cells.
type participantGenerator struct {
	CellSelector api.CellSelector
}

func (g participantGenerator) Generate(
	identifier string,
	vpaProperties map[dmodel.VPAType]any,
	properties map[string]any,
	cells []api.Cell,
	dProfiles []api.DataspaceProfile) (*api.ParticipantProfile, error) {

	// TODO process vpaProperties properties - decompose properties into VPA properties

	cell, err := g.CellSelector(dmodel.VpaDeploymentType, cells, dProfiles)
	if err != nil {
		return nil, err
	}

	connector := g.generateVPA(dmodel.ConnectorType, cell)
	cservice := g.generateVPA(dmodel.CredentialServiceType, cell)
	dplane := g.generateVPA(dmodel.DataPlaneType, cell)
	vpas := []api.VirtualParticipantAgent{connector, cservice, dplane}

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
func (g participantGenerator) generateVPA(vpaType dmodel.VPAType, cell *api.Cell) api.VirtualParticipantAgent {
	connector := api.VirtualParticipantAgent{
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
	return connector
}
