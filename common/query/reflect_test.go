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

import (
	"testing"
)

// Test structures for GetFieldValue tests
type Person struct {
	Name string
	Age  int
}

type Company struct {
	Name    string
	CEO     *Person
	Founded int
}

type Address struct {
	Street string
	City   string
	Zip    string
}

type Employee struct {
	Name    string
	Address Address
	Manager *Employee
}

// TestGetFieldValue tests extraction of simple fields
func TestGetFieldValue_SimpleFields(t *testing.T) {
	person := Person{Name: "John", Age: 30}

	tests := []struct {
		name      string
		obj       any
		fieldPath string
		expected  any
		wantErr   bool
	}{
		{
			name:      "string field",
			obj:       person,
			fieldPath: "Name",
			expected:  "John",
			wantErr:   false,
		},
		{
			name:      "int field",
			obj:       person,
			fieldPath: "Age",
			expected:  30,
			wantErr:   false,
		},
		{
			name:      "non-existent field",
			obj:       person,
			fieldPath: "Email",
			expected:  nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetFieldValue(tt.obj, tt.fieldPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("GetFieldValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetFieldValue_NestedFields tests extraction of nested fields with dot notation
func TestGetFieldValue_NestedFields(t *testing.T) {
	ceo := &Person{Name: "Alice", Age: 50}
	company := Company{Name: "TechCorp", CEO: ceo, Founded: 2010}

	tests := []struct {
		name      string
		obj       any
		fieldPath string
		expected  any
		wantErr   bool
	}{
		{
			name:      "nested pointer field",
			obj:       company,
			fieldPath: "CEO.Name",
			expected:  "Alice",
			wantErr:   false,
		},
		{
			name:      "nested pointer int field",
			obj:       company,
			fieldPath: "CEO.Age",
			expected:  50,
			wantErr:   false,
		},
		{
			name:      "nested struct field",
			obj:       company,
			fieldPath: "CEO",
			expected:  ceo,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetFieldValue(tt.obj, tt.fieldPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("GetFieldValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetFieldValue_DeepNesting tests deeply nested field access
func TestGetFieldValue_DeepNesting(t *testing.T) {
	emp1 := &Employee{Name: "Bob"}
	emp2 := &Employee{Name: "Charlie", Manager: emp1}
	emp3 := Employee{Name: "David", Manager: emp2}

	tests := []struct {
		name      string
		obj       any
		fieldPath string
		expected  any
		wantErr   bool
	}{
		{
			name:      "two levels deep",
			obj:       emp3,
			fieldPath: "Manager.Name",
			expected:  "Charlie",
			wantErr:   false,
		},
		{
			name:      "three levels deep",
			obj:       emp3,
			fieldPath: "Manager.Manager.Name",
			expected:  "Bob",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetFieldValue(tt.obj, tt.fieldPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("GetFieldValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetFieldValue_InvalidCases tests error conditions
func TestGetFieldValue_InvalidCases(t *testing.T) {
	tests := []struct {
		name      string
		obj       any
		fieldPath string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "nil pointer dereference",
			obj:       Company{Name: "NoCEO", CEO: nil},
			fieldPath: "CEO.Name",
			wantErr:   true,
			errMsg:    "invalid Value at Field",
		},
		{
			name:      "non-struct type",
			obj:       "string",
			fieldPath: "Name",
			wantErr:   true,
			errMsg:    "cannot access Field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetFieldValue(tt.obj, tt.fieldPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCompareValues_Equality tests equal and not equal operators
func TestCompareValues_Equality(t *testing.T) {
	tests := []struct {
		name         string
		op           Operator
		fieldValue   any
		compareValue any
		expected     bool
	}{
		{
			name:         "equal strings",
			op:           OpEqual,
			fieldValue:   "hello",
			compareValue: "hello",
			expected:     true,
		},
		{
			name:         "not equal strings",
			op:           OpEqual,
			fieldValue:   "hello",
			compareValue: "world",
			expected:     false,
		},
		{
			name:         "not equal - equal values",
			op:           OpNotEqual,
			fieldValue:   "hello",
			compareValue: "hello",
			expected:     false,
		},
		{
			name:         "not equal - different values",
			op:           OpNotEqual,
			fieldValue:   "hello",
			compareValue: "world",
			expected:     true,
		},
		{
			name:         "equal ints",
			op:           OpEqual,
			fieldValue:   42,
			compareValue: 42,
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

// TestCompareValues_Numeric tests numeric comparison operators
func TestCompareValues_Numeric(t *testing.T) {
	tests := []struct {
		name         string
		op           Operator
		fieldValue   any
		compareValue any
		expected     bool
	}{
		{
			name:         "greater than",
			op:           OpGreater,
			fieldValue:   10,
			compareValue: 5,
			expected:     true,
		},
		{
			name:         "greater than false",
			op:           OpGreater,
			fieldValue:   5,
			compareValue: 10,
			expected:     false,
		},
		{
			name:         "less than",
			op:           OpLess,
			fieldValue:   5,
			compareValue: 10,
			expected:     true,
		},
		{
			name:         "less than false",
			op:           OpLess,
			fieldValue:   10,
			compareValue: 5,
			expected:     false,
		},
		{
			name:         "greater equal",
			op:           OpGreaterEqual,
			fieldValue:   10,
			compareValue: 10,
			expected:     true,
		},
		{
			name:         "less equal",
			op:           OpLessEqual,
			fieldValue:   10,
			compareValue: 10,
			expected:     true,
		},
		{
			name:         "greater with float",
			op:           OpGreater,
			fieldValue:   10.5,
			compareValue: 10.2,
			expected:     true,
		},
		{
			name:         "float64 vs int",
			op:           OpEqual,
			fieldValue:   float64(42),
			compareValue: 42,
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

// TestCompareValues_StringOperations tests string-based operators
func TestCompareValues_StringOperations(t *testing.T) {
	tests := []struct {
		name         string
		op           Operator
		fieldValue   any
		compareValue any
		expected     bool
	}{
		{
			name:         "contains true",
			op:           OpContains,
			fieldValue:   "hello world",
			compareValue: "world",
			expected:     true,
		},
		{
			name:         "contains false",
			op:           OpContains,
			fieldValue:   "hello world",
			compareValue: "xyz",
			expected:     false,
		},
		{
			name:         "like operator",
			op:           OpLike,
			fieldValue:   "test string",
			compareValue: "string",
			expected:     true,
		},
		{
			name:         "starts with true",
			op:           OpStartsWith,
			fieldValue:   "hello world",
			compareValue: "hello",
			expected:     true,
		},
		{
			name:         "starts with false",
			op:           OpStartsWith,
			fieldValue:   "hello world",
			compareValue: "world",
			expected:     false,
		},
		{
			name:         "ends with true",
			op:           OpEndsWith,
			fieldValue:   "hello world",
			compareValue: "world",
			expected:     true,
		},
		{
			name:         "ends with false",
			op:           OpEndsWith,
			fieldValue:   "hello world",
			compareValue: "hello",
			expected:     false,
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

// TestCompareValues_NullOperations tests IS NULL and IS NOT NULL operators
func TestCompareValues_NullOperations(t *testing.T) {
	tests := []struct {
		name         string
		op           Operator
		fieldValue   any
		compareValue any
		expected     bool
	}{
		{
			name:         "is null - nil value",
			op:           OpIsNull,
			fieldValue:   nil,
			compareValue: nil,
			expected:     true,
		},
		{
			name:         "is null - non-nil value",
			op:           OpIsNull,
			fieldValue:   "something",
			compareValue: nil,
			expected:     false,
		},
		{
			name:         "is not null - non-nil value",
			op:           OpIsNotNull,
			fieldValue:   "something",
			compareValue: nil,
			expected:     true,
		},
		{
			name:         "is not null - nil value",
			op:           OpIsNotNull,
			fieldValue:   nil,
			compareValue: nil,
			expected:     false,
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

// TestCompareValues_InOperators tests IN and NOT IN operators
func TestCompareValues_InOperators(t *testing.T) {
	tests := []struct {
		name         string
		op           Operator
		fieldValue   any
		compareValue any
		expected     bool
	}{
		{
			name:         "in - found",
			op:           OpIn,
			fieldValue:   "apple",
			compareValue: []any{"apple", "banana", "orange"},
			expected:     true,
		},
		{
			name:         "in - not found",
			op:           OpIn,
			fieldValue:   "grape",
			compareValue: []any{"apple", "banana", "orange"},
			expected:     false,
		},
		{
			name:         "in - with ints",
			op:           OpIn,
			fieldValue:   5,
			compareValue: []any{1, 5, 10},
			expected:     true,
		},
		{
			name:         "not in - found",
			op:           OpNotIn,
			fieldValue:   "apple",
			compareValue: []any{"apple", "banana", "orange"},
			expected:     false,
		},
		{
			name:         "not in - not found",
			op:           OpNotIn,
			fieldValue:   "grape",
			compareValue: []any{"apple", "banana", "orange"},
			expected:     true,
		},
		{
			name:         "in - empty slice",
			op:           OpIn,
			fieldValue:   "apple",
			compareValue: []any{},
			expected:     false,
		},
		{
			name:         "in - non-slice returns false",
			op:           OpIn,
			fieldValue:   "apple",
			compareValue: "not a slice",
			expected:     false,
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

// TestCompareValues_UnsupportedOperator tests unknown operator
func TestCompareValues_UnsupportedOperator(t *testing.T) {
	result := CompareValues(Operator("UNKNOWN"), "value", "value")
	if result != false {
		t.Errorf("CompareValues with unknown operator should return false, got %v", result)
	}
}

// TestToFloat64 tests numeric type conversions
func TestToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected float64
		ok       bool
	}{
		{
			name:     "float64",
			value:    3.14,
			expected: 3.14,
			ok:       true,
		},
		{
			name:     "float32",
			value:    float32(2.71),
			expected: 2.7100000381469727,
			ok:       true,
		},
		{
			name:     "int",
			value:    42,
			expected: 42.0,
			ok:       true,
		},
		{
			name:     "int32",
			value:    int32(100),
			expected: 100.0,
			ok:       true,
		},
		{
			name:     "int64",
			value:    int64(9999),
			expected: 9999.0,
			ok:       true,
		},
		{
			name:     "uint",
			value:    uint(50),
			expected: 50.0,
			ok:       true,
		},
		{
			name:     "uint64",
			value:    uint64(12345),
			expected: 12345.0,
			ok:       true,
		},
		{
			name:     "string (not convertible)",
			value:    "not a number",
			expected: 0,
			ok:       false,
		},
		{
			name:     "nil (not convertible)",
			value:    nil,
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := toFloat64(tt.value)
			if ok != tt.ok {
				t.Errorf("toFloat64() ok = %v, want %v", ok, tt.ok)
				return
			}
			if ok && result != tt.expected {
				t.Errorf("toFloat64() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestInSlice tests slice membership
func TestInSlice(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		slice    any
		expected bool
	}{
		{
			name:     "string in slice",
			value:    "apple",
			slice:    []any{"apple", "banana", "orange"},
			expected: true,
		},
		{
			name:     "string not in slice",
			value:    "grape",
			slice:    []any{"apple", "banana", "orange"},
			expected: false,
		},
		{
			name:     "int in slice",
			value:    5,
			slice:    []any{1, 2, 5, 10},
			expected: true,
		},
		{
			name:     "empty slice",
			value:    "anything",
			slice:    []any{},
			expected: false,
		},
		{
			name:     "not a slice",
			value:    "value",
			slice:    "not a slice",
			expected: false,
		},
		{
			name:     "nil value in slice",
			value:    nil,
			slice:    []any{nil, "value"},
			expected: true,
		},
		{
			name:     "int slice",
			value:    3,
			slice:    []int{1, 2, 3, 4},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inSlice(tt.value, tt.slice)
			if result != tt.expected {
				t.Errorf("inSlice(%v, %v) = %v, want %v", tt.value, tt.slice, result, tt.expected)
			}
		})
	}
}

// TestToString tests conversion to string
func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{
			name:     "string value",
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "int value",
			value:    42,
			expected: "42",
		},
		{
			name:     "float value",
			value:    3.14,
			expected: "3.14",
		},
		{
			name:     "bool value",
			value:    true,
			expected: "true",
		},
		{
			name:     "nil value",
			value:    nil,
			expected: "<nil>",
		},
		{
			name:     "slice value",
			value:    []int{1, 2, 3},
			expected: "[1 2 3]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toString(tt.value)
			if result != tt.expected {
				t.Errorf("toString(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

// TestCompare tests the internal compare function
func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		a        any
		b        any
		expected int
	}{
		{
			name:     "equal numbers",
			a:        10,
			b:        10,
			expected: 0,
		},
		{
			name:     "a > b (numbers)",
			a:        20,
			b:        10,
			expected: 1,
		},
		{
			name:     "a < b (numbers)",
			a:        10,
			b:        20,
			expected: -1,
		},
		{
			name:     "equal strings",
			a:        "hello",
			b:        "hello",
			expected: 0,
		},
		{
			name:     "a > b (strings)",
			a:        "world",
			b:        "hello",
			expected: 1,
		},
		{
			name:     "a < b (strings)",
			a:        "hello",
			b:        "world",
			expected: -1,
		},
		{
			name:     "string vs string comparison fallback",
			a:        "100",
			b:        "20",
			expected: -1, // String comparison, not numeric
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compare(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compare(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestCompareValues_NumericTypeConversions tests numeric comparisons with mixed types
func TestCompareValues_NumericTypeConversions(t *testing.T) {
	tests := []struct {
		name         string
		op           Operator
		fieldValue   any
		compareValue any
		expected     bool
	}{
		{
			name:         "int32 greater than int64",
			op:           OpGreater,
			fieldValue:   int32(100),
			compareValue: int64(50),
			expected:     true,
		},
		{
			name:         "uint8 less than uint32",
			op:           OpLess,
			fieldValue:   uint8(10),
			compareValue: uint32(20),
			expected:     true,
		},
		{
			name:         "int greater or equal to float",
			op:           OpGreaterEqual,
			fieldValue:   100,
			compareValue: float64(100.0),
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

// TestGetFieldValue_MapAccess tests extraction of values from map[string]any fields
func TestGetFieldValue_MapAccess(t *testing.T) {
	type DataContainer struct {
		ID    string
		Props map[string]any
	}

	container := DataContainer{
		ID: "123",
		Props: map[string]any{
			"Color":  "blue",
			"Size":   42,
			"Active": true,
			"Score":  95.5,
		},
	}

	tests := []struct {
		name      string
		obj       any
		fieldPath string
		expected  any
		wantErr   bool
	}{
		{
			name:      "map string value",
			obj:       container,
			fieldPath: "Props.Color",
			expected:  "blue",
			wantErr:   false,
		},
		{
			name:      "map int value",
			obj:       container,
			fieldPath: "Props.Size",
			expected:  42,
			wantErr:   false,
		},
		{
			name:      "map bool value",
			obj:       container,
			fieldPath: "Props.Active",
			expected:  true,
			wantErr:   false,
		},
		{
			name:      "map float value",
			obj:       container,
			fieldPath: "Props.Score",
			expected:  95.5,
			wantErr:   false,
		},
		{
			name:      "map key not found",
			obj:       container,
			fieldPath: "Props.NonExistent",
			expected:  nil,
			wantErr:   true,
		},
		{
			name:      "regular field access before map",
			obj:       container,
			fieldPath: "ID",
			expected:  "123",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetFieldValue(tt.obj, tt.fieldPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("GetFieldValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestGetFieldValue_NestedMapAccess tests deeply nested map access with struct fields
func TestGetFieldValue_NestedMapAccess(t *testing.T) {
	type MetaInfo struct {
		Metadata map[string]any
	}

	type Item struct {
		Name string
		Meta MetaInfo
	}

	item := Item{
		Name: "Widget",
		Meta: MetaInfo{
			Metadata: map[string]any{
				"Version": "1.0",
				"Author":  "John",
				"Status":  "active",
			},
		},
	}

	tests := []struct {
		name      string
		obj       any
		fieldPath string
		expected  any
		wantErr   bool
	}{
		{
			name:      "struct to map access",
			obj:       item,
			fieldPath: "Meta.Metadata.Version",
			expected:  "1.0",
			wantErr:   false,
		},
		{
			name:      "struct to map access string",
			obj:       item,
			fieldPath: "Meta.Metadata.Author",
			expected:  "John",
			wantErr:   false,
		},
		{
			name:      "struct to map access another field",
			obj:       item,
			fieldPath: "Meta.Metadata.Status",
			expected:  "active",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetFieldValue(tt.obj, tt.fieldPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("GetFieldValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}
