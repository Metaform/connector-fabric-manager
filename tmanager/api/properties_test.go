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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestProperties_Value(t *testing.T) {
	t.Run("nil properties", func(t *testing.T) {
		var p Properties

		value, err := p.Value()

		require.NoError(t, err)
		require.Nil(t, value)
	})

	t.Run("empty properties", func(t *testing.T) {
		p := make(Properties)

		value, err := p.Value()

		require.NoError(t, err)
		require.Nil(t, value)
	})

	t.Run("properties with string value", func(t *testing.T) {
		p := Properties{
			"environment": "production",
		}

		value, err := p.Value()

		require.NoError(t, err)
		require.NotNil(t, value)

		bytes, ok := value.([]byte)
		require.True(t, ok)

		var result map[string]any
		err = json.Unmarshal(bytes, &result)
		require.NoError(t, err)
		require.Equal(t, "production", result["environment"])
	})

	t.Run("properties with multiple types", func(t *testing.T) {
		p := Properties{
			"environment": "production",
			"capacity":    100,
			"enabled":     true,
			"tags":        []string{"critical", "monitored"},
			"metadata": map[string]any{
				"owner": "platform-team",
			},
		}

		value, err := p.Value()

		require.NoError(t, err)
		require.NotNil(t, value)

		bytes, ok := value.([]byte)
		require.True(t, ok)

		var result map[string]any
		err = json.Unmarshal(bytes, &result)
		require.NoError(t, err)
		require.Equal(t, "production", result["environment"])
		require.Equal(t, float64(100), result["capacity"]) // JSON unmarshals numbers as float64
		require.Equal(t, true, result["enabled"])
		require.Len(t, result["tags"], 2)

		metadata, ok := result["metadata"].(map[string]any)
		require.True(t, ok)
		require.Equal(t, "platform-team", metadata["owner"])
	})
}

func TestProperties_Scan(t *testing.T) {
	t.Run("scan nil value", func(t *testing.T) {
		var p Properties

		err := p.Scan(nil)

		require.NoError(t, err)
		require.Nil(t, p)
	})

	t.Run("scan empty byte slice", func(t *testing.T) {
		var p Properties

		err := p.Scan([]byte{})

		require.NoError(t, err)
		require.NotNil(t, p)
		require.Len(t, p, 0)
	})

	t.Run("scan valid JSON bytes", func(t *testing.T) {
		var p Properties
		jsonData := []byte(`{"environment":"production","capacity":100,"enabled":true}`)

		err := p.Scan(jsonData)

		require.NoError(t, err)
		require.NotNil(t, p)
		require.Equal(t, "production", p["environment"])
		require.Equal(t, float64(100), p["capacity"])
		require.Equal(t, true, p["enabled"])
	})

	t.Run("scan valid JSON string", func(t *testing.T) {
		var p Properties
		jsonString := `{"region":"us-east-1","tags":["critical","monitored"]}`

		err := p.Scan(jsonString)

		require.NoError(t, err)
		require.NotNil(t, p)
		require.Equal(t, "us-east-1", p["region"])

		tags, ok := p["tags"].([]any)
		require.True(t, ok)
		require.Len(t, tags, 2)
		require.Equal(t, "critical", tags[0])
		require.Equal(t, "monitored", tags[1])
	})

	t.Run("scan complex nested JSON", func(t *testing.T) {
		var p Properties
		jsonData := []byte(`{
			"metadata": {
				"owner": "platform-team",
				"cost_center": "engineering"
			},
			"monitoring": {
				"enabled": true,
				"interval": "30s",
				"alerts": ["error", "warning"]
			}
		}`)

		err := p.Scan(jsonData)

		require.NoError(t, err)
		require.NotNil(t, p)

		metadata, ok := p["metadata"].(map[string]any)
		require.True(t, ok)
		require.Equal(t, "platform-team", metadata["owner"])
		require.Equal(t, "engineering", metadata["cost_center"])

		monitoring, ok := p["monitoring"].(map[string]any)
		require.True(t, ok)
		require.Equal(t, true, monitoring["enabled"])
		require.Equal(t, "30s", monitoring["interval"])

		alerts, ok := monitoring["alerts"].([]any)
		require.True(t, ok)
		require.Len(t, alerts, 2)
	})

	t.Run("scan invalid JSON", func(t *testing.T) {
		var p Properties
		invalidJSON := []byte(`{"invalid": json}`)

		err := p.Scan(invalidJSON)

		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid character")
	})

	t.Run("scan unsupported type", func(t *testing.T) {
		var p Properties

		err := p.Scan(123)

		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot scan int into Properties")
	})

	t.Run("scan time value", func(t *testing.T) {
		var p Properties

		err := p.Scan(time.Now())

		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot scan time.Time into Properties")
	})
}

