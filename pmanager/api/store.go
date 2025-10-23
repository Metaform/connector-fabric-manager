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

import "github.com/metaform/connector-fabric-manager/common/model"

// DefinitionStore manages DeploymentDefinitions and ActivityDefinitions.
type DefinitionStore interface {

	// FindDeploymentDefinition retrieves the DeploymentDefinition associated with the given type.
	// Returns the DeploymentDefinition object or store.ErrNotFound if the definition cannot be found.
	FindDeploymentDefinition(deploymentType model.DeploymentType) (*DeploymentDefinition, error)

	// FindActivityDefinition retrieves the ActivityDefinition associated with the given type.
	// Returns the ActivityDefinition object or store.ErrNotFound if the definition cannot be found.
	FindActivityDefinition(activityType ActivityType) (*ActivityDefinition, error)

	ExistsDeploymentDefinition(deploymentType model.DeploymentType) (bool, error)

	ExistsActivityDefinition(activityType ActivityType) (bool, error)

	// StoreDeploymentDefinition saves or updates a DeploymentDefinition
	StoreDeploymentDefinition(definition *DeploymentDefinition) (*DeploymentDefinition, error)

	// StoreActivityDefinition saves or updates a ActivityDefinition
	StoreActivityDefinition(definition *ActivityDefinition) (*ActivityDefinition, error)

	// DeleteDeploymentDefinition removes a DeploymentDefinition for the given type, returning true if successful.
	DeleteDeploymentDefinition(deploymentType model.DeploymentType) (bool, error)

	// DeleteActivityDefinition removes an ActivityDefinition for the given type, returning true if successful.
	DeleteActivityDefinition(activityType ActivityType) (bool, error)

	// ListDeploymentDefinitions returns DeploymentDefinition instances with pagination support
	ListDeploymentDefinitions(offset, limit int) ([]*DeploymentDefinition, bool, error)

	// ListActivityDefinitions returns ActivityDefinition instances with pagination support
	ListActivityDefinitions(offset, limit int) ([]*ActivityDefinition, bool, error)
}
