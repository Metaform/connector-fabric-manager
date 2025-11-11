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

package sqlstore

import (
	"fmt"
	"strings"

	"github.com/metaform/connector-fabric-manager/common/query"
)

// SQLBuilder converts predicates to SQL
// This allows different SQL dialects or custom SQL generation strategies
type SQLBuilder interface {
	// BuildSQL converts a predicate to SQL WHERE clause and arguments
	BuildSQL(predicate query.Predicate) (string, []any)
}

// buildSQL is the core SQL generation logic
func buildSQL(predicate query.Predicate) (string, []any) {
	switch p := predicate.(type) {
	case *query.AtomicPredicate:
		return buildSimplePredicateSQL(p)
	case *query.CompoundPredicate:
		return buildCompoundPredicateSQL(p)
	default:
		return "true", nil
	}
}

func buildSimplePredicateSQL(predicate *query.AtomicPredicate) (string, []any) {
	switch predicate.Operator {
	case query.OpIsNull:
		return fmt.Sprintf("%s IS NULL", predicate.Field), nil
	case query.OpIsNotNull:
		return fmt.Sprintf("%s IS NOT NULL", predicate.Field), nil
	case query.OpIn, query.OpNotIn:
		values := predicate.Value.([]any)
		placeholders := make([]string, len(values))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		sql := fmt.Sprintf("%s %s (%s)", predicate.Field, predicate.Operator, strings.Join(placeholders, ","))
		return sql, values
	default:
		return fmt.Sprintf("%s %s ?", predicate.Field, predicate.Operator), []any{predicate.Value}
	}
}

func buildCompoundPredicateSQL(predicate *query.CompoundPredicate) (string, []any) {
	if len(predicate.Predicates) == 0 {
		return "true", nil
	}
	if len(predicate.Predicates) == 1 {
		return buildSQL(predicate.Predicates[0])
	}

	parts := make([]string, len(predicate.Predicates))
	var args []any

	for i, pred := range predicate.Predicates {
		sql, predArgs := buildSQL(pred)
		parts[i] = fmt.Sprintf("(%s)", sql)
		args = append(args, predArgs...)
	}

	return strings.Join(parts, fmt.Sprintf(" %s ", predicate.Operator)), args
}

// DefaultSQLBuilder is the default SQL builder using ANSI SQL
type DefaultSQLBuilder struct{}

func (b *DefaultSQLBuilder) BuildSQL(predicate query.Predicate) (string, []any) {
	return buildSQL(predicate)
}
