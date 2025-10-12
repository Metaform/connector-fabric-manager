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
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/dmodel"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParticipantProfileGenerator_Generate(t *testing.T) {
	now := time.Now().UTC()

	t.Run("successful generation", func(t *testing.T) {
		mockCellSelector := func(deploymentType dmodel.DeploymentType, cells []api.Cell, dProfiles []api.DataspaceProfile) (*api.Cell, error) {
			return &api.Cell{
				DeployableEntity: api.DeployableEntity{
					Entity: api.Entity{
						ID:      "cell-123",
						Version: 1,
					},
					State:          api.DeploymentStateActive,
					StateTimestamp: now,
				},
				Properties: make(api.Properties),
			}, nil
		}

		generator := participantGenerator{
			CellSelector: mockCellSelector,
		}

		identifier := "participant-abc"
		deploymentProperties := api.Properties{
			"department": "IT",
			"region":     "US-West",
		}

		extensionProperties := api.Properties{
			"test": "value",
		}

		cells := []api.Cell{
			{
				DeployableEntity: api.DeployableEntity{
					Entity: api.Entity{
						ID:      "cell-123",
						Version: 1,
					},
					State:          api.DeploymentStateActive,
					StateTimestamp: now,
				},
			},
		}

		dProfiles := []api.DataspaceProfile{
			{
				Entity: api.Entity{
					ID:      "profile-456",
					Version: 1,
				},
				DeploymentProperties: make(api.Properties),
			},
		}

		profile, err := generator.Generate(identifier, deploymentProperties, extensionProperties, cells, dProfiles)

		require.NoError(t, err)
		require.NotNil(t, profile)

		// Validate basic profile structure
		assert.NotEmpty(t, profile.ID)
		_, err = uuid.Parse(profile.ID)
		assert.NoError(t, err, "ID should be a valid UUID")
		assert.Equal(t, int64(0), profile.Version)
		assert.Equal(t, identifier, profile.Identifier)
		assert.Equal(t, deploymentProperties, profile.DeploymentProperties)
		assert.Equal(t, extensionProperties, profile.ExtensionProperties)
		assert.Equal(t, dProfiles, profile.DataSpaceProfiles)

		// Validate VPAs
		assert.Len(t, profile.VPAs, 1)
		vpa := profile.VPAs[0]
		assert.NotEmpty(t, vpa.ID)
		_, err = uuid.Parse(vpa.ID)
		assert.NoError(t, err, "VPA ID should be a valid UUID")
		assert.Equal(t, int64(0), vpa.Version)
		assert.Equal(t, dmodel.ConnectorType, vpa.Type)
		assert.Equal(t, api.DeploymentStateActive, vpa.State)
		assert.Equal(t, "cell-123", vpa.Cell.ID)
		assert.NotNil(t, vpa.Properties)
		assert.NotNil(t, vpa.StateTimestamp)
	})

	t.Run("error when cell selector fails", func(t *testing.T) {
		mockCellSelector := func(deploymentType dmodel.DeploymentType, cells []api.Cell, dProfiles []api.DataspaceProfile) (*api.Cell, error) {
			return nil, assert.AnError
		}

		generator := participantGenerator{
			CellSelector: mockCellSelector,
		}

		profile, err := generator.Generate(
			"test-participant",
			make(map[string]any),
			make(map[string]any),
			[]api.Cell{},
			[]api.DataspaceProfile{})

		require.Error(t, err)
		require.Nil(t, profile)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("cell selector receives correct deployment type", func(t *testing.T) {
		var receivedDeploymentType dmodel.DeploymentType
		mockCellSelector := func(deploymentType dmodel.DeploymentType, cells []api.Cell, dProfiles []api.DataspaceProfile) (*api.Cell, error) {
			receivedDeploymentType = deploymentType
			return &api.Cell{
				DeployableEntity: api.DeployableEntity{
					Entity: api.Entity{
						ID:      "cell-123",
						Version: 1,
					},
					State:          api.DeploymentStateActive,
					StateTimestamp: now,
				},
				Properties: make(api.Properties),
			}, nil
		}

		generator := participantGenerator{
			CellSelector: mockCellSelector,
		}

		_, err := generator.Generate(
			"test",
			make(map[string]any),
			make(map[string]any),
			[]api.Cell{},
			[]api.DataspaceProfile{})

		require.NoError(t, err)
		assert.Equal(t, dmodel.VpaDeploymentType, receivedDeploymentType)
	})

	t.Run("cell selector receives correct parameters", func(t *testing.T) {
		var receivedCells []api.Cell
		var receivedProfiles []api.DataspaceProfile

		mockCellSelector := func(deploymentType dmodel.DeploymentType, cells []api.Cell, dProfiles []api.DataspaceProfile) (*api.Cell, error) {
			receivedCells = cells
			receivedProfiles = dProfiles
			return &api.Cell{
				DeployableEntity: api.DeployableEntity{
					Entity: api.Entity{
						ID:      "cell-123",
						Version: 1,
					},
					State:          api.DeploymentStateActive,
					StateTimestamp: now,
				},
				Properties: make(api.Properties),
			}, nil
		}

		generator := participantGenerator{
			CellSelector: mockCellSelector,
		}

		inputCells := []api.Cell{
			{
				DeployableEntity: api.DeployableEntity{
					Entity: api.Entity{ID: "cell-1"},
				},
			},
		}
		inputProfiles := []api.DataspaceProfile{
			{
				Entity: api.Entity{ID: "profile-1"},
			},
		}

		_, err := generator.Generate(
			"test",
			make(map[string]any),
			make(map[string]any),
			inputCells,
			inputProfiles)

		require.NoError(t, err)
		assert.Equal(t, inputCells, receivedCells)
		assert.Equal(t, inputProfiles, receivedProfiles)
	})

	t.Run("multiple dataspace profiles are correctly assigned", func(t *testing.T) {
		mockCellSelector := func(deploymentType dmodel.DeploymentType, cells []api.Cell, dProfiles []api.DataspaceProfile) (*api.Cell, error) {
			return &api.Cell{
				DeployableEntity: api.DeployableEntity{
					Entity: api.Entity{
						ID:      "cell-123",
						Version: 1,
					},
					State:          api.DeploymentStateActive,
					StateTimestamp: now,
				},
				Properties: make(api.Properties),
			}, nil
		}

		generator := participantGenerator{
			CellSelector: mockCellSelector,
		}

		dProfiles := []api.DataspaceProfile{
			{
				Entity: api.Entity{
					ID:      "profile-1",
					Version: 1,
				},
			},
			{
				Entity: api.Entity{
					ID:      "profile-2",
					Version: 2,
				},
			},
			{
				Entity: api.Entity{
					ID:      "profile-3",
					Version: 1,
				},
			},
		}

		profile, err := generator.Generate(
			"multi-profile-test",
			make(map[string]any),
			make(map[string]any),
			[]api.Cell{},
			dProfiles)

		require.NoError(t, err)
		require.NotNil(t, profile)
		assert.Len(t, profile.DataSpaceProfiles, 3)
		assert.Equal(t, dProfiles, profile.DataSpaceProfiles)
	})

}

func TestParticipantProfileGenerator_generateConnector(t *testing.T) {
	now := time.Now().UTC()

	t.Run("generates connector", func(t *testing.T) {
		generator := participantGenerator{}

		cellProperties := api.Properties{
			"environment": "production",
			"region":      "eu-west-1",
			"capacity":    500,
			"metadata": map[string]any{
				"owner": "platform-team",
				"tags":  []string{"critical", "production"},
			},
		}

		inputCell := &api.Cell{
			DeployableEntity: api.DeployableEntity{
				Entity: api.Entity{
					ID:      "prop-test-cell",
					Version: 3,
				},
				State:          api.DeploymentStateActive,
				StateTimestamp: now,
			},
			Properties: cellProperties,
		}

		connector := generator.generateConnector(inputCell)

		assert.Equal(t, cellProperties, connector.Cell.Properties)
		assert.NotSame(t, &cellProperties, &connector.Cell.Properties, "Properties should be a copy, not the same reference")
	})

	t.Run("generates unique connector IDs", func(t *testing.T) {
		generator := participantGenerator{}

		inputCell := &api.Cell{
			DeployableEntity: api.DeployableEntity{
				Entity: api.Entity{
					ID:      "test-cell",
					Version: 1,
				},
				State:          api.DeploymentStateActive,
				StateTimestamp: now,
			},
			Properties: make(api.Properties),
		}

		connector1 := generator.generateConnector(inputCell)
		connector2 := generator.generateConnector(inputCell)
		connector3 := generator.generateConnector(inputCell)

		ids := map[string]bool{
			connector1.ID: true,
			connector2.ID: true,
			connector3.ID: true,
		}
		assert.Len(t, ids, 3, "All connector IDs should be unique")
	})

}
