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
	"encoding/json"
	"fmt"
	"testing"
)

// TestAtomicPredicate_JSON tests JSON serialization/deserialization of AtomicPredicate
func TestAtomicPredicate_JSON(t *testing.T) {
	tests := []struct {
		name      string
		predicate *AtomicPredicate
	}{
		{
			name:      "Eq predicate",
			predicate: Eq("Name", "Alice"),
		},
		{
			name:      "Gt predicate",
			predicate: Gt("Age", 30),
		},
		{
			name:      "In predicate",
			predicate: In("Status", "active", "pending", "approved"),
		},
		{
			name:      "Contains predicate",
			predicate: Contains("Email", "@example.com"),
		},
		{
			name:      "IsNull predicate",
			predicate: IsNull("DeletedAt"),
		},
		{
			name:      "Like predicate",
			predicate: Like("Name", "john"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(tt.predicate)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			// Unmarshal back
			var unmarshaled AtomicPredicate
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			// Verify fields match
			if unmarshaled.Field != tt.predicate.Field {
				t.Errorf("Field mismatch: got %v, want %v", unmarshaled.Field, tt.predicate.Field)
			}
			if unmarshaled.Operator != tt.predicate.Operator {
				t.Errorf("Operator mismatch: got %v, want %v", unmarshaled.Operator, tt.predicate.Operator)
			}
		})
	}
}

// TestCompoundPredicate_JSON tests JSON serialization/deserialization of CompoundPredicate
func TestCompoundPredicate_JSON(t *testing.T) {
	tests := []struct {
		name      string
		predicate *CompoundPredicate
	}{
		{
			name: "Simple AND",
			predicate: And(
				Eq("Name", "Alice"),
				Eq("Active", true),
			),
		},
		{
			name: "Simple OR",
			predicate: Or(
				Eq("Status", "pending"),
				Eq("Status", "approved"),
			),
		},
		{
			name: "Nested AND-OR",
			predicate: And(
				Eq("Active", true),
				Or(
					Eq("Role", "admin"),
					Eq("Role", "moderator"),
				),
			),
		},
		{
			name: "Complex nested",
			predicate: Or(
				And(
					Gte("Age", 18),
					Lte("Age", 65),
				),
				And(
					Eq("Status", "student"),
					Eq("Active", true),
				),
			),
		},
		{
			name: "Deep nesting",
			predicate: And(
				Eq("Active", true),
				Or(
					And(
						Gte("Age", 18),
						Lte("Age", 65),
					),
					And(
						Eq("Status", "student"),
						Contains("Email", "@school.edu"),
					),
				),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(tt.predicate)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			// Unmarshal back
			var unmarshaled CompoundPredicate
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			// Verify basic properties
			if unmarshaled.Operator != tt.predicate.Operator {
				t.Errorf("Operator mismatch: got %s, want %s", unmarshaled.Operator, tt.predicate.Operator)
			}
			if len(unmarshaled.Predicates) != len(tt.predicate.Predicates) {
				t.Errorf("Predicates length mismatch: got %d, want %d", len(unmarshaled.Predicates), len(tt.predicate.Predicates))
			}
		})
	}
}

// TestUnmarshalPredicate tests the UnmarshalPredicate helper function with both types
func TestUnmarshalPredicate(t *testing.T) {
	tests := []struct {
		name        string
		predicate   Predicate
		testEntity  TestEntity
		shouldMatch bool
	}{
		{
			name:        "Atomic predicate",
			predicate:   Eq("Name", "Alice"),
			shouldMatch: true,
			testEntity:  TestEntity{Name: "Alice"},
		},
		{
			name: "Compound predicate",
			predicate: And(
				Eq("Name", "Alice"),
				Gt("Age", 25),
			),
			shouldMatch: true,
			testEntity:  TestEntity{Name: "Alice", Age: 30},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the predicate
			jsonData, err := json.Marshal(tt.predicate)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			// Use helper function to unmarshal
			unmarshaled, err := UnmarshalPredicate(jsonData)
			if err != nil {
				t.Fatalf("UnmarshalPredicate() error = %v", err)
			}

			// Verify it works with matching
			if result := unmarshaled.Matches(tt.testEntity, nil); result != tt.shouldMatch {
				t.Errorf("Matches() = %v, want %v", result, tt.shouldMatch)
			}
		})
	}
}

