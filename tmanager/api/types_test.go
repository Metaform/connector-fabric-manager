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

package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeploymentState_String(t *testing.T) {
	tests := []struct {
		name  string
		state DeploymentState
		want  string
	}{
		{"Initial", DeploymentStateInitial, "initial"},
		{"Pending", DeploymentStatePending, "pending"},
		{"Active", DeploymentStateActive, "active"},
		{"Offline", DeploymentStateOffline, "offline"},
		{"Error", DeploymentStateError, "error"},
		{"Locked", DeploymentStateLocked, "locked"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.state.String())
		})
	}
}

func TestDeploymentState_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		state DeploymentState
		want  bool
	}{
		{"Initial is valid", DeploymentStateInitial, true},
		{"Pending is valid", DeploymentStatePending, true},
		{"Active is valid", DeploymentStateActive, true},
		{"Offline is valid", DeploymentStateOffline, true},
		{"Error is valid", DeploymentStateError, true},
		{"Locked is valid", DeploymentStateLocked, true},
		{"Empty string is invalid", DeploymentState(""), false},
		{"Invalid state is invalid", DeploymentState("invalid"), false},
		{"Case sensitive - uppercase is invalid", DeploymentState("ACTIVE"), false},
		{"Case sensitive - mixed case is invalid", DeploymentState("Active"), false},
		{"Whitespace state is invalid", DeploymentState(" active "), false},
		{"Similar but wrong state is invalid", DeploymentState("activ"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.state.IsValid())
		})
	}
}

func TestDeploymentState_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		state   DeploymentState
		want    string
		wantErr bool
	}{
		{"Initial marshals correctly", DeploymentStateInitial, `"initial"`, false},
		{"Pending marshals correctly", DeploymentStatePending, `"pending"`, false},
		{"Active marshals correctly", DeploymentStateActive, `"active"`, false},
		{"Offline marshals correctly", DeploymentStateOffline, `"offline"`, false},
		{"Error marshals correctly", DeploymentStateError, `"error"`, false},
		{"Locked marshals correctly", DeploymentStateLocked, `"locked"`, false},
		{"Invalid state marshals as string", DeploymentState("invalid"), `"invalid"`, false},
		{"Empty state marshals as empty string", DeploymentState(""), `""`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.state.MarshalJSON()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, string(got))
			}
		})
	}
}

func TestDeploymentState_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    DeploymentState
		wantErr bool
	}{
		{"Initial unmarshals correctly", `"initial"`, DeploymentStateInitial, false},
		{"Pending unmarshals correctly", `"pending"`, DeploymentStatePending, false},
		{"Active unmarshals correctly", `"active"`, DeploymentStateActive, false},
		{"Offline unmarshals correctly", `"offline"`, DeploymentStateOffline, false},
		{"Error unmarshals correctly", `"error"`, DeploymentStateError, false},
		{"Locked unmarshals correctly", `"locked"`, DeploymentStateLocked, false},
		{"Invalid state returns error", `"invalid"`, DeploymentState(""), true},
		{"Empty string returns zero", `""`, DeploymentState(""), false},
		{"Uppercase returns error", `"ACTIVE"`, DeploymentState(""), true},
		{"Mixed case returns error", `"Active"`, DeploymentState(""), true},
		{"Whitespace returns error", `" active "`, DeploymentState(""), true},
		{"Invalid JSON returns error", `invalid`, DeploymentState(""), true},
		{"Number returns error", `123`, DeploymentState(""), true},
		{"Boolean returns error", `true`, DeploymentState(""), true},
		{"Object returns error", `{"state":"active"}`, DeploymentState(""), true},
		{"Array returns error", `["active"]`, DeploymentState(""), true},
		{"Null returns zero", `null`, DeploymentState(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var state DeploymentState
			err := json.Unmarshal([]byte(tt.json), &state)
			if tt.wantErr {
				require.Error(t, err)
				if tt.json == `"invalid"` || tt.json == `""` || tt.json == `"ACTIVE"` ||
					tt.json == `"Active"` || tt.json == `" active "` {
					assert.Contains(t, err.Error(), "invalid deployment state")
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, state)
			}
		})
	}
}

