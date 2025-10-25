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
	"time"

	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/agent"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

const (
	ActivityType = "onboard-service-activity"
	countKey     = "ob.count"
)

func LaunchAndWaitSignal(shutdown <-chan struct{}) {
	config := agent.LauncherConfig{
		AgentName:    "Onboarding Agent",
		ConfigPrefix: "obagent",
		ActivityType: ActivityType,
		NewProcessor: func(monitor system.LogMonitor) api.ActivityProcessor {
			return &ConnectorActivityProcessor{monitor}
		},
	}
	agent.LaunchAgent(shutdown, config)
}

type ConnectorActivityProcessor struct {
	monitor system.LogMonitor
}

func (t ConnectorActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	count, found := ctx.Value(countKey)
	if (found) && (count.(float64) > 0) {
		t.monitor.Infof("Onboarding and credential setup complete")
		ctx.Delete(countKey)
		return api.ActivityResult{Result: api.ActivityResultComplete}
	}
	t.monitor.Infof("Onboarding Initiated")
	ctx.SetValue(countKey, 1)
	return api.ActivityResult{Result: api.ActivityResultSchedule, WaitOnReschedule: 2 * time.Second}
}
