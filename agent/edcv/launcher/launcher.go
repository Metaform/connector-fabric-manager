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
	"github.com/metaform/connector-fabric-manager/agent/edcv/activity"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsagent"
)

const (
	ActivityType = "edcv-activity"
)

func LaunchAndWaitSignal(shutdown <-chan struct{}) {
	config := natsagent.LauncherConfig{
		AgentName:    "EDC-V Agent",
		ConfigPrefix: "edcvagent",
		ActivityType: ActivityType,
		NewProcessor: func(monitor system.LogMonitor) api.ActivityProcessor {
			return &activity.EDCVActivityProcessor{Monitor: monitor}
		},
	}
	natsagent.LaunchAgent(shutdown, config)
}

type ConnectorActivityProcessor struct {
	monitor system.LogMonitor
}
