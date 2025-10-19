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

package tmstore

import (
	"github.com/metaform/connector-fabric-manager/common/system"
)

type InMemoryServiceAssembly struct {
	system.DefaultServiceAssembly
}

func (a *InMemoryServiceAssembly) Name() string {
	return "Tenant Manager In-Memory Store"
}

func (a *InMemoryServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{TManagerStoreKey}
}

func (a *InMemoryServiceAssembly) Init(context *system.InitContext) error {
	context.Registry.Register(TManagerStoreKey, NewInMemoryTManagerStore(true))
	return nil
}