func TestProperties_Get(t *testing.T) {
	t.Run("get from nil properties", func(t *testing.T) {
		var p Properties

		value, exists := p.Get("key")

		require.False(t, exists)
		require.Nil(t, value)
	})

	t.Run("get existing key", func(t *testing.T) {
		p := Properties{
			"environment": "production",
			"capacity":    100,
		}

		value, exists := p.Get("environment")

		require.True(t, exists)
		require.Equal(t, "production", value)
	})

	t.Run("get non-existing key", func(t *testing.T) {
		p := Properties{
			"environment": "production",
		}

		value, exists := p.Get("nonexistent")

		require.False(t, exists)
		require.Nil(t, value)
	})

	t.Run("get nil value", func(t *testing.T) {
		p := Properties{
			"nullable": nil,
		}

		value, exists := p.Get("nullable")

		require.True(t, exists)
		require.Nil(t, value)
	})
}

func TestProperties_GetString(t *testing.T) {
	t.Run("get string value", func(t *testing.T) {
		p := Properties{
			"environment": "production",
		}

		value, ok := p.GetString("environment")

		require.True(t, ok)
		require.Equal(t, "production", value)
	})

	t.Run("get non-existing key", func(t *testing.T) {
		p := Properties{
			"environment": "production",
		}

		_, ok := p.GetString("nonexistent")

		require.False(t, ok)
	})

	t.Run("get non-string value", func(t *testing.T) {
		p := Properties{
			"capacity": 100,
			"enabled":  true,
		}

		capacityValue, capacityOk := p.GetString("capacity")
		enabledValue, enabledOk := p.GetString("enabled")

		require.False(t, capacityOk)
		require.Equal(t, "", capacityValue)
		require.False(t, enabledOk)
		require.Equal(t, "", enabledValue)
	})

	t.Run("get nil value", func(t *testing.T) {
		p := Properties{
			"nullable": nil,
		}

		value, ok := p.GetString("nullable")

		require.False(t, ok)
		require.Equal(t, "", value)
	})

	t.Run("get from nil properties", func(t *testing.T) {
		var p Properties

		value, ok := p.GetString("key")

		require.False(t, ok)
		require.Equal(t, "", value)
	})
}

func TestProperties_GetInt(t *testing.T) {
	t.Run("get int value", func(t *testing.T) {
		p := Properties{
			"capacity": 100,
		}

		value, ok := p.GetInt("capacity")

		require.True(t, ok)
		require.Equal(t, 100, value)
	})

	t.Run("get non-existing key", func(t *testing.T) {
		p := Properties{
			"environment": "production",
		}

		value, ok := p.GetInt("nonexistent")

		require.False(t, ok)
		require.Equal(t, 0, value)
	})

	t.Run("get non-numeric value", func(t *testing.T) {
		p := Properties{
			"environment": "production",
			"enabled":     true,
		}

		envValue, envOk := p.GetInt("environment")
		enabledValue, enabledOk := p.GetInt("enabled")

		require.False(t, envOk)
		require.Equal(t, 0, envValue)
		require.False(t, enabledOk)
		require.Equal(t, 0, enabledValue)
	})

	t.Run("get nil value", func(t *testing.T) {
		p := Properties{
			"nullable": nil,
		}

		value, ok := p.GetInt("nullable")

		require.False(t, ok)
		require.Equal(t, 0, value)
	})

	t.Run("get from nil properties", func(t *testing.T) {
		var p Properties

		value, ok := p.GetInt("key")

		require.False(t, ok)
		require.Equal(t, 0, value)
	})
}

