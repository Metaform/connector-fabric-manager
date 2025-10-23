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

package common

import (
	"context"
	"fmt"
	"time"

	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/spf13/viper"
)

const (
	uriKey    = "uri"
	bucketKey = "bucket"
	streamKey = "stream"
	timeout   = 10 * time.Second
)

func Launch(shutdown <-chan struct{}, aConfig AgentConfig, assemblies ...system.ServiceAssembly) {
	mode := runtime.LoadMode()

	monitor := runtime.LoadLogMonitor(aConfig.LogPrefix, mode)
	//goland:noinspection GoUnhandledErrorResult
	defer monitor.Sync()

	assembler := system.NewServiceAssembler(monitor, aConfig.VConfig, mode)
	for _, assembly := range assemblies {
		assembler.Register(assembly)
	}
	runtime.AssembleAndLaunch(assembler, aConfig.Name, monitor, shutdown)
}

func LoadAgentConfig(name string, logPrefix string, configPrefix string) *AgentConfig {
	vConfig := system.LoadConfigOrPanic(configPrefix)
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
	return &AgentConfig{
		Name:       name,
		LogPrefix:  logPrefix,
		VConfig:    vConfig,
		URI:        uri,
		Bucket:     bucketValue,
		StreamName: streamValue,
	}
}

type AgentConfig struct {
	Name       string
	LogPrefix  string
	VConfig    *viper.Viper
	URI        string
	Bucket     string
	StreamName string
}

func SetupConsumer(natsClient *natsclient.NatsClient, streamName string, activityType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stream, err := natsclient.SetupStream(ctx, natsClient, streamName)

	if err != nil {
		return fmt.Errorf("error setting up agent stream: %w", err)
	}

	_, err = natsclient.SetupConsumer(ctx, stream, activityType)

	if err != nil {
		return fmt.Errorf("error setting up agent consumer: %w", err)
	}

	return nil
}
