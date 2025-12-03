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
)

const (
	cfmTenantsTable             = "tenants"
	cfmCellsTable               = "cells"
	cfmParticipantProfilesTable = "participant_profiles"
	cfmDataspaceProfilesTable   = "dataspace_profiles"
)

func createTenantsTable(db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id TEXT PRIMARY KEY,
			version INT DEFAULT 1,
			properties JSONB
		)
	`, cfmTenantsTable))
	return err
}

func createParticipantProfilesTable(db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id VARCHAR(255) PRIMARY KEY,
				version BIGINT NOT NULL,
				identifier VARCHAR(255),
				tenantId VARCHAR(255),
				dataspaceProfileIds JSONB,
				vpas JSONB,
				error BOOLEAN,
				errorDetail TEXT,
				properties JSONB
			);
			CREATE INDEX IF NOT EXISTS idx_participant_tenant ON %s(tenantid)
	`, cfmParticipantProfilesTable, cfmParticipantProfilesTable))
	return err
}

func createDataspaceProfilesTable(db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id TEXT PRIMARY KEY,
			version INT DEFAULT 1,
			artifacts JSONB,
			deployments JSONB,
			properties JSONB
		)
	`, cfmDataspaceProfilesTable))
	return err
}

func createCellsTable(db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id TEXT PRIMARY KEY,
			version INT DEFAULT 1,
			STATE TEXT NOT NULL,
			state_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			properties JSONB
		)
	`, cfmCellsTable))
	return err
}
