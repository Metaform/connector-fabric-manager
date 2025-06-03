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
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
	"testing"
)

func TestMappingEntry_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    MappingEntry
		wantErr bool
	}{
		{
			name: "string input",
			json: `"test_value"`,
			want: MappingEntry{
				Source: "test_value",
				Target: "test_value",
			},
			wantErr: false,
		},
		{
			name: "object input",
			json: `{"source": "src_value", "target": "tgt_value"}`,
			want: MappingEntry{
				Source: "src_value",
				Target: "tgt_value",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			json:    `{"invalid`,
			want:    MappingEntry{},
			wantErr: true,
		},
		{
			name:    "invalid type",
			json:    `42`,
			want:    MappingEntry{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got MappingEntry
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("MappingEntry.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.DeepEqual(t, got, tt.want)
		})
	}
}

func TestParseDeploymentDefinition(t *testing.T) {
	result, err := ParseDeploymentDefinition([]byte(testDefinition))
	require.NoError(t, err)

	expected := &DeploymentDefinition{
		Type:       "tenant.example.com",
		ApiVersion: "1.0",
		Resource: Resource{
			Group:       "example.com",
			Singular:    "tenant",
			Plural:      "tenants",
			Description: "Deploys infrastructure and configuration required to support a tenant",
		},
		Versions: []Version{
			{
				Version: "1.0.0",
				Active:  true,
				Schema: map[string]any{
					"openAPIV3Schema": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"cell": map[string]any{
								"type": "string",
							},
						},
						"required": []any{"cell"},
					},
				},
				OrchestrationDefinition: OrchestrationDefinition{
					{
						Parallel: false,
						Activities: []Activity{
							{
								Type: "dns.example.com",
								Inputs: []MappingEntry{
									MappingEntry{
										Source: "cell",
										Target: "cell",
									},
									MappingEntry{
										Source: "baseUrl",
										Target: "baseUrl",
									},
								},
							},
							{
								Type: "ihtenant.example.com",
								Inputs: []MappingEntry{
									MappingEntry{
										Source: "cell",
										Target: "cell",
									}, MappingEntry{
										Source: "test.dataspaces",
										Target: "dataspaces",
									},
								},
							},
						},
					},

					{
						Parallel:   true,
						Activities: []Activity{},
					},
				},
			},
		},
	}

	assert.DeepEqual(t, result, expected)
}

const testDefinition = `{
  "type": "tenant.example.com",
  "apiVersion": "1.0",
  "resource": {
    "group": "example.com",
    "singular": "tenant",
    "plural": "tenants",
    "description": "Deploys infrastructure and configuration required to support a tenant"
  },
  "versions": [
    {
      "version": "1.0.0",
      "active": true,
      "schema": {
        "openAPIV3Schema": {
          "type": "object",
          "properties": {
            "cell": {
              "type": "string"
            }
          },
          "required": [
            "cell"
          ]
        }
      },
      "orchestration": [
        {
          "parallel": false,
          "activities": [
            {
              "type": "dns.example.com",
              "inputs": [
                "cell",
                "baseUrl"
              ]
            },
            {
              "type": "ihtenant.example.com",
              "inputs": [
                "cell",
                {
                  "source": "test.dataspaces",
                  "target": "dataspaces"
                }
              ]
            }
          ]
        },
        {
          "parallel": true,
          "activities": []
        }
      ]
    }
  ]
}`
