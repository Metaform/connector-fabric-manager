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
	"testing"
	"time"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToParticipantProfile(t *testing.T) {
	testTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.Local)

	input := &api.ParticipantProfile{
		Entity: api.Entity{
			ID:      "participant-123",
			Version: 1,
		},
		Identifier: "test-participant",
		VPAs: []api.VirtualParticipantAgent{
			{
				DeployableEntity: api.DeployableEntity{
					Entity: api.Entity{
						ID:      "vpa-123",
						Version: 2,
					},
					State:          api.DeploymentStateActive,
					StateTimestamp: testTime,
				},
				Type: model.ConnectorType,
				Cell: api.Cell{
					DeployableEntity: api.DeployableEntity{
						Entity: api.Entity{
							ID:      "cell-123",
							Version: 1,
						},
						State:          api.DeploymentStateActive,
						StateTimestamp: testTime,
					},
					Properties: api.Properties{"cell-key": "cell-value"},
				},
				Properties: api.Properties{"vpa-key": "vpa-value"},
			},
		},
		Properties: api.Properties{"participant-key": "participant-value"},
	}

	result := ToParticipantProfile(input)

	require.NotNil(t, result)
	assert.Equal(t, "participant-123", result.ID)
	assert.Equal(t, int64(1), result.Version)
	assert.Equal(t, "test-participant", result.Identifier)
	assert.Len(t, result.VPAs, 1)
	assert.Equal(t, map[string]any{"participant-key": "participant-value"}, result.Properties)
}

func TestToVPACollection(t *testing.T) {
	testTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	input := &api.ParticipantProfile{
		VPAs: []api.VirtualParticipantAgent{
			{
				DeployableEntity: api.DeployableEntity{
					Entity: api.Entity{
						ID:      "vpa-1",
						Version: 1,
					},
					State:          api.DeploymentStateActive,
					StateTimestamp: testTime,
				},
				Type: model.ConnectorType,
				Cell: api.Cell{
					DeployableEntity: api.DeployableEntity{
						Entity: api.Entity{
							ID:      "cell-1",
							Version: 1,
						},
						State:          api.DeploymentStateActive,
						StateTimestamp: testTime,
					},
					Properties: api.Properties{},
				},
				Properties: api.Properties{"key1": "value1"},
			},
			{
				DeployableEntity: api.DeployableEntity{
					Entity: api.Entity{
						ID:      "vpa-2",
						Version: 2,
					},
					State:          api.DeploymentStatePending,
					StateTimestamp: testTime,
				},
				Type: model.CredentialServiceType,
				Cell: api.Cell{
					DeployableEntity: api.DeployableEntity{
						Entity: api.Entity{
							ID:      "cell-2",
							Version: 1,
						},
						State:          api.DeploymentStatePending,
						StateTimestamp: testTime,
					},
					Properties: api.Properties{},
				},
				Properties: api.Properties{"key2": "value2"},
			},
		},
	}

	result := ToVPACollection(input)

	require.Len(t, result, 2)

	// First VPA
	assert.Equal(t, "vpa-1", result[0].ID)
	assert.Equal(t, int64(1), result[0].Version)
	assert.Equal(t, "active", result[0].State)
	assert.Equal(t, testTime, result[0].StateTimestamp)
	assert.Equal(t, model.ConnectorType, result[0].Type)
	assert.Equal(t, "cell-1", result[0].Cell.ID)
	assert.Equal(t, map[string]any{"key1": "value1"}, result[0].Properties)

	// Second VPA
	assert.Equal(t, "vpa-2", result[1].ID)
	assert.Equal(t, int64(2), result[1].Version)
	assert.Equal(t, "pending", result[1].State)
	assert.Equal(t, testTime, result[1].StateTimestamp)
	assert.Equal(t, model.CredentialServiceType, result[1].Type)
	assert.Equal(t, "cell-2", result[1].Cell.ID)
	assert.Equal(t, map[string]any{"key2": "value2"}, result[1].Properties)
}

