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
	"fmt"
	"reflect"
	"strings"
)

// GetFieldValue extracts a Value from an object by Field name
// Supports nested fields with dot notation (e.g., "Entity.ID")
// Supports map access where a field is a map[string]any (e.g., "Properties.Foo")
// Supports slice traversal with recursive field path evaluation (e.g., "Entities.Entity.ID")
func GetFieldValue(obj any, fieldPath string) (any, error) {
	parts := strings.Split(fieldPath, ".")
	val := reflect.ValueOf(obj)

	for i, part := range parts {
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if !val.IsValid() {
			return nil, fmt.Errorf("invalid Value at Field %s", part)
		}

		// Check if this is a map[string]any for map access
		if val.Kind() == reflect.Map && val.Type().Key().Kind() == reflect.String {
			mapVal := val.MapIndex(reflect.ValueOf(part))
			if !mapVal.IsValid() {
				return nil, fmt.Errorf("key %s not found in map", part)
			}
			val = mapVal
			continue
		}

		if val.Kind() != reflect.Struct {
			return nil, fmt.Errorf("cannot access Field %s on non-struct type %v", part, val.Type())
		}
		var err error
		val, err = getFieldValueCaseInsensitive(val, part)
		if err != nil {
			return nil, fmt.Errorf("error getting Field %s not found %w", part, err)
		}

		// After getting the field, check if we have a slice and more parts to traverse
		if i < len(parts)-1 {
			// Dereference pointer if needed
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}

			// If this field is a slice, recursively process remaining path on each element
			if val.Kind() == reflect.Slice {
				remainingPath := strings.Join(parts[i+1:], ".")
				var results []any
				for j := 0; j < val.Len(); j++ {
					sliceElem := val.Index(j)
					elemValue, err := GetFieldValue(sliceElem.Interface(), remainingPath)
					if err == nil && elemValue != nil {
						results = append(results, elemValue)
					}
				}
				// Return the collected results as a slice
				if len(results) > 0 {
					return results, nil
				}
				return nil, fmt.Errorf("no matching values found in slice for path %s", remainingPath)
			}
		}
	}
	return val.Interface(), nil
}

func getFieldValueCaseInsensitive(val reflect.Value, fieldName string) (reflect.Value, error) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("not a struct")
	}

	// Search case-insensitively
	field, found := val.Type().FieldByNameFunc(func(s string) bool {
		return strings.EqualFold(s, fieldName)
	})

	if !found {
		return reflect.Value{}, fmt.Errorf("field %s not found", fieldName)
	}

	return val.FieldByIndex(field.Index), nil
}

// CompareValues compares two values based on an Operator
// Works with various types: strings, numbers, slices, etc.
// When fieldValue is a slice (from nested slice traversal), checks if ANY element matches
func CompareValues(op Operator, fieldValue, compareValue any) bool {
	// If fieldValue is a slice, check if any element matches
	fieldValueSlice := reflect.ValueOf(fieldValue)
	if fieldValueSlice.Kind() == reflect.Slice {
		for i := 0; i < fieldValueSlice.Len(); i++ {
			elemValue := fieldValueSlice.Index(i).Interface()
			if compareValueForSingleElement(op, elemValue, compareValue) {
				return true
			}
		}
		return false
	}

	return compareValueForSingleElement(op, fieldValue, compareValue)
}

