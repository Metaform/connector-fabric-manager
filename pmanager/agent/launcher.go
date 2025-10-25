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

package agent

import (
	"fmt"

	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/spf13/viper"
)

const (
	uriKey    = "uri"
	bucketKey = "bucket"
	streamKey = "stream"
)

type LauncherConfig struct {
	AgentName    string
	ConfigPrefix string
	ActivityType string
	NewProcessor func(monitor system.LogMonitor) api.ActivityProcessor
}

type agentConfig struct {
	Name       string
	URI        string
	Bucket     string
	StreamName string
	VConfig    *viper.Viper
}

func LaunchAgent(shutdown <-chan struct{}, config LauncherConfig) {
	cfg := loadAgentConfig(config.AgentName, config.ConfigPrefix)

	assembly := &AgentServiceAssembly{
		agentName:    config.AgentName,
		activityType: config.ActivityType,
		uri:          cfg.URI,
		bucket:       cfg.Bucket,
		streamName:   cfg.StreamName,
		newProcessor: config.NewProcessor,
	}

	mode := runtime.LoadMode()

	monitor := runtime.LoadLogMonitor(config.ConfigPrefix, mode)
	//goland:noinspection GoUnhandledErrorResult
	defer monitor.Sync()

	assembler := system.NewServiceAssembler(monitor, cfg.VConfig, mode)
	assembler.Register(assembly)
	runtime.AssembleAndLaunch(assembler, cfg.Name, monitor, shutdown)
}

func loadAgentConfig(name string, configPrefix string) *agentConfig {
	vConfig := system.LoadConfigOrPanic(configPrefix)
	uri := vConfig.GetString(uriKey)
	bucketValue := vConfig.GetString(bucketKey)
	streamValue := vConfig.GetString(streamKey)

	err := runtime.CheckRequiredParams(
		fmt.Sprintf("%s.%s", configPrefix, uriKey), uri,
		fmt.Sprintf("%s.%s", configPrefix, bucketKey), bucketValue,
		fmt.Sprintf("%s.%s", configPrefix, streamKey), streamValue)
	if err != nil {
		panic(fmt.Errorf("error loading agent configuration: %w", err))
	}
	return &agentConfig{
		Name:       name,
		URI:        uri,
		Bucket:     bucketValue,
		StreamName: streamValue,
		VConfig:    vConfig,
	}
}
