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
	"context"
	"time"

	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/mvd/common"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsorchestration"
)

const (
	ActivityType = "dns-activity"
)

func LaunchAndWaitSignal(shutdown <-chan struct{}) {
	aConfig := common.LoadAgentConfig("DNS Agent", "dns-agent", "dnsagent")
	common.Launch(shutdown, *aConfig, &dnsAgentServiceAssembly{
		uri:                    aConfig.URI,
		bucket:                 aConfig.Bucket,
		streamName:             aConfig.StreamName,
		DefaultServiceAssembly: system.DefaultServiceAssembly{},
	})
}

type dnsAgentServiceAssembly struct {
	uri        string
	bucket     string
	streamName string
	system.DefaultServiceAssembly
}

func (d dnsAgentServiceAssembly) Name() string {
	return "DNS Agent"
}

func (d dnsAgentServiceAssembly) Start(startCtx *system.StartContext) error {
	natsClient, err := natsclient.NewNatsClient(d.uri, d.bucket)
	if err != nil {
		return err
	}

	if err = common.SetupConsumer(natsClient, d.streamName, ActivityType); err != nil {
		return err
	}

	executor := &natsorchestration.NatsActivityExecutor{
		Client:            natsclient.NewMsgClient(natsClient),
		StreamName:        d.streamName,
		ActivityType:      ActivityType,
		ActivityProcessor: DNSActivityProcessor{startCtx.LogMonitor},
		Monitor:           system.NoopMonitor{},
	}

	ctx := context.Background()
	return executor.Execute(ctx)
}

type DNSActivityProcessor struct {
	monitor system.LogMonitor
}

func (t DNSActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	count, found := ctx.Value("dns.count")
	if (found) && (count.(float64) > 0) {
		t.monitor.Infof("DNS provisioning complete")
		ctx.Delete("dns.count")
		return api.ActivityResult{Result: api.ActivityResultComplete}
	}
	t.monitor.Infof("DNS provisioning requested")
	ctx.SetValue("dns.count", 1)
	return api.ActivityResult{Result: api.ActivityResultSchedule, WaitOnReschedule: 1 * time.Second}
}
