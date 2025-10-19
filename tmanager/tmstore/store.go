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

package tmstore

import (
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

const (
	TManagerStoreKey system.ServiceType = "tmstore:TManagerStore"
)

// TManagerStore manages tenant entities.
type TManagerStore interface {

	// GetCells returns all cells in the system.
	GetCells() ([]api.Cell, error)

	// GetDataspaceProfiles GetCells returns all dataspace profiles in the system.
	GetDataspaceProfiles() ([]api.DataspaceProfile, error)

	// FindDeployment returns a deployment record by ID. If not found, returns errors.NotFound.
	FindDeployment(id string) (*api.DeploymentRecord, error)

	// DeploymentExists returns true if a deployment record exists with the given ID.
	DeploymentExists(id string) (bool, error)

	// CreateDeployment creates a new deployment record.
	CreateDeployment(record api.DeploymentRecord) (*api.DeploymentRecord, error)

	// UpdateDeployment updates an existing deployment record.
	UpdateDeployment(record api.DeploymentRecord) error
}
