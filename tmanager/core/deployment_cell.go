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

package core

import (
	"context"

	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

type cellDeployer struct {
	trxContext store.TransactionContext
	store      api.EntityStore[api.Cell]
}

func (d cellDeployer) RecordExternalDeployment(ctx context.Context, cell api.Cell) (*api.Cell, error) {
	return store.Tx[api.Cell](d.trxContext).AndReturn(ctx, func(ctx context.Context) (*api.Cell, error) {
		return d.store.Create(ctx, &cell)
	})
}
