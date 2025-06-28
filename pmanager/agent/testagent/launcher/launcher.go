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
	"github.com/metaform/connector-fabric-manager/common/config"
	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/common/system"
)

const (
	configPrefix = "testagent"
)

func Launch() {
	mode := runtime.LoadMode()

	logMonitor := runtime.LoadLogMonitor(mode)
	//goland:noinspection GoUnhandledErrorResult
	defer logMonitor.Sync()

	vConfig := config.LoadConfigOrPanic(configPrefix)

	assembler := system.NewServiceAssembler(logMonitor, vConfig, mode)
	assembler.Register(&testAgenServiceAssemby{})
	runtime.AssembleAndLaunch(assembler, "Test Agent", vConfig, logMonitor)
}

type testAgenServiceAssemby struct {
	system.DefaultServiceAssembly
}

func (t testAgenServiceAssemby) Name() string {
	return "Test Agent"
}
