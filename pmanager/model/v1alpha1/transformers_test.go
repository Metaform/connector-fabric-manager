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
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/stretchr/testify/assert"
)

func TestToAPIActivityDefinition(t *testing.T) {
	tests := []struct {
		name       string
		definition *ActivityDefinition
		expected   *api.ActivityDefinition
	}{
		{
			name: "complete activity definition",
			definition: &ActivityDefinition{
				Type:         "http-request",
				Description:  "Makes HTTP requests",
				InputSchema:  map[string]any{"url": "string"},
				OutputSchema: map[string]any{"response": "object"},
			},
			expected: &api.ActivityDefinition{
				Type:         api.ActivityType("http-request"),
				Description:  "Makes HTTP requests",
				InputSchema:  map[string]any{"url": "string"},
				OutputSchema: map[string]any{"response": "object"},
			},
		},
		{
			name: "minimal activity definition",
			definition: &ActivityDefinition{
				Type: "basic-task",
			},
			expected: &api.ActivityDefinition{
				Type: api.ActivityType("basic-task"),
			},
		},
		{
			name:       "empty activity definition",
			definition: &ActivityDefinition{},
			expected: &api.ActivityDefinition{
				Type: api.ActivityType(""),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToAPIActivityDefinition(tt.definition)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToAPIActivityDefinition_NilInput(t *testing.T) {
	// Test that the function handles nil input gracefully
	assert.NotPanics(t, func() {
		result := ToAPIActivityDefinition(nil)
		assert.Empty(t, result.Type)
		assert.Empty(t, result.Description)
		assert.Nil(t, result.InputSchema)
		assert.Nil(t, result.OutputSchema)
	})
}

func TestToAPIOrchestrationDefinition(t *testing.T) {
	tests := []struct {
		name                    string
		orchestrationDefinition *OrchestrationDefinition
		expected                *api.OrchestrationDefinition
	}{
		{
			name: "complete orchestration definition",
			orchestrationDefinition: &OrchestrationDefinition{
				Type:        "kubernetes",
				Description: "Test",
				Schema:      map[string]any{"version": "v1"},
				Activities: []Activity{
					{
						ID:            "activity-1",
						Type:          "http-request",
						Discriminator: "test-discriminator",
						Inputs: []MappingEntry{
							{Source: "input.url", Target: "request.url"},
							{Source: "input.method", Target: "request.method"},
						},
						DependsOn: []string{"activity-0"},
					},
					{
						ID:        "activity-2",
						Type:      "data-transform",
						Inputs:    []MappingEntry{},
						DependsOn: []string{"activity-1"},
					},
				},
			},
			expected: &api.OrchestrationDefinition{
				Type:        model.OrchestrationType("kubernetes"),
				Description: "Test",
				Active:      true,
				Schema:      map[string]any{"version": "v1"},
				Activities: []api.Activity{
					{
						ID:            "activity-1",
						Type:          api.ActivityType("http-request"),
						Discriminator: "test-discriminator",
						Inputs: []api.MappingEntry{
							{Source: "input.url", Target: "request.url"},
							{Source: "input.method", Target: "request.method"},
						},
						DependsOn: []string{"activity-0"},
					},
					{
						ID:        "activity-2",
						Type:      api.ActivityType("data-transform"),
						Inputs:    []api.MappingEntry{},
						DependsOn: []string{"activity-1"},
					},
				},
			},
		},
		{
			name: "minimal orchestration definition",
			orchestrationDefinition: &OrchestrationDefinition{
				Type:       "docker",
				Activities: []Activity{},
			},
			expected: &api.OrchestrationDefinition{
				Type:       model.OrchestrationType("docker"),
				Active:     true,
				Activities: []api.Activity{},
			},
		},
		{
			name:                    "empty orchestration definition",
			orchestrationDefinition: &OrchestrationDefinition{},
			expected: &api.OrchestrationDefinition{
				Type:       model.OrchestrationType(""),
				Active:     true,
				Activities: []api.Activity{},
			},
		},
		{
			name: "single activity without dependencies",
			orchestrationDefinition: &OrchestrationDefinition{
				Type: "local",
				Activities: []Activity{
					{
						ID:   "standalone-activity",
						Type: "file-processor",
						Inputs: []MappingEntry{
							{Source: "file.path", Target: "processor.input"},
						},
					},
				},
			},
			expected: &api.OrchestrationDefinition{
				Type:   model.OrchestrationType("local"),
				Active: true,
				Activities: []api.Activity{
					{
						ID:   "standalone-activity",
						Type: api.ActivityType("file-processor"),
						Inputs: []api.MappingEntry{
							{Source: "file.path", Target: "processor.input"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToAPIOrchestrationDefinition(tt.orchestrationDefinition)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToAPIOrchestrationDefinition_NilInput(t *testing.T) {
	// Test that the function handles nil input gracefully
	assert.NotPanics(t, func() {
		result := ToAPIOrchestrationDefinition(nil)
		assert.Empty(t, result.Type)
		assert.Nil(t, result.Schema)
		assert.Len(t, result.Activities, 0)
	})
}

func TestToAPIMappingEntries(t *testing.T) {
	tests := []struct {
		name     string
		entries  []MappingEntry
		expected []api.MappingEntry
	}{
		{
			name: "multiple mapping entries",
			entries: []MappingEntry{
				{Source: "input.name", Target: "output.fullName"},
				{Source: "input.age", Target: "output.yearsOld"},
				{Source: "input.email", Target: "output.contactEmail"},
			},
			expected: []api.MappingEntry{
				{Source: "input.name", Target: "output.fullName"},
				{Source: "input.age", Target: "output.yearsOld"},
				{Source: "input.email", Target: "output.contactEmail"},
			},
		},
		{
			name: "single mapping entry",
			entries: []MappingEntry{
				{Source: "data.value", Target: "result.processed"},
			},
			expected: []api.MappingEntry{
				{Source: "data.value", Target: "result.processed"},
			},
		},
		{
			name:     "empty mapping entries",
			entries:  []MappingEntry{},
			expected: []api.MappingEntry{},
		},
		{
			name: "mapping entries with empty strings",
			entries: []MappingEntry{
				{Source: "", Target: ""},
				{Source: "valid.source", Target: ""},
				{Source: "", Target: "valid.target"},
			},
			expected: []api.MappingEntry{
				{Source: "", Target: ""},
				{Source: "valid.source", Target: ""},
				{Source: "", Target: "valid.target"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToAPIMappingEntries(tt.entries)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToAPIMappingEntries_NilInput(t *testing.T) {
	result := ToAPIMappingEntries(nil)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestToOrchestrationEntry_VerifiesInputs(t *testing.T) {
	testTime := time.Now()
	input := api.OrchestrationEntry{
		ID:                "test-id-123",
		CorrelationID:     "corr-id-456",
		State:             5,
		StateTimestamp:    testTime,
		CreatedTimestamp:  testTime.Add(-time.Hour),
		OrchestrationType: model.OrchestrationType("TestType"),
	}

	result := ToOrchestrationEntry(&input)

	assert.Equal(t, input.ID, result.ID)
	assert.Equal(t, input.CorrelationID, result.CorrelationID)
	assert.Equal(t, int(input.State), result.State)
	assert.Equal(t, input.StateTimestamp, result.StateTimestamp)
	assert.Equal(t, input.CreatedTimestamp, result.CreatedTimestamp)
	assert.Equal(t, input.OrchestrationType, result.OrchestrationType)
}

func TestToOrchestration(t *testing.T) {
	now := time.Now()
	apiOrchestration := &api.Orchestration{
		ID:                "test-orch-1",
		CorrelationID:     "corr-123",
		State:             api.OrchestrationStateRunning,
		StateTimestamp:    now,
		CreatedTimestamp:  now.Add(-1 * time.Hour),
		OrchestrationType: "test-type",
		ProcessingData:    map[string]any{"key1": "value1"},
		OutputData:        map[string]any{"key2": "value2"},
		Completed:         map[string]struct{}{"activity1": {}},
		Steps: []api.OrchestrationStep{
			{
				Activities: []api.Activity{
					{
						ID:            "activity-1",
						Type:          "test.activity",
						Discriminator: "test-disc",
						Inputs: []api.MappingEntry{
							{Source: "src1", Target: "tgt1"},
						},
						DependsOn: []string{"activity-0"},
					},
				},
			},
		},
	}

	result := ToOrchestration(apiOrchestration)

	assert.Equal(t, "test-orch-1", result.ID)
	assert.Equal(t, "corr-123", result.CorrelationID)
	assert.Equal(t, int(api.OrchestrationStateRunning), result.State)
	assert.Equal(t, now, result.StateTimestamp)
	assert.Equal(t, now.Add(-1*time.Hour), result.CreatedTimestamp)
	assert.Equal(t, model.OrchestrationType("test-type"), result.OrchestrationType)
	assert.Equal(t, map[string]any{"key1": "value1"}, result.ProcessingData)
	assert.Equal(t, map[string]any{"key2": "value2"}, result.OutputData)
	assert.Equal(t, 1, len(result.Steps))
	assert.Equal(t, 1, len(result.Steps[0].Activities))
	assert.Equal(t, "activity-1", result.Steps[0].Activities[0].ID)
	assert.Equal(t, "test.activity", result.Steps[0].Activities[0].Type)
	assert.Equal(t, 1, len(result.Steps[0].Activities[0].Inputs))
	assert.Equal(t, "src1", result.Steps[0].Activities[0].Inputs[0].Source)
}
