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

package v1alpha1

import (
	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

func ToParticipantProfile(input *api.ParticipantProfile) *ParticipantProfile {
	return &ParticipantProfile{
		Entity: Entity{
			ID:      input.ID,
			Version: input.Version,
		},
		Identifier:  input.Identifier,
		VPAs:        ToVPACollection(input),
		Properties:  input.Properties,
		Error:       input.Error,
		ErrorDetail: input.ErrorDetail,
	}
}

func ToVPACollection(input *api.ParticipantProfile) []VirtualParticipantAgent {
	vpas := make([]VirtualParticipantAgent, len(input.VPAs))
	for i, vpa := range input.VPAs {
		vpas[i] = *ToVPA(vpa)
	}
	return vpas
}

func ToVPA(input api.VirtualParticipantAgent) *VirtualParticipantAgent {
	return &VirtualParticipantAgent{
		DeployableEntity: DeployableEntity{
			Entity: Entity{
				ID:      input.ID,
				Version: input.Version,
			},
			State:          input.State.String(),
			StateTimestamp: input.StateTimestamp,
		},
		Type:       input.Type,
		Cell:       *ToCell(input.Cell),
		Properties: input.Properties,
	}
}

func ToCell(input api.Cell) *Cell {
	return &Cell{
		Entity: Entity{
			ID:      input.ID,
			Version: input.Version,
		},
		NewCell: NewCell{
			State:          input.State.String(),
			StateTimestamp: input.StateTimestamp,
			Properties:     input.Properties,
		},
	}
}

func ToAPIParticipantProfile(input *ParticipantProfile) *api.ParticipantProfile {
	return &api.ParticipantProfile{
		Entity: api.Entity{
			ID:      input.ID,
			Version: input.Version,
		},
		Identifier:  input.Identifier,
		VPAs:        ToAPIVPACollection(input.VPAs),
		Properties:  api.ToProperties(input.Properties),
		Error:       input.Error,
		ErrorDetail: input.ErrorDetail,
	}
}

func ToAPIVPACollection(vpas []VirtualParticipantAgent) []api.VirtualParticipantAgent {
	apiVPAs := make([]api.VirtualParticipantAgent, len(vpas))
	for i, vpa := range vpas {
		apiVPAs[i] = *ToAPIVPA(vpa)
	}
	return apiVPAs
}

func ToAPIVPA(input VirtualParticipantAgent) *api.VirtualParticipantAgent {
	state, _ := api.ToDeploymentState(input.State)
	return &api.VirtualParticipantAgent{
		DeployableEntity: api.DeployableEntity{
			Entity: api.Entity{
				ID:      input.ID,
				Version: input.Version,
			},
			State:          state,
			StateTimestamp: input.StateTimestamp.UTC(), // Force UTC
		},
		Type:       input.Type,
		Cell:       *ToAPICell(input.Cell),
		Properties: api.ToProperties(input.Properties),
	}
}

func ToAPICell(input Cell) *api.Cell {
	state, _ := api.ToDeploymentState(input.State)
	return &api.Cell{
		DeployableEntity: api.DeployableEntity{
			Entity: api.Entity{
				ID:      input.ID,
				Version: input.Version,
			},
			State:          state,
			StateTimestamp: input.StateTimestamp.UTC(), // Force UTC
		},
		Properties: api.ToProperties(input.Properties),
	}
}

func NewAPICell(input NewCell) *api.Cell {
	state, _ := api.ToDeploymentState(input.State)
	return &api.Cell{
		DeployableEntity: api.DeployableEntity{
			Entity: api.Entity{
				ID:      uuid.New().String(),
				Version: 0,
			},
			State:          state,
			StateTimestamp: input.StateTimestamp.UTC(), // Force UTC
		},
		Properties: api.ToProperties(input.Properties),
	}
}

func ToDataspaceProfile(input *api.DataspaceProfile) *DataspaceProfile {
	deployments := make([]DataspaceDeployment, len(input.Deployments))
	for i, deployment := range input.Deployments {
		deployments[i] = DataspaceDeployment{
			DeployableEntity: DeployableEntity{
				Entity: Entity{
					ID:      deployment.ID,
					Version: deployment.Version,
				},
				State:          deployment.State.String(),
				StateTimestamp: deployment.StateTimestamp.UTC(), // Convert to UTC
			},
			CellID:     deployment.Cell.ID,
			Properties: deployment.Properties,
		}
	}

	return &DataspaceProfile{
		Entity: Entity{
			ID:      input.ID,
			Version: input.Version,
		},
		Artifacts:   input.Artifacts,
		Deployments: deployments,
		Properties:  input.Properties,
	}
}
