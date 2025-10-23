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

//go:generate mockery

package api

import (
	"context"

	"time"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/system"
)

const (
	ProvisionManagerKey       system.ServiceType = "pmapi:ProvisionManager"
	DefinitionStoreKey        system.ServiceType = "pmapi:DefinitionStore"
	DeploymentOrchestratorKey system.ServiceType = "pmapi:DeploymentOrchestrator"
	DefinitionManagerKey      system.ServiceType = "pmapi:DefinitionManager"
)

// ProvisionManager handles deployments to the system.
type ProvisionManager interface {

	// Start initializes a new orchestration and starts its execution.
	// If a recoverable error is encountered one of model.RecoverableError, model.ClientError, or model.FatalError will be returned.
	Start(ctx context.Context, manifest *model.DeploymentManifest) (*Orchestration, error)

	// Cancel terminates an orchestration execution.
	// If a recoverable error is encountered one of model.RecoverableError, model.ClientError, or model.FatalError will be returned.
	Cancel(ctx context.Context, deploymentID string) error

	// GetOrchestration returns an orchestration by its ID or nil if not found.
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
// defined by WaitOnReschedule.
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
	Result           ActivityResultType
	WaitOnReschedule time.Duration
	Error            error
}

type ActivityContext interface {
	OID() string
	ID() string
	InputData() ImmutableMap
	SetValue(key string, value any)
	Value(key string) (any, bool)
	Values() map[string]any
	Delete(key string)
	SetOutputValue(key string, value any)
	OutputValues() map[string]any
	Context() context.Context
}

type ImmutableMap interface {
	Get(key string) (any, bool)
	Keys() []string
	Size() int
}

type DefinitionManager interface {
	CreateDeploymentDefinition(ctx context.Context, definition *DeploymentDefinition) (*DeploymentDefinition, error)
	CreateActivityDefinition(ctx context.Context, definition *ActivityDefinition) (*ActivityDefinition, error)
}
