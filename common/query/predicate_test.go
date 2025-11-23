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

// TestEntity is a simple test struct for predicate matching
type TestEntity struct {
	ID       string
	Name     string
	Age      int
	Score    float64
	Active   bool
	Email    string
	Category string
	Count    int
}

// TestAtomicPredicate_Matches_Equality tests equality predicates
func TestAtomicPredicate_Matches_Equality(t *testing.T) {
	entity := TestEntity{
		ID:       "123",
		Name:     "Alice",
		Age:      30,
		Score:    95.5,
		Active:   true,
		Email:    "alice@example.com",
		Category: "Premium",
		Count:    5,
	}

	tests := []struct {
		name      string
		predicate *AtomicPredicate
		expected  bool
	}{
		{
			name:      "equal string - match",
			predicate: Eq("Name", "Alice"),
			expected:  true,
		},
		{
			name:      "equal string - no match",
			predicate: Eq("Name", "Bob"),
			expected:  false,
		},
		{
			name:      "equal int - match",
			predicate: Eq("Age", 30),
			expected:  true,
		},
		{
			name:      "equal int - no match",
			predicate: Eq("Age", 25),
			expected:  false,
		},
		{
			name:      "equal bool - true",
			predicate: Eq("Active", true),
			expected:  true,
		},
		{
			name:      "equal bool - false no match",
			predicate: Eq("Active", false),
			expected:  false,
		},
		{
			name:      "equal float - match",
			predicate: Eq("Score", 95.5),
			expected:  true,
		},
		{
			name:      "equal float - no match",
			predicate: Eq("Score", 90.0),
			expected:  false,
		},
		{
			name:      "not equal - match",
			predicate: Neq("Category", "Standard"),
			expected:  true,
		},
		{
			name:      "not equal - no match",
			predicate: Neq("Category", "Premium"),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(entity, nil)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAtomicPredicate_Matches_Comparison tests comparison operators
func TestAtomicPredicate_Matches_Comparison(t *testing.T) {
	entity := TestEntity{
		ID:    "123",
		Name:  "Alice",
		Age:   30,
		Score: 95.5,
	}

	tests := []struct {
		name      string
		predicate *AtomicPredicate
		expected  bool
	}{
		{
			name:      "greater than - true",
			predicate: Gt("Age", 25),
			expected:  true,
		},
		{
			name:      "greater than - false",
			predicate: Gt("Age", 30),
			expected:  false,
		},
		{
			name:      "greater than - equal",
			predicate: Gt("Age", 35),
			expected:  false,
		},
		{
			name:      "greater or equal - equal",
			predicate: Gte("Age", 30),
			expected:  true,
		},
		{
			name:      "greater or equal - true",
			predicate: Gte("Age", 25),
			expected:  true,
		},
		{
			name:      "greater or equal - false",
			predicate: Gte("Age", 35),
			expected:  false,
		},
		{
			name:      "less than - true",
			predicate: Lt("Age", 35),
			expected:  true,
		},
		{
			name:      "less than - false",
			predicate: Lt("Age", 30),
			expected:  false,
		},
		{
			name:      "less or equal - equal",
			predicate: Lte("Age", 30),
			expected:  true,
		},
		{
			name:      "less or equal - true",
			predicate: Lte("Age", 35),
			expected:  true,
		},
		{
			name:      "less or equal - false",
			predicate: Lte("Age", 25),
			expected:  false,
		},
		{
			name:      "greater than float",
			predicate: Gt("Score", 90.0),
			expected:  true,
		},
		{
			name:      "less than float",
			predicate: Lt("Score", 100.0),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(entity, nil)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAtomicPredicate_Matches_In tests IN and NOT IN operators
func TestAtomicPredicate_Matches_In(t *testing.T) {
	entity := TestEntity{
		ID:       "123",
		Name:     "Alice",
		Age:      30,
		Category: "Premium",
	}

	tests := []struct {
		name      string
		predicate *AtomicPredicate
		expected  bool
	}{
		{
			name:      "IN - string match",
			predicate: In("Category", "Premium", "Gold", "Platinum"),
			expected:  true,
		},
		{
			name:      "IN - string no match",
			predicate: In("Category", "Standard", "Gold"),
			expected:  false,
		},
		{
			name:      "IN - int match",
			predicate: In("Age", 25, 30, 35),
			expected:  true,
		},
		{
			name:      "IN - int no match",
			predicate: In("Age", 20, 25, 35),
			expected:  false,
		},
		{
			name:      "IN - single value match",
			predicate: In("Name", "Alice"),
			expected:  true,
		},
		{
			name:      "IN - single value no match",
			predicate: In("Name", "Bob"),
			expected:  false,
		},
		{
			name:      "NOT IN - match",
			predicate: NotIn("Category", "Standard", "Gold"),
			expected:  true,
		},
		{
			name:      "NOT IN - no match",
			predicate: NotIn("Category", "Premium", "Platinum"),
			expected:  false,
		},
		{
			name:      "NOT IN - int match",
			predicate: NotIn("Age", 20, 25, 35),
			expected:  true,
		},
		{
			name:      "NOT IN - int no match",
			predicate: NotIn("Age", 30, 35),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(entity, nil)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAtomicPredicate_Matches_StringPatterns tests string pattern matching
func TestAtomicPredicate_Matches_StringPatterns(t *testing.T) {
	entity := TestEntity{
		ID:    "test-123",
		Name:  "Alice Johnson",
		Email: "alice@example.com",
	}

	tests := []struct {
		name      string
		predicate *AtomicPredicate
		expected  bool
	}{
		{
			name:      "LIKE - match",
			predicate: Like("Name", "ice"),
			expected:  true,
		},
		{
			name:      "LIKE - no match",
			predicate: Like("Name", "xyz"),
			expected:  false,
		},
		{
			name:      "CONTAINS - match",
			predicate: Contains("Email", "@example"),
			expected:  true,
		},
		{
			name:      "CONTAINS - no match",
			predicate: Contains("Email", "@test"),
			expected:  false,
		},
		{
			name:      "STARTS_WITH - match",
			predicate: StartsWith("Name", "Alice"),
			expected:  true,
		},
		{
			name:      "STARTS_WITH - no match",
			predicate: StartsWith("Name", "Bob"),
			expected:  false,
		},
		{
			name:      "ENDS_WITH - match",
			predicate: EndsWith("Name", "Johnson"),
			expected:  true,
		},
		{
			name:      "ENDS_WITH - no match",
			predicate: EndsWith("Name", "Smith"),
			expected:  false,
		},
		{
			name:      "CONTAINS case sensitive",
			predicate: Contains("Name", "alice"),
			expected:  false,
		},
		{
			name:      "STARTS_WITH partial",
			predicate: StartsWith("ID", "test"),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(entity, nil)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAtomicPredicate_Matches_Null tests NULL checking operators
func TestAtomicPredicate_Matches_Null(t *testing.T) {
	entity := TestEntity{
		ID:   "123",
		Name: "Alice",
		Age:  30,
	}

	tests := []struct {
		name      string
		predicate *AtomicPredicate
		expected  bool
	}{
		{
			name:      "IS NOT NULL - non-nil value",
			predicate: IsNotNull("Name"),
			expected:  true,
		},
		{
			name:      "IS NOT NULL - zero value",
			predicate: IsNotNull("Age"),
			expected:  true,
		},
		{
			name:      "IS NULL - non-existent field",
			predicate: IsNull("NonExistent"),
			expected:  true, // Field doesn't exist, so it's nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(entity, nil)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAtomicPredicate_Matches_WithCustomMatcher tests with custom FieldMatcher
func TestAtomicPredicate_Matches_WithCustomMatcher(t *testing.T) {
	entity := TestEntity{
		ID:   "123",
		Name: "Alice",
	}

	customMatcher := &DefaultFieldMatcher{}
	predicate := Eq("Name", "Alice")

	result := predicate.Matches(entity, customMatcher)
	if !result {
		t.Errorf("Matches with custom matcher failed, expected true, got %v", result)
	}
}

// TestAtomicPredicate_Matches_InvalidField tests handling of invalid fields
func TestAtomicPredicate_Matches_InvalidField(t *testing.T) {
	entity := TestEntity{
		Name: "Alice",
	}

	tests := []struct {
		name      string
		predicate *AtomicPredicate
		expected  bool
	}{
		{
			name:      "non-existent field",
			predicate: Eq("NonExistent", "value"),
			expected:  false,
		},
		{
			name:      "non-existent field with comparison",
			predicate: Gt("NonExistent", 10),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(entity, nil)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCompoundPredicate_Matches_AND tests AND conjunction logic
func TestCompoundPredicate_Matches_AND(t *testing.T) {
	entity := TestEntity{
		ID:       "123",
		Name:     "Alice",
		Age:      30,
		Active:   true,
		Category: "Premium",
	}

	tests := []struct {
		name      string
		predicate *CompoundPredicate
		expected  bool
	}{
		{
			name: "AND - all true",
			predicate: And(
				Eq("Name", "Alice"),
				Eq("Active", true),
			),
			expected: true,
		},
		{
			name: "AND - one false",
			predicate: And(
				Eq("Name", "Alice"),
				Eq("Active", false),
			),
			expected: false,
		},
		{
			name: "AND - all false",
			predicate: And(
				Eq("Name", "Bob"),
				Eq("Active", false),
			),
			expected: false,
		},
		{
			name: "AND - multiple conditions all true",
			predicate: And(
				Eq("Name", "Alice"),
				Gte("Age", 25),
				Eq("Active", true),
				Eq("Category", "Premium"),
			),
			expected: true,
		},
		{
			name: "AND - multiple conditions one false",
			predicate: And(
				Eq("Name", "Alice"),
				Gte("Age", 35),
				Eq("Active", true),
			),
			expected: false,
		},
		{
			name: "AND - comparison operators",
			predicate: And(
				Gte("Age", 25),
				Lte("Age", 35),
			),
			expected: true,
		},
		{
			name:      "AND - empty predicates",
			predicate: And(),
			expected:  true,
		},
		{
			name: "AND - single predicate true",
			predicate: And(
				Eq("Name", "Alice"),
			),
			expected: true,
		},
		{
			name: "AND - single predicate false",
			predicate: And(
				Eq("Name", "Bob"),
			),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(entity, nil)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCompoundPredicate_Matches_OR tests OR conjunction logic
func TestCompoundPredicate_Matches_OR(t *testing.T) {
	entity := TestEntity{
		ID:       "123",
		Name:     "Alice",
		Age:      30,
		Category: "Premium",
	}

	tests := []struct {
		name      string
		predicate *CompoundPredicate
		expected  bool
	}{
		{
			name: "OR - first true",
			predicate: Or(
				Eq("Name", "Alice"),
				Eq("Name", "Bob"),
			),
			expected: true,
		},
		{
			name: "OR - second true",
			predicate: Or(
				Eq("Name", "Bob"),
				Eq("Name", "Alice"),
			),
			expected: true,
		},
		{
			name: "OR - all true",
			predicate: Or(
				Eq("Name", "Alice"),
				Eq("Category", "Premium"),
			),
			expected: true,
		},
		{
			name: "OR - all false",
			predicate: Or(
				Eq("Name", "Bob"),
				Eq("Category", "Standard"),
			),
			expected: false,
		},
		{
			name: "OR - multiple conditions one true",
			predicate: Or(
				Eq("Name", "Bob"),
				Eq("Name", "Charlie"),
				Eq("Category", "Premium"),
			),
			expected: true,
		},
		{
			name: "OR - multiple conditions all false",
			predicate: Or(
				Eq("Name", "Bob"),
				Eq("Category", "Standard"),
				Eq("Age", 25),
			),
			expected: false,
		},
		{
			name:      "OR - empty predicates",
			predicate: Or(),
			expected:  true,
		},
		{
			name: "OR - single predicate true",
			predicate: Or(
				Eq("Name", "Alice"),
			),
			expected: true,
		},
		{
			name: "OR - single predicate false",
			predicate: Or(
				Eq("Name", "Bob"),
			),
			expected: false,
		},
		{
			name: "OR - comparison operators",
			predicate: Or(
				Lt("Age", 25),
				Gt("Age", 35),
			),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(entity, nil)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCompoundPredicate_Matches_Nested tests nested compound predicates
func TestCompoundPredicate_Matches_Nested(t *testing.T) {
	entity := TestEntity{
		ID:       "123",
		Name:     "Alice",
		Age:      30,
		Active:   true,
		Category: "Premium",
		Score:    95.5,
	}

	tests := []struct {
		name      string
		predicate *CompoundPredicate
		expected  bool
	}{
		{
			name: "AND with nested OR - match",
			predicate: And(
				Eq("Active", true),
				Or(
					Eq("Category", "Premium"),
					Eq("Category", "Gold"),
				),
			),
			expected: true,
		},
		{
			name: "AND with nested OR - no match",
			predicate: And(
				Eq("Active", true),
				Or(
					Eq("Category", "Standard"),
					Eq("Category", "Gold"),
				),
			),
			expected: false,
		},
		{
			name: "OR with nested AND - match",
			predicate: Or(
				And(
					Eq("Name", "Bob"),
					Eq("Age", 25),
				),
				And(
					Eq("Name", "Alice"),
					Gte("Age", 25),
				),
			),
			expected: true,
		},
		{
			name: "OR with nested AND - no match",
			predicate: Or(
				And(
					Eq("Name", "Bob"),
					Eq("Age", 25),
				),
				And(
					Eq("Name", "Charlie"),
					Gte("Age", 40),
				),
			),
			expected: false,
		},
		{
			name: "Complex nested - (Active AND (Premium OR Gold))",
			predicate: And(
				Eq("Active", true),
				Or(
					Eq("Category", "Premium"),
					Eq("Category", "Gold"),
				),
			),
			expected: true,
		},
		{
			name: "Complex nested - ((Age >= 25 AND Age <= 35) AND Active)",
			predicate: And(
				And(
					Gte("Age", 25),
					Lte("Age", 35),
				),
				Eq("Active", true),
			),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(entity, nil)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCompoundPredicate_Matches_MixedOperators tests compound predicates with various operators
func TestCompoundPredicate_Matches_MixedOperators(t *testing.T) {
	entity := TestEntity{
		ID:       "123",
		Name:     "Alice Johnson",
		Age:      30,
		Active:   true,
		Category: "Premium",
	}

	tests := []struct {
		name      string
		predicate *CompoundPredicate
		expected  bool
	}{
		{
			name: "AND with comparison and string match",
			predicate: And(
				Gte("Age", 25),
				Contains("Name", "Alice"),
			),
			expected: true,
		},
		{
			name: "AND with IN operator",
			predicate: And(
				In("Category", "Premium", "Gold"),
				Eq("Active", true),
			),
			expected: true,
		},
		{
			name: "OR with string patterns",
			predicate: Or(
				StartsWith("Name", "Bob"),
				EndsWith("Name", "Johnson"),
			),
			expected: true,
		},
		{
			name: "AND with NOT IN",
			predicate: And(
				NotIn("Category", "Standard", "Basic"),
				Gte("Age", 25),
			),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(entity, nil)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAtomicPredicate_Matches_CaseInsensitiveFieldNames tests case-insensitive field access in predicates
func Test_Matches_CaseInsensitiveFieldNames(t *testing.T) {
	entity := TestEntity{
		ID:       "123",
		Name:     "Alice",
		Age:      30,
		Score:    95.5,
		Active:   true,
		Email:    "alice@example.com",
		Category: "Premium",
		Count:    5,
	}

	tests := []struct {
		name      string
		predicate *AtomicPredicate
		expected  bool
	}{
		{
			name:      "lowercase field name - string equality",
			predicate: Eq("name", "Alice"),
			expected:  true,
		},
		{
			name:      "uppercase field name - string equality",
			predicate: Eq("NAME", "Alice"),
			expected:  true,
		},
		{
			name:      "mixed case field name - string equality",
			predicate: Eq("NaMe", "Alice"),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(entity, nil)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAtomicPredicate_Matches_TypeConversion tests type conversion in matching
func TestAtomicPredicate_Matches_TypeConversion(t *testing.T) {
	// Using struct with pointer to test type handling
	type TestEntityWithPointer struct {
		Count *int
		Score *float64
	}

	count := 5
	score := 95.5
	entity := TestEntityWithPointer{
		Count: &count,
		Score: &score,
	}

	tests := []struct {
		name      string
		predicate *AtomicPredicate
		expected  bool
	}{
		{
			name:      "pointer int equality",
			predicate: Eq("Count", 5),
			expected:  true,
		},
		{
			name:      "pointer float equality",
			predicate: Eq("Score", 95.5),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(entity, nil)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAtomicPredicate_Matches_MapProperties tests predicates against map[string]any properties
func TestAtomicPredicate_Matches_MapProperties(t *testing.T) {
	type MapEntity struct {
		ID         string
		Properties map[string]any
	}

	entity := MapEntity{
		ID: "entity-1",
		Properties: map[string]any{
			"Foo":      "bar",
			"Status":   "active",
			"Count":    42,
			"Rating":   4.5,
			"Verified": true,
		},
	}

	tests := []struct {
		name      string
		predicate *AtomicPredicate
		expected  bool
	}{
		{
			name:      "map string equality match",
			predicate: Eq("Properties.Foo", "bar"),
			expected:  true,
		},
		{
			name:      "map string equality no match",
			predicate: Eq("Properties.Foo", "baz"),
			expected:  false,
		},
		{
			name:      "map int equality match",
			predicate: Eq("Properties.Count", 42),
			expected:  true,
		},
		{
			name:      "map int greater than",
			predicate: Gt("Properties.Count", 40),
			expected:  true,
		},
		{
			name:      "map int less than",
			predicate: Lt("Properties.Count", 50),
			expected:  true,
		},
		{
			name:      "map float equality",
			predicate: Eq("Properties.Rating", 4.5),
			expected:  true,
		},
		{
			name:      "map bool equality true",
			predicate: Eq("Properties.Verified", true),
			expected:  true,
		},
		{
			name:      "map bool equality false",
			predicate: Eq("Properties.Verified", false),
			expected:  false,
		},
		{
			name:      "map string contains",
			predicate: Contains("Properties.Status", "ctiv"),
			expected:  true,
		},
		{
			name:      "map string IN operator",
			predicate: In("Properties.Status", "active", "pending"),
			expected:  true,
		},
		{
			name:      "map nonexistent key",
			predicate: Eq("Properties.Missing", "value"),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(entity, nil)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCompoundPredicate_Matches_MapPropertiesWithLogic tests compound predicates with map properties
func TestCompoundPredicate_Matches_MapPropertiesWithLogic(t *testing.T) {
	type Resource struct {
		Name       string
		Attributes map[string]any
	}

	resource := Resource{
		Name: "TestResource",
		Attributes: map[string]any{
			"Type":     "compute",
			"Replicas": 3,
			"Region":   "us-east-1",
			"Enabled":  true,
		},
	}

	tests := []struct {
		name      string
		predicate *CompoundPredicate
		expected  bool
	}{
		{
			name: "AND with two map properties",
			predicate: And(
				Eq("Attributes.Type", "compute"),
				Eq("Attributes.Enabled", true),
			),
			expected: true,
		},
		{
			name: "AND with one false map property",
			predicate: And(
				Eq("Attributes.Type", "compute"),
				Eq("Attributes.Enabled", false),
			),
			expected: false,
		},
		{
			name: "OR with map properties",
			predicate: Or(
				Eq("Attributes.Region", "us-west-2"),
				Eq("Attributes.Type", "compute"),
			),
			expected: true,
		},
		{
			name: "AND with map comparison and regular field",
			predicate: And(
				Eq("Name", "TestResource"),
				Gt("Attributes.Replicas", 2),
			),
			expected: true,
		},
		{
			name: "nested AND/OR with map properties",
			predicate: And(
				Eq("Attributes.Enabled", true),
				Or(
					Eq("Attributes.Type", "storage"),
					Eq("Attributes.Type", "compute"),
				),
			),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(resource, nil)
			if result != tt.expected {
				t.Errorf("Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestAtomicPredicate_StringAliasNormalization verifies that string type aliases are properly normalized during predicate matching
func TestAtomicPredicate_StringAliasNormalization(t *testing.T) {
	type Status string

	// Create a test object with a plain string field
	testObj := struct {
		Name   string
		Status string
	}{
		Name:   "test",
		Status: "active",
	}

	matcher := &DefaultFieldMatcher{}

	tests := []struct {
		name      string
		predicate *AtomicPredicate
		expected  bool
	}{
		{
			name: "string alias matches plain string field",
			predicate: &AtomicPredicate{
				Field:    "Status",
				Operator: OpEqual,
				Value:    Status("active"), // String alias
			},
			expected: true,
		},
		{
			name: "string alias does not match different value",
			predicate: &AtomicPredicate{
				Field:    "Status",
				Operator: OpEqual,
				Value:    Status("inactive"), // String alias with different value
			},
			expected: false,
		},
		{
			name: "string alias with contains operator",
			predicate: &AtomicPredicate{
				Field:    "Status",
				Operator: OpContains,
				Value:    Status("act"), // String alias substring
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.predicate.Matches(testObj, matcher)
			if result != tt.expected {
				t.Errorf("predicate.Matches() = %v, want %v", result, tt.expected)
			}
		})
	}
}
