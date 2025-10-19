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

package tmstore

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryTManagerStore_GetCells(t *testing.T) {
	t.Run("returns empty list when no cells", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)

		cells, err := store.GetCells()

		require.NoError(t, err)
		assert.Empty(t, cells)
	})

	t.Run("returns all stored cells", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)
		now := time.Now()

		// Add test cells
		cell1 := api.Cell{
			DeployableEntity: api.DeployableEntity{
				Entity: api.Entity{
					ID:      "cell-1",
					Version: 1,
				},
				State:          api.DeploymentStateActive,
				StateTimestamp: now,
			},
			Properties: api.Properties{"type": "k8s"},
		}

		cell2 := api.Cell{
			DeployableEntity: api.DeployableEntity{
				Entity: api.Entity{
					ID:      "cell-2",
					Version: 1,
				},
				State:          api.DeploymentStatePending,
				StateTimestamp: now,
			},
			Properties: api.Properties{"type": "docker"},
		}

		err := store.cellStorage.Create(cell1)
		require.NoError(t, err)
		err = store.cellStorage.Create(cell2)
		require.NoError(t, err)

		cells, err := store.GetCells()

		require.NoError(t, err)
		assert.Len(t, cells, 2)

		// Verify both cells are returned
		assert.ElementsMatch(t, []string{"cell-1", "cell-2"}, []string{cells[0].ID, cells[1].ID})
		assert.Equal(t, "k8s", cells[0].Properties["type"])
		assert.Equal(t, "docker", cells[1].Properties["type"])
	})
}

func TestInMemoryTManagerStore_GetDataspaceProfiles(t *testing.T) {
	t.Run("returns empty list when no profiles", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)

		profiles, err := store.GetDataspaceProfiles()

		require.NoError(t, err)
		assert.Empty(t, profiles)
	})

	t.Run("returns all stored profiles", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)

		// Add test profiles
		profile1 := api.DataspaceProfile{
			Entity: api.Entity{
				ID:      "profile-1",
				Version: 1,
			},
			Artifacts:  []string{"artifact-1"},
			Properties: api.Properties{"env": "prod"},
		}

		profile2 := api.DataspaceProfile{
			Entity: api.Entity{
				ID:      "profile-2",
				Version: 2,
			},
			Artifacts:  []string{"artifact-2", "artifact-3"},
			Properties: api.Properties{"env": "dev"},
		}

		err := store.dProfileStorage.Create(profile1)
		require.NoError(t, err)
		err = store.dProfileStorage.Create(profile2)
		require.NoError(t, err)

		profiles, err := store.GetDataspaceProfiles()

		require.NoError(t, err)
		assert.Len(t, profiles, 2)

		// Verify profiles are returned (order may vary)
		profileIDs := make(map[string]api.DataspaceProfile)
		for _, profile := range profiles {
			profileIDs[profile.ID] = profile
		}

		assert.Contains(t, profileIDs, "profile-1")
		assert.Contains(t, profileIDs, "profile-2")
		assert.Equal(t, "prod", profileIDs["profile-1"].Properties["env"])
		assert.Equal(t, "dev", profileIDs["profile-2"].Properties["env"])
		assert.Len(t, profileIDs["profile-1"].Artifacts, 1)
		assert.Len(t, profileIDs["profile-2"].Artifacts, 2)
	})
}

func TestInMemoryTManagerStore_FindDeployment(t *testing.T) {
	t.Run("returns error when deployment not found", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)

		record, err := store.FindDeployment("non-existent-id")

		require.Error(t, err)
		assert.Nil(t, record)
		assert.Equal(t, model.ErrNotFound, err)
	})

	t.Run("returns deployment when found", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)
		now := time.Now()

		// Create test deployment record
		testRecord := api.DeploymentRecord{
			ID:            "deployment-123",
			CorrelationID: "corr-456",
			State:         api.ProcessingStateRunning,
			Timestamp:     now,
			TenantID:      "tenant-789",
			ManifestID:    "manifest-abc",
			Success:       false,
		}

		err := store.deploymentStorage.Create(testRecord)
		require.NoError(t, err)

		record, err := store.FindDeployment("deployment-123")

		require.NoError(t, err)
		require.NotNil(t, record)
		assert.Equal(t, "deployment-123", record.ID)
		assert.Equal(t, "corr-456", record.CorrelationID)
		assert.Equal(t, api.ProcessingStateRunning, record.State)
		assert.Equal(t, "tenant-789", record.TenantID)
		assert.Equal(t, "manifest-abc", record.ManifestID)
		assert.False(t, record.Success)
	})

	t.Run("returns error for empty ID", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)

		record, err := store.FindDeployment("")

		require.Error(t, err)
		assert.Nil(t, record)
		assert.Equal(t, model.ErrNotFound, err)
	})
}

