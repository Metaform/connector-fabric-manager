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

package dmodel

// DeploymentManifest represents the configuration details for a system deployment. An Orchestration is instantiated
// from the manifest and executed.
//
// The manifest includes a unique identifier, the type of deployment specified by a DeploymentDefinition, and a payload
// of deployment-specific data, which will be passed as input to the Orchestration.
type DeploymentManifest struct {
	ID             string         `json:"id"`
	DeploymentType DeploymentType `json:"deploymentType"`
	Payload        map[string]any `json:"payload"`
}

// VPAManifest represents the configuration details for a VPA deployment.
type VPAManifest struct {
	ID         string         `json:"id"`
	VPAType    VPAType        `json:"vpaType"`
	Cell       string         `json:"cell"`
	Properties map[string]any `json:"properties"`
}

type DeploymentType string

func (dt DeploymentType) String() string {
	return string(dt)
}

type VPAType string

func (dt VPAType) String() string {
	return string(dt)
}

const (
	VpaDeploymentType     DeploymentType = "cfm.vpa"
	ConnectorType         VPAType        = "cfm.connector"
	CredentialServiceType VPAType        = "cfm.credentialservice"
	DataPlaneType         VPAType        = "cfm.dataplane"
	VpaPayloadType                       = "cfm.vpas"
)
