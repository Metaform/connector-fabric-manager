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
	"database/sql"
	"errors"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDBTable(t *testing.T) {
	_, err := testDB.Exec(`
		DROP TABLE IF EXISTS test_table CASCADE;
		CREATE TABLE test_table (
			id SERIAL PRIMARY KEY,
			value TEXT NOT NULL
		)
	`)
	require.NoError(t, err)
}

func TestDBTransactionContext(t *testing.T) {
	setupTestDBTable(t)
	defer cleanupTestData(t, testDB)

	ctx := context.Background()
	trxContext := NewDBTransactionContext(testDB)

	t.Run("Successful transaction", func(t *testing.T) {
		err := trxContext.Execute(ctx, func(ctx context.Context) error {
			trx := ctx.Value(SQLTransactionKey).(*sql.Tx)
			_, err := trx.Exec("INSERT INTO test_table (value) VALUES ($1)", "test1")
			return err
		})

		assert.NoError(t, err)

		// Verify the data was inserted
		var count int
		err = testDB.QueryRow("SELECT COUNT(*) FROM test_table WHERE value = $1", "test1").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("Failed transaction should rollback", func(t *testing.T) {
		initialCount := 0
		err := testDB.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&initialCount)
		assert.NoError(t, err)

		err = trxContext.Execute(ctx, func(ctx context.Context) error {
			trx := ctx.Value(SQLTransactionKey).(*sql.Tx)

			_, err := trx.Exec("INSERT INTO test_table (value) VALUES ($1)", "test2")
			if err != nil {
				return err
			}
			return errors.New("forced error")
		})

		assert.Error(t, err)

		// Verify the data was rolled back
		var count int
		err = testDB.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, initialCount, count)
	})

	t.Run("Panic should rollback", func(t *testing.T) {
		initialCount := 0
		err := testDB.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&initialCount)
		assert.NoError(t, err)

		assert.Panics(t, func() {
			_ = trxContext.Execute(ctx, func(ctx context.Context) error {
				trx := ctx.Value(SQLTransactionKey).(*sql.Tx)
				_, err := trx.Exec("INSERT INTO test_table (value) VALUES ($1)", "test3")
				if err != nil {
					return err
				}
				panic("forced panic")
			})
		})

		// Verify the data was rolled back
		var count int
		err = testDB.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, initialCount, count)
	})
}
