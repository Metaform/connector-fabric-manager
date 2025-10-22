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

package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockResponseWriter implements http.ResponseWriter for testing
type mockResponseWriter struct {
	headers    http.Header
	statusCode int
	body       *bytes.Buffer
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{
		headers: make(http.Header),
		body:    &bytes.Buffer{},
	}
}

func (m *mockResponseWriter) Header() http.Header {
	return m.headers
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	return m.body.Write(data)
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

func TestWriteErrorWithID(t *testing.T) {
	t.Run("writes error response with ID", func(t *testing.T) {
		w := newMockResponseWriter()
		handler := HttpHandler{Monitor: system.NoopMonitor{}}
		handler.WriteErrorWithID(w, "Test error", http.StatusBadRequest, "error-123")

		assert.Equal(t, http.StatusBadRequest, w.statusCode)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response ErrorResponse
		err := json.Unmarshal(w.body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Bad Request", response.Error)
		assert.Equal(t, "Test error", response.Message)
		assert.Equal(t, 400, response.Code)
		assert.Equal(t, "error-123", response.ID)
	})

	t.Run("writes error response without ID", func(t *testing.T) {
		w := newMockResponseWriter()

		handler := HttpHandler{Monitor: system.NoopMonitor{}}
		handler.WriteErrorWithID(w, "Server error", http.StatusInternalServerError, "")

		assert.Equal(t, http.StatusInternalServerError, w.statusCode)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response ErrorResponse
		err := json.Unmarshal(w.body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Internal Server Error", response.Error)
		assert.Equal(t, "Server error", response.Message)
		assert.Equal(t, 500, response.Code)
		assert.Empty(t, response.ID)
	})

	t.Run("handles different status codes", func(t *testing.T) {
		testCases := []struct {
			statusCode    int
			expectedError string
		}{
			{http.StatusNotFound, "Not Found"},
			{http.StatusUnauthorized, "Unauthorized"},
			{http.StatusForbidden, "Forbidden"},
		}

		handler := HttpHandler{Monitor: system.NoopMonitor{}}

		for _, tc := range testCases {
			w := newMockResponseWriter()
			handler.WriteErrorWithID(w, "Test message", tc.statusCode, "test-id")

			assert.Equal(t, tc.statusCode, w.statusCode)

			var response ErrorResponse
			err := json.Unmarshal(w.body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedError, response.Error)
		}
	})
}

func TestWriteError(t *testing.T) {
	t.Run("delegates to WriteErrorWithID with empty ID", func(t *testing.T) {
		w := newMockResponseWriter()
		handler := HttpHandler{Monitor: system.NoopMonitor{}}

		handler.WriteError(w, "Test message", http.StatusInternalServerError)

		assert.Equal(t, http.StatusInternalServerError, w.statusCode)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response ErrorResponse
		err := json.Unmarshal(w.body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Internal Server Error", response.Error)
		assert.Equal(t, "Test message", response.Message)
		assert.Equal(t, 500, response.Code)
		assert.Empty(t, response.ID)
	})
}
