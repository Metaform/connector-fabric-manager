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
	"github.com/metaform/connector-fabric-manager/pmanager/pmcore"
	"github.com/metaform/connector-fabric-manager/tmanager/tmhandler"
)

const (
	defaultPort  = 8181
	configPrefix = "pmconfig"
	key          = "httpPort"
)

func Launch() {
	mode := runtime.LoadMode()

	logMonitor := runtime.LoadLogMonitor(mode)
	//goland:noinspection GoUnhandledErrorResult
	defer logMonitor.Sync()

	vConfig := config.LoadConfigOrPanic(configPrefix)
	vConfig.SetDefault(key, defaultPort)

	pManager := pmcore.NewProvisioningManager(logMonitor, vConfig, mode)
	pManager.ServiceAssembler.Register(&routing.RouterServiceAssembly{})
	pManager.ServiceAssembler.Register(&tmhandler.HandlerServiceAssembly{})

	runtime.AssembleAndLaunch(&pManager.ServiceAssembler, "Provisioning Manager", vConfig, logMonitor)
}
