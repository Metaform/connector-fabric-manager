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
	"time"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsagent"
)

const (
	ActivityType = "dns-activity"
	key          = "dns.count"
)

func LaunchAndWaitSignal(shutdown <-chan struct{}) {
	config := natsagent.LauncherConfig{
		AgentName:    "DNS Agent",
		ConfigPrefix: "dnsagent",
		ActivityType: ActivityType,
		NewProcessor: func(ctx *natsagent.AgentContext) api.ActivityProcessor {
			return &DNSActivityProcessor{ctx.Monitor}
		},
	}
	natsagent.LaunchAgent(shutdown, config)
}

type DNSActivityProcessor struct {
	monitor system.LogMonitor
}

func (t DNSActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	identifier, found := ctx.Value(model.ParticipantIdentifier)
	if !found {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("missing participant identifier")}
	}
	if ctx.Discriminator() == api.DisposeDiscriminator {
		t.monitor.Infof("DNS disposal complete: %s", identifier)
		return api.ActivityResult{Result: api.ActivityResultComplete}
	}
	count, found := ctx.Value(key)
	if (found) && (count.(float64) > 0) {
		t.monitor.Infof("DNS provisioning complete: %s", identifier)
		ctx.Delete(key)
		return api.ActivityResult{Result: api.ActivityResultComplete}
	}
	t.monitor.Infof("DNS provisioning requested: %s", identifier)
	ctx.SetValue(key, 1)
	return api.ActivityResult{Result: api.ActivityResultSchedule, WaitOnReschedule: 1 * time.Second}
}