func TestProperties_Set(t *testing.T) {
	t.Run("set on nil properties", func(t *testing.T) {
		var p Properties

		p.Set("environment", "production")

		require.NotNil(t, p)
		require.Equal(t, "production", p["environment"])
	})

	t.Run("set on initialized properties", func(t *testing.T) {
		p := make(Properties)

		p.Set("environment", "production")
		p.Set("capacity", 100)
		p.Set("enabled", true)

		assert.Equal(t, "production", p["environment"])
		assert.Equal(t, 100, p["capacity"])
		assert.Equal(t, true, p["enabled"])
	})

	t.Run("overwrite existing value", func(t *testing.T) {
		p := Properties{
			"environment": "staging",
		}

		p.Set("environment", "production")

		assert.Equal(t, "production", p["environment"])
	})

	t.Run("set nil value", func(t *testing.T) {
		p := make(Properties)

		p.Set("nullable", nil)

		value, exists := p.Get("nullable")
		require.True(t, exists)
		require.Nil(t, value)
	})
}

func TestProperties_Roundtrip(t *testing.T) {
	t.Run("round trip through database interfaces", func(t *testing.T) {

		original := Properties{
			"environment": "production",
			"capacity":    100,
			"enabled":     true,
			"tags":        []string{"critical", "monitored"},
			"metadata": map[string]any{
				"owner": "platform-team",
			},
		}

		// Convert to database value
		value, err := original.Value()
		require.NoError(t, err)
		require.NotNil(t, value)

		// Scan from database value
		var restored Properties
		err = restored.Scan(value)
		require.NoError(t, err)

		// Verify deserialization
		assert.Equal(t, "production", restored["environment"])
		assert.Equal(t, 100.0, restored["capacity"])

		enabledValue, exists := restored.Get("enabled")
		require.True(t, exists)
		require.Equal(t, true, enabledValue)

		tags, exists := restored.Get("tags")
		require.True(t, exists)
		tagsSlice, ok := tags.([]any)
		require.True(t, ok)
		require.Len(t, tagsSlice, 2)
	})

	t.Run("round trip through JSON marshaling", func(t *testing.T) {
		original := Properties{
			"region": "us-east-1",
			"count":  5,
		}

		bytes, err := json.Marshal(original)
		require.NoError(t, err)

		var restored Properties
		err = json.Unmarshal(bytes, &restored)
		require.NoError(t, err)

		require.Equal(t, "us-east-1", restored["region"])
		require.Equal(t, 5.0, restored["count"])
	})
}

func TestProperties_EdgeCases(t *testing.T) {
	t.Run("properties with special characters in keys", func(t *testing.T) {
		p := Properties{
			"key-with-dashes":    "value1",
			"key.with.dots":      "value2",
			"key_with_undercore": "value3",
			"key with spaces":    "value4",
		}

		value, err := p.Value()
		require.NoError(t, err)

		var restored Properties
		err = restored.Scan(value)
		require.NoError(t, err)

		assert.Equal(t, "value1", restored["key-with-dashes"])
		assert.Equal(t, "value2", restored["key.with.dots"])
		assert.Equal(t, "value3", restored["key_with_undercore"])
		assert.Equal(t, "value4", restored["key with spaces"])
	})

	t.Run("deeply nested properties", func(t *testing.T) {
		p := Properties{
			"level1": map[string]any{
				"level2": map[string]any{
					"level3": map[string]any{
						"level4": "deep value",
					},
				},
			},
		}

		value, err := p.Value()
		require.NoError(t, err)

		var restored Properties
		err = restored.Scan(value)
		require.NoError(t, err)

		level1, exists := restored.Get("level1")
		require.True(t, exists)

		level1Map, ok := level1.(map[string]any)
		require.True(t, ok)

		_, exists = level1Map["level2"]
		require.True(t, exists)
	})
}