func TestToVPA(t *testing.T) {
	testTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.FixedZone("EST", -5*60*60))

	input := api.VirtualParticipantAgent{
		DeployableEntity: api.DeployableEntity{
			Entity: api.Entity{
				ID:      "vpa-456",
				Version: 3,
			},
			State:          api.DeploymentStateError,
			StateTimestamp: testTime,
		},
		Type: model.DataPlaneType,
		Cell: api.Cell{
			DeployableEntity: api.DeployableEntity{
				Entity: api.Entity{
					ID:      "cell-456",
					Version: 2,
				},
				State:          api.DeploymentStateOffline,
				StateTimestamp: testTime,
			},
			Properties: api.Properties{"cell-prop": "cell-val"},
		},
		Properties: api.Properties{"vpa-prop": "vpa-val"},
	}

	result := ToVPA(input)

	require.NotNil(t, result)
	assert.Equal(t, "vpa-456", result.ID)
	assert.Equal(t, int64(3), result.Version)
	assert.Equal(t, "error", result.State)
	assert.Equal(t, testTime, result.StateTimestamp)
	assert.Equal(t, model.DataPlaneType, result.Type)
	assert.Equal(t, "cell-456", result.Cell.ID)
	assert.Equal(t, int64(2), result.Cell.Version)
	assert.Equal(t, map[string]any{"vpa-prop": "vpa-val"}, result.Properties)
}

func TestToCell(t *testing.T) {
	testTime := time.Date(2025, 6, 15, 10, 30, 45, 123456789, time.UTC)

	input := api.Cell{
		DeployableEntity: api.DeployableEntity{
			Entity: api.Entity{
				ID:      "cell-789",
				Version: 5,
			},
			State:          api.DeploymentStateLocked,
			StateTimestamp: testTime,
		},
		Properties: api.Properties{
			"environment": "production",
			"region":      "us-west-2",
			"capacity":    100,
		},
	}

	result := ToCell(input)

	require.NotNil(t, result)
	assert.Equal(t, "cell-789", result.ID)
	assert.Equal(t, int64(5), result.Version)
	assert.Equal(t, "locked", result.State)
	assert.Equal(t, testTime, result.StateTimestamp)
	assert.Equal(t, map[string]any{
		"environment": "production",
		"region":      "us-west-2",
		"capacity":    100,
	}, result.Properties)
}

func TestToAPIParticipantProfile(t *testing.T) {
	testTime := time.Date(2025, 3, 10, 14, 25, 0, 0, time.Local)

	input := &ParticipantProfile{
		Entity: Entity{
			ID:      "api-participant-123",
			Version: 4,
		},
		Identifier: "api-test-participant",
		VPAs: []VirtualParticipantAgent{
			{
				DeployableEntity: DeployableEntity{
					Entity: Entity{
						ID:      "api-vpa-123",
						Version: 1,
					},
					State:          "active",
					StateTimestamp: testTime,
				},
				Type: model.ConnectorType,
				Cell: Cell{
					Entity: Entity{
						ID:      "api-cell-123",
						Version: 1,
					},
					NewCell: NewCell{
						State:          "active",
						StateTimestamp: testTime,
						Properties:     map[string]any{"cell-key": "cell-value"},
					},
				},
				Properties: map[string]any{"vpa-key": "vpa-value"},
			},
		},
		Properties: map[string]any{"participant-key": "participant-value"},
	}

	result := ToAPIParticipantProfile(input)

	require.NotNil(t, result)
	assert.Equal(t, "api-participant-123", result.ID)
	assert.Equal(t, int64(4), result.Version)
	assert.Equal(t, "api-test-participant", result.Identifier)
	assert.Len(t, result.VPAs, 1)
	assert.Contains(t, result.Properties, "participant-key")
	assert.Equal(t, "participant-value", result.Properties["participant-key"])
}

