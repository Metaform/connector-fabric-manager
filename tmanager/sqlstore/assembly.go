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
	"database/sql"
	"fmt"

	"github.com/metaform/connector-fabric-manager/common/sqlstore"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/tmanager/api"

	_ "github.com/lib/pq" // Register PostgreSQL driver
)

const (
	driverName = "postgres"
	dsnKey     = "dsn"
)

type PostgresServiceAssembly struct {
	system.DefaultServiceAssembly
	db *sql.DB
}

func (a *PostgresServiceAssembly) Name() string {
	return "Tenant Manager Postgres"
}

func (a *PostgresServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{api.CellStoreKey, api.DataspaceProfileStoreKey, api.ParticipantProfileStoreKey, store.TransactionContextKey}
}

func (a *PostgresServiceAssembly) Init(ictx *system.InitContext) error {
	if !ictx.Config.IsSet(dsnKey) {
		return fmt.Errorf("missing Postgres DSN configuration: %s", dsnKey)
	}
	dsn := ictx.Config.GetString(dsnKey)

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return fmt.Errorf("error connecting to DB at %s: %w", dsn, err)
	}

	a.db = db

	createTables(db)

	cellStore := newCellStore()
	dataspaceStore := newDataspaceProfileStore()
	participantStore := newParticipantProfileStore()
	tenantStore := newTenantStore()

	ictx.Registry.Register(api.TenantStoreKey, tenantStore)
	ictx.Registry.Register(api.ParticipantProfileStoreKey, participantStore)
	ictx.Registry.Register(api.DataspaceProfileStoreKey, dataspaceStore)
	ictx.Registry.Register(api.CellStoreKey, cellStore)

	txContext := sqlstore.NewDBTransactionContext(db)
	ictx.Registry.Register(store.TransactionContextKey, txContext)

	return nil
}

func (a *PostgresServiceAssembly) Finalize() error {
	if a.db != nil {
		a.db.Close()
	}
	return nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS cells (
			id TEXT PRIMARY KEY,
			version INT DEFAULT 1,
			STATE TEXT NOT NULL,
			state_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			properties JSONB
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS dataspace_profiles (
			id TEXT PRIMARY KEY,
			version INT DEFAULT 1,
			artifacts JSONB,
			deployments JSONB,
			properties JSONB
		)
	`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS participant_profiles (
				id VARCHAR(255) PRIMARY KEY,
				version BIGINT NOT NULL,
				identifier VARCHAR(255),
				tenantId VARCHAR(255),
				dataspaceProfileIds JSONB,
				vpas JSONB,
				error BOOLEAN,
				errorDetail TEXT,
				properties JSONB
			)
	`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tenants (
			id TEXT PRIMARY KEY,
			version INT DEFAULT 1,
			properties JSONB
		)
	`)

	if err != nil {
		return err
	}

	return nil
}
