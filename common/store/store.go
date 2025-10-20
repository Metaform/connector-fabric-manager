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
	"github.com/metaform/connector-fabric-manager/common/types"
)

const (
	TransactionContextKey system.ServiceType = "store:TransactionContext"
)

// TransactionContext defines an interface for managing transactional operations.
type TransactionContext interface {
	Execute(ctx context.Context, callback func(ctx context.Context) error) error
}

// TrxFunc represents a transactional function wrapper allowing execution with a TransactionContext.
// For example:
//
//	store.Trx(ctx).AndReturn(ctx, func(ctx context.Context) (*MyType, error) {
//		return store.FindById(ctx, "my-id")
//	})
type TrxFunc[T any] struct {
	ctx TransactionContext
}

func Trx[T any](ctx TransactionContext) TrxFunc[T] {
	return TrxFunc[T]{ctx: ctx}
}

func (tf TrxFunc[T]) AndReturn(ctx context.Context, callback func(context.Context) (*T, error)) (*T, error) {
	var result *T
	var callbackErr error

	err := tf.ctx.Execute(ctx, func(ctx context.Context) error {
		result, callbackErr = callback(ctx)
		return callbackErr
	})

	if err != nil {
		return nil, err
	}

	return result, callbackErr
}

type NoOpTransactionContext struct{}

func (n NoOpTransactionContext) Execute(ctx context.Context, callback func(ctx context.Context) error) error {
	return callback(ctx)
}

var ErrNotFound = &types.BadRequestError{Message: "not found"}

type NoOpTrxAssembly struct {
	system.DefaultServiceAssembly
}

func (n NoOpTrxAssembly) Name() string {
	return "NoOpTrxAssembly"
}

func (n NoOpTrxAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{TransactionContextKey}
}

func (n *NoOpTrxAssembly) Init(context *system.InitContext) error {
	context.Registry.Register(TransactionContextKey, NoOpTransactionContext{})
	return nil
}
