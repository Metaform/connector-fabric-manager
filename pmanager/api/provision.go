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

import "github.com/metaform/connector-fabric-manager/common/system"

const (
	ProvisionManagerKey    system.ServiceType = "pmapi:ProvisionManager"
	ProvisionerRegistryKey system.ServiceType = "pmapi:ProvisionerRegistry"
)

type ProvisionerRegistry interface {
	RegisterProvisioner(provisioner Provisioner)
}

type ProvisionerBase interface {
	Start(manifest *DeploymentManifest) (string, error)
	Cancel(id string) error
}

type ProvisionManager interface {
	ProvisionerBase
}

type Provisioner interface {
	CanProcess(manifest *DeploymentManifest) bool
	ProvisionerBase
}