func TestToAPIVPACollection(t *testing.T) {
	testTime := time.Date(2025, 2, 14, 8, 15, 30, 0, time.UTC)

	input := []VirtualParticipantAgent{
		{
			DeployableEntity: DeployableEntity{
				Entity: Entity{
					ID:      "vpa-collection-1",
					Version: 1,
				},
				State:          "pending",
				StateTimestamp: testTime,
			},
			Type: model.ConnectorType,
			Cell: Cell{
				Entity: Entity{
					ID:      "cell-collection-1",
					Version: 1,
				},
				NewCell: NewCell{
					State:          "pending",
					StateTimestamp: testTime,
					Properties:     map[string]any{},
				},
			},
			Properties: map[string]any{"prop1": "val1"},
		},
		{
			DeployableEntity: DeployableEntity{
				Entity: Entity{
					ID:      "vpa-collection-2",
					Version: 2,
				},
				State:          "offline",
				StateTimestamp: testTime,
			},
			Type: model.DataPlaneType,
			Cell: Cell{
				Entity: Entity{
					ID:      "cell-collection-2",
					Version: 1,
				},
				NewCell: NewCell{
					State:          "offline",
					StateTimestamp: testTime,
					Properties:     map[string]any{},
				},
			},
			Properties: map[string]any{"prop2": "val2"},
		},
	}

	result := ToAPIVPACollection(input)

	require.Len(t, result, 2)

	// First VPA
	assert.Equal(t, "vpa-collection-1", result[0].ID)
	assert.Equal(t, int64(1), result[0].Version)
	assert.Equal(t, api.DeploymentStatePending, result[0].State)
	assert.Equal(t, testTime.UTC(), result[0].StateTimestamp) // Should be converted to UTC
	assert.Equal(t, model.ConnectorType, result[0].Type)

	// Second VPA
	assert.Equal(t, "vpa-collection-2", result[1].ID)
	assert.Equal(t, int64(2), result[1].Version)
	assert.Equal(t, api.DeploymentStateOffline, result[1].State)
	assert.Equal(t, testTime.UTC(), result[1].StateTimestamp) // Should be converted to UTC
	assert.Equal(t, model.DataPlaneType, result[1].Type)
}

func TestToAPIVPA(t *testing.T) {
	// Test with non-UTC timezone to verify UTC conversion
	testTime := time.Date(2025, 4, 20, 16, 45, 0, 0, time.FixedZone("JST", 9*60*60))

	input := VirtualParticipantAgent{
		DeployableEntity: DeployableEntity{
			Entity: Entity{
				ID:      "api-vpa-456",
				Version: 6,
			},
			State:          "locked",
			StateTimestamp: testTime,
		},
		Type: model.CredentialServiceType,
		Cell: Cell{
			Entity: Entity{
				ID:      "api-cell-456",
				Version: 3,
			},
			NewCell: NewCell{
				State:          "locked",
				StateTimestamp: testTime,
				Properties:     map[string]any{"cell-env": "staging"},
			},
		},
		Properties: map[string]any{"vpa-env": "staging"},
	}

	result := ToAPIVPA(input)

	require.NotNil(t, result)
	assert.Equal(t, "api-vpa-456", result.ID)
	assert.Equal(t, int64(6), result.Version)
	assert.Equal(t, api.DeploymentStateLocked, result.State)
	assert.Equal(t, testTime.UTC(), result.StateTimestamp)      // Should be converted to UTC
	assert.Equal(t, time.UTC, result.StateTimestamp.Location()) // Verify timezone is UTC
	assert.Equal(t, model.CredentialServiceType, result.Type)
	assert.Equal(t, "api-cell-456", result.Cell.ID)
	assert.Contains(t, result.Properties, "vpa-env")
}

func TestToAPICell(t *testing.T) {
	// Test with different timezone to verify UTC conversion
	testTime := time.Date(2025, 7, 4, 9, 0, 0, 0, time.FixedZone("PST", -8*60*60))

	input := Cell{
		Entity: Entity{
			ID:      "api-cell-789",
			Version: 7,
		},
		NewCell: NewCell{
			State:          "error",
			StateTimestamp: testTime,
			Properties: map[string]any{
				"cluster":   "prod-cluster-1",
				"nodes":     5,
				"cpu_cores": 32,
				"memory_gb": 128,
			},
		},
	}

	result := ToAPICell(input)

	require.NotNil(t, result)
	assert.Equal(t, "api-cell-789", result.ID)
	assert.Equal(t, int64(7), result.Version)
	assert.Equal(t, api.DeploymentStateError, result.State)
	assert.Equal(t, testTime.UTC(), result.StateTimestamp)      // Should be converted to UTC
	assert.Equal(t, time.UTC, result.StateTimestamp.Location()) // Verify timezone is UTC
	assert.Contains(t, result.Properties, "cluster")
	assert.Equal(t, "prod-cluster-1", result.Properties["cluster"])
	assert.Contains(t, result.Properties, "nodes")
	assert.Equal(t, 5, result.Properties["nodes"])
}

