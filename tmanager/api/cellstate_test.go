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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCellState_String(t *testing.T) {
	tests := []struct {
		name  string
		state CellState
		want  string
	}{
		{"Initial", CellStateInitial, "initial"},
		{"Pending", CellStatePending, "pending"},
		{"Active", CellStateActive, "active"},
		{"Offline", CellStateOffline, "offline"},
		{"Error", CellStateError, "error"},
		{"Locked", CellStateLocked, "locked"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.state.String())
		})
	}
}

func TestCellState_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		state CellState
		want  bool
	}{
		{"Initial is valid", CellStateInitial, true},
		{"Pending is valid", CellStatePending, true},
		{"Active is valid", CellStateActive, true},
		{"Offline is valid", CellStateOffline, true},
		{"Error is valid", CellStateError, true},
		{"Locked is valid", CellStateLocked, true},
		{"Empty string is invalid", CellState(""), false},
		{"Invalid state is invalid", CellState("invalid"), false},
		{"Case sensitive - uppercase is invalid", CellState("ACTIVE"), false},
		{"Case sensitive - mixed case is invalid", CellState("Active"), false},
		{"Whitespace state is invalid", CellState(" active "), false},
		{"Similar but wrong state is invalid", CellState("activ"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.state.IsValid())
		})
	}
}

func TestCellState_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		state   CellState
		want    string
		wantErr bool
	}{
		{"Initial marshals correctly", CellStateInitial, `"initial"`, false},
		{"Pending marshals correctly", CellStatePending, `"pending"`, false},
		{"Active marshals correctly", CellStateActive, `"active"`, false},
		{"Offline marshals correctly", CellStateOffline, `"offline"`, false},
		{"Error marshals correctly", CellStateError, `"error"`, false},
		{"Locked marshals correctly", CellStateLocked, `"locked"`, false},
		{"Invalid state marshals as string", CellState("invalid"), `"invalid"`, false},
		{"Empty state marshals as empty string", CellState(""), `""`, false},
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

func TestCellState_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    CellState
		wantErr bool
	}{
		{"Initial unmarshals correctly", `"initial"`, CellStateInitial, false},
		{"Pending unmarshals correctly", `"pending"`, CellStatePending, false},
		{"Active unmarshals correctly", `"active"`, CellStateActive, false},
		{"Offline unmarshals correctly", `"offline"`, CellStateOffline, false},
		{"Error unmarshals correctly", `"error"`, CellStateError, false},
		{"Locked unmarshals correctly", `"locked"`, CellStateLocked, false},
		{"Invalid state returns error", `"invalid"`, CellState(""), true},
		{"Empty string returns error", `""`, CellState(""), true},
		{"Uppercase returns error", `"ACTIVE"`, CellState(""), true},
		{"Mixed case returns error", `"Active"`, CellState(""), true},
		{"Whitespace returns error", `" active "`, CellState(""), true},
		{"Invalid JSON returns error", `invalid`, CellState(""), true},
		{"Number returns error", `123`, CellState(""), true},
		{"Boolean returns error", `true`, CellState(""), true},
		{"Object returns error", `{"state":"active"}`, CellState(""), true},
		{"Array returns error", `["active"]`, CellState(""), true},
		{"Null returns error", `null`, CellState(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var state CellState
			err := json.Unmarshal([]byte(tt.json), &state)
			if tt.wantErr {
				require.Error(t, err)
				if tt.json == `"invalid"` || tt.json == `""` || tt.json == `"ACTIVE"` ||
					tt.json == `"Active"` || tt.json == `" active "` {
					assert.Contains(t, err.Error(), "invalid cell state")
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, state)
			}
		})
	}
}

