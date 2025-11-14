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
	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsagent"
)

const (
	agentName    = "Test Agent"
	activityType = "test.activity"
	configPrefix = "testagent"
)

func LaunchAndWaitSignal() {
	Launch(runtime.CreateSignalShutdownChan())
}

func Launch(shutdown <-chan struct{}) {
	config := natsagent.LauncherConfig{
		AgentName:    agentName,
		ConfigPrefix: configPrefix,
		ActivityType: activityType,
		NewProcessor: func(monitor system.LogMonitor) api.ActivityProcessor {
			return &TestActivityProcessor{monitor}
		},
	}
	natsagent.LaunchAgent(shutdown, config)
}

type TestActivityProcessor struct {
	monitor system.LogMonitor
}

func (t TestActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	if ctx.Discriminator() == api.DisposeDiscriminator {
		// a disposal request
		t.monitor.Infof("Processed dispose")
		return api.ActivityResult{Result: api.ActivityResultComplete}
	}
	ctx.SetOutputValue("agent.test.output", "test output")
	t.monitor.Infof("Processed deploy")
	return api.ActivityResult{Result: api.ActivityResultComplete}
}