// TestUnmarshalPredicates tests the UnmarshalPredicates helper function
func TestUnmarshalPredicates(t *testing.T) {
	predicates := []Predicate{
		Eq("Name", "Alice"),
		Gt("Age", 25),
		And(Eq("Active", true), Eq("Role", "admin")),
	}

	// Marshal array of predicates
	jsonData, err := json.Marshal(predicates)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Unmarshal using helper function
	unmarshaled, err := UnmarshalPredicates(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalPredicates() error = %v", err)
	}

	// Verify count
	if len(unmarshaled) != len(predicates) {
		t.Errorf("Predicates length mismatch: got %d, want %d", len(unmarshaled), len(predicates))
	}
}

// TestPredicate_JSONRoundTrip_WithMatching tests that JSON roundtrip preserves matching behavior
func TestPredicate_JSONRoundTrip_WithMatching(t *testing.T) {
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

	predicates := []Predicate{
		Eq("Name", "Alice"),
		And(Eq("Name", "Alice"), Gte("Age", 25)),
		Or(Eq("Category", "Premium"), Eq("Category", "Gold")),
		And(
			Eq("Active", true),
			Or(Gte("Age", 25), Eq("Category", "Student")),
		),
	}

	for i, pred := range predicates {
		t.Run(fmt.Sprintf("Predicate_%d", i), func(t *testing.T) {
			// Get original match result
			originalMatch := pred.Matches(entity, nil)

			// Marshal and unmarshal
			jsonData, err := json.Marshal(pred)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			unmarshaled, err := UnmarshalPredicate(jsonData)
			if err != nil {
				t.Fatalf("UnmarshalPredicate() error = %v", err)
			}

			// Get unmarshaled match result
			unmarshaledMatch := unmarshaled.Matches(entity, nil)

			// Should match
			if originalMatch != unmarshaledMatch {
				t.Errorf("Match behavior changed after roundtrip: original=%v, unmarshaled=%v", originalMatch, unmarshaledMatch)
			}
		})
	}
}

// TestPredicate_JSON_WithVariousTypes tests JSON serialization with various value types
func TestPredicate_JSON_WithVariousTypes(t *testing.T) {
	tests := []struct {
		name      string
		predicate *AtomicPredicate
	}{
		{
			name:      "String value",
			predicate: Eq("Name", "Alice"),
		},
		{
			name:      "Int value",
			predicate: Gt("Age", 30),
		},
		{
			name:      "Float value",
			predicate: Lt("Score", 95.5),
		},
		{
			name:      "Bool value",
			predicate: Eq("Active", true),
		},
		{
			name:      "Nil value",
			predicate: IsNull("DeletedAt"),
		},
		{
			name:      "Array value",
			predicate: In("Status", "active", "pending"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.predicate)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			var unmarshaled AtomicPredicate
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			if unmarshaled.Field != tt.predicate.Field {
				t.Errorf("Field mismatch after JSON roundtrip")
			}
			if unmarshaled.Operator != tt.predicate.Operator {
				t.Errorf("Operator mismatch after JSON roundtrip")
			}
		})
	}
}

// TestAtomicPredicate_JSON_Structure tests that the JSON output is clean
func TestAtomicPredicate_JSON_Structure(t *testing.T) {
	predicate := Eq("Name", "Alice")
	jsonData, err := json.Marshal(predicate)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var result map[string]any
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	// Verify only expected fields
	if result["field"] != "Name" {
		t.Errorf("field = %v, want Name", result["field"])
	}
	if result["operator"] != "=" {
		t.Errorf("operator = %v, want =", result["operator"])
	}
	if result["value"] != "Alice" {
		t.Errorf("value = %v, want Alice", result["value"])
	}
}

// TestCompoundPredicate_JSON_Structure tests that compound predicates have predicates field
func TestCompoundPredicate_JSON_Structure(t *testing.T) {
	predicate := And(Eq("Name", "Alice"), Gt("Age", 30))
	jsonData, err := json.Marshal(predicate)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var result map[string]any
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Fatalf("Unmarshal to map failed: %v", err)
	}

	// Verify structure has predicates field
	if _, hasPredicates := result["predicates"]; !hasPredicates {
		t.Errorf("CompoundPredicate must have 'predicates' field")
	}
	if result["operator"] != "AND" {
		t.Errorf("operator = %v, want AND", result["operator"])
	}
}

