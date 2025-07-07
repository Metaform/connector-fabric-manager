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
	"github.com/metaform/connector-fabric-manager/assembly/httpclient"
	"github.com/metaform/connector-fabric-manager/assembly/routing"
	"github.com/metaform/connector-fabric-manager/common/config"
	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/memorystore"
	"github.com/metaform/connector-fabric-manager/pmanager/pmhandler"
)

const (
	defaultPort  = 8181
	configPrefix = "pmconfig"
	httpKey      = "httpPort"
	storeKey     = "sql"
)

func LaunchAndWaitSignal() {
	Launch(runtime.CreateSignalShutdownChan())
}

func Launch(shutdown <-chan struct{}) {
	mode := runtime.LoadMode()

	logMonitor := runtime.LoadLogMonitor(mode)
	//goland:noinspection GoUnhandledErrorResult
	defer logMonitor.Sync()

	vConfig := config.LoadConfigOrPanic(configPrefix)
	vConfig.SetDefault(httpKey, defaultPort)

	assembler := system.NewServiceAssembler(logMonitor, vConfig, mode)
	if vConfig.IsSet(storeKey) {
		// TODO add SQL assembly
		panic("SQL storage not yet implemented")
	} else {
		assembler.Register(&memorystore.MemoryStoreServiceAssembly{})
	}

	assembler.Register(&httpclient.HttpClientServiceAssembly{})
	assembler.Register(&routing.RouterServiceAssembly{})
	assembler.Register(&pmhandler.HandlerServiceAssembly{})

	runtime.AssembleAndLaunch(assembler, "Provision Manager", logMonitor, shutdown)
}