func TestDeploymentState_Value(t *testing.T) {
	tests := []struct {
		name    string
		state   DeploymentState
		want    string
		wantErr bool
	}{
		{"Initial returns correct value", DeploymentStateInitial, "initial", false},
		{"Pending returns correct value", DeploymentStatePending, "pending", false},
		{"Active returns correct value", DeploymentStateActive, "active", false},
		{"Offline returns correct value", DeploymentStateOffline, "offline", false},
		{"Error returns correct value", DeploymentStateError, "error", false},
		{"Locked returns correct value", DeploymentStateLocked, "locked", false},
		{"Invalid state returns error", DeploymentState("invalid"), "", true},
		{"Empty state returns error", DeploymentState(""), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.state.Value()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid cell state")
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDeploymentState_Scan(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		want    DeploymentState
		wantErr bool
	}{
		{"String initial scans correctly", "initial", DeploymentStateInitial, false},
		{"String pending scans correctly", "pending", DeploymentStatePending, false},
		{"String active scans correctly", "active", DeploymentStateActive, false},
		{"String offline scans correctly", "offline", DeploymentStateOffline, false},
		{"String error scans correctly", "error", DeploymentStateError, false},
		{"String locked scans correctly", "locked", DeploymentStateLocked, false},
		{"Byte slice initial scans correctly", []byte("initial"), DeploymentStateInitial, false},
		{"Byte slice active scans correctly", []byte("active"), DeploymentStateActive, false},
		{"Nil value scans to empty", nil, DeploymentState(""), false},
		{"Invalid string returns error", "invalid", DeploymentState(""), true},
		{"Empty string returns error", "", DeploymentState(""), true},
		{"Uppercase string returns error", "ACTIVE", DeploymentState(""), true},
		{"Mixed case string returns error", "Active", DeploymentState(""), true},
		{"Whitespace string returns error", " active ", DeploymentState(""), true},
		{"Integer value returns error", 123, DeploymentState(""), true},
		{"Boolean value returns error", true, DeploymentState(""), true},
		{"Float value returns error", 3.14, DeploymentState(""), true},
		{"Struct value returns error", struct{}{}, DeploymentState(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var state DeploymentState
			err := state.Scan(tt.value)
			if tt.wantErr {
				require.Error(t, err)
				if tt.value != nil {
					switch v := tt.value.(type) {
					case string, []byte:
						if v != "" {
							assert.Contains(t, err.Error(), "invalid cell state")
						}
					default:
						assert.Contains(t, err.Error(), "cannot scan")
					}
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, state)
			}
		})
	}
}

func TestDeploymentState_JSONRoundTrip(t *testing.T) {
	states := []DeploymentState{
		DeploymentStateInitial,
		DeploymentStatePending,
		DeploymentStateActive,
		DeploymentStateOffline,
		DeploymentStateError,
		DeploymentStateLocked,
	}

	for _, originalState := range states {
		t.Run(fmt.Sprintf("Round trip for %s", originalState), func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(originalState)
			require.NoError(t, err)

			// Unmarshal from JSON
			var unmarshaledState DeploymentState
			err = json.Unmarshal(jsonData, &unmarshaledState)
			require.NoError(t, err)

			// Verify they're equal
			assert.Equal(t, originalState, unmarshaledState)
		})
	}
}

func TestDeploymentState_DatabaseRoundTrip(t *testing.T) {
	states := []DeploymentState{
		DeploymentStateInitial,
		DeploymentStatePending,
		DeploymentStateActive,
		DeploymentStateOffline,
		DeploymentStateError,
		DeploymentStateLocked,
	}

	for _, originalState := range states {
		t.Run(fmt.Sprintf("Database round trip for %s", originalState), func(t *testing.T) {
			dbValue, err := originalState.Value()
			require.NoError(t, err)

			var scannedState DeploymentState
			err = scannedState.Scan(dbValue)
			require.NoError(t, err)

			assert.Equal(t, originalState, scannedState)
		})
	}
}

func TestDeploymentState_WithStruct(t *testing.T) {
	t.Run("Cell struct with state", func(t *testing.T) {
		cell := Cell{
			DeployableEntity: DeployableEntity{
				Entity: Entity{
					ID:      "cell-123",
					Version: 0,
				},
				State:          DeploymentStateActive,
				StateTimestamp: time.Now(),
			},
			Properties: Properties{
				"key1": "value1",
			},
		}

		jsonData, err := json.Marshal(cell)
		require.NoError(t, err)
		assert.Contains(t, string(jsonData), `"active"`)

		var unmarshaledCell Cell
		err = json.Unmarshal(jsonData, &unmarshaledCell)
		require.NoError(t, err)
		assert.Equal(t, DeploymentStateActive, unmarshaledCell.State)
		assert.Equal(t, "cell-123", unmarshaledCell.ID)
	})
}

func TestDeploymentState_AllStatesValid(t *testing.T) {
	t.Run("All defined constants are valid", func(t *testing.T) {
		states := []DeploymentState{
			DeploymentStateInitial,
			DeploymentStatePending,
			DeploymentStateActive,
			DeploymentStateOffline,
			DeploymentStateError,
			DeploymentStateLocked,
		}

		for _, state := range states {
			assert.True(t, state.IsValid(), "State %s should be valid", state)
		}
	})
}

func TestDeploymentState_EdgeCases(t *testing.T) {
	t.Run("Zero value behavior", func(t *testing.T) {
		var state DeploymentState
		assert.False(t, state.IsValid())
		assert.Equal(t, "", state.String())
	})

	t.Run("Pointer handling", func(t *testing.T) {
		state := DeploymentStateActive
		ptr := &state

		// Verify pointer works with methods
		assert.True(t, ptr.IsValid())
		assert.Equal(t, "active", ptr.String())

		// Verify JSON marshaling with a pointer
		jsonData, err := json.Marshal(ptr)
		require.NoError(t, err)
		assert.Equal(t, `"active"`, string(jsonData))
	})

	t.Run("Enum constants are distinct", func(t *testing.T) {
		states := []DeploymentState{
			DeploymentStateInitial,
			DeploymentStatePending,
			DeploymentStateActive,
			DeploymentStateOffline,
			DeploymentStateError,
			DeploymentStateLocked,
		}

		// Verify states are unique
		for i, state1 := range states {
			for j, state2 := range states {
				if i != j {
					assert.NotEqual(t, state1, state2, "States %s and %s should be different", state1, state2)
				}
			}
		}
	})

}

func TestToDeploymentState(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      DeploymentState
		expectError   bool
		errorContains string
	}{
		{
			name:        "valid initial state",
			input:       "initial",
			expected:    DeploymentStateInitial,
			expectError: false,
		},
		{
			name:        "valid pending state",
			input:       "pending",
			expected:    DeploymentStatePending,
			expectError: false,
		},
		{
			name:        "valid active state",
			input:       "active",
			expected:    DeploymentStateActive,
			expectError: false,
		},
		{
			name:        "valid locked state",
			input:       "locked",
			expected:    DeploymentStateLocked,
			expectError: false,
		},
		{
			name:        "valid offline state",
			input:       "offline",
			expected:    DeploymentStateOffline,
			expectError: false,
		},
		{
			name:        "valid error state",
			input:       "error",
			expected:    DeploymentStateError,
			expectError: false,
		},
		{
			name:          "invalid state",
			input:         "invalid",
			expected:      "",
			expectError:   true,
			errorContains: "invalid deployment state: invalid",
		},
		{
			name:          "empty state",
			input:         "",
			expected:      "",
			expectError:   true,
			errorContains: "invalid deployment state:",
		},
		{
			name:        "uppercase state",
			input:       "ACTIVE",
			expected:    "active",
			expectError: false,
		},
		{
			name:        "mixed case state",
			input:       "Initial",
			expected:    "initial",
			expectError: false,
		},
		{
			name:          "whitespace state",
			input:         " active ",
			expected:      "",
			expectError:   true,
			errorContains: "invalid deployment state:  active ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToDeploymentState(tt.input)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Empty(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
				assert.True(t, result.IsValid())
			}
		})
	}
}

