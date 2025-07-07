package sqlstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestContainer(t *testing.T) (testcontainers.Container, string) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())
	return container, dsn
}

func setupTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS test_table (
			id SERIAL PRIMARY KEY,
			VALUE TEXT NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}
}

func TestDBTransactionContext(t *testing.T) {
	// Setup
	ctx := context.Background()
	container, dsn := setupTestContainer(t)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate container: %v", err)
		}
	}()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	setupTestDB(t, db)
	trxContext := NewDBTransactionContext(db)

	t.Run("Successful transaction", func(t *testing.T) {
		err := trxContext.Execute(ctx, func(ctx context.Context) error {
			trx := ctx.Value(SQLTransactionKey).(*sql.Tx)
			_, err := trx.Exec("INSERT INTO test_table (value) VALUES ($1)", "test1")
			return err
		})

		assert.NoError(t, err)

		// Verify the data was inserted
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM test_table WHERE value = $1", "test1").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("Failed transaction should rollback", func(t *testing.T) {
		initialCount := 0
		err := db.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&initialCount)
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
		err = db.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, initialCount, count)
	})

	t.Run("Panic should rollback", func(t *testing.T) {
		initialCount := 0
		err := db.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&initialCount)
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
		err = db.QueryRow("SELECT COUNT(*) FROM test_table").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, initialCount, count)
	})
}
