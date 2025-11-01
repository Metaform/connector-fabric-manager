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

package core

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/collection"
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

type participantDeployer struct {
	participantGenerator participantGenerator
	deploymentClient     api.DeploymentClient
	trxContext           store.TransactionContext
	participantStore     api.EntityStore[api.ParticipantProfile]
	cellStore            api.EntityStore[api.Cell]
	dataspaceStore       api.EntityStore[api.DataspaceProfile]
}

func (d participantDeployer) GetProfile(ctx context.Context, profileID string) (*api.ParticipantProfile, error) {
	return store.Trx[api.ParticipantProfile](d.trxContext).AndReturn(ctx, func(ctx context.Context) (*api.ParticipantProfile, error) {
		return d.participantStore.FindById(ctx, profileID)
	})
}

func (d participantDeployer) DeployProfile(
	ctx context.Context,
	identifier string,
	vpaProperties api.VPAPropMap,
	properties map[string]any) (*api.ParticipantProfile, error) {

	// TODO perform property validation against a custom schema
	return store.Trx[api.ParticipantProfile](d.trxContext).AndReturn(ctx, func(ctx context.Context) (*api.ParticipantProfile, error) {
		cells, err := collection.CollectAll(d.cellStore.GetAll(ctx))
		if err != nil {
			return nil, err
		}

		dProfiles, err := collection.CollectAll(d.dataspaceStore.GetAll(ctx))
		if err != nil {
			return nil, err
		}

		participantProfile, err := d.participantGenerator.Generate(
			identifier,
			vpaProperties,
			properties,
			cells,
			dProfiles)
		if err != nil {
			return nil, err
		}

		dManifest := model.DeploymentManifest{
			ID:             uuid.New().String(),
			CorrelationID:  participantProfile.ID,
			DeploymentType: model.VPADeploymentType,
			Payload:        make(map[string]any),
		}

		dManifest.Payload[model.ParticipantIdentifier] = participantProfile.Identifier

		vpaManifests := make([]model.VPAManifest, 0, len(participantProfile.VPAs))
		for _, vpa := range participantProfile.VPAs {
			vpaManifest := model.VPAManifest{
				ID:         vpa.ID,
				VPAType:    vpa.Type,
				Cell:       vpa.Cell.ID,
				Properties: vpa.Properties,
			}
			vpaManifests = append(vpaManifests, vpaManifest)
		}

		dManifest.Payload[model.VPAPayloadType] = vpaManifests
		result, err := d.participantStore.Create(ctx, participantProfile)
		if err != nil {
			return nil, fmt.Errorf("error creating participant %s: %w", identifier, err)
		}

		// Only send the deployment message if the storage operation succeeded. If the deployment fails, the transaction
		// will be rolled back.
		err = d.deploymentClient.Send(ctx, dManifest)
		if err != nil {
			return nil, fmt.Errorf("error deploying participant %s: %w", identifier, err)
		}

		return result, nil
	})
}

func (d participantDeployer) DisposeProfile(ctx context.Context, identifier string) error {
	//TODO implement me
	panic("implement me")
}

type vpaDeploymentCallbackHandler struct {
	participantStore api.EntityStore[api.ParticipantProfile]
	trxContext       store.TransactionContext
	monitor          system.LogMonitor
}

func (h vpaDeploymentCallbackHandler) handle(ctx context.Context, response model.DeploymentResponse) error {
	return h.trxContext.Execute(ctx, func(c context.Context) error {
		// Note de-duplication does not need to be performed as this operation is idempotent
		profile, err := h.participantStore.FindById(c, response.CorrelationID)
		if err != nil {
			h.monitor.Infof("Error retrieving participant profile '%s' for manifest %s: %w", response.CorrelationID, response.ManifestID, err)
			// Do not return error as this is fatal and the message must be acked
			return nil
		}
		switch {
		case response.Success:
			// Place all output values under VPStateData key
			vpaProps := make(map[string]any)
			for key, value := range response.Properties {
				vpaProps[key] = value
			}
			profile.Properties[model.VPAStateData] = vpaProps

			for i, vpa := range profile.VPAs {
				vpa.State = api.DeploymentStateActive
				// TODO update timestamp based on returned data
				profile.VPAs[i] = vpa // Use range index because vpa is a copy
			}
		default:
			// TODO update VPA status
			profile.Error = true
			profile.ErrorDetail = response.ErrorDetail
		}
		err = h.participantStore.Update(c, profile)
		if err != nil {
			return fmt.Errorf("error updating participant profile %s processing response for manifest %s: %w", response.CorrelationID, response.ManifestID, err)
		}
		return nil
	})
}
