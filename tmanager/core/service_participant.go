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
	"strings"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/collection"
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/common/types"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

type participantService struct {
	participantGenerator participantGenerator
	provisionClient      api.ProvisionClient
	trxContext           store.TransactionContext
	participantStore     store.EntityStore[*api.ParticipantProfile]
	cellStore            store.EntityStore[*api.Cell]
	dataspaceStore       store.EntityStore[*api.DataspaceProfile]
	monitor              system.LogMonitor
}

func (d participantService) GetProfile(ctx context.Context, tenantID string, participantID string) (*api.ParticipantProfile, error) {
	return store.Trx[api.ParticipantProfile](d.trxContext).AndReturn(ctx, func(ctx context.Context) (*api.ParticipantProfile, error) {
		profile, err := d.participantStore.FindById(ctx, participantID)
		if err != nil {
			return nil, err
		}
		if profile.TenantID != tenantID {
			return nil, types.ErrNotFound
		}
		return profile, nil
	})
}

func (d participantService) DeployProfile(
	ctx context.Context,
	tenantID string,
	identifier string,
	vpaProperties api.VPAPropMap,
	properties map[string]any) (*api.ParticipantProfile, error) {

	// TODO perform property validation against a custom schema
	return store.Trx[api.ParticipantProfile](d.trxContext).AndReturn(ctx, func(ctx context.Context) (*api.ParticipantProfile, error) {
		cells, err := collection.CollectAllDeref(d.cellStore.GetAll(ctx))
		if err != nil {
			return nil, err
		}

		dProfiles, err := collection.CollectAllDeref(d.dataspaceStore.GetAll(ctx))
		if err != nil {
			return nil, err
		}

		participantProfile, err := d.participantGenerator.Generate(
			identifier,
			tenantID,
			vpaProperties,
			properties,
			cells,
			dProfiles)
		if err != nil {
			return nil, err
		}

		oManifest := model.OrchestrationManifest{
			ID:                uuid.New().String(),
			CorrelationID:     participantProfile.ID,
			OrchestrationType: model.VPADeployType,
			Payload:           make(map[string]any),
		}

		oManifest.Payload[model.ParticipantIdentifier] = participantProfile.Identifier

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

		oManifest.Payload[model.VPAData] = vpaManifests
		result, err := d.participantStore.Create(ctx, participantProfile)
		if err != nil {
			return nil, fmt.Errorf("error creating participant %s: %w", identifier, err)
		}

		// Only send the orchestration message if the storage operation succeeded. If the send fails, the transaction
		// will be rolled back.
		err = d.provisionClient.Send(ctx, oManifest)
		if err != nil {
			return nil, fmt.Errorf("error deploying participant %s: %w", identifier, err)
		}

		return result, nil
	})
}

func (d participantService) DisposeProfile(ctx context.Context, tenantID string, participantID string) error {
	return d.trxContext.Execute(ctx, func(c context.Context) error {
		profile, err := d.participantStore.FindById(c, participantID)
		if err != nil {
			return err
		}
		if profile.TenantID != tenantID {
			return types.ErrNotFound
		}
		states := make([]string, 0, len(profile.VPAs))
		for _, vpa := range profile.VPAs {
			if vpa.State != api.DeploymentStateActive {
				states = append(states, vpa.ID+":"+vpa.State.String())
			}
		}
		if len(states) > 0 {
			return fmt.Errorf("cannot dispose VPAs %s in states: %s", participantID, strings.Join(states, ","))
		}
		stateData, found := profile.Properties[model.VPAStateData]
		if !found {
			return fmt.Errorf("profile is not deployed or is missing state data: %s", participantID)
		}

		oManifest := model.OrchestrationManifest{
			ID:                uuid.New().String(),
			CorrelationID:     participantID,
			OrchestrationType: model.VPADisposeType,
			Payload:           make(map[string]any),
		}

		oManifest.Payload[model.ParticipantIdentifier] = profile.Identifier
		oManifest.Payload[model.VPAStateData] = stateData

		vpaManifests := make([]model.VPAManifest, 0, len(profile.VPAs))
		for _, vpa := range profile.VPAs {
			vpaManifest := model.VPAManifest{
				ID:         vpa.ID,
				VPAType:    vpa.Type,
				Cell:       vpa.Cell.ID,
				Properties: vpa.Properties,
			}
			vpaManifests = append(vpaManifests, vpaManifest)

			// Set to disposing
			vpa.State = api.DeploymentStateDisposing
		}

		oManifest.Payload[model.VPAData] = vpaManifests

		err = d.participantStore.Update(c, profile)
		if err != nil {
			return fmt.Errorf("error disposing participant %s: %w", participantID, err)
		}

		// Only send the orchestration message if the storage operation succeeded. If the send fails, the transaction
		// will be rolled back.
		err = d.provisionClient.Send(ctx, oManifest)
		if err != nil {
			return fmt.Errorf("error disposing participant %s: %w", participantID, err)
		}

		return nil
	})
}

type vpaCallbackHandler struct {
	participantStore store.EntityStore[*api.ParticipantProfile]
	trxContext       store.TransactionContext
	monitor          system.LogMonitor
}

func (h vpaCallbackHandler) handleDeploy(ctx context.Context, response model.OrchestrationResponse) error {
	return h.handle(ctx, response, func(profile *api.ParticipantProfile, resp model.OrchestrationResponse) {
		// Place all output values under VPStateData key
		vpaProps := make(map[string]any)
		for key, value := range resp.Properties {
			vpaProps[key] = value
		}
		profile.Properties[model.VPAStateData] = vpaProps

		for i, vpa := range profile.VPAs {
			vpa.State = api.DeploymentStateActive
			// TODO update timestamp based on returned data
			profile.VPAs[i] = vpa // Use range index because vpa is a copy
		}
	})
}

func (h vpaCallbackHandler) handleDispose(ctx context.Context, response model.OrchestrationResponse) error {
	return h.handle(ctx, response, func(profile *api.ParticipantProfile, resp model.OrchestrationResponse) {
		for i, vpa := range profile.VPAs {
			// Update state
			vpa.State = api.DeploymentStateDisposed
			// TODO update timestamp based on returned data
			profile.VPAs[i] = vpa // Use range index because vpa is a copy
		}
	})
}

// handle processes the asynchronous response to participant VPA deployment request.
func (h vpaCallbackHandler) handle(
	ctx context.Context,
	response model.OrchestrationResponse,
	handler func(profile *api.ParticipantProfile, resp model.OrchestrationResponse)) error {

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
			handler(profile, response)
		default:
			// TODO update VPA status
			profile.Error = true
			profile.ErrorDetail = response.ErrorDetail
		}
		err = h.participantStore.Update(c, profile)
		if err != nil {
			return fmt.Errorf("error updating participant profile %s processing VPA response for manifest %s: %w", response.CorrelationID, response.ManifestID, err)
		}
		return nil
	})
}
