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

package tmcore

import (
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type TManager struct {
	ServiceAssembler system.ServiceAssembler
}

func NewTenantManager(logger *zap.Logger, viper *viper.Viper, mode system.RuntimeMode) *TManager {
	return &TManager{
		ServiceAssembler: *system.NewServiceAssembler(logger, viper, mode),
	}
}
