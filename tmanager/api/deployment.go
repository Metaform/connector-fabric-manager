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

func ToVPAMap(vpaProperties map[string]map[string]any) *VPAPropMap {
	vpaPropsMap := make(VPAPropMap)
	for vpaTypeStr, props := range vpaProperties {
		vpaType := model.VPAType(vpaTypeStr)
		vpaPropsMap[vpaType] = props
	}
	return &vpaPropsMap
}
