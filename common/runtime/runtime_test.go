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
	"strings"
	"testing"
)

func TestCheckRequiredParams(t *testing.T) {
	tests := []struct {
		name        string
		params      []interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name:        "empty parameters should pass",
			params:      []interface{}{},
			expectError: false,
		},
		{
			name:        "valid pair of parameters should pass",
			params:      []interface{}{"key", "value"},
			expectError: false,
		},
		{
			name:        "multiple valid pairs should pass",
			params:      []interface{}{"key1", "value1", "key2", "value2"},
			expectError: false,
		},
		{
			name:        "odd number of parameters should fail",
			params:      []interface{}{"key"},
			expectError: true,
			errorMsg:    "arguments must be even, got 1",
		},
		{
			name:        "three parameters should fail",
			params:      []interface{}{"key1", "value1", "key2"},
			expectError: true,
			errorMsg:    "arguments must be even, got 3",
		},
		{
			name:        "nil value at odd index should fail",
			params:      []interface{}{"key", nil},
			expectError: true,
			errorMsg:    "key not specified",
		},
		{
			name:        "nil value in second pair should fail",
			params:      []interface{}{"key1", "value1", "key2", nil},
			expectError: true,
			errorMsg:    "key2 not specified",
		},
		{
			name:        "multiple nil values should fail with multiple errors",
			params:      []interface{}{"key1", nil, "key2", nil},
			expectError: true,
			errorMsg:    "key1 not specified, key2 not specified",
		},
		{
			name:        "nil key at even index is allowed",
			params:      []interface{}{nil, "value"},
			expectError: false,
		},
		{
			name:        "mix of valid and invalid parameters",
			params:      []interface{}{"key1", "value1", "key2", nil, "key3", "value3"},
			expectError: true,
			errorMsg:    "key2 not specified",
		},
		{
			name:        "empty string values are not allowed",
			params:      []interface{}{"key1", 0, "key2", "", "key3", false},
			expectError: true,
			errorMsg:    "key2 is empty",
		},
		{
			name:        "zero integer value is allowed",
			params:      []interface{}{"key", 0},
			expectError: false,
		},
		{
			name:        "false boolean value is allowed",
			params:      []interface{}{"key", false},
			expectError: false,
		},
		{
			name:        "multiple nil keys with valid values",
			params:      []interface{}{nil, "value1", nil, "value2"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckRequiredParams(tt.params...)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}

				if !strings.Contains(err.Error(), "missing parameters:") {
					t.Errorf("expected error to contain 'missing parameters:', got: %v", err)
				}

				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error to contain '%s', got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateParamsOddNumberCheck(t *testing.T) {
	// Test specifically for odd number validation
	tests := []struct {
		name   string
		params []interface{}
		count  int
	}{
		{"single param", []interface{}{"key"}, 1},
		{"three params", []interface{}{"key1", "value1", "key2"}, 3},
		{"five params", []interface{}{"key1", "value1", "key2", "value2", "key3"}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckRequiredParams(tt.params...)

			if err == nil {
				t.Error("expected error for odd number of parameters")
				return
			}

			expectedMsg := "arguments must be even, got " + string(rune(tt.count+'0'))
			if !strings.Contains(err.Error(), expectedMsg) {
				t.Errorf("expected error to contain '%s', got: %v", expectedMsg, err)
			}
		})
	}
}

func TestValidateParamsNilValueCheck(t *testing.T) {
	// Test specifically for nil value validation at odd indices
	tests := []struct {
		name     string
		params   []interface{}
		errorMsg string
	}{
		{
			name:     "nil at index 1",
			params:   []interface{}{"key", nil},
			errorMsg: "key not specified",
		},
		{
			name:     "nil at index 3",
			params:   []interface{}{"key1", "value1", "key2", nil},
			errorMsg: "key2 not specified",
		},
		{
			name:     "nil at index 5",
			params:   []interface{}{"key1", "value1", "key2", "value2", "key3", nil},
			errorMsg: "key3 not specified",
		},
		{
			name:     "multiple nils at odd indices",
			params:   []interface{}{"key1", nil, "key2", nil, "key3", nil},
			errorMsg: "key1 not specified, key2 not specified, key3 not specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckRequiredParams(tt.params...)

			if err == nil {
				t.Error("expected error for nil values")
				return
			}

			if !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("expected error to contain '%s', got: %v", tt.errorMsg, err)
			}
		})
	}
}