func TestToDeploymentState_AllValidStates(t *testing.T) {
	validStates := map[string]DeploymentState{
		"initial": DeploymentStateInitial,
		"pending": DeploymentStatePending,
		"active":  DeploymentStateActive,
		"locked":  DeploymentStateLocked,
		"offline": DeploymentStateOffline,
		"error":   DeploymentStateError,
	}

	for input, expected := range validStates {
		t.Run("valid_"+input, func(t *testing.T) {
			result, err := ToDeploymentState(input)

			require.NoError(t, err)
			assert.Equal(t, expected, result)
			assert.Equal(t, input, result.String())
		})
	}
}

func TestToDeploymentState_ConsistencyWithEnum(t *testing.T) {
	// Test that our function is consistent with the enum's validation
	testCases := []string{
		"initial", "pending", "active", "locked", "offline", "error",
		"invalid", "", "unknown", "ACTIVE", "Initial",
	}

	for _, input := range testCases {
		t.Run("consistency_"+input, func(t *testing.T) {
			result, err := ToDeploymentState(input)
			enumState := DeploymentState(strings.ToLower(input)) // Enum states are lowercase

			if enumState.IsValid() {
				require.NoError(t, err)
				assert.Equal(t, enumState, result)
			} else {
				require.Error(t, err)
				assert.Empty(t, result)
			}
		})
	}
}

