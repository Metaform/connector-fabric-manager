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
	"context"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/system"
)

const (
	DeploymentHandlerRegistryKey system.ServiceType = "tmapi:DeploymentHandlerRegistry"
	DeploymentClientKey          system.ServiceType = "tmapi:DeploymentClient"
	ParticipantProfileServiceKey system.ServiceType = "tmapi:ParticipantProfileService"
	DataspaceProfileServiceKey   system.ServiceType = "tmapi:DataspaceProfileService"
	CellServiceKey               system.ServiceType = "tmapi:CellService"
)

type VPAPropMap = map[model.VPAType]map[string]any

// DeploymentClient asynchronously deploys a manifest to the provision manager. Implementations may use different wire protocols.
type DeploymentClient interface {
	// Send deploys the specified manifest.
	// If a recoverable error is encountered one of model.RecoverableError, model.ClientError, or model.FatalError will be returned.
	Send(ctx context.Context, manifest model.DeploymentManifest) error
}

// DeploymentCallbackHandler is called when a deployment is complete.
// If a recoverable error is encountered one of model.RecoverableError, model.ClientError, or model.FatalError will be returned.
type DeploymentCallbackHandler func(context.Context, model.DeploymentResponse) error

// DeploymentHandlerRegistry registers deployment handlers by deployment type.
type DeploymentHandlerRegistry interface {
	RegisterDeploymentHandler(deploymentType model.DeploymentType, handler DeploymentCallbackHandler)
}

// ParticipantProfileService performs participant profile operations, including deploying associated VPAs.
type ParticipantProfileService interface {
	DeployProfile(ctx context.Context, identifier string, vpaProperties VPAPropMap, properties map[string]any) (*ParticipantProfile, error)
	DisposeProfile(ctx context.Context, identifier string) error
	GetProfile(ctx context.Context, id string) (*ParticipantProfile, error)
}

// DataspaceProfileService performs dataspace profile operations.
type DataspaceProfileService interface {
	CreateProfile(ctx context.Context, artifacts []string, properties map[string]any) (*DataspaceProfile, error)
	DeployProfile(ctx context.Context, profileID string, cellID string) error
	GetProfile(ctx context.Context, profileID string) (*DataspaceProfile, error)
}

// CellService performs cell operations.
type CellService interface {
	RecordExternalDeployment(ctx context.Context, cell Cell) (*Cell, error)
}

func ToVPAMap(vpaProperties map[string]map[string]any) *VPAPropMap {
	vpaPropsMap := make(VPAPropMap)
	for vpaTypeStr, props := range vpaProperties {
		vpaType := model.VPAType(vpaTypeStr)
		vpaPropsMap[vpaType] = props
	}
	return &vpaPropsMap
}
