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
	Type         string         `json:"type"`
	Provider     string         `json:"provider"`
	Description  string         `json:"description"`
	InputSchema  map[string]any `json:"inputSchema"`
	OutputSchema map[string]any `json:"outputSchema"`
}

type Activity struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	Inputs    []MappingEntry `json:"inputs"`
	DependsOn []string       `json:"dependsOn"`
}

type MappingEntry struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type DeploymentDefinition struct {
	Type       string         `json:"type"`
	Schema     map[string]any `json:"schema"`
	Activities []Activity     `json:"activities"`
}
