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
	ProvisionManagerKey          system.ServiceType = "pmapi:ProvisionManager"
	ActivityProcessorRegistryKey system.ServiceType = "pmapi:ActivityProcessorRegistry"
)

// ProvisionManager handles deployments to the system.
type ProvisionManager interface {
	Start(manifest *DeploymentManifest) (string, error)
	Cancel(id string) error
}

type ActivityProcessorRegistry interface {
	RegisterProcessor(processor ActivityProcessor)
}

// ActivityProcessor executes activities for a given type.
//
// If the execution completes successfully, the processor returns ActivityResultContinue.
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
	ActivityResultContinue   = 1
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
