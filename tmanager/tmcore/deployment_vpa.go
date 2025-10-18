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
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/dmodel"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/metaform/connector-fabric-manager/tmanager/tmstore"
)

type participantDeployer struct {
	participantGenerator participantGenerator
	deploymentClient     api.DeploymentClient
	trxContext           store.TransactionContext
}

func (d participantDeployer) Deploy(
	ctx context.Context,
	identifier string,
	vpaProperties api.VpaPropMap,
	properties map[string]any) error {

	// TODO perform property validation against a custom schema
	return d.trxContext.Execute(ctx, func(ctx context.Context) error {
		// FIXME get cells and dataspace profiles from store
		cells, dProfiles := seedData()

		participantProfile, err := d.participantGenerator.Generate(
			identifier,
			vpaProperties,
			properties,
			cells,
			dProfiles)
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
				Properties: vpa.Properties,
			}
			vpaManifests = append(vpaManifests, vpaManifest)
		}

		dManifest.Payload[dmodel.VpaPayloadType] = vpaManifests

		err = d.deploymentClient.Deploy(ctx, dManifest)
		if err != nil {
			return fmt.Errorf("error deploying participant %s: %w", identifier, err)
		}

		// TODO persist
		return nil
	})
}

type vpaDeploymentCallbackHandler struct {
	TenantStore tmstore.TenantStore
}

func (h vpaDeploymentCallbackHandler) handle(_ context.Context, response dmodel.DeploymentResponse) error {
	if !response.Success {
		fmt.Println("Deployment failed:" + response.ErrorDetail)
		// TODO move to error state
		return nil
	}
	fmt.Println("Deployment succeeded:" + response.ManifestID)
	return nil
}

// seedData temporary function to initialize and return sample cells and dataspace profiles for use in deployment workflows.
func seedData() ([]api.Cell, []api.DataspaceProfile) {
	cells := []api.Cell{
		{
			DeployableEntity: api.DeployableEntity{
				Entity: api.Entity{
					ID:      "cell-001",
					Version: 1,
				},
				State:          api.DeploymentStateActive,
				StateTimestamp: time.Now(),
			},
			Properties: api.Properties{
				"region": "us-east-1",
				"type":   "kubernetes",
			},
		},
	}

	dProfiles := []api.DataspaceProfile{
		{
			Entity: api.Entity{
				ID:      "dataspace-profile-001",
				Version: 1,
			},
			Artifacts: []string{"connector-runtime", "policy-engine"},
			Deployments: []api.DataspaceDeployment{
				{
					DeployableEntity: api.DeployableEntity{
						Entity: api.Entity{
							ID:      "deployment-001",
							Version: 1,
						},
						State:          api.DeploymentStateActive,
						StateTimestamp: time.Now(),
					},
					Cell:       cells[0], // Reference to the first cell
					Properties: api.Properties{},
				},
			},
			Properties: api.Properties{},
		},
	}
	return cells, dProfiles
}
