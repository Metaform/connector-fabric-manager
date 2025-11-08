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

package launcher

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsagent"
)

const (
	ActivityType = "connector-activity"
)

func LaunchAndWaitSignal(shutdown <-chan struct{}) {
	config := natsagent.LauncherConfig{
		AgentName:    "Connector Agent",
		ConfigPrefix: "cagent",
		ActivityType: ActivityType,
		NewProcessor: func(monitor system.LogMonitor) api.ActivityProcessor {
			return &ConnectorActivityProcessor{monitor}
		},
	}
	natsagent.LaunchAgent(shutdown, config)
}

type ConnectorActivityProcessor struct {
	monitor system.LogMonitor
}

func (t ConnectorActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	identifier, found := ctx.InputData().Get(model.ParticipantIdentifier)
	if !found {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("missing participant identifier")}
	}
	_, found = ctx.InputData().Get(model.VPADispose)
	if found {
		// disposal request
		ctx.SetOutputValue(model.VPADispose, true)
		return api.ActivityResult{Result: api.ActivityResultComplete}
	}
	// Return state data
	ctx.SetOutputValue(model.ConnectorId, uuid.New().String())

	t.monitor.Infof("Connector provisioning complete: %s", identifier)
	return api.ActivityResult{Result: api.ActivityResultComplete}
}
