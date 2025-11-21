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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivityDefinitionTypeValidation(t *testing.T) {

	tests := []struct {
		name    string
		actDef  ActivityDefinition
		wantErr bool
	}{
		{
			name: "valid type with alphanumeric and allowed chars",
			actDef: ActivityDefinition{
				Type: "my-activity_v1.0",
			},
			wantErr: false,
		},
		{
			name: "valid type with only alphanumeric",
			actDef: ActivityDefinition{
				Type: "SimpleActivity",
			},
			wantErr: false,
		},
		{
			name: "valid type with dots and underscores",
			actDef: ActivityDefinition{
				Type: "com.example.activity_name",
			},
			wantErr: false,
		},
		{
			name: "valid type with hyphens",
			actDef: ActivityDefinition{
				Type: "my-activity-type",
			},
			wantErr: false,
		},
		{
			name: "valid type with numbers",
			actDef: ActivityDefinition{
				Type: "activity123",
			},
			wantErr: false,
		},
		{
			name: "invalid type with space",
			actDef: ActivityDefinition{
				Type: "invalid activity",
			},
			wantErr: true,
		},
		{
			name: "invalid type with special chars @",
			actDef: ActivityDefinition{
				Type: "invalid@activity",
			},
			wantErr: true,
		},
		{
			name: "invalid type with special chars #",
			actDef: ActivityDefinition{
				Type: "invalid#activity",
			},
			wantErr: true,
		},
		{
			name: "invalid type with special chars $",
			actDef: ActivityDefinition{
				Type: "invalid$activity",
			},
			wantErr: true,
		},
		{
			name: "invalid type with special chars %",
			actDef: ActivityDefinition{
				Type: "invalid%activity",
			},
			wantErr: true,
		},
		{
			name: "invalid type with special chars &",
			actDef: ActivityDefinition{
				Type: "invalid&activity",
			},
			wantErr: true,
		},
		{
			name: "invalid type empty string",
			actDef: ActivityDefinition{
				Type: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := model.Validator.Var(tt.actDef.Type, "required,modeltype")
			if tt.wantErr {
				require.Error(t, err, "expected validation error")
			} else {
				require.NoError(t, err, "expected no validation error")
			}
		})
	}
}

func TestActivityTypeValidation(t *testing.T) {

	tests := []struct {
		name     string
		activity Activity
		wantErr  bool
	}{
		{
			name: "valid activity type",
			activity: Activity{
				ID:   "act-1",
				Type: "process-data",
			},
			wantErr: false,
		},
		{
			name: "valid activity type with dots",
			activity: Activity{
				ID:   "act-2",
				Type: "com.example.ProcessActivity",
			},
			wantErr: false,
		},
		{
			name: "valid activity type with underscores",
			activity: Activity{
				ID:   "act-3",
				Type: "process_data_activity",
			},
			wantErr: false,
		},
		{
			name: "valid activity type with numbers",
			activity: Activity{
				ID:   "act-4",
				Type: "activity2process3",
			},
			wantErr: false,
		},
		{
			name: "valid activity type with version",
			activity: Activity{
				ID:   "act-5",
				Type: "activity.v1.0",
			},
			wantErr: false,
		},
		{
			name: "invalid activity type with space",
			activity: Activity{
				ID:   "act-6",
				Type: "invalid type",
			},
			wantErr: true,
		},
		{
			name: "invalid activity type with @",
			activity: Activity{
				ID:   "act-7",
				Type: "invalid@type",
			},
			wantErr: true,
		},
		{
			name: "invalid activity type with #",
			activity: Activity{
				ID:   "act-8",
				Type: "invalid#type",
			},
			wantErr: true,
		},
		{
			name: "invalid activity type with $",
			activity: Activity{
				ID:   "act-9",
				Type: "invalid$type",
			},
			wantErr: true,
		},
		{
			name: "invalid activity type with %",
			activity: Activity{
				ID:   "act-10",
				Type: "invalid%type",
			},
			wantErr: true,
		},
		{
			name: "invalid activity type with &",
			activity: Activity{
				ID:   "act-11",
				Type: "invalid&type",
			},
			wantErr: true,
		},
		{
			name: "invalid activity type empty",
			activity: Activity{
				ID:   "act-12",
				Type: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := model.Validator.Var(tt.activity.Type, "required,modeltype")
			if tt.wantErr {
				assert.Error(t, err, "expected validation error for type: %s", tt.activity.Type)
			} else {
				assert.NoError(t, err, "expected no validation error for type: %s", tt.activity.Type)
			}
		})
	}
}

func TestOrchestrationDefinitionTypeValidation(t *testing.T) {
	tests := []struct {
		name    string
		orchDef OrchestrationDefinition
		wantErr bool
	}{
		{
			name: "valid orchestration type",
			orchDef: OrchestrationDefinition{
				Type: "sequential-orchestration",
			},
			wantErr: false,
		},
		{
			name: "valid orchestration type with version",
			orchDef: OrchestrationDefinition{
				Type: "orchestration.v1.0",
			},
			wantErr: false,
		},
		{
			name: "valid orchestration type with namespace",
			orchDef: OrchestrationDefinition{
				Type: "com.company.orchestration_type",
			},
			wantErr: false,
		},
		{
			name: "valid orchestration type with numbers",
			orchDef: OrchestrationDefinition{
				Type: "orch123-type",
			},
			wantErr: false,
		},
		{
			name: "valid orchestration type simple",
			orchDef: OrchestrationDefinition{
				Type: "Orchestration",
			},
			wantErr: false,
		},
		{
			name: "invalid orchestration type with space",
			orchDef: OrchestrationDefinition{
				Type: "invalid orchestration",
			},
			wantErr: true,
		},
		{
			name: "invalid orchestration type with @",
			orchDef: OrchestrationDefinition{
				Type: "invalid@orchestration",
			},
			wantErr: true,
		},
		{
			name: "invalid orchestration type with #",
			orchDef: OrchestrationDefinition{
				Type: "invalid#orchestration",
			},
			wantErr: true,
		},
		{
			name: "invalid orchestration type with $",
			orchDef: OrchestrationDefinition{
				Type: "invalid$orchestration",
			},
			wantErr: true,
		},
		{
			name: "invalid orchestration type with %",
			orchDef: OrchestrationDefinition{
				Type: "invalid%orchestration",
			},
			wantErr: true,
		},
		{
			name: "invalid orchestration type with &",
			orchDef: OrchestrationDefinition{
				Type: "invalid&orchestration",
			},
			wantErr: true,
		},
		{
			name: "invalid orchestration type empty",
			orchDef: OrchestrationDefinition{
				Type: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := model.Validator.Var(tt.orchDef.Type, "required,modeltype")
			if tt.wantErr {
				assert.Error(t, err, "expected validation error for type: %s", tt.orchDef.Type)
			} else {
				assert.NoError(t, err, "expected no validation error for type: %s", tt.orchDef.Type)
			}
		})
	}
}