// compareValueForSingleElement performs the actual comparison for a single value
func compareValueForSingleElement(op Operator, fieldValue, compareValue any) bool {
	// Normalize string aliases to plain strings before comparison
	fieldValue = normalizeTypeAlias(fieldValue)
	compareValue = normalizeTypeAlias(compareValue)

	switch op {
	case OpEqual:
		// Try numeric comparison for numbers
		aNum, aOk := toFloat64(fieldValue)
		bNum, bOk := toFloat64(compareValue)
		if aOk && bOk {
			return aNum == bNum
		}
		return fieldValue == compareValue
	case OpNotEqual:
		// Try numeric comparison for numbers
		aNum, aOk := toFloat64(fieldValue)
		bNum, bOk := toFloat64(compareValue)
		if aOk && bOk {
			return aNum != bNum
		}
		return fieldValue != compareValue
	case OpGreater:
		return compare(fieldValue, compareValue) > 0
	case OpGreaterEqual:
		return compare(fieldValue, compareValue) >= 0
	case OpLess:
		return compare(fieldValue, compareValue) < 0
	case OpLessEqual:
		return compare(fieldValue, compareValue) <= 0
	case OpIn:
		return inSlice(fieldValue, compareValue)
	case OpNotIn:
		return !inSlice(fieldValue, compareValue)
	case OpLike, OpContains:
		return strings.Contains(toString(fieldValue), toString(compareValue))
	case OpStartsWith:
		return strings.HasPrefix(toString(fieldValue), toString(compareValue))
	case OpEndsWith:
		return strings.HasSuffix(toString(fieldValue), toString(compareValue))
	case OpIsNull:
		return fieldValue == nil
	case OpIsNotNull:
		return fieldValue != nil
	}
	return false
}

// compare handles comparison of different types
func compare(a, b any) int {
	// Try numeric comparison first
	aNum, aOk := toFloat64(a)
	bNum, bOk := toFloat64(b)
	if aOk && bOk {
		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}
		return 0
	}

	// Fall back to string comparison
	aStr := toString(a)
	bStr := toString(b)
	if aStr < bStr {
		return -1
	}
	if aStr > bStr {
		return 1
	}
	return 0
}

// toFloat64 attempts to convert an any to float64
func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case *float64:
		if val != nil {
			return *val, true
		}
	case *float32:
		if val != nil {
			return float64(*val), true
		}
	case *int:
		if val != nil {
			return float64(*val), true
		}
	case *int8:
		if val != nil {
			return float64(*val), true
		}
	case *int16:
		if val != nil {
			return float64(*val), true
		}
	case *int32:
		if val != nil {
			return float64(*val), true
		}
	case *int64:
		if val != nil {
			return float64(*val), true
		}
	case *uint:
		if val != nil {
			return float64(*val), true
		}
	case *uint8:
		if val != nil {
			return float64(*val), true
		}
	case *uint16:
		if val != nil {
			return float64(*val), true
		}
	case *uint32:
		if val != nil {
			return float64(*val), true
		}
	case *uint64:
		if val != nil {
			return float64(*val), true
		}
	}
	return 0, false
}

// inSlice checks if a Value is in a slice
func inSlice(value, slice any) bool {
	// Normalize the search value to handle type aliases
	value = normalizeTypeAlias(value)

	sliceVal := reflect.ValueOf(slice)
	if sliceVal.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < sliceVal.Len(); i++ {
		sliceItem := sliceVal.Index(i).Interface()
		// Normalize the slice item to handle type aliases
		sliceItem = normalizeTypeAlias(sliceItem)

		// Try numeric comparison first
		aNum, aOk := toFloat64(value)
		bNum, bOk := toFloat64(sliceItem)
		if aOk && bOk {
			if aNum == bNum {
				return true
			}
			continue
		}
		if value == sliceItem {
			return true
		}
	}
	return false
}

// toString converts any Value to string
func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// normalizeTypeAlias converts type aliases to their underlying types for comparison.
// Handles string aliases, numeric aliases (int*, uint*, float*), and iota-based enum types.
func normalizeTypeAlias(value any) any {
	t := reflect.TypeOf(value)
	if t == nil {
		return value
	}

	kind := t.Kind()
	typeName := t.String()

	// If the type name matches the kind name, it's a base type (not an alias)
	if typeName == kind.String() {
		return value
	}

	// Convert aliases to their underlying types
	switch kind {
	case reflect.String:
		return reflect.ValueOf(value).String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(value).Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(value).Uint()
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(value).Float()
	default:
		return value
	}
}
