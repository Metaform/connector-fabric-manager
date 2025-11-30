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
	"time"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

type dataspaceProfileService struct {
	trxContext   store.TransactionContext
	profileStore store.EntityStore[*api.DataspaceProfile]
	cellStore    store.EntityStore[*api.Cell]
}

func (d dataspaceProfileService) GetProfile(ctx context.Context, profileID string) (*api.DataspaceProfile, error) {
	return d.profileStore.FindByID(ctx, profileID)
}

func (d dataspaceProfileService) CreateProfile(ctx context.Context, artifacts []string, properties map[string]any) (*api.DataspaceProfile, error) {
	return store.Trx[api.DataspaceProfile](d.trxContext).AndReturn(ctx, func(ctx context.Context) (*api.DataspaceProfile, error) {
		profile := &api.DataspaceProfile{
			Entity: api.Entity{
				ID:      uuid.New().String(),
				Version: 0,
			},
			Artifacts:   artifacts,
			Deployments: make([]api.DataspaceDeployment, 0),
			Properties:  properties,
		}
		return d.profileStore.Create(ctx, profile)
	})
}

func (d dataspaceProfileService) DeployProfile(ctx context.Context, profileID string, cellID string) error {
	return d.trxContext.Execute(ctx, func(_ context.Context) error {
		profile, err := d.profileStore.FindByID(ctx, profileID)
		if err != nil {
			return err
		}

		cell, err := d.cellStore.FindByID(ctx, cellID)
		if err != nil {
			return err
		}

		// TODO validate not already deployed and handle deployment
		profile.Deployments = append(profile.Deployments, api.DataspaceDeployment{
			DeployableEntity: api.DeployableEntity{
				Entity: api.Entity{
					ID:      uuid.New().String(),
					Version: 0,
				},
				State:          api.DeploymentStateActive,
				StateTimestamp: time.Time{}.UTC(),
			},
			Cell:       *cell,
			Properties: make(map[string]any),
		})
		err = d.profileStore.Update(ctx, profile)
		if err != nil {
			return err
		}
		return nil

	})

}
