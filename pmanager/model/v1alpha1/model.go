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

type ActivityDefinition struct {
	Type         string         `json:"type" validate:"required"`
	Description  string         `json:"description omitempty"`
	InputSchema  map[string]any `json:"inputSchema omitempty"`
	OutputSchema map[string]any `json:"outputSchema omitempty"`
}

type Activity struct {
	ID            string         `json:"id" validate:"required"`
	Type          string         `json:"type" validate:"required"`
	Discriminator string         `json:"discriminator" validate:"false"`
	Inputs        []MappingEntry `json:"inputs omitempty"`
	DependsOn     []string       `json:"dependsOn omitempty"`
}

type MappingEntry struct {
	Source string `json:"source" validate:"required"`
	Target string `json:"target" validate:"required"`
}

type OrchestrationDefinition struct {
	Type        string         `json:"type" validate:"required"`
	Description string         `json:"description omitempty"`
	Schema      map[string]any `json:"schema omitempty"`
	Activities  []Activity     `json:"activities" validate:"required,min=1"`
}
