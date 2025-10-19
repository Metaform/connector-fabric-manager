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

package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeploymentManifest_JSONRoundTrip(t *testing.T) {
	originalManifest := DeploymentManifest{
		ID:             "manifest-1",
		DeploymentType: VpaDeploymentType,
		Payload: map[string]any{
			"environment": "production",
			"capacity":    100,
			"enabled":     true,
		},
	}

	jsonData, err := json.Marshal(originalManifest)
	require.NoError(t, err)

	var unmarshaledManifest DeploymentManifest
	err = json.Unmarshal(jsonData, &unmarshaledManifest)
	require.NoError(t, err)

	assert.Equal(t, originalManifest.ID, unmarshaledManifest.ID)
	assert.Equal(t, originalManifest.DeploymentType, unmarshaledManifest.DeploymentType)

	require.NotNil(t, unmarshaledManifest.Payload)
	comparePayload(t, originalManifest.Payload, unmarshaledManifest.Payload)
}

func TestDeploymentResponse_JSONRoundTrip(t *testing.T) {
	originalResponse := DeploymentResponse{
		ID:             "response-1",
		Success:        true,
		ErrorDetail:    "",
		ManifestID:     "manifest-1",
		DeploymentType: VpaDeploymentType,
		Properties: map[string]any{
			"endpoint": "https://example.com",
			"status":   "running",
			"port":     8080,
		},
	}

	jsonData, err := json.Marshal(originalResponse)
	require.NoError(t, err)

	var unmarshaledResponse DeploymentResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	require.NoError(t, err)

	assert.Equal(t, originalResponse.ID, unmarshaledResponse.ID)
	assert.Equal(t, originalResponse.Success, unmarshaledResponse.Success)
	assert.Equal(t, originalResponse.ErrorDetail, unmarshaledResponse.ErrorDetail)
	assert.Equal(t, originalResponse.ManifestID, unmarshaledResponse.ManifestID)
	assert.Equal(t, originalResponse.DeploymentType, unmarshaledResponse.DeploymentType)

	require.NotNil(t, unmarshaledResponse.Properties)
	comparePayload(t, originalResponse.Properties, unmarshaledResponse.Properties)
}

func TestVPAManifest_JSONRoundTrip(t *testing.T) {
	originalManifest := VPAManifest{
		ID:      "vpa-manifest-1",
		VPAType: ConnectorType,
		Cell:    "cell-west-1",
		Properties: map[string]any{
			"region":    "us-west-1",
			"replicas":  3,
			"enabled":   true,
			"endpoints": []string{"http://api1.example.com", "http://api2.example.com"},
		},
	}

	jsonData, err := json.Marshal(originalManifest)
	require.NoError(t, err)

	var unmarshaledManifest VPAManifest
	err = json.Unmarshal(jsonData, &unmarshaledManifest)
	require.NoError(t, err)

	assert.Equal(t, originalManifest.ID, unmarshaledManifest.ID)
	assert.Equal(t, originalManifest.VPAType, unmarshaledManifest.VPAType)
	assert.Equal(t, originalManifest.Cell, unmarshaledManifest.Cell)

	// Handle properties comparison (JSON unmarshaling converts numbers to float64)
	require.NotNil(t, unmarshaledManifest.Properties)
	comparePayload(t, originalManifest.Properties, unmarshaledManifest.Properties)
}

// comparePayload is a helper function to compare payloads/properties accounting for JSON type conversions
func comparePayload(t *testing.T, original, unmarshaled map[string]any) {
	for key, originalValue := range original {
		unmarshaledValue, exists := unmarshaled[key]
		require.True(t, exists, "Key %s should exist in unmarshaled payload", key)

		switch v := originalValue.(type) {
		case int:
			assert.Equal(t, float64(v), unmarshaledValue, "Value for key %s should match", key)
		case []string:
			unmarshaledSlice, ok := unmarshaledValue.([]any)
			require.True(t, ok, "Value for key %s should be a slice", key)
			require.Len(t, unmarshaledSlice, len(v))
			for i, str := range v {
				assert.Equal(t, str, unmarshaledSlice[i])
			}
		case map[string]any:
			// Recursively compare nested maps
			unmarshaledMap, ok := unmarshaledValue.(map[string]any)
			require.True(t, ok, "Value for key %s should be a map", key)
			comparePayload(t, v, unmarshaledMap)
		default:
			assert.Equal(t, originalValue, unmarshaledValue, "Value for key %s should match", key)
		}
	}
}