func TestInMemoryTManagerStore_DeploymentExists(t *testing.T) {
	t.Run("returns false when deployment does not exist", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)

		exists, err := store.DeploymentExists("non-existent-id")

		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("returns true when deployment exists", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)

		testRecord := api.DeploymentRecord{
			ID:        "deployment-123",
			TenantID:  "tenant-789",
			Timestamp: time.Now(),
		}

		err := store.deploymentStorage.Create(testRecord)
		require.NoError(t, err)

		exists, err := store.DeploymentExists("deployment-123")

		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("returns false for empty ID", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)

		exists, err := store.DeploymentExists("")

		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestInMemoryTManagerStore_CreateDeployment(t *testing.T) {
	t.Run("creates deployment with generated ID", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)
		now := time.Now()

		inputRecord := api.DeploymentRecord{
			CorrelationID: "corr-123",
			State:         api.ProcessingStateInitialized,
			Timestamp:     now,
			TenantID:      "tenant-456",
			ManifestID:    "manifest-789",
			Success:       false,
		}

		createdRecord, err := store.CreateDeployment(inputRecord)

		require.NoError(t, err)
		require.NotNil(t, createdRecord)

		// Verify ID was generated
		assert.NotEmpty(t, createdRecord.ID)
		_, err = uuid.Parse(createdRecord.ID)
		assert.NoError(t, err, "ID should be a valid UUID")

		// Verify other fields are preserved
		assert.Equal(t, "corr-123", createdRecord.CorrelationID)
		assert.Equal(t, api.ProcessingStateInitialized, createdRecord.State)
		assert.Equal(t, now, createdRecord.Timestamp)
		assert.Equal(t, "tenant-456", createdRecord.TenantID)
		assert.Equal(t, "manifest-789", createdRecord.ManifestID)
		assert.False(t, createdRecord.Success)

		// Verify record is stored
		storedRecord, err := store.FindDeployment(createdRecord.ID)
		require.NoError(t, err)
		assert.Equal(t, createdRecord, storedRecord)
	})

	t.Run("creates multiple deployments with unique IDs", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)

		record1, err := store.CreateDeployment(api.DeploymentRecord{
			CorrelationID: "corr-1",
			TenantID:      "tenant-1",
		})
		require.NoError(t, err)

		record2, err := store.CreateDeployment(api.DeploymentRecord{
			CorrelationID: "corr-2",
			TenantID:      "tenant-2",
		})
		require.NoError(t, err)

		assert.NotEqual(t, record1.ID, record2.ID)
		assert.NotEmpty(t, record1.ID)
		assert.NotEmpty(t, record2.ID)
	})

	t.Run("preserves complex fields", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)

		responseProps := map[string]any{
			"key1": "value1",
			"key2": 42,
			"nested": map[string]any{
				"inner": "value",
			},
		}

		inputRecord := api.DeploymentRecord{
			CorrelationID:      "corr-complex",
			TenantID:           "tenant-complex",
			ErrorDetail:        "some error message",
			ResponseProperties: responseProps,
		}

		createdRecord, err := store.CreateDeployment(inputRecord)

		require.NoError(t, err)
		assert.Equal(t, "some error message", createdRecord.ErrorDetail)
		assert.Equal(t, responseProps, createdRecord.ResponseProperties)
	})
}

