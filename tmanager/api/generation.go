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
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/dmodel"
)

const (
	CellSelectorKey system.ServiceType = "tmapi:CellSelector"
)

// CellSelector selects a cell for resource deployment.
type CellSelector func(dmodel.DeploymentType, []Cell, []DataspaceProfile) (*Cell, error)
