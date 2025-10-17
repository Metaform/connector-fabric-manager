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

package tmcore

import (
	"context"

	"github.com/metaform/connector-fabric-manager/common/dmodel"
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

// deploymentCallbackService registers api.DeploymentCallbackHandler instances and dispatches deployment responses.
type deploymentCallbackService struct {
	handlers map[string]api.DeploymentCallbackHandler
}

func (d deploymentCallbackService) RegisterDeploymentHandler(deploymentType dmodel.DeploymentType, handler api.DeploymentCallbackHandler) {
	d.handlers[deploymentType.String()] = handler
}

func (d deploymentCallbackService) Dispatch(ctx context.Context, response dmodel.DeploymentResponse) error {
	handler, found := d.handlers[response.DeploymentType.String()]
	if !found {
		return model.NewFatalError("deployment handler not found for type: %s", response.DeploymentType)
	}
	return handler(ctx, response)
}