func TestNewAPICell(t *testing.T) {
	// Test with non-UTC timezone to verify UTC conversion
	testTime := time.Date(2025, 8, 12, 20, 30, 45, 123000000, time.FixedZone("CET", 1*60*60))

	input := NewCell{
		State:          "initial",
		StateTimestamp: testTime,
		Properties: map[string]any{
			"type":       "kubernetes",
			"version":    "1.28",
			"auto_scale": true,
		},
	}

	result := NewAPICell(input)

	require.NotNil(t, result)
	assert.NotEmpty(t, result.ID)             // Should generate a new UUID
	assert.Equal(t, int64(0), result.Version) // Should be 0 for new cells
	assert.Equal(t, api.DeploymentStateInitial, result.State)
	assert.Equal(t, testTime.UTC(), result.StateTimestamp)      // Should be converted to UTC
	assert.Equal(t, time.UTC, result.StateTimestamp.Location()) // Verify timezone is UTC
	assert.Contains(t, result.Properties, "type")
	assert.Equal(t, "kubernetes", result.Properties["type"])
	assert.Contains(t, result.Properties, "auto_scale")
	assert.Equal(t, true, result.Properties["auto_scale"])
}

func TestTimestampUTCConversion(t *testing.T) {
	// Test various timezones to ensure they all convert to UTC properly
	timezones := []struct {
		name string
		zone *time.Location
	}{
		{"EST", time.FixedZone("EST", -5*60*60)},
		{"PST", time.FixedZone("PST", -8*60*60)},
		{"JST", time.FixedZone("JST", 9*60*60)},
		{"CET", time.FixedZone("CET", 1*60*60)},
		{"Local", time.Local},
		{"UTC", time.UTC},
	}

	baseTime := time.Date(2025, 5, 15, 12, 30, 45, 123456789, time.UTC)

	for _, tz := range timezones {
		t.Run(tz.name, func(t *testing.T) {
			testTime := baseTime.In(tz.zone)

			// Test ToAPIVPA
			vpaInput := VirtualParticipantAgent{
				DeployableEntity: DeployableEntity{
					Entity:         Entity{ID: "vpa", Version: 1},
					State:          "active",
					StateTimestamp: testTime,
				},
				Type: model.ConnectorType,
				Cell: Cell{
					Entity: Entity{ID: "cell", Version: 1},
					NewCell: NewCell{
						State:          "active",
						StateTimestamp: testTime,
						Properties:     map[string]any{},
					},
				},
				Properties: map[string]any{},
			}

			vpaResult := ToAPIVPA(vpaInput)
			assert.Equal(t, time.UTC, vpaResult.StateTimestamp.Location())
			assert.Equal(t, baseTime.UTC(), vpaResult.StateTimestamp)

			// Test ToAPICell
			cellInput := Cell{
				Entity: Entity{ID: "cell", Version: 1},
				NewCell: NewCell{
					State:          "active",
					StateTimestamp: testTime,
					Properties:     map[string]any{},
				},
			}

			cellResult := ToAPICell(cellInput)
			assert.Equal(t, time.UTC, cellResult.StateTimestamp.Location())
			assert.Equal(t, baseTime.UTC(), cellResult.StateTimestamp)

			// Test NewAPICell
			newCellInput := NewCell{
				State:          "active",
				StateTimestamp: testTime,
				Properties:     map[string]any{},
			}

			newCellResult := NewAPICell(newCellInput)
			assert.Equal(t, time.UTC, newCellResult.StateTimestamp.Location())
			assert.Equal(t, baseTime.UTC(), newCellResult.StateTimestamp)
		})
	}
}

func TestEmptyAndNilInputs(t *testing.T) {
	t.Run("ToParticipantProfile with nil", func(t *testing.T) {
		// This would panic, so we test with minimal valid input
		input := &api.ParticipantProfile{}
		result := ToParticipantProfile(input)
		require.NotNil(t, result)
		assert.Empty(t, result.ID)
		assert.Equal(t, int64(0), result.Version)
	})

	t.Run("ToVPACollection with empty VPAs", func(t *testing.T) {
		input := &api.ParticipantProfile{VPAs: []api.VirtualParticipantAgent{}}
		result := ToVPACollection(input)
		assert.Len(t, result, 0)
	})

	t.Run("ToAPIVPACollection with empty slice", func(t *testing.T) {
		result := ToAPIVPACollection([]VirtualParticipantAgent{})
		assert.Len(t, result, 0)
	})

	t.Run("Properties handling", func(t *testing.T) {
		// Test nil properties
		input := &api.ParticipantProfile{
			Properties: nil,
		}
		result := ToParticipantProfile(input)
		assert.Nil(t, result.Properties)

		// Test empty properties
		input.Properties = api.Properties{}
		result = ToParticipantProfile(input)
		assert.Empty(t, result.Properties)
	})
}
