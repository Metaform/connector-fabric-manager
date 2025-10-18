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

	"github.com/metaform/connector-fabric-manager/common/dmodel"
	"github.com/metaform/connector-fabric-manager/common/system"
)

const (
	DeploymentHandlerRegistryKey system.ServiceType = "tmapi:DeploymentHandlerRegistry"
	ParticipantDeployerKey       system.ServiceType = "tmapi:ParticipantDeployer"
	DeploymentClientKey          system.ServiceType = "tmapi:DeploymentClient"
)

type VpaPropMap = map[dmodel.VPAType]map[string]any

// ParticipantDeployer creates a participant profile and deploys its associated VPAs.
type ParticipantDeployer interface {
	Deploy(ctx context.Context, identifier string, vpaProperties VpaPropMap, properties map[string]any) error
}

// DeploymentClient asynchronously deploys a manifest to the provision manager. Implementations may use different wire protocols.
type DeploymentClient interface {

	// Deploy deploys the specified manifest.
	// If a recoverable error is encountered one of model.RecoverableError, model.ClientError, or model.FatalError will be returned.
	Deploy(ctx context.Context, manifest dmodel.DeploymentManifest) error
}

// DeploymentCallbackHandler is called when a deployment is complete.
// If a recoverable error is encountered one of model.RecoverableError, model.ClientError, or model.FatalError will be returned.
type DeploymentCallbackHandler func(context.Context, dmodel.DeploymentResponse) error

// DeploymentHandlerRegistry registers deployment handlers by deployment type.
type DeploymentHandlerRegistry interface {
	RegisterDeploymentHandler(deploymentType dmodel.DeploymentType, handler DeploymentCallbackHandler)
}
