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
	store                tmstore.TManagerStore
}

func (d participantDeployer) Deploy(
	ctx context.Context,
	identifier string,
	vpaProperties api.VpaPropMap,
	properties map[string]any) error {

	// TODO perform property validation against a custom schema
	return d.trxContext.Execute(ctx, func(ctx context.Context) error {
		cells, err := d.store.GetCells()
		if err != nil {
			return err
		}
		dProfiles, err := d.store.GetDataspaceProfiles()
		if err != nil {
			return err
		}

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
	TenantStore tmstore.TManagerStore
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