func TestToProperties(t *testing.T) {
	t.Run("convert nil map", func(t *testing.T) {
		var m map[string]any

		result := ToProperties(m)

		require.NotNil(t, result)
		require.Len(t, result, 0)
	})

	t.Run("convert empty map", func(t *testing.T) {
		m := make(map[string]any)

		result := ToProperties(m)

		require.NotNil(t, result)
		require.Len(t, result, 0)
	})

	t.Run("convert map with string values", func(t *testing.T) {
		m := map[string]any{
			"environment": "production",
			"region":      "us-east-1",
		}

		result := ToProperties(m)

		require.NotNil(t, result)
		require.Len(t, result, 2)
		require.Equal(t, "production", result["environment"])
		require.Equal(t, "us-east-1", result["region"])
	})

	t.Run("convert map with multiple types", func(t *testing.T) {
		m := map[string]any{
			"environment": "production",
			"capacity":    100,
			"enabled":     true,
			"tags":        []string{"critical", "monitored"},
			"metadata": map[string]any{
				"owner": "platform-team",
			},
		}

		result := ToProperties(m)

		require.NotNil(t, result)
		require.Len(t, result, 5)
		require.Equal(t, "production", result["environment"])
		require.Equal(t, 100, result["capacity"])
		require.Equal(t, true, result["enabled"])

		tags, ok := result["tags"].([]string)
		require.True(t, ok)
		require.Len(t, tags, 2)
		require.Equal(t, "critical", tags[0])
		require.Equal(t, "monitored", tags[1])

		metadata, ok := result["metadata"].(map[string]any)
		require.True(t, ok)
		require.Equal(t, "platform-team", metadata["owner"])
	})

	t.Run("convert map preserves original", func(t *testing.T) {
		m := map[string]any{
			"environment": "production",
			"capacity":    100,
		}

		result := ToProperties(m)

		// Modify result to ensure original is not affected
		result["environment"] = "staging"
		result["new_key"] = "new_value"

		require.Equal(t, "production", m["environment"])
		require.Equal(t, 100, m["capacity"])
		_, exists := m["new_key"]
		require.False(t, exists)

		// Verify result has the modifications
		require.Equal(t, "staging", result["environment"])
		require.Equal(t, "new_value", result["new_key"])
	})
}
