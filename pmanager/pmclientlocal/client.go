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

package pmclientlocal

import (
	"github.com/metaform/connector-fabric-manager/assembly/httpclient"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

type localProvisionManagerClient struct{}

func (l localProvisionManagerClient) Provision(manifest *api.DeploymentManifest) error {
	return nil
}

type LocalPmClientServiceAssembly struct {
	system.DefaultServiceAssembly
}

func (l LocalPmClientServiceAssembly) Name() string {
	return "Local Provision Manager Client"
}

func (l LocalPmClientServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{api.ProvisionManagerClientKey}
}

func (l LocalPmClientServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{httpclient.HttpClientKey}
}

func (l LocalPmClientServiceAssembly) Init(context *system.InitContext) error {
	context.Registry.Register(api.ProvisionManagerClientKey, localProvisionManagerClient{})
	return nil
}
