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
	"context"
	"iter"

	"github.com/metaform/connector-fabric-manager/common/query"
	"github.com/metaform/connector-fabric-manager/common/store"
)

func NewPostgresEntityStore[T store.EntityType]() *PostgresEntityStore[T] {
	estore := &PostgresEntityStore[T]{
		matcher: &query.DefaultFieldMatcher{},
	}
	return estore
}

type PostgresEntityStore[T store.EntityType] struct {
	matcher query.FieldMatcher
}

func (p PostgresEntityStore[T]) FindByID(ctx context.Context, id string) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresEntityStore[T]) Exists(ctx context.Context, id string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresEntityStore[T]) Create(ctx context.Context, entity T) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresEntityStore[T]) Update(ctx context.Context, entity T) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresEntityStore[T]) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresEntityStore[T]) GetAll(ctx context.Context) iter.Seq2[T, error] {
	//TODO implement me
	panic("implement me")
}

func (p PostgresEntityStore[T]) GetAllCount(ctx context.Context) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresEntityStore[T]) GetAllPaginated(ctx context.Context, opts store.PaginationOptions) iter.Seq2[T, error] {
	//TODO implement me
	panic("implement me")
}

func (p PostgresEntityStore[T]) FindByPredicate(ctx context.Context, predicate query.Predicate) iter.Seq2[T, error] {
	//TODO implement me
	panic("implement me")
}

func (p PostgresEntityStore[T]) FindByPredicatePaginated(ctx context.Context, predicate query.Predicate, opts store.PaginationOptions) iter.Seq2[T, error] {
	//TODO implement me
	panic("implement me")
}

func (p PostgresEntityStore[T]) FindFirstByPredicate(ctx context.Context, predicate query.Predicate) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresEntityStore[T]) CountByPredicate(ctx context.Context, predicate query.Predicate) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresEntityStore[T]) DeleteByPredicate(ctx context.Context, predicate query.Predicate) error {
	//TODO implement me
	panic("implement me")
}
