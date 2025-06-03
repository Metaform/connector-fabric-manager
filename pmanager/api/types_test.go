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

func TestOrchestration_CanProceedToNextActivity(t *testing.T) {
	tests := []struct {
		name          string
		orchestration *Orchestration
		activityID    string
		validator     func([]string) bool
		want          bool
		wantErr       bool
	}{
		{
			name: "single step sequential orchestration",
			orchestration: &Orchestration{
				Steps: []OrchestrationStep{
					{
						Parallel: false,
						Activities: []Activity{
							{ID: "act1", Type: "test"},
						},
					},
				},
			},
			activityID: "act1",
			validator:  func([]string) bool { return true },
			want:       true,
			wantErr:    false,
		},
		{
			name: "multiple sequential steps",
			orchestration: &Orchestration{
				Steps: []OrchestrationStep{
					{
						Parallel: false,
						Activities: []Activity{
							{ID: "act1", Type: "test"},
							{ID: "act2", Type: "test"},
						},
					},
					{
						Parallel: false,
						Activities: []Activity{
							{ID: "act3", Type: "test"},
						},
					},
				},
			},
			activityID: "act2",
			validator:  func([]string) bool { return true },
			want:       true,
			wantErr:    false,
		},
		{
			name: "parallel step with all activities completed",
			orchestration: &Orchestration{
				Steps: []OrchestrationStep{
					{
						Parallel: true,
						Activities: []Activity{
							{ID: "act1", Type: "test"},
							{ID: "act2", Type: "test"},
							{ID: "act3", Type: "test"},
						},
					},
				},
			},
			activityID: "act2",
			validator:  func([]string) bool { return true },
			want:       true,
			wantErr:    false,
		},
		{
			name: "parallel step with pending activities",
			orchestration: &Orchestration{
				Steps: []OrchestrationStep{
					{
						Parallel: true,
						Activities: []Activity{
							{ID: "act1", Type: "test"},
							{ID: "act2", Type: "test"},
							{ID: "act3", Type: "test"},
						},
					},
				},
			},
			activityID: "act2",
			validator:  func([]string) bool { return false },
			want:       false,
			wantErr:    false,
		},
		{
			name: "activity not found",
			orchestration: &Orchestration{
				Steps: []OrchestrationStep{
					{
						Parallel: false,
						Activities: []Activity{
							{ID: "act1", Type: "test"},
						},
					},
				},
			},
			activityID: "non-existent",
			validator:  func([]string) bool { return true },
			want:       true,
			wantErr:    true,
		},
		{
			name: "mixed parallel and sequential steps",
			orchestration: &Orchestration{
				Steps: []OrchestrationStep{
					{
						Parallel: false,
						Activities: []Activity{
							{ID: "act1", Type: "test"},
						},
					},
					{
						Parallel: true,
						Activities: []Activity{
							{ID: "act2", Type: "test"},
							{ID: "act3", Type: "test"},
							{ID: "act4", Type: "test"},
						},
					},
					{
						Parallel: false,
						Activities: []Activity{
							{ID: "act5", Type: "test"},
						},
					},
				},
			},
			activityID: "act3",
			validator:  func([]string) bool { return false },
			want:       false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.orchestration.CanProceedToNextActivity(tt.activityID, tt.validator)
			if (err != nil) != tt.wantErr {
				t.Errorf("%v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("%v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStepForActivity(t *testing.T) {
	t.Run("single step orchestration - activity found", func(t *testing.T) {
		// Setup
		activity := Activity{ID: "activity1"}
		orchestration := &Orchestration{
			Steps: []OrchestrationStep{
				{
					Activities: []Activity{activity},
				},
			},
		}

		step, err := orchestration.GetStepForActivity("activity1")

		require.NoError(t, err)
		require.NotNil(t, step)
		require.Equal(t, activity, step.Activities[0])
	})

	t.Run("two step orchestration - activity found in second step", func(t *testing.T) {
		// Setup
		activity1 := Activity{ID: "activity1"}
		activity2 := Activity{ID: "activity2"}
		orchestration := &Orchestration{
			Steps: []OrchestrationStep{
				{
					Activities: []Activity{activity1},
				},
				{
					Activities: []Activity{activity2},
				},
			},
		}

		step, err := orchestration.GetStepForActivity("activity2")

		// Assert
		require.NoError(t, err)
		require.NotNil(t, step)
		require.Equal(t, activity2, step.Activities[0])
	})

	t.Run("activity not found", func(t *testing.T) {
		// Setup
		activity := Activity{ID: "activity1"}
		orchestration := &Orchestration{
			Steps: []OrchestrationStep{
				{
					Activities: []Activity{activity},
				},
			},
		}

		step, err := orchestration.GetStepForActivity("nonexistent")

		require.Error(t, err)
		require.Nil(t, step)
		require.Contains(t, err.Error(), "step not found for activity: nonexistent")
	})

	t.Run("empty orchestration", func(t *testing.T) {
		// Setup
		orchestration := &Orchestration{
			Steps: []OrchestrationStep{},
		}

		step, err := orchestration.GetStepForActivity("activity1")

		require.Error(t, err)
		require.Nil(t, step)
		require.Contains(t, err.Error(), "step not found for activity: activity1")
	})
}

func TestGetNextActivities(t *testing.T) {
	t.Run("single step with sequential activities", func(t *testing.T) {
		orch := &Orchestration{
			Steps: []OrchestrationStep{
				{
					Parallel: false,
					Activities: []Activity{
						{ID: "a1"},
						{ID: "a2"},
						{ID: "a3"},
					},
				},
			},
		}

		activities, next := orch.GetNextActivities("a1")
		require.Equal(t, 1, len(activities))
		require.Equal(t, "a2", activities[0].ID)
		require.False(t, next)

		activities, next = orch.GetNextActivities("a2")
		require.Equal(t, 1, len(activities))
		require.Equal(t, "a3", activities[0].ID)
		require.False(t, next)

		activities, next = orch.GetNextActivities("a3")
		require.Empty(t, activities)
		require.False(t, next)
	})

	t.Run("single step with parallel activities", func(t *testing.T) {
		orch := &Orchestration{
			Steps: []OrchestrationStep{
				{
					Parallel: true,
					Activities: []Activity{
						{ID: "a1"},
						{ID: "a2"},
						{ID: "a3"},
					},
				},
			},
		}

		activities, next := orch.GetNextActivities("a1")
		require.Equal(t, 0, len(activities))
		require.False(t, next)
	})

	t.Run("multiple steps - sequential to parallel", func(t *testing.T) {
		orch := &Orchestration{
			Steps: []OrchestrationStep{
				{
					Parallel: false,
					Activities: []Activity{
						{ID: "a1"},
						{ID: "a2"},
					},
				},
				{
					Parallel: true,
					Activities: []Activity{
						{ID: "b1"},
						{ID: "b2"},
						{ID: "b3"},
					},
				},
			},
		}

		activities, next := orch.GetNextActivities("a1")
		require.Equal(t, 1, len(activities))
		require.Equal(t, "a2", activities[0].ID)
		require.False(t, next)

		activities, next = orch.GetNextActivities("a2")
		require.Equal(t, 2, len(activities))
		require.Equal(t, "b1", activities[0].ID)
		require.Equal(t, "b2", activities[1].ID)
		require.True(t, next)
	})

	t.Run("multiple steps - parallel to sequential", func(t *testing.T) {
		orch := &Orchestration{
			Steps: []OrchestrationStep{
				{
					Parallel: true,
					Activities: []Activity{
						{ID: "a1"},
						{ID: "a2"},
					},
				},
				{
					Parallel: false,
					Activities: []Activity{
						{ID: "b1"},
						{ID: "b2"},
					},
				},
			},
		}

		activities, next := orch.GetNextActivities("a2")
		require.Equal(t, 1, len(activities))
		require.Equal(t, "b1", activities[0].ID)
		require.False(t, next)
	})

	t.Run("empty step in between", func(t *testing.T) {
		orch := &Orchestration{
			Steps: []OrchestrationStep{
				{
					Parallel: false,
					Activities: []Activity{
						{ID: "a1"},
					},
				},
				{
					Parallel:   false,
					Activities: []Activity{},
				},
				{
					Parallel: false,
					Activities: []Activity{
						{ID: "c1"},
					},
				},
			},
		}

		activities, next := orch.GetNextActivities("a1")
		require.Empty(t, activities)
		require.False(t, next)
	})

	t.Run("non-existent activity", func(t *testing.T) {
		orch := &Orchestration{
			Steps: []OrchestrationStep{
				{
					Parallel: false,
					Activities: []Activity{
						{ID: "a1"},
					},
				},
			},
		}

		activities, next := orch.GetNextActivities("non-existent")
		require.Empty(t, activities)
		require.False(t, next)
	})

	t.Run("last activity in parallel step", func(t *testing.T) {
		orch := &Orchestration{
			Steps: []OrchestrationStep{
				{
					Parallel: true,
					Activities: []Activity{
						{ID: "a1"},
						{ID: "a2"},
						{ID: "a3"},
					},
				},
			},
		}

		activities, next := orch.GetNextActivities("a3")
		require.Empty(t, activities)
		require.False(t, next)
	})
}

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
