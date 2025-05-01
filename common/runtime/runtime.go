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

package runtime

import (
	"flag"
	"fmt"
	"github.com/metaform/connector-fabric-manager/common/system"
	"go.uber.org/zap"
)

const mode = "mode"

func LoadLogger(mode system.RuntimeMode) *zap.Logger {
	var logger *zap.Logger
	var err error
	if mode == system.DevelopmentMode || mode == system.DebugMode {
		logger, err = zap.NewDevelopment(zap.WithCaller(false))
	} else {
		logger, err = zap.NewProduction(zap.WithCaller(false))
	}
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}
	return logger
}

func LoadMode() system.RuntimeMode {
	modeFlag := flag.String(mode, system.ProductionMode, "Runtime mode: development, production, or debug")
	flag.Parse()

	mode, err := system.ParseRuntimeMode(*modeFlag)
	if err != nil {
		panic(fmt.Errorf("error parsing runtime mode: %w", err))
	}
	return mode
}
