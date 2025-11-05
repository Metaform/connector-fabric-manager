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

// DefinitionStore manages OrchestrationDefinition and ActivityDefinitions.
type DefinitionStore interface {

	// FindOrchestrationDefinition retrieves the OrchestrationDefinition associated with the given type.
	// Returns the OrchestrationDefinition object or store.ErrNotFound if the definition cannot be found.
	FindOrchestrationDefinition(orchestrationType model.OrchestrationType) (*OrchestrationDefinition, error)

	// FindActivityDefinition retrieves the ActivityDefinition associated with the given type.
	// Returns the ActivityDefinition object or store.ErrNotFound if the definition cannot be found.
	FindActivityDefinition(activityType ActivityType) (*ActivityDefinition, error)

	ExistsOrchestrationDefinition(orchestrationType model.OrchestrationType) (bool, error)

	ExistsActivityDefinition(activityType ActivityType) (bool, error)

	// StoreOrchestrationDefinition saves or updates a OrchestrationDefinition
	StoreOrchestrationDefinition(definition *OrchestrationDefinition) (*OrchestrationDefinition, error)

	// StoreActivityDefinition saves or updates a ActivityDefinition
	StoreActivityDefinition(definition *ActivityDefinition) (*ActivityDefinition, error)

	// DeleteOrchestrationDefinition removes a OrchestrationDefinition for the given type, returning true if successful.
	DeleteOrchestrationDefinition(orchestrationType model.OrchestrationType) (bool, error)

	// DeleteActivityDefinition removes an ActivityDefinition for the given type, returning true if successful.
	DeleteActivityDefinition(activityType ActivityType) (bool, error)

	// ListOrchestrationDefinitions returns OrchestrationDefinition instances with pagination support
	ListOrchestrationDefinitions(offset, limit int) ([]*OrchestrationDefinition, bool, error)

	// ListActivityDefinitions returns ActivityDefinition instances with pagination support
	ListActivityDefinitions(offset, limit int) ([]*ActivityDefinition, bool, error)
}
