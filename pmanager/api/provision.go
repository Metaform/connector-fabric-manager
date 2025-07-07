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
	"time"
)

const (
	ProvisionManagerKey system.ServiceType = "pmapi:ProvisionManager"
	DefinitionStoreKey  system.ServiceType = "pmapi:DefinitionStore"
)

// ProvisionManager handles deployments to the system.
type ProvisionManager interface {
	Start(ctx context.Context, manifest *DeploymentManifest) (*Orchestration, error)
	Cancel(ctx context.Context, deploymentID string) error
	GetOrchestration(ctx context.Context, deploymentID string) (*Orchestration, error)
}

// DeploymentOrchestrator orchestrates deployments.
// Implementations must support idempotent behavior.
type DeploymentOrchestrator interface {

	// ExecuteOrchestration starts the execution of the specified orchestration, processing its steps and activities.
	ExecuteOrchestration(ctx context.Context, orchestration *Orchestration) error

	// GetOrchestration retrieves an Orchestration by its ID or nil if not found.
	GetOrchestration(ctx context.Context, id string) (*Orchestration, error)
}

// ActivityProcessor executes activities for a given type.
//
// If the execution completes successfully, the processor returns ActivityResultComplete.
//
// If the processor returns ActivityResultWait, the activity will remain outstanding until completion is asynchronously signaled.
//
// If the processor returns ActivityResultSchedule, the orchestration engine will reschedule message delivery in the duration
// defined by WaitMillis.
//
// If the processor encounters an error, it returns an ActivityResultRetryError or an ActivityResultFatalError.
type ActivityProcessor interface {
	Process(activityContext ActivityContext) ActivityResult
}

type ActivityResultType int

const (
	ActivityResultWait       = 0
	ActivityResultComplete   = 1
	ActivityResultSchedule   = 2
	ActivityResultRetryError = -1
	ActivityResultFatalError = -2
)

type ActivityResult struct {
	Result     ActivityResultType
	WaitMillis time.Duration
	Error      error
}

type ActivityContext interface {
	OID() string
	ID() string
	SetValue(key string, value any)
	Value(key string) (any, bool)
	Values() map[string]any
	Context() context.Context
}

// DefinitionStore manages DeploymentDefinitions and ActivityDefinitions.
type DefinitionStore interface {

	// FindDeploymentDefinition retrieves the DeploymentDefinition associated with the given id.
	// Returns the DeploymentDefinition object or store.ErrNotFound if the definition cannot be found.
	FindDeploymentDefinition(id string) (*DeploymentDefinition, error)

	// FindActivityDefinition retrieves the ActivityDefinition associated with the given id.
	// Returns the ActivityDefinition object or store.ErrNotFound if the definition cannot be found.
	FindActivityDefinition(id string) (*ActivityDefinition, error)

	// StoreDeploymentDefinition saves or updates a DeploymentDefinition
	StoreDeploymentDefinition(id string, definition *DeploymentDefinition)

	// StoreActivityDefinition saves or updates a ActivityDefinition
	StoreActivityDefinition(id string, definition *ActivityDefinition)

	// DeleteDeploymentDefinition removes a DeploymentDefinition for the given id, returning true if successful.
	DeleteDeploymentDefinition(id string) bool

	// DeleteActivityDefinition removes an ActivityDefinition for the given id, returning true if successful.
	DeleteActivityDefinition(id string) bool

	// ListDeploymentDefinitions returns DeploymentDefinition instances with pagination support
	ListDeploymentDefinitions(offset, limit int) ([]*DeploymentDefinition, bool, error)

	// ListActivityDefinitions returns ActivityDefinition instances with pagination support
	ListActivityDefinitions(offset, limit int) ([]*ActivityDefinition, bool, error)
}
