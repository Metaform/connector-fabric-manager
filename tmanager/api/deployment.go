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

	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/dmodel"
)

const (
	DeploymentHandlerRegistryKey    system.ServiceType = "tmapi:DeploymentHandlerRegistry"
	DeploymentCallbackDispatcherKey system.ServiceType = "tmapi:DeploymentCallbackDispatcher"
	ParticipantDeployerKey          system.ServiceType = "tmapi:ParticipantDeployer"
	DeploymentClientKey             system.ServiceType = "tmapi:DeploymentClient"
)

// ParticipantDeployer creates a participant profile and deploys its associated VPAs.
type ParticipantDeployer interface {
	Deploy(ctx context.Context, identifier string, properties map[string]any) error
}

// DeploymentClient asynchronously deploys a manifest to the provision manager. Implementations may use different wire protocols.
type DeploymentClient interface {

	// Deploy deploys the specified manifest.
	// If a recoverable error is encountered one of model.RecoverableError, model.ClientError, or model.FatalError will be returned.
	Deploy(ctx context.Context, manifest dmodel.DeploymentManifest) error
}

// DeploymentResponse is asynchronously returned by the deployment client.
type DeploymentResponse struct {
	ID             string
	Success        bool
	ErrorDetail    string
	ManifestID     string
	DeploymentType string
	Properties     map[string]any
}

// DeploymentCallbackHandler is called when a deployment is complete.
// If a recoverable error is encountered one of model.RecoverableError, model.ClientError, or model.FatalError will be returned.
type DeploymentCallbackHandler func(context.Context, DeploymentResponse) error

// DeploymentHandlerRegistry registers deployment handlers by deployment type.
type DeploymentHandlerRegistry interface {
	RegisterDeploymentHandler(deploymentType string, handler DeploymentCallbackHandler)
}

// DeploymentCallbackDispatcher routes deployment responses to the associated handler.
type DeploymentCallbackDispatcher interface {

	// Dispatch is invoked when a deployment is complete.
	// If a recoverable error is encountered one of model.RecoverableError, model.ClientError, or model.FatalError will be returned.
	Dispatch(ctx context.Context, response DeploymentResponse) error
}
