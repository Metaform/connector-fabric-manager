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

import "time"

type Entity struct {
	ID      string `json:"id"`
	Version int64  `json:"version"`
}
type NewCell struct {
	State          string         `json:"state"`
	StateTimestamp time.Time      `json:"stateTimestamp"`
	Properties     map[string]any `json:"properties,omitempty"`
}

type Cell struct {
	Entity
	NewCell
}

type NewDataspaceProfile struct {
	Artifacts  []string       ` json:"artifacts,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
}

type NewDataspaceProfileDeployment struct {
	ProfileID string `json:"profileId"`
	CellID    string `json:"cellId"`
}

type DataspaceDeployment struct {
	Entity
	State          string         `json:"state"`
	StateTimestamp time.Time      `json:"stateTimestamp"`
	CellID         string         `json:"cellId"`
	Properties     map[string]any `json:"properties,omitempty"`
}
type DataspaceProfile struct {
	Entity
	Artifacts   []string
	Deployments []DataspaceDeployment
	Properties  map[string]any `json:"properties,omitempty"`
}