// TestUnmarshalPredicate_Discrimination tests correct type discrimination
func TestUnmarshalPredicate_Discrimination(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		wantType string
	}{
		{
			name:     "Atomic (no predicates field)",
			json:     `{"field":"Name","operator":"=","value":"Alice"}`,
			wantType: "atomic",
		},
		{
			name:     "Compound (has predicates field)",
			json:     `{"operator":"AND","predicates":[{"field":"A","operator":"=","value":"B"}]}`,
			wantType: "compound",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pred, err := UnmarshalPredicate([]byte(tt.json))
			if err != nil {
				t.Fatalf("UnmarshalPredicate() error = %v", err)
			}

			switch tt.wantType {
			case "atomic":
				if _, ok := pred.(*AtomicPredicate); !ok {
					t.Errorf("Expected *AtomicPredicate, got %T", pred)
				}
			case "compound":
				if _, ok := pred.(*CompoundPredicate); !ok {
					t.Errorf("Expected *CompoundPredicate, got %T", pred)
				}
			}
		})
	}
}

// TestUnmarshalPredicate_InvalidJSON tests error handling
func TestUnmarshalPredicate_InvalidJSON(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    []byte
		expectError bool
	}{
		{
			name:        "Invalid JSON",
			jsonData:    []byte("{invalid}"),
			expectError: true,
		},
		{
			name:        "Malformed nested predicates",
			jsonData:    []byte(`{"operator":"AND","predicates":["{bad}"]}`),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UnmarshalPredicate(tt.jsonData)
			if (err != nil) != tt.expectError {
				t.Errorf("UnmarshalPredicate() error = %v, expectError = %v", err, tt.expectError)
			}
		})
	}
}

// Product is a simple test struct for JSON predicate testing
type Product struct {
	ID       string
	Name     string
	Price    float64
	InStock  bool
	Quantity int
	Supplier string
	SKU      string
}

// TestAtomicPredicate_JSON_CaseInsensitiveOperators tests that operators work with any case
func TestAtomicPredicate_JSON_CaseInsensitiveOperators(t *testing.T) {
	tests := []struct {
		name             string
		json             string
		expectedOperator Operator
	}{
		{
			name:             "Lowercase operator",
			json:             `{"field":"Name","operator":"=","value":"Widget"}`,
			expectedOperator: OpEqual,
		},
		{
			name:             "Uppercase operator",
			json:             `{"field":"Name","operator":"=","value":"Widget"}`,
			expectedOperator: OpEqual,
		},
		{
			name:             "Lowercase 'in' operator",
			json:             `{"field":"Supplier","operator":"in","value":["ACME","TechCorp"]}`,
			expectedOperator: OpIn,
		},
		{
			name:             "Uppercase 'IN' operator",
			json:             `{"field":"Supplier","operator":"IN","value":["ACME","TechCorp"]}`,
			expectedOperator: OpIn,
		},
		{
			name:             "Lowercase 'like' operator",
			json:             `{"field":"Name","operator":"like","value":"widget%"}`,
			expectedOperator: OpLike,
		},
		{
			name:             "Uppercase 'LIKE' operator",
			json:             `{"field":"Name","operator":"LIKE","value":"widget%"}`,
			expectedOperator: OpLike,
		},
		{
			name:             "Mixed case 'Like' operator",
			json:             `{"field":"Name","operator":"Like","value":"widget%"}`,
			expectedOperator: OpLike,
		},
		{
			name:             "Lowercase 'not in' operator",
			json:             `{"field":"Supplier","operator":"not in","value":["Discontinued"]}`,
			expectedOperator: OpNotIn,
		},
		{
			name:             "Uppercase 'NOT IN' operator",
			json:             `{"field":"Supplier","operator":"NOT IN","value":["Discontinued"]}`,
			expectedOperator: OpNotIn,
		},
		{
			name:             "Lowercase 'starts_with' operator",
			json:             `{"field":"SKU","operator":"starts_with","value":"PROD"}`,
			expectedOperator: OpStartsWith,
		},
		{
			name:             "Uppercase 'STARTS_WITH' operator",
			json:             `{"field":"SKU","operator":"STARTS_WITH","value":"PROD"}`,
			expectedOperator: OpStartsWith,
		},
		{
			name:             "Lowercase 'is null' operator",
			json:             `{"field":"SKU","operator":"is null","value":null}`,
			expectedOperator: OpIsNull,
		},
		{
			name:             "Uppercase 'IS NULL' operator",
			json:             `{"field":"SKU","operator":"IS NULL","value":null}`,
			expectedOperator: OpIsNull,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pred AtomicPredicate
			err := json.Unmarshal([]byte(tt.json), &pred)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			if pred.Operator != tt.expectedOperator {
				t.Errorf("Operator mismatch: got %v, want %v", pred.Operator, tt.expectedOperator)
			}
		})
	}
}

