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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestModel struct {
	ID   string
	Name string
}

func TestAndReturn_Success(t *testing.T) {
	ctx := context.Background()
	expectedModel := &TestModel{ID: "test-id", Name: "test-name"}
	noOpTrx := NoOpTransactionContext{}

	result, err := Trx[TestModel](noOpTrx).AndReturn(ctx, func(ctx context.Context) (*TestModel, error) {
		return expectedModel, nil
	})

	require.NoError(t, err)
	assert.Equal(t, expectedModel, result)
}

func TestAndReturn_CallbackError(t *testing.T) {
	ctx := context.Background()
	callbackErr := errors.New("callback error")
	noOpTrx := NoOpTransactionContext{}

	result, err := Trx[TestModel](noOpTrx).AndReturn(ctx, func(ctx context.Context) (*TestModel, error) {
		return nil, callbackErr
	})

	require.Error(t, err)
	assert.Equal(t, callbackErr, err)
	assert.Nil(t, result)
}

func TestAndReturn_NilResult(t *testing.T) {
	ctx := context.Background()
	noOpTrx := NoOpTransactionContext{}

	result, err := Trx[TestModel](noOpTrx).AndReturn(ctx, func(ctx context.Context) (*TestModel, error) {
		return nil, nil
	})

	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestAndReturn_ContextPropagation(t *testing.T) {
	parentCtx := context.Background()
	ctxKey := "test-key"
	ctxValue := "test-value"
	ctx := context.WithValue(parentCtx, ctxKey, ctxValue)
	noOpTrx := NoOpTransactionContext{}

	_, err := Trx[TestModel](noOpTrx).AndReturn(ctx, func(ctx context.Context) (*TestModel, error) {
		assert.Equal(t, ctxValue, ctx.Value(ctxKey))
		return &TestModel{}, nil
	})

	require.NoError(t, err)
}

func TestAndReturn_WithStructPointer(t *testing.T) {
	ctx := context.Background()
	noOpTrx := NoOpTransactionContext{}

	originalModel := TestModel{ID: "original-id", Name: "original-name"}

	result, err := Trx[TestModel](noOpTrx).AndReturn(ctx, func(ctx context.Context) (*TestModel, error) {
		// Simulate modifying the model
		model := originalModel
		model.Name = "modified-name"
		return &model, nil
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "original-id", result.ID)
	assert.Equal(t, "modified-name", result.Name)
}

func TestNoOpTransactionContext_Execute(t *testing.T) {
	ctx := context.Background()
	noOpTrx := NoOpTransactionContext{}

	executed := false
	err := noOpTrx.Execute(ctx, func(ctx context.Context) error {
		executed = true
		return nil
	})

	require.NoError(t, err)
	assert.True(t, executed)
}

func TestNoOpTransactionContext_WithError(t *testing.T) {
	ctx := context.Background()
	noOpTrx := NoOpTransactionContext{}
	expectedErr := errors.New("callback error")

	err := noOpTrx.Execute(ctx, func(ctx context.Context) error {
		return expectedErr
	})

	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestTxFunction_Creation(t *testing.T) {
	noOpTrx := NoOpTransactionContext{}

	txFunc := Trx[TestModel](noOpTrx)

	assert.NotNil(t, txFunc)
	assert.Equal(t, noOpTrx, txFunc.ctx)
}

func TestAndReturn_CallbackReturnsNewInstance(t *testing.T) {
	ctx := context.Background()
	noOpTrx := NoOpTransactionContext{}

	result, err := Trx[TestModel](noOpTrx).AndReturn(ctx, func(ctx context.Context) (*TestModel, error) {
		// Create a new instance each time
		return &TestModel{
			ID:   "dynamic-id",
			Name: "dynamic-name",
		}, nil
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "dynamic-id", result.ID)
	assert.Equal(t, "dynamic-name", result.Name)
}
