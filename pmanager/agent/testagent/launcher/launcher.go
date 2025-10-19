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
	"fmt"
	"time"

	"github.com/metaform/connector-fabric-manager/common/config"
	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsorchestration"
)

const (
	logPrefix    = "test-agent"
	activityType = "test.activity"
	configPrefix = "testagent"
	uriKey       = "uri"
	bucketKey    = "bucket"
	streamKey    = "stream"
	timeout      = 10 * time.Second
)

func LaunchAndWaitSignal() {
	Launch(runtime.CreateSignalShutdownChan())
}

func Launch(shutdown <-chan struct{}) {
	mode := runtime.LoadMode()

	logMonitor := runtime.LoadLogMonitor(logPrefix, mode)
	//goland:noinspection GoUnhandledErrorResult
	defer logMonitor.Sync()

	vConfig := config.LoadConfigOrPanic(configPrefix)

	assembler := system.NewServiceAssembler(logMonitor, vConfig, mode)

	uri := vConfig.GetString(uriKey)
	bucketValue := vConfig.GetString(bucketKey)
	streamValue := vConfig.GetString(streamKey)

	err := runtime.CheckRequiredParams(
		fmt.Sprintf("%s.%s", configPrefix, uriKey), uri,
		fmt.Sprintf("%s.%s", configPrefix, bucketKey), bucketValue,
		fmt.Sprintf("%s.%s", configPrefix, streamKey), streamValue)
	if err != nil {
		panic(fmt.Errorf("error launching test agent: %w", err))
	}

	assembler.Register(&testAgentServiceAssembly{uri: uri, bucket: bucketValue, streamName: streamValue})
	runtime.AssembleAndLaunch(assembler, "Test Agent", logMonitor, shutdown)
}

type testAgentServiceAssembly struct {
	uri        string
	bucket     string
	streamName string
	system.DefaultServiceAssembly
}

func (t testAgentServiceAssembly) Name() string {
	return "Test Agent"
}

func (t testAgentServiceAssembly) Start(startCtx *system.StartContext) error {

	natsClient, err := natsclient.NewNatsClient(t.uri, t.bucket)
	if err != nil {
		return err
	}

	if err = SetupConsumer(natsClient, t.streamName); err != nil {
		return err
	}

	executor := &natsorchestration.NatsActivityExecutor{
		Client:            natsclient.NewMsgClient(natsClient),
		StreamName:        t.streamName,
		ActivityType:      activityType,
		ActivityProcessor: TestActivityProcessor{startCtx.LogMonitor},
		Monitor:           system.NoopMonitor{},
	}

	ctx := context.Background()
	return executor.Execute(ctx)
}

type TestActivityProcessor struct {
	monitor system.LogMonitor
}

func (t TestActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	t.monitor.Infof("Processed activity")
	return api.ActivityResult{Result: api.ActivityResultComplete}
}

func SetupConsumer(natsClient *natsclient.NatsClient, streamName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stream, err := natsclient.SetupStream(ctx, natsClient, streamName)

	if err != nil {
		return fmt.Errorf("error setting up NATS test agent stream: %w", err)
	}

	_, err = natsclient.SetupConsumer(ctx, stream, activityType)

	if err != nil {
		return fmt.Errorf("error setting up NATS test agent consumer: %w", err)
	}

	return nil
}
