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
	"github.com/metaform/connector-fabric-manager/assembly/routing"
	"github.com/metaform/connector-fabric-manager/common/config"
	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/tmanager/natsdeployment"
	"github.com/metaform/connector-fabric-manager/tmanager/tmcore"
	"github.com/metaform/connector-fabric-manager/tmanager/tmhandler"
	"github.com/metaform/connector-fabric-manager/tmanager/tmstore"
)

const (
	logPrefix    = "tmanager"
	defaultPort  = 8080
	configPrefix = "tm"
	key          = "httpPort"

	uriKey    = "uri"
	bucketKey = "bucket"
	streamKey = "stream"
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
	vConfig.SetDefault(key, defaultPort)

	uri := vConfig.GetString(uriKey)
	bucketValue := vConfig.GetString(bucketKey)
	streamValue := vConfig.GetString(streamKey)

	assembler := system.NewServiceAssembler(logMonitor, vConfig, mode)
	assembler.Register(&routing.RouterServiceAssembly{})
	assembler.Register(&tmhandler.HandlerServiceAssembly{})
	assembler.Register(&tmcore.TMCoreServiceAssembly{})
	assembler.Register(&store.NoOpTrxAssembly{})
	assembler.Register(&tmstore.InMemoryServiceAssembly{})
	assembler.Register(natsdeployment.NewNatsDeploymentServiceAssembly(uri, bucketValue, streamValue))

	runtime.AssembleAndLaunch(assembler, "Tenant Manager", logMonitor, shutdown)
}
