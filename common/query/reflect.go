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
func GetFieldValue(obj any, fieldPath string) (any, error) {
	parts := strings.Split(fieldPath, ".")
	val := reflect.ValueOf(obj)

	for _, part := range parts {
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if !val.IsValid() {
			return nil, fmt.Errorf("invalid Value at Field %s", part)
		}
		if val.Kind() != reflect.Struct {
			return nil, fmt.Errorf("cannot access Field %s on non-struct type %v", part, val.Type())
		}
		val = val.FieldByName(part)
		if !val.IsValid() {
			return nil, fmt.Errorf("Field %s not found", part)
		}
	}
	return val.Interface(), nil
}

// CompareValues compares two values based on an Operator
// Works with various types: strings, numbers, slices, etc.
func CompareValues(op Operator, fieldValue, compareValue any) bool {
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
	sliceVal := reflect.ValueOf(slice)
	if sliceVal.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < sliceVal.Len(); i++ {
		sliceItem := sliceVal.Index(i).Interface()
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
