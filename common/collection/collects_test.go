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

package collection

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectAll(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		seq := func(yield func(int, error) bool) {
			// Empty sequence - don't yield anything
		}

		result, err := CollectAll(seq)

		require.NoError(t, err)
		assert.Empty(t, result)
		assert.NotNil(t, result) // Should be empty slice, not nil
	})

	t.Run("multiple items sequence", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5}
		seq := func(yield func(int, error) bool) {
			for _, item := range items {
				if !yield(item, nil) {
					break
				}
			}
		}

		result, err := CollectAll(seq)

		require.NoError(t, err)
		require.Len(t, result, 5)
		assert.Equal(t, items, result)
	})

	t.Run("sequence with error at beginning", func(t *testing.T) {
		expectedErr := errors.New("test error")
		seq := func(yield func(string, error) bool) {
			yield("", expectedErr)
		}

		result, err := CollectAll(seq)

		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
	})

	t.Run("sequence with error in middle", func(t *testing.T) {
		expectedErr := errors.New("middle error")
		seq := func(yield func(string, error) bool) {
			if !yield("item1", nil) {
				return
			}
			if !yield("item2", nil) {
				return
			}
			if !yield("", expectedErr) {
				return
			} // Error in the middle
			if !yield("item4", nil) {
				return
			} // This should not be reached
		}

		result, err := CollectAll(seq)

		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result) // Should return nil on error
	})

	t.Run("sequence with error at end", func(t *testing.T) {
		expectedErr := errors.New("end error")
		seq := func(yield func(int, error) bool) {
			yield(1, nil)
			yield(2, nil)
			yield(3, nil)
			yield(0, expectedErr) // Error at end
		}

		result, err := CollectAll(seq)

		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, result)
	})

	t.Run("sequence with nil pointer values", func(t *testing.T) {
		type TestEntity struct {
			Value string
		}

		seq := func(yield func(*TestEntity, error) bool) {
			yield(nil, nil)                        // Nil pointer
			yield(&TestEntity{Value: "test"}, nil) // Valid pointer
			yield(nil, nil)                        // Nil pointer again
		}

		result, err := CollectAll(seq)

		require.NoError(t, err)
		require.Len(t, result, 3)
		assert.Nil(t, result[0])
		assert.NotNil(t, result[1])
		assert.Equal(t, "test", result[1].Value)
		assert.Nil(t, result[2])
	})
}