// TestAtomicPredicate_JSON_OperatorNormalization_WithMatching tests that normalized operators work correctly
func TestAtomicPredicate_JSON_OperatorNormalization_WithMatching(t *testing.T) {
	product := Product{
		ID:       "P001",
		Name:     "Premium Widget",
		Price:    29.99,
		InStock:  true,
		Quantity: 100,
		Supplier: "ACME Corp",
		SKU:      "PROD-12345",
	}

	tests := []struct {
		name        string
		json        string
		shouldMatch bool
	}{
		{
			name:        "Lowercase '=' operator matches",
			json:        `{"field":"Name","operator":"=","value":"Premium Widget"}`,
			shouldMatch: true,
		},
		{
			name:        "Lowercase 'in' operator matches",
			json:        `{"field":"Supplier","operator":"in","value":["ACME Corp","TechCorp"]}`,
			shouldMatch: true,
		},
		{
			name:        "Lowercase 'in' operator doesn't match",
			json:        `{"field":"Supplier","operator":"in","value":["Discontinued","Old Supplier"]}`,
			shouldMatch: false,
		},
		{
			name:        "Uppercase 'IN' operator matches",
			json:        `{"field":"Supplier","operator":"IN","value":["ACME Corp","TechCorp"]}`,
			shouldMatch: true,
		},
		{
			name:        "Lowercase 'contains' operator matches",
			json:        `{"field":"Name","operator":"contains","value":"Widget"}`,
			shouldMatch: true,
		},
		{
			name:        "Uppercase 'CONTAINS' operator matches",
			json:        `{"field":"Name","operator":"CONTAINS","value":"Widget"}`,
			shouldMatch: true,
		},
		{
			name:        "Lowercase 'starts_with' matches",
			json:        `{"field":"SKU","operator":"starts_with","value":"PROD"}`,
			shouldMatch: true,
		},
		{
			name:        "Uppercase 'STARTS_WITH' matches",
			json:        `{"field":"SKU","operator":"STARTS_WITH","value":"PROD"}`,
			shouldMatch: true,
		},
		{
			name:        "Lowercase '>' operator matches",
			json:        `{"field":"Price","operator":">","value":20.0}`,
			shouldMatch: true,
		},
		{
			name:        "Lowercase 'is not null' matches",
			json:        `{"field":"SKU","operator":"is not null","value":null}`,
			shouldMatch: true,
		},
		{
			name:        "Uppercase 'IS NOT NULL' matches",
			json:        `{"field":"SKU","operator":"IS NOT NULL","value":null}`,
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pred AtomicPredicate
			err := json.Unmarshal([]byte(tt.json), &pred)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			result := pred.Matches(product, nil)
			if result != tt.shouldMatch {
				t.Errorf("Matches() = %v, want %v", result, tt.shouldMatch)
			}
		})
	}
}

// TestMatchAllPredicateJSON tests JSON marshaling and unmarshaling of MatchAllPredicate
func TestMatchAllPredicateJSON(t *testing.T) {
	tests := []struct {
		name            string
		predicate       Predicate
		expectedJSON    string
		shouldUnmarshal bool
	}{
		{
			name:            "MatchAllPredicate marshaling",
			predicate:       &MatchAllPredicate{},
			expectedJSON:    `{"type":"matchAll"}`,
			shouldUnmarshal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			jsonBytes, err := json.Marshal(tt.predicate)
			if err != nil {
				t.Fatalf("failed to marshal predicate: %v", err)
			}

			jsonStr := string(jsonBytes)
			if jsonStr != tt.expectedJSON {
				t.Errorf("marshal result = %q, want %q", jsonStr, tt.expectedJSON)
			}

			// Test unmarshaling
			if tt.shouldUnmarshal {
				unmarshaled, err := UnmarshalPredicate(jsonBytes)
				if err != nil {
					t.Fatalf("failed to unmarshal predicate: %v", err)
				}

				// Verify it's a MatchAllPredicate
				matchAll, ok := unmarshaled.(*MatchAllPredicate)
				if !ok {
					t.Errorf("unmarshaled predicate type = %T, want *MatchAllPredicate", unmarshaled)
					return
				}

				// Verify it still matches everything
				testEntity := struct {
					Name string
					Age  int
				}{Name: "test", Age: 30}

				if !matchAll.Matches(testEntity, &DefaultFieldMatcher{}) {
					t.Errorf("expected MatchAllPredicate to match all objects after unmarshaling")
				}
			}
		})
	}
}
