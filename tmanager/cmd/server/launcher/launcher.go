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
	"github.com/metaform/connector-fabric-manager/tmanager/tmcore"
	"github.com/metaform/connector-fabric-manager/tmanager/tmhandler"
)

const (
	defaultPort  = 8080
	configPrefix = "tmconfig"
	key          = "httpPort"
)

func Launch() {
	mode := runtime.LoadMode()

	logMonitor := runtime.LoadLogMonitor(mode)
	//goland:noinspection GoUnhandledErrorResult
	defer logMonitor.Sync()

	vConfig := config.LoadConfigOrPanic(configPrefix)
	vConfig.SetDefault(key, defaultPort)

	tManager := tmcore.NewTenantManager(logMonitor, vConfig, mode)
	tManager.ServiceAssembler.Register(&routing.RouterServiceAssembly{})
	tManager.ServiceAssembler.Register(&tmhandler.HandlerServiceAssembly{})

	runtime.AssembleAndLaunch(&tManager.ServiceAssembler, "Tenant Manager", vConfig, logMonitor)

}
