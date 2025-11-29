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

// Test nested slice traversal with dot notation
func TestGetFieldValue_NestedSliceTraversal(t *testing.T) {
	type Cell struct {
		ID   string
		Name string
	}

	type VPA struct {
		ID   string
		Cell Cell
	}

	type Profile struct {
		ID   string
		VPAs []VPA
	}

	profile := Profile{
		ID: "profile-1",
		VPAs: []VPA{
			{
				ID: "vpa-1",
				Cell: Cell{
					ID:   "cell1",
					Name: "Cell One",
				},
			},
			{
				ID: "vpa-2",
				Cell: Cell{
					ID:   "cell2",
					Name: "Cell Two",
				},
			},
		},
	}

	tests := []struct {
		name        string
		fieldPath   string
		expected    []string
		expectError bool
	}{
		{
			name:      "Single level nested slice - VPAs.ID",
			fieldPath: "VPAs.ID",
			expected:  []string{"vpa-1", "vpa-2"},
		},
		{
			name:      "Two level nested slice - VPAs.Cell.ID",
			fieldPath: "VPAs.Cell.ID",
			expected:  []string{"cell1", "cell2"},
		},
		{
			name:      "Two level nested slice - VPAs.Cell.Name",
			fieldPath: "VPAs.Cell.Name",
			expected:  []string{"Cell One", "Cell Two"},
		},
		{
			name:        "Invalid path in slice",
			fieldPath:   "VPAs.InvalidField",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetFieldValue(profile, tt.fieldPath)
			if (err != nil) != tt.expectError {
				t.Errorf("GetFieldValue() error = %v, expectError = %v", err, tt.expectError)
				return
			}

			if err == nil {
				resultSlice, ok := result.([]any)
				if !ok {
					t.Errorf("Expected []any, got %T", result)
					return
				}

				if len(resultSlice) != len(tt.expected) {
					t.Errorf("Expected %d results, got %d", len(tt.expected), len(resultSlice))
					return
				}

				for i, val := range resultSlice {
					if val != tt.expected[i] {
						t.Errorf("Result[%d] = %v, want %v", i, val, tt.expected[i])
					}
				}
			}
		})
	}
}

// Test slice comparison with CompareValues
func TestCompareValues_SliceWithAny(t *testing.T) {
	tests := []struct {
		name        string
		operator    Operator
		fieldValue  any
		compareVal  any
		expectedRes bool
	}{
		{
			name:        "Slice with OpEqual - match found",
			operator:    OpEqual,
			fieldValue:  []any{"cell1", "cell2", "cell3"},
			compareVal:  "cell1",
			expectedRes: true,
		},
		{
			name:        "Slice with OpEqual - no match",
			operator:    OpEqual,
			fieldValue:  []any{"cell1", "cell2", "cell3"},
			compareVal:  "cell4",
			expectedRes: false,
		},
		{
			name:        "Slice with OpContains - match found",
			operator:    OpContains,
			fieldValue:  []any{"cell-alpha", "cell-beta", "cell-gamma"},
			compareVal:  "beta",
			expectedRes: true,
		},
		{
			name:        "Slice with OpStartsWith - match found",
			operator:    OpStartsWith,
			fieldValue:  []any{"vpa-1", "vpa-2", "connector-1"},
			compareVal:  "vpa",
			expectedRes: true,
		},
		{
			name:        "Slice with OpEndsWith - match found",
			operator:    OpEndsWith,
			fieldValue:  []any{"cell1", "cell2", "data-cell3"},
			compareVal:  "cell3",
			expectedRes: true,
		},
		{
			name:        "Single value (not slice) equality",
			operator:    OpEqual,
			fieldValue:  "cell1",
			compareVal:  "cell1",
			expectedRes: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareValues(tt.operator, tt.fieldValue, tt.compareVal)
			if result != tt.expectedRes {
				t.Errorf("CompareValues() = %v, want %v", result, tt.expectedRes)
			}
		})
	}
}

// Test predicate matching with nested slice fields
func TestAtomicPredicate_NestedSliceMatching(t *testing.T) {
	type Cell struct {
		ID   string
		Name string
	}

	type VPA struct {
		ID   string
		Cell Cell
	}

	type Profile struct {
		ID   string
		VPAs []VPA
	}

	profile := Profile{
		ID: "profile-1",
		VPAs: []VPA{
			{
				ID: "vpa-1",
				Cell: Cell{
					ID:   "cell1",
					Name: "Cell One",
				},
			},
			{
				ID: "vpa-2",
				Cell: Cell{
					ID:   "cell2",
					Name: "Cell Two",
				},
			},
		},
	}

	tests := []struct {
		name     string
		pred     *AtomicPredicate
		expected bool
	}{
		{
			name: "Match nested slice field with ID equality",
			pred: &AtomicPredicate{
				Field:    "VPAs.Cell.ID",
				Operator: OpEqual,
				Value:    "cell1",
			},
			expected: true,
		},
		{
			name: "No match nested slice field",
			pred: &AtomicPredicate{
				Field:    "VPAs.Cell.ID",
				Operator: OpEqual,
				Value:    "cell99",
			},
			expected: false,
		},
		{
			name: "Match nested slice field with CONTAINS",
			pred: &AtomicPredicate{
				Field:    "VPAs.Cell.Name",
				Operator: OpContains,
				Value:    "One",
			},
			expected: true,
		},
		{
			name: "Match nested slice parent field",
			pred: &AtomicPredicate{
				Field:    "VPAs.ID",
				Operator: OpEqual,
				Value:    "vpa-2",
			},
			expected: true,
		},
	}

	matcher := &DefaultFieldMatcher{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pred.Matches(profile, matcher)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}
