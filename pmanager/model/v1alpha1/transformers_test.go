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
				Provider:     "http-provider",
				Description:  "Makes HTTP requests",
				InputSchema:  map[string]any{"url": "string"},
				OutputSchema: map[string]any{"response": "object"},
			},
			expected: &api.ActivityDefinition{
				Type:         api.ActivityType("http-request"),
				Provider:     "http-provider",
				Description:  "Makes HTTP requests",
				InputSchema:  map[string]interface{}{"url": "string"},
				OutputSchema: map[string]interface{}{"response": "object"},
			},
		},
		{
			name: "minimal activity definition",
			definition: &ActivityDefinition{
				Type:     "basic-task",
				Provider: "basic-provider",
			},
			expected: &api.ActivityDefinition{
				Type:     api.ActivityType("basic-task"),
				Provider: "basic-provider",
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
		assert.Empty(t, result.Provider)
		assert.Empty(t, result.Description)
		assert.Nil(t, result.InputSchema)
		assert.Nil(t, result.OutputSchema)
	})
}

func TestToAPIDeploymentDefinition(t *testing.T) {
	tests := []struct {
		name                 string
		deploymentDefinition *DeploymentDefinition
		expected             *api.DeploymentDefinition
	}{
		{
			name: "complete deployment definition",
			deploymentDefinition: &DeploymentDefinition{
				Type:   "kubernetes",
				Schema: map[string]interface{}{"version": "v1"},
				Activities: []Activity{
					{
						ID:   "activity-1",
						Type: "http-request",
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
			expected: &api.DeploymentDefinition{
				Type:   model.DeploymentType("kubernetes"),
				Active: true,
				Schema: map[string]interface{}{"version": "v1"},
				Activities: []api.Activity{
					{
						ID:   "activity-1",
						Type: api.ActivityType("http-request"),
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
			name: "minimal deployment definition",
			deploymentDefinition: &DeploymentDefinition{
				Type:       "docker",
				Activities: []Activity{},
			},
			expected: &api.DeploymentDefinition{
				Type:       model.DeploymentType("docker"),
				Active:     true,
				Activities: []api.Activity{},
			},
		},
		{
			name:                 "empty deployment definition",
			deploymentDefinition: &DeploymentDefinition{},
			expected: &api.DeploymentDefinition{
				Type:       model.DeploymentType(""),
				Active:     true,
				Activities: []api.Activity{},
			},
		},
		{
			name: "single activity without dependencies",
			deploymentDefinition: &DeploymentDefinition{
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
			expected: &api.DeploymentDefinition{
				Type:   model.DeploymentType("local"),
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
			result := ToAPIDeploymentDefinition(tt.deploymentDefinition)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToAPIDeploymentDefinition_NilInput(t *testing.T) {
	// Test that the function handles nil input gracefully
	assert.NotPanics(t, func() {
		result := ToAPIDeploymentDefinition(nil)
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

// Benchmark tests
func BenchmarkToAPIActivityDefinition(b *testing.B) {
	definition := &ActivityDefinition{
		Type:         "http-request",
		Provider:     "http-provider",
		Description:  "Makes HTTP requests",
		InputSchema:  map[string]interface{}{"url": "string"},
		OutputSchema: map[string]interface{}{"response": "object"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToAPIActivityDefinition(definition)
	}
}

func BenchmarkToAPIDeploymentDefinition(b *testing.B) {
	deploymentDefinition := &DeploymentDefinition{
		Type:   "kubernetes",
		Schema: map[string]interface{}{"version": "v1"},
		Activities: []Activity{
			{
				ID:   "activity-1",
				Type: "http-request",
				Inputs: []MappingEntry{
					{Source: "input.url", Target: "request.url"},
					{Source: "input.method", Target: "request.method"},
				},
				DependsOn: []string{"activity-0"},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToAPIDeploymentDefinition(deploymentDefinition)
	}
}

func BenchmarkToAPIMappingEntries(b *testing.B) {
	entries := []MappingEntry{
		{Source: "input.name", Target: "output.fullName"},
		{Source: "input.age", Target: "output.yearsOld"},
		{Source: "input.email", Target: "output.contactEmail"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToAPIMappingEntries(entries)
	}
}