func TestInMemoryTManagerStore_UpdateDeployment(t *testing.T) {
	t.Run("updates existing deployment", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)
		originalTime := time.Now()

		// Create initial deployment
		initialRecord := api.DeploymentRecord{
			ID:        "deployment-update-test",
			TenantID:  "tenant-123",
			State:     api.ProcessingStateInitialized,
			Timestamp: originalTime,
			Success:   false,
		}

		err := store.deploymentStorage.Create(initialRecord)
		require.NoError(t, err)

		// Update the deployment
		updatedTime := originalTime.Add(time.Hour)
		updateRecord := api.DeploymentRecord{
			ID:          "deployment-update-test",
			TenantID:    "tenant-123",
			State:       api.ProcessingStateCompleted,
			Timestamp:   updatedTime,
			Success:     true,
			ErrorDetail: "",
		}

		err = store.UpdateDeployment(updateRecord)
		require.NoError(t, err)

		// Verify updates were applied
		storedRecord, err := store.FindDeployment("deployment-update-test")
		require.NoError(t, err)
		assert.Equal(t, api.ProcessingStateCompleted, storedRecord.State)
		assert.Equal(t, updatedTime, storedRecord.Timestamp)
		assert.True(t, storedRecord.Success)
		assert.Empty(t, storedRecord.ErrorDetail)
	})

	t.Run("returns error when deployment not found", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)

		updateRecord := api.DeploymentRecord{
			ID:       "non-existent-deployment",
			TenantID: "tenant-123",
		}

		err := store.UpdateDeployment(updateRecord)

		require.Error(t, err)
		assert.Equal(t, model.ErrNotFound, err)
	})

	t.Run("returns error for empty ID", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)

		updateRecord := api.DeploymentRecord{
			ID:       "",
			TenantID: "tenant-123",
		}

		err := store.UpdateDeployment(updateRecord)

		require.Error(t, err)
		assert.Equal(t, model.ErrInvalidInput, err)
	})

	t.Run("updates complex fields correctly", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)

		initialRecord := api.DeploymentRecord{
			ID:                 "deployment-complex-update",
			TenantID:           "tenant-123",
			ResponseProperties: map[string]any{"old": "value"},
		}

		err := store.deploymentStorage.Create(initialRecord)
		require.NoError(t, err)

		// Update with new complex data
		newResponseProps := map[string]any{
			"updated": "value",
			"nested": map[string]any{
				"data": []string{"item1", "item2"},
			},
		}

		updateRecord := api.DeploymentRecord{
			ID:                 "deployment-complex-update",
			TenantID:           "tenant-123",
			ResponseProperties: newResponseProps,
			ErrorDetail:        "updated error message",
		}

		err = store.UpdateDeployment(updateRecord)
		require.NoError(t, err)

		// Verify complex fields were updated
		storedRecord, err := store.FindDeployment("deployment-complex-update")
		require.NoError(t, err)
		assert.Equal(t, newResponseProps, storedRecord.ResponseProperties)
		assert.Equal(t, "updated error message", storedRecord.ErrorDetail)
	})
}

func TestInMemoryTManagerStore_ConcurrentAccess(t *testing.T) {
	t.Run("handles concurrent reads and writes", func(t *testing.T) {
		store := NewInMemoryTManagerStore(false)

		// This test verifies thread safety of the store
		initialRecord := api.DeploymentRecord{
			ID:       "concurrent-test",
			TenantID: "tenant-concurrent",
		}

		err := store.deploymentStorage.Create(initialRecord)
		require.NoError(t, err)

		// Perform concurrent operations
		done := make(chan bool, 2)

		// Concurrent reader
		go func() {
			defer func() { done <- true }()
			for i := 0; i < 100; i++ {
				_, err := store.FindDeployment("concurrent-test")
				assert.NoError(t, err)

				exists, err := store.DeploymentExists("concurrent-test")
				assert.NoError(t, err)
				assert.True(t, exists)
			}
		}()

		// Concurrent writer
		go func() {
			defer func() { done <- true }()
			for i := 0; i < 100; i++ {
				updateRecord := api.DeploymentRecord{
					ID:       "concurrent-test",
					TenantID: "tenant-concurrent",
					State:    api.ProcessingStateRunning,
				}
				err := store.UpdateDeployment(updateRecord)
				assert.NoError(t, err)
			}
		}()

		// Wait for both goroutines to complete
		<-done
		<-done

		// Verify final state
		finalRecord, err := store.FindDeployment("concurrent-test")
		require.NoError(t, err)
		assert.Equal(t, "concurrent-test", finalRecord.ID)
	})
}
