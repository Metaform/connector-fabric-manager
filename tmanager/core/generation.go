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

package core

import (
	"errors"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

// defaultVPASelector iterates through cells and dataspace profiles to find and return the first active cell; returns an
// error if none are found.
func defaultVPASelector(_ model.DeploymentType, cells []api.Cell, dProfiles []api.DataspaceProfile) (*api.Cell, error) {
	for _, cell := range cells {
		if cell.State == api.DeploymentStateActive {
			for _, dProfile := range dProfiles {
				for _, deployment := range dProfile.Deployments {
					if deployment.State == api.DeploymentStateActive && deployment.Cell.ID == cell.ID {
						return &cell, nil
					}
				}
			}
		}
	}
	return nil, errors.New("no active cell found")
}
