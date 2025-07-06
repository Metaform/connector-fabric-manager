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

package store

import (
	"context"
	"github.com/metaform/connector-fabric-manager/common/system"
)

const (
	TransactionContextKey system.ServiceType = "store:TransactionContext"
)

// TransactionContext defines an interface for managing transactional operations.
type TransactionContext interface {
	Execute(ctx context.Context, callback func(ctx context.Context) error) error
}

type NoOpTransactionContext struct{}

func (n NoOpTransactionContext) Execute(ctx context.Context, callback func(ctx context.Context) error) error {
	return callback(ctx)
}