func TestCellState_Value(t *testing.T) {
	tests := []struct {
		name    string
		state   CellState
		want    string
		wantErr bool
	}{
		{"Initial returns correct value", CellStateInitial, "initial", false},
		{"Pending returns correct value", CellStatePending, "pending", false},
		{"Active returns correct value", CellStateActive, "active", false},
		{"Offline returns correct value", CellStateOffline, "offline", false},
		{"Error returns correct value", CellStateError, "error", false},
		{"Locked returns correct value", CellStateLocked, "locked", false},
		{"Invalid state returns error", CellState("invalid"), "", true},
		{"Empty state returns error", CellState(""), "", true},
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

func TestCellState_Scan(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    CellState
		wantErr bool
	}{
		{"String initial scans correctly", "initial", CellStateInitial, false},
		{"String pending scans correctly", "pending", CellStatePending, false},
		{"String active scans correctly", "active", CellStateActive, false},
		{"String offline scans correctly", "offline", CellStateOffline, false},
		{"String error scans correctly", "error", CellStateError, false},
		{"String locked scans correctly", "locked", CellStateLocked, false},
		{"Byte slice initial scans correctly", []byte("initial"), CellStateInitial, false},
		{"Byte slice active scans correctly", []byte("active"), CellStateActive, false},
		{"Nil value scans to empty", nil, CellState(""), false},
		{"Invalid string returns error", "invalid", CellState(""), true},
		{"Empty string returns error", "", CellState(""), true},
		{"Uppercase string returns error", "ACTIVE", CellState(""), true},
		{"Mixed case string returns error", "Active", CellState(""), true},
		{"Whitespace string returns error", " active ", CellState(""), true},
		{"Integer value returns error", 123, CellState(""), true},
		{"Boolean value returns error", true, CellState(""), true},
		{"Float value returns error", 3.14, CellState(""), true},
		{"Struct value returns error", struct{}{}, CellState(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var state CellState
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

func TestCellState_JSONRoundTrip(t *testing.T) {
	states := []CellState{
		CellStateInitial,
		CellStatePending,
		CellStateActive,
		CellStateOffline,
		CellStateError,
		CellStateLocked,
	}

	for _, originalState := range states {
		t.Run(fmt.Sprintf("Round trip for %s", originalState), func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(originalState)
			require.NoError(t, err)

			// Unmarshal from JSON
			var unmarshaledState CellState
			err = json.Unmarshal(jsonData, &unmarshaledState)
			require.NoError(t, err)

			// Verify they're equal
			assert.Equal(t, originalState, unmarshaledState)
		})
	}
}

func TestCellState_DatabaseRoundTrip(t *testing.T) {
	states := []CellState{
		CellStateInitial,
		CellStatePending,
		CellStateActive,
		CellStateOffline,
		CellStateError,
		CellStateLocked,
	}

	for _, originalState := range states {
		t.Run(fmt.Sprintf("Database round trip for %s", originalState), func(t *testing.T) {
			dbValue, err := originalState.Value()
			require.NoError(t, err)

			var scannedState CellState
			err = scannedState.Scan(dbValue)
			require.NoError(t, err)

			assert.Equal(t, originalState, scannedState)
		})
	}
}

func TestCellState_WithStruct(t *testing.T) {
	t.Run("Cell struct with state", func(t *testing.T) {
		cell := Cell{
			Entity: Entity{
				ID: "cell-123",
			},
			State: CellStateActive,
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
		assert.Equal(t, CellStateActive, unmarshaledCell.State)
		assert.Equal(t, "cell-123", unmarshaledCell.ID)
	})
}

func TestCellState_AllStatesValid(t *testing.T) {
	t.Run("All defined constants are valid", func(t *testing.T) {
		states := []CellState{
			CellStateInitial,
			CellStatePending,
			CellStateActive,
			CellStateOffline,
			CellStateError,
			CellStateLocked,
		}

		for _, state := range states {
			assert.True(t, state.IsValid(), "State %s should be valid", state)
		}
	})
}

func TestCellState_EdgeCases(t *testing.T) {
	t.Run("Zero value behavior", func(t *testing.T) {
		var state CellState
		assert.False(t, state.IsValid())
		assert.Equal(t, "", state.String())
	})

	t.Run("Pointer handling", func(t *testing.T) {
		state := CellStateActive
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
		states := []CellState{
			CellStateInitial,
			CellStatePending,
			CellStateActive,
			CellStateOffline,
			CellStateError,
			CellStateLocked,
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
