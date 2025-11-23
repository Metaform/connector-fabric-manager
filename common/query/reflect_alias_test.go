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

package query

import "testing"

// TestNormalizeTypeAlias_StringAliases tests normalization of string type aliases
func TestNormalizeTypeAlias_StringAliases(t *testing.T) {
	type CustomString string

	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "string alias converts to string",
			input:    CustomString("hello"),
			expected: "hello",
		},
		{
			name:     "plain string remains unchanged",
			input:    "world",
			expected: "world",
		},
		{
			name:     "empty string alias",
			input:    CustomString(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeTypeAlias(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeTypeAlias(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestNormalizeTypeAlias_IntAliases tests normalization of integer type aliases
func TestNormalizeTypeAlias_IntAliases(t *testing.T) {
	type CustomInt int
	type CustomInt8 int8
	type CustomInt16 int16
	type CustomInt32 int32
	type CustomInt64 int64

	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "int alias converts to int64",
			input:    CustomInt(42),
			expected: int64(42),
		},
		{
			name:     "int8 alias converts to int64",
			input:    CustomInt8(10),
			expected: int64(10),
		},
		{
			name:     "int16 alias converts to int64",
			input:    CustomInt16(100),
			expected: int64(100),
		},
		{
			name:     "int32 alias converts to int64",
			input:    CustomInt32(1000),
			expected: int64(1000),
		},
		{
			name:     "int64 alias converts to int64",
			input:    CustomInt64(10000),
			expected: int64(10000),
		},
		{
			name:     "negative int alias",
			input:    CustomInt(-50),
			expected: int64(-50),
		},
		{
			name:     "zero int alias",
			input:    CustomInt(0),
			expected: int64(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeTypeAlias(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeTypeAlias(%v) = %v (type %T), want %v (type %T)", tt.input, result, result, tt.expected, tt.expected)
			}
		})
	}
}

// TestNormalizeTypeAlias_UintAliases tests normalization of unsigned integer type aliases
func TestNormalizeTypeAlias_UintAliases(t *testing.T) {
	type CustomUint uint
	type CustomUint8 uint8
	type CustomUint16 uint16
	type CustomUint32 uint32
	type CustomUint64 uint64

	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "uint alias converts to uint64",
			input:    CustomUint(42),
			expected: uint64(42),
		},
		{
			name:     "uint8 alias converts to uint64",
			input:    CustomUint8(10),
			expected: uint64(10),
		},
		{
			name:     "uint16 alias converts to uint64",
			input:    CustomUint16(100),
			expected: uint64(100),
		},
		{
			name:     "uint32 alias converts to uint64",
			input:    CustomUint32(1000),
			expected: uint64(1000),
		},
		{
			name:     "uint64 alias converts to uint64",
			input:    CustomUint64(10000),
			expected: uint64(10000),
		},
		{
			name:     "zero uint alias",
			input:    CustomUint(0),
			expected: uint64(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeTypeAlias(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeTypeAlias(%v) = %v (type %T), want %v (type %T)", tt.input, result, result, tt.expected, tt.expected)
			}
		})
	}
}

// TestNormalizeTypeAlias_FloatAliases tests normalization of float type aliases
func TestNormalizeTypeAlias_FloatAliases(t *testing.T) {
	type CustomFloat32 float32
	type CustomFloat64 float64

	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "float32 alias converts to float64",
			input:    CustomFloat32(3.14),
			expected: float64(CustomFloat32(3.14)),
		},
		{
			name:     "float64 alias converts to float64",
			input:    CustomFloat64(2.71),
			expected: float64(2.71),
		},
		{
			name:     "plain float64 stays float64",
			input:    float64(1.5),
			expected: float64(1.5),
		},
		{
			name:     "zero float32 alias",
			input:    CustomFloat32(0.0),
			expected: float64(0.0),
		},
		{
			name:     "negative float alias",
			input:    CustomFloat64(-99.99),
			expected: float64(-99.99),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeTypeAlias(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeTypeAlias(%v) = %v (type %T), want %v (type %T)", tt.input, result, result, tt.expected, tt.expected)
			}
		})
	}
}

// TestNormalizeTypeAlias_IotaEnums tests normalization of iota-based enum types
func TestNormalizeTypeAlias_IotaEnums(t *testing.T) {
	type Status int
	const (
		StatusUnknown Status = iota
		StatusActive
		StatusInactive
		StatusPending
	)

	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "iota enum converts to int64",
			input:    StatusActive,
			expected: int64(1),
		},
		{
			name:     "iota enum zero value",
			input:    StatusUnknown,
			expected: int64(0),
		},
		{
			name:     "iota enum high value",
			input:    StatusPending,
			expected: int64(3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeTypeAlias(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeTypeAlias(%v) = %v (type %T), want %v (type %T)", tt.input, result, result, tt.expected, tt.expected)
			}
		})
	}
}

// TestNormalizeTypeAlias_NilAndNonAliasTypes tests nil values and non-alias base types
func TestNormalizeTypeAlias_NilAndNonAliasTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: nil,
		},
		{
			name:     "base string type",
			input:    "base",
			expected: "base",
		},
		{
			name:     "base bool type",
			input:    true,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeTypeAlias(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeTypeAlias(%v) = %v (type %T), want %v (type %T)", tt.input, result, result, tt.expected, tt.expected)
			}
		})
	}
}

// TestNormalizeTypeAlias_IntegrationWithCompareValues tests that normalizeTypeAlias works correctly within CompareValues
func TestNormalizeTypeAlias_IntegrationWithCompareValues(t *testing.T) {
	type CustomStatus int
	const (
		StatusA CustomStatus = iota
		StatusB
		StatusC
	)

	type StatusName string

	tests := []struct {
		name         string
		op           Operator
		fieldValue   any
		compareValue any
		expected     bool
	}{
		{
			name:         "string alias equality",
			op:           OpEqual,
			fieldValue:   StatusName("active"),
			compareValue: "active",
			expected:     true,
		},
		{
			name:         "int alias equality",
			op:           OpEqual,
			fieldValue:   StatusB,
			compareValue: 1,
			expected:     true,
		},
		{
			name:         "int alias not equal",
			op:           OpNotEqual,
			fieldValue:   StatusC,
			compareValue: 1,
			expected:     true,
		},
		{
			name:         "int alias greater than",
			op:           OpGreater,
			fieldValue:   StatusC,
			compareValue: StatusB,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareValues(tt.op, tt.fieldValue, tt.compareValue)
			if result != tt.expected {
				t.Errorf("CompareValues(%v, %v, %v) = %v, want %v", tt.op, tt.fieldValue, tt.compareValue, result, tt.expected)
			}
		})
	}
}
