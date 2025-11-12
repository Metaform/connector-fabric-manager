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
	"strings"
)

// Predicate defines a query condition that can be evaluated in memory or converted into a query language.
type Predicate interface {
	// Matches evaluates the predicate against an object using the provided matcher
	Matches(obj any, matcher FieldMatcher) bool
	// String returns a readable representation
	String() string
}

// Field represents a queryable Field
type Field string

// Operator defines comparison operators
type Operator string

const (
	OpEqual        Operator = "="
	OpNotEqual     Operator = "!="
	OpGreater      Operator = ">"
	OpGreaterEqual Operator = ">="
	OpLess         Operator = "<"
	OpLessEqual    Operator = "<="
	OpIn           Operator = "IN"
	OpNotIn        Operator = "NOT IN"
	OpLike         Operator = "LIKE"
	OpNotLike      Operator = "NOT LIKE"
	OpIsNull       Operator = "IS NULL"
	OpIsNotNull    Operator = "IS NOT NULL"
	OpContains     Operator = "CONTAINS"
	OpStartsWith   Operator = "STARTS_WITH"
	OpEndsWith     Operator = "ENDS_WITH"
)

// FieldMatcher is a strategy for extracting and comparing Field values from objects
// Implementations provide custom logic for specific types
type FieldMatcher interface {
	// GetFieldValue extracts a Value from an object by Field name
	GetFieldValue(obj any, fieldName string) (any, error)
	// CompareValues compares two values based on the Operator
	CompareValues(op Operator, fieldValue, compareValue any) bool
}

// DefaultFieldMatcher provides generic reflection-based Field matching
type DefaultFieldMatcher struct{}

func (m *DefaultFieldMatcher) GetFieldValue(obj any, fieldName string) (any, error) {
	return GetFieldValue(obj, fieldName)
}

func (m *DefaultFieldMatcher) CompareValues(op Operator, fieldValue, compareValue any) bool {
	return CompareValues(op, fieldValue, compareValue)
}

// AtomicPredicate is a basic Field comparison predicate
type AtomicPredicate struct {
	Field    Field    `json:"field"`
	Operator Operator `json:"operator"`
	Value    any      `json:"value"`
}

// UnmarshalJSON deserializes an AtomicPredicate from JSON, normalizing the operator to uppercase
func (p *AtomicPredicate) UnmarshalJSON(data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if field, ok := m["field"].(string); ok {
		p.Field = Field(field)
	}
	if op, ok := m["operator"].(string); ok {
		p.Operator = Operator(strings.ToUpper(op))
	}
	p.Value = m["value"]

	return nil
}

// Eq creates a predicate for equality (syntactic sugar)
func Eq(field Field, value any) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: OpEqual,
		Value:    value,
	}
}

// Neq creates a not-equal predicate
func Neq(field Field, value any) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: OpNotEqual,
		Value:    value,
	}
}

// Gt creates a greater-than predicate
func Gt(field Field, value any) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: OpGreater,
		Value:    value,
	}
}

// Gte creates a greater-than-or-equal predicate
func Gte(field Field, value any) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: OpGreaterEqual,
		Value:    value,
	}
}

// Lt creates a less-than predicate
func Lt(field Field, value any) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: OpLess,
		Value:    value,
	}
}

// Lte creates a less-than-or-equal predicate
func Lte(field Field, value any) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: OpLessEqual,
		Value:    value,
	}
}

// In creates an IN predicate
func In(field Field, values ...any) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: OpIn,
		Value:    values,
	}
}

// NotIn creates a NOT IN predicate
func NotIn(field Field, values ...any) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: OpNotIn,
		Value:    values,
	}
}

// Like creates a LIKE predicate
func Like(field Field, pattern string) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: OpLike,
		Value:    pattern,
	}
}

// Contains creates a CONTAINS predicate (for in-memory substring matching)
func Contains(field Field, substring string) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: OpContains,
		Value:    substring,
	}
}

// StartsWith creates a STARTS_WITH predicate
func StartsWith(field Field, prefix string) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: OpStartsWith,
		Value:    prefix,
	}
}

// EndsWith creates an ENDS_WITH predicate
func EndsWith(field Field, suffix string) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: OpEndsWith,
		Value:    suffix,
	}
}

// IsNull creates an IS NULL predicate
func IsNull(field Field) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: OpIsNull,
		Value:    nil,
	}
}

