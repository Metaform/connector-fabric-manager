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

// ActivityProcessor executes activities for a given type. If the execution completes successfully, the processor
// returns true.
//
// If the process requires time to complete, the processor should return false to indicate that the orchestration should
// wait for the activity to complete. At that point, the processor is responsible for signaling completion of the
// activity.
//
// If the processor encounters an error, it returns an error to indicate that the orchestration should fail.
type ActivityProcessor interface {
	Process(activityContext ActivityContext) (bool, error)
}

type ActivityContext interface {
	OID() string
	ID() string
	SetValue(key string, value any)
	Value(key string) any
	Context() context.Context
}
