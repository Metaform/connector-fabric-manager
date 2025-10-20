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

package api

import (
	"context"
	"iter"

	"github.com/metaform/connector-fabric-manager/common/system"
)

const (
	ParticipantProfileStoreKey system.ServiceType = "tmstore:ParticipantProfileStore"
	CellStoreKey               system.ServiceType = "tmstore:CellStore"
	DataspaceProfileStoreKey   system.ServiceType = "tmstore:DataspaceProfileStore"
)

// PaginationOptions defines pagination parameters for entity retrieval.
type PaginationOptions struct {
	// Offset is the number of items to skip from the beginning.
	Offset int
	// Limit is the maximum number of items to return. If 0, returns all items.
	Limit int
	// Cursor is an optional cursor for cursor-based pagination (implementation-specific).
	Cursor string
}

// DefaultPaginationOptions returns default pagination settings (no pagination).
func DefaultPaginationOptions() PaginationOptions {
	return PaginationOptions{
		Offset: 0,
		Limit:  0, // 0 means no limit
		Cursor: "",
	}
}

// EntityStore defines the interface for entity storage.
type EntityStore[T any] interface {
	FindById(ctx context.Context, id string) (*T, error)
	Exists(ctx context.Context, id string) (bool, error)
	Create(ctx context.Context, entity *T) (*T, error)
	Update(ctx context.Context, entity *T) error
	GetAll(ctx context.Context) iter.Seq2[T, error]
	GetAllPaginated(ctx context.Context, opts PaginationOptions) iter.Seq2[T, error]
}
