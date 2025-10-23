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
	"time"

	"github.com/metaform/connector-fabric-manager/common/model"
)

type Entity struct {
	ID      string `json:"id" required:"true"`
	Version int64  `json:"version" required:"true"`
}
type NewCell struct {
	State          string         `json:"state" required:"true"`
	StateTimestamp time.Time      `json:"stateTimestamp" required:"true"`
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
	ProfileID string `json:"profileId" required:"true"`
	CellID    string `json:"cellId,omitempty"`
}

type DataspaceDeployment struct {
	DeployableEntity
	CellID     string         `json:"cellId,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
}
type DataspaceProfile struct {
	Entity
	Artifacts   []string              `json:"artifacts,omitempty"`
	Deployments []DataspaceDeployment `json:"deployments,omitempty"`
	Properties  map[string]any        `json:"properties,omitempty"`
}

type NewParticipantProfileDeployment struct {
	Identifier    string                    `json:"identifier" required:"true"`
	CellID        string                    `json:"cellId" required:"true"`
	VPAProperties map[string]map[string]any `json:"vpaProperties,omitempty"`
	Properties    map[string]any            `json:"properties,omitempty"`
}

type ParticipantProfile struct {
	Entity
	Identifier  string                    `json:"identifier" required:"true"`
	VPAs        []VirtualParticipantAgent `json:"vpas,omitempty"`
	Properties  map[string]any            `json:"properties,omitempty"`
	Error       bool                      `json:"error"`
	ErrorDetail string                    `json:"errorDetail,omitempty"`
}

type VirtualParticipantAgent struct {
	DeployableEntity
	Type       model.VPAType  `json:"type" required:"true"`
	Cell       Cell           `json:"cell" required:"true"`
	Properties map[string]any `json:"properties,omitempty"`
}

type DeployableEntity struct {
	Entity
	State          string    `json:"state" required:"true"`
	StateTimestamp time.Time `json:"stateTimestamp" required:"true"`
}
