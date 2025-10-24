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
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/mvd/common"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

const (
	ActivityType = "connector-activity"
)

func LaunchAndWaitSignal(shutdown <-chan struct{}) {
	config := common.LauncherConfig{
		AgentName:    "Connector Agent",
		ConfigPrefix: "cagent",
		ActivityType: ActivityType,
		NewProcessor: func(monitor system.LogMonitor) api.ActivityProcessor {
			return &ConnectorActivityProcessor{monitor}
		},
	}
	common.LaunchAgent(shutdown, config)
}

type ConnectorActivityProcessor struct {
	monitor system.LogMonitor
}

func (t ConnectorActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	t.monitor.Infof("Connector provisioning complete")
	return api.ActivityResult{Result: api.ActivityResultComplete}
}
