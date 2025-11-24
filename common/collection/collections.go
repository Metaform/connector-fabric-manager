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

import "iter"

// CollectAll collects all sequence elements into a slice.
func CollectAll[T any](seq iter.Seq2[T, error]) ([]T, error) {
	var result []T
	for item, err := range seq {
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	if result == nil {
		return []T{}, nil
	}
	return result, nil
}

func CollectAllDeref[T any](seq iter.Seq2[*T, error]) ([]T, error) {
	var result []T
	for item, err := range seq {
		if err != nil {
			return nil, err
		}
		result = append(result, *item)
	}
	if result == nil {
		return []T{}, nil
	}
	return result, nil
}

// DerefSlice dereferences slice entry points.
func DerefSlice[T any](ptrs []*T) []T {
	values := make([]T, 0, len(ptrs))
	for _, ptr := range ptrs {
		if ptr != nil {
			values = append(values, *ptr)
		}
	}
	return values
}
