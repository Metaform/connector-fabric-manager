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
	"iter"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/query"
)

// DefinitionStore manages OrchestrationDefinition and ActivityDefinitions.
type DefinitionStore interface {

	// FindOrchestrationDefinition retrieves the OrchestrationDefinition associated with the given type.
	// Returns the OrchestrationDefinition object or store.ErrNotFound if the definition cannot be found.
	FindOrchestrationDefinition(ctx context.Context, orchestrationType model.OrchestrationType) (*OrchestrationDefinition, error)

	// FindOrchestrationDefinitionsByPredicate retrieves OrchestrationDefinition instances matching the given predicate.
	FindOrchestrationDefinitionsByPredicate(ctx context.Context, predicate query.Predicate) iter.Seq2[OrchestrationDefinition, error]

	// ExistsOrchestrationDefinition returns true if an OrchestrationDefinition exists for the given type.
	ExistsOrchestrationDefinition(ctx context.Context, orchestrationType model.OrchestrationType) (bool, error)

	// GetOrchestrationDefinitionCount returns the number of OrchestrationDefinitions matching the given predicate.
	GetOrchestrationDefinitionCount(_ context.Context, predicate query.Predicate) (int, error)

	// FindActivityDefinition retrieves the ActivityDefinition associated with the given type.
	// Returns the ActivityDefinition object or store.ErrNotFound if the definition cannot be found.
	FindActivityDefinition(ctx context.Context, activityType ActivityType) (*ActivityDefinition, error)

	// FindActivityDefinitionsByPredicate retrieves ActivityDefinition instances matching the given predicate.
	FindActivityDefinitionsByPredicate(ctx context.Context, predicate query.Predicate) iter.Seq2[ActivityDefinition, error]

	// ExistsActivityDefinition returns true if an ActivityDefinition exists for the given type.
	ExistsActivityDefinition(ctx context.Context, activityType ActivityType) (bool, error)

	// GetActivityDefinitionCount returns the number of ActivityDefinitions matching the given predicate.
	GetActivityDefinitionCount(_ context.Context, predicate query.Predicate) (int, error)

	// StoreOrchestrationDefinition saves or updates a OrchestrationDefinition
	StoreOrchestrationDefinition(ctx context.Context, definition *OrchestrationDefinition) (*OrchestrationDefinition, error)

	// StoreActivityDefinition saves or updates a ActivityDefinition
	StoreActivityDefinition(ctx context.Context, definition *ActivityDefinition) (*ActivityDefinition, error)

	// DeleteOrchestrationDefinition removes a OrchestrationDefinition for the given type, returning true if successful.
	DeleteOrchestrationDefinition(ctx context.Context, orchestrationType model.OrchestrationType) (bool, error)

	ActivityDefinitionReferences(ctx context.Context, activityType ActivityType) ([]string, error)

	// DeleteActivityDefinition removes an ActivityDefinition for the given type, returning true if successful.
	DeleteActivityDefinition(ctx context.Context, activityType ActivityType) (bool, error)

	// ListOrchestrationDefinitions returns OrchestrationDefinition instances with pagination support
	ListOrchestrationDefinitions(ctx context.Context, offset int, limit int) ([]*OrchestrationDefinition, bool, error)

	// ListActivityDefinitions returns ActivityDefinition instances with pagination support
	ListActivityDefinitions(ctx context.Context, offset int, limit int) ([]*ActivityDefinition, bool, error)
}