// IsNotNull creates an IS NOT NULL predicate
func IsNotNull(field Field) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: OpIsNotNull,
		Value:    nil,
	}
}

// NewComparison creates a predicate with a specific Operator (for advanced cases)
func NewComparison(field Field, op Operator, value any) *AtomicPredicate {
	return &AtomicPredicate{
		Field:    field,
		Operator: op,
		Value:    value,
	}
}

func (p *AtomicPredicate) Matches(obj any, matcher FieldMatcher) bool {
	if matcher == nil {
		matcher = &DefaultFieldMatcher{}
	}

	fieldValue, err := matcher.GetFieldValue(obj, string(p.Field))

	// For NULL checks, treat non-existent fields (errors) as nil
	if err != nil {
		if p.Operator == OpIsNull {
			return true
		}
		if p.Operator == OpIsNotNull {
			return false
		}
		return false
	}

	return matcher.CompareValues(p.Operator, fieldValue, p.Value)
}

func (p *AtomicPredicate) String() string {
	switch p.Operator {
	case OpIsNull, OpIsNotNull:
		return fmt.Sprintf("%s %s", p.Field, p.Operator)
	default:
		return fmt.Sprintf("%s %s %v", p.Field, p.Operator, p.Value)
	}
}

// CompoundPredicate combines multiple Predicates with AND/OR logic
type CompoundPredicate struct {
	Operator   string      `json:"operator"`
	Predicates []Predicate `json:"predicates"`
}

// And creates an AND conjunction of Predicates
func And(predicates ...Predicate) *CompoundPredicate {
	return &CompoundPredicate{
		Predicates: predicates,
		Operator:   "AND",
	}
}

// Or creates an OR conjunction of Predicates
func Or(predicates ...Predicate) *CompoundPredicate {
	return &CompoundPredicate{
		Predicates: predicates,
		Operator:   "OR",
	}
}

func (p *CompoundPredicate) Matches(obj any, matcher FieldMatcher) bool {
	if len(p.Predicates) == 0 {
		return true
	}

	for _, pred := range p.Predicates {
		matches := pred.Matches(obj, matcher)
		if p.Operator == "AND" && !matches {
			return false
		}
		if p.Operator == "OR" && matches {
			return true
		}
	}

	if p.Operator == "AND" {
		return true
	}
	return false
}

func (p *CompoundPredicate) String() string {
	parts := make([]string, len(p.Predicates))
	for i, pred := range p.Predicates {
		parts[i] = pred.String()
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, fmt.Sprintf(" %s ", p.Operator)))
}

// UnmarshalJSON deserializes a CompoundPredicate from JSON
func (p *CompoundPredicate) UnmarshalJSON(data []byte) error {
	var aux struct {
		Operator   string            `json:"operator"`
		Predicates []json.RawMessage `json:"predicates"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	p.Operator = aux.Operator
	p.Predicates = make([]Predicate, len(aux.Predicates))

	for i, rawPred := range aux.Predicates {
		pred, err := UnmarshalPredicate(rawPred)
		if err != nil {
			return fmt.Errorf("failed to unmarshal predicate at index %d: %w", i, err)
		}
		p.Predicates[i] = pred
	}

	return nil
}

// UnmarshalPredicate unmarshals a Predicate from JSON bytes into the Predicate interface
func UnmarshalPredicate(data []byte) (Predicate, error) {
	// Parse JSON into a generic object to peek at structure
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal predicate: %w", err)
	}

	// If predicates field exists, it's a compound predicate
	if _, hasPredicates := obj["predicates"]; hasPredicates {
		var compound CompoundPredicate
		if err := json.Unmarshal(data, &compound); err != nil {
			return nil, err
		}
		return &compound, nil
	}

	// Otherwise, it's an atomic predicate
	var atomic AtomicPredicate
	if err := json.Unmarshal(data, &atomic); err != nil {
		return nil, err
	}
	return &atomic, nil
}

// UnmarshalPredicates unmarshals a slice of Predicates from JSON bytes
func UnmarshalPredicates(data []byte) ([]Predicate, error) {
	var rawPredicates []json.RawMessage
	if err := json.Unmarshal(data, &rawPredicates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal predicate array: %w", err)
	}

	predicates := make([]Predicate, len(rawPredicates))
	for i, rawPred := range rawPredicates {
		pred, err := UnmarshalPredicate(rawPred)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal predicate at index %d: %w", i, err)
		}
		predicates[i] = pred
	}

	return predicates, nil
}
