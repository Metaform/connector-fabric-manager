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
	"time"

	"github.com/stretchr/testify/assert"
)

// TestTimePredicateSupport verifies that predicates work with time.Time fields
type TimeTestEntity struct {
	ID               string
	StateTimestamp   time.Time
	CreatedTimestamp time.Time
}

func TestTimeEqualityPredicate(t *testing.T) {
	referenceTime := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)

	entity := TimeTestEntity{
		ID:               "test-1",
		StateTimestamp:   referenceTime,
		CreatedTimestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	matcher := &DefaultFieldMatcher{}

	// Test equality
	predicate := Eq("StateTimestamp", referenceTime)
	assert.True(t, predicate.Matches(entity, matcher), "Equality predicate should match exact time")

	// Test inequality with different time
	differentTime := time.Date(2025, 1, 16, 10, 30, 0, 0, time.UTC)
	predicate = Eq("StateTimestamp", differentTime)
	assert.False(t, predicate.Matches(entity, matcher), "Equality predicate should not match different time")
}

func TestTimeGreaterThanPredicate(t *testing.T) {
	baseTime := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	laterTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)
	earlierTime := time.Date(2025, 1, 10, 10, 30, 0, 0, time.UTC)

	entity := TimeTestEntity{
		ID:               "test-1",
		StateTimestamp:   baseTime,
		CreatedTimestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	matcher := &DefaultFieldMatcher{}

	// Test greater than with earlier time (should match)
	predicate := Gt("StateTimestamp", earlierTime)
	assert.True(t, predicate.Matches(entity, matcher), "Time should be greater than earlier time")

	// Test greater than with later time (should not match)
	predicate = Gt("StateTimestamp", laterTime)
	assert.False(t, predicate.Matches(entity, matcher), "Time should not be greater than later time")

	// Test greater than with same time (should not match)
	predicate = Gt("StateTimestamp", baseTime)
	assert.False(t, predicate.Matches(entity, matcher), "Time should not be greater than itself")
}

func TestTimeGreaterOrEqualPredicate(t *testing.T) {
	baseTime := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	laterTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)
	earlierTime := time.Date(2025, 1, 10, 10, 30, 0, 0, time.UTC)

	entity := TimeTestEntity{
		ID:               "test-1",
		StateTimestamp:   baseTime,
		CreatedTimestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	matcher := &DefaultFieldMatcher{}

	// Test >= with earlier time (should match)
	predicate := Gte("StateTimestamp", earlierTime)
	assert.True(t, predicate.Matches(entity, matcher), "Time should be >= earlier time")

	// Test >= with same time (should match)
	predicate = Gte("StateTimestamp", baseTime)
	assert.True(t, predicate.Matches(entity, matcher), "Time should be >= itself")

	// Test >= with later time (should not match)
	predicate = Gte("StateTimestamp", laterTime)
	assert.False(t, predicate.Matches(entity, matcher), "Time should not be >= later time")
}

func TestTimeLessThanPredicate(t *testing.T) {
	baseTime := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	laterTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)
	earlierTime := time.Date(2025, 1, 10, 10, 30, 0, 0, time.UTC)

	entity := TimeTestEntity{
		ID:               "test-1",
		StateTimestamp:   baseTime,
		CreatedTimestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	matcher := &DefaultFieldMatcher{}

	// Test less than with later time (should match)
	predicate := Lt("StateTimestamp", laterTime)
	assert.True(t, predicate.Matches(entity, matcher), "Time should be < later time")

	// Test less than with earlier time (should not match)
	predicate = Lt("StateTimestamp", earlierTime)
	assert.False(t, predicate.Matches(entity, matcher), "Time should not be < earlier time")

	// Test less than with same time (should not match)
	predicate = Lt("StateTimestamp", baseTime)
	assert.False(t, predicate.Matches(entity, matcher), "Time should not be < itself")
}

func TestTimeLessOrEqualPredicate(t *testing.T) {
	baseTime := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	laterTime := time.Date(2025, 1, 20, 10, 30, 0, 0, time.UTC)
	earlierTime := time.Date(2025, 1, 10, 10, 30, 0, 0, time.UTC)

	entity := TimeTestEntity{
		ID:               "test-1",
		StateTimestamp:   baseTime,
		CreatedTimestamp: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	matcher := &DefaultFieldMatcher{}

	// Test <= with later time (should match)
	predicate := Lte("StateTimestamp", laterTime)
	assert.True(t, predicate.Matches(entity, matcher), "Time should be <= later time")

	// Test <= with same time (should match)
	predicate = Lte("StateTimestamp", baseTime)
	assert.True(t, predicate.Matches(entity, matcher), "Time should be <= itself")

	// Test <= with earlier time (should not match)
	predicate = Lte("StateTimestamp", earlierTime)
	assert.False(t, predicate.Matches(entity, matcher), "Time should not be <= earlier time")
}

func TestTimeCompoundPredicates(t *testing.T) {
	startTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	midTime := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	endTime := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	entity := TimeTestEntity{
		ID:               "test-1",
		StateTimestamp:   midTime,
		CreatedTimestamp: startTime,
	}

	matcher := &DefaultFieldMatcher{}

	// Test AND: StateTimestamp >= startTime AND StateTimestamp <= endTime
	predicate := And(
		Gte("StateTimestamp", startTime),
		Lte("StateTimestamp", endTime),
	)
	assert.True(t, predicate.Matches(entity, matcher), "Time should be within range using AND")

	// Test AND with failing condition
	predicate = And(
		Gte("StateTimestamp", startTime),
		Lt("StateTimestamp", midTime),
	)
	assert.False(t, predicate.Matches(entity, matcher), "AND should fail if any condition fails")

	// Test OR: StateTimestamp > endTime OR StateTimestamp < startTime
	predicate = Or(
		Gt("StateTimestamp", endTime),
		Lt("StateTimestamp", startTime),
	)
	assert.False(t, predicate.Matches(entity, matcher), "OR should fail if all conditions fail")

	// Test OR with passing condition
	predicate = Or(
		Gt("StateTimestamp", endTime),
		Lt("StateTimestamp", endTime),
	)
	assert.True(t, predicate.Matches(entity, matcher), "OR should pass if any condition passes")
}
