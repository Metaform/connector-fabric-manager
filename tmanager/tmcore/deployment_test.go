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
	"context"
	"testing"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeploymentCallbackService_RegisterDeploymentHandler(t *testing.T) {

	t.Run("register handler", func(t *testing.T) {
		service := deploymentCallbackService{
			handlers: make(map[string]api.DeploymentCallbackHandler),
		}

		handler := func(ctx context.Context, response api.DeploymentResponse) error {
			return nil
		}

		service.RegisterDeploymentHandler("test-type", handler)

		require.Contains(t, service.handlers, "test-type")
		assert.NotNil(t, service.handlers["test-type"])
	})

	t.Run("register multiple handlers", func(t *testing.T) {
		service := deploymentCallbackService{
			handlers: make(map[string]api.DeploymentCallbackHandler),
		}

		handler1 := func(ctx context.Context, response api.DeploymentResponse) error {
			return nil
		}
		handler2 := func(ctx context.Context, response api.DeploymentResponse) error {
			return model.NewClientError("test error")
		}

		service.RegisterDeploymentHandler("type1", handler1)
		service.RegisterDeploymentHandler("type2", handler2)

		require.Contains(t, service.handlers, "type1")
		require.Contains(t, service.handlers, "type2")
		assert.Len(t, service.handlers, 2)
	})

	t.Run("overwrite existing handler", func(t *testing.T) {
		service := deploymentCallbackService{
			handlers: make(map[string]api.DeploymentCallbackHandler),
		}

		originalHandler := func(ctx context.Context, response api.DeploymentResponse) error {
			return model.NewClientError("original")
		}
		newHandler := func(ctx context.Context, response api.DeploymentResponse) error {
			return model.NewClientError("new")
		}

		service.RegisterDeploymentHandler("test-type", originalHandler)
		service.RegisterDeploymentHandler("test-type", newHandler)

		require.Contains(t, service.handlers, "test-type")
		assert.Len(t, service.handlers, 1)
	})
}

func TestDeploymentCallbackService_Dispatch(t *testing.T) {
	t.Run("dispatch to registered handler", func(t *testing.T) {
		service := deploymentCallbackService{
			handlers: make(map[string]api.DeploymentCallbackHandler),
		}

		var receivedResponse api.DeploymentResponse
		handler := func(ctx context.Context, response api.DeploymentResponse) error {
			receivedResponse = response
			return nil
		}

		service.RegisterDeploymentHandler("vpa", handler)

		response := api.DeploymentResponse{
			Success:        true,
			ErrorDetail:    "",
			ManifestID:     "manifest-123",
			DeploymentType: "vpa",
			Properties: map[string]any{
				"version": "1.0.0",
			},
		}

		err := service.Dispatch(context.Background(), response)

		require.NoError(t, err)
		assert.Equal(t, response, receivedResponse)
	})

	t.Run("dispatch with handler returning error", func(t *testing.T) {
		service := deploymentCallbackService{
			handlers: make(map[string]api.DeploymentCallbackHandler),
		}

		expectedError := model.NewClientError("deployment failed")
		handler := func(ctx context.Context, response api.DeploymentResponse) error {
			return expectedError
		}

		service.RegisterDeploymentHandler("vpa", handler)

		response := api.DeploymentResponse{
			Success:        false,
			ErrorDetail:    "connection timeout",
			DeploymentType: "vpa",
		}

		err := service.Dispatch(context.Background(), response)

		require.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("dispatch to unregistered deployment type", func(t *testing.T) {
		service := deploymentCallbackService{
			handlers: make(map[string]api.DeploymentCallbackHandler),
		}

		response := api.DeploymentResponse{
			DeploymentType: "nonexistent-type",
		}

		err := service.Dispatch(context.Background(), response)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "deployment handler not found for type: nonexistent-type")

		require.True(t, model.IsFatal(err))

	})

	t.Run("dispatch with context cancellation", func(t *testing.T) {
		service := deploymentCallbackService{
			handlers: make(map[string]api.DeploymentCallbackHandler),
		}

		handler := func(ctx context.Context, response api.DeploymentResponse) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		}

		service.RegisterDeploymentHandler("vpa", handler)

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		response := api.DeploymentResponse{
			DeploymentType: "vpa",
		}

		err := service.Dispatch(ctx, response)

		require.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})

}

func TestDeploymentCallbackService_Integration(t *testing.T) {
	t.Run("multiple handlers", func(t *testing.T) {
		service := deploymentCallbackService{
			handlers: make(map[string]api.DeploymentCallbackHandler),
		}

		var connectorCalls int
		var dataspaceCalls int

		connectorHandler := func(ctx context.Context, response api.DeploymentResponse) error {
			connectorCalls++
			return nil
		}

		dataspaceHandler := func(ctx context.Context, response api.DeploymentResponse) error {
			dataspaceCalls++
			return model.NewRecoverableError("temporary failure")
		}

		service.RegisterDeploymentHandler("vpa", connectorHandler)
		service.RegisterDeploymentHandler("dprofile", dataspaceHandler)

		connectorResponse := api.DeploymentResponse{DeploymentType: "vpa"}
		err := service.Dispatch(context.Background(), connectorResponse)
		require.NoError(t, err)

		dataspaceResponse := api.DeploymentResponse{DeploymentType: "dprofile"}
		err = service.Dispatch(context.Background(), dataspaceResponse)
		require.Error(t, err)

		// Verify only correct handlers were called
		assert.Equal(t, 1, connectorCalls)
		assert.Equal(t, 1, dataspaceCalls)

		require.True(t, model.IsRecoverable(err))
	})

}
