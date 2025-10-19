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

package _type

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewRecoverableWrappedError(t *testing.T) {
	tests := []struct {
		name     string
		cause    error
		message  string
		args     []any
		expected string
	}{
		{
			name:     "with cause and simple message",
			cause:    errors.New("original error"),
			message:  "temporary failure",
			args:     nil,
			expected: "temporary failure: original error",
		},
		{
			name:     "with cause and formatted message",
			cause:    errors.New("connection failed"),
			message:  "retry attempt %d failed",
			args:     []any{3},
			expected: "retry attempt 3 failed: connection failed",
		},
		{
			name:     "with nil cause",
			cause:    nil,
			message:  "temporary failure",
			args:     nil,
			expected: "temporary failure",
		},
		{
			name:     "with nested wrapped error",
			cause:    fmt.Errorf("wrapped: %w", errors.New("root cause")),
			message:  "retry failed",
			args:     nil,
			expected: "retry failed: wrapped: root cause",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewRecoverableWrappedError(tt.cause, tt.message, tt.args...)

			if err.Error() != tt.expected {
				t.Errorf("NewRecoverableWrappedError() = %v, want %v", err.Error(), tt.expected)
			}

			if !IsRecoverable(err) {
				t.Error("NewRecoverableWrappedError() should return a RecoverableError")
			}

			// Test unwrapping
			if tt.cause != nil && errors.Unwrap(err) != tt.cause {
				t.Error("NewRecoverableWrappedError() should preserve the cause")
			}
		})
	}
}

func TestNewClientWrappedError(t *testing.T) {
	tests := []struct {
		name     string
		cause    error
		message  string
		args     []any
		expected string
	}{
		{
			name:     "with cause and simple message",
			cause:    errors.New("validation failed"),
			message:  "invalid request",
			args:     nil,
			expected: "invalid request: validation failed",
		},
		{
			name:     "with cause and formatted message",
			cause:    errors.New("missing field"),
			message:  "field %s is required",
			args:     []any{"email"},
			expected: "field email is required: missing field",
		},
		{
			name:     "with nil cause",
			cause:    nil,
			message:  "bad request",
			args:     nil,
			expected: "bad request",
		},
		{
			name:     "with nested wrapped error",
			cause:    fmt.Errorf("validation: %w", errors.New("empty value")),
			message:  "input validation failed",
			args:     nil,
			expected: "input validation failed: validation: empty value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewClientWrappedError(tt.cause, tt.message, tt.args...)

			if err.Error() != tt.expected {
				t.Errorf("NewClientWrappedError() = %v, want %v", err.Error(), tt.expected)
			}

			if !IsClientError(err) {
				t.Error("NewClientWrappedError() should return a ClientError")
			}

			// Test unwrapping
			if tt.cause != nil && errors.Unwrap(err) != tt.cause {
				t.Error("NewClientWrappedError() should preserve the cause")
			}
		})
	}
}

func TestNewFatalWrappedError(t *testing.T) {
	tests := []struct {
		name     string
		cause    error
		message  string
		args     []any
		expected string
	}{
		{
			name:     "with cause and simple message",
			cause:    errors.New("database connection lost"),
			message:  "system error",
			args:     nil,
			expected: "system error: database connection lost",
		},
		{
			name:     "with cause and formatted message",
			cause:    errors.New("out of memory"),
			message:  "critical failure in %s",
			args:     []any{"authentication service"},
			expected: "critical failure in authentication service: out of memory",
		},
		{
			name:     "with nil cause",
			cause:    nil,
			message:  "fatal error",
			args:     nil,
			expected: "fatal error",
		},
		{
			name:     "with nested wrapped error",
			cause:    fmt.Errorf("startup: %w", errors.New("config missing")),
			message:  "service initialization failed",
			args:     nil,
			expected: "service initialization failed: startup: config missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewFatalWrappedError(tt.cause, tt.message, tt.args...)

			if err.Error() != tt.expected {
				t.Errorf("NewFatalWrappedError() = %v, want %v", err.Error(), tt.expected)
			}

			var fatalErr FatalError
			if !errors.As(err, &fatalErr) {
				t.Error("NewFatalWrappedError() should return a FatalError")
			}

			// Test unwrapping
			if tt.cause != nil && errors.Unwrap(err) != tt.cause {
				t.Error("NewFatalWrappedError() should preserve the cause")
			}
		})
	}
}

func TestWrappedErrorsWithErrorsIs(t *testing.T) {
	originalErr := errors.New("original error")

	recoverableErr := NewRecoverableWrappedError(originalErr, "recoverable wrapper")
	clientErr := NewClientWrappedError(originalErr, "client wrapper")
	fatalErr := NewFatalWrappedError(originalErr, "fatal wrapper")

	if !errors.Is(recoverableErr, originalErr) {
		t.Error("errors.Is should find the original error in recoverable wrapped error")
	}

	if !errors.Is(clientErr, originalErr) {
		t.Error("errors.Is should find the original error in client wrapped error")
	}

	if !errors.Is(fatalErr, originalErr) {
		t.Error("errors.Is should find the original error in fatal wrapped error")
	}
}
