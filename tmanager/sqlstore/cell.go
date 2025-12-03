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
	"encoding/json"
	"fmt"
	"time"

	"github.com/metaform/connector-fabric-manager/common/sqlstore"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
)

func newCellStore() store.EntityStore[*api.Cell] {
	columnNames := []string{"id", "version", "state", "stateTimestamp", "properties"}
	builder := sqlstore.NewPostgresJSONBBuilder().WithJSONBFieldTypes(map[string]sqlstore.JSONBFieldType{
		"properties": sqlstore.JSONBFieldTypeScalar,
	})

	estore := sqlstore.NewPostgresEntityStore[*api.Cell](
		"cells",
		columnNames,
		recordToCellEntity,
		cellEntityToRecord,
		builder,
	)

	return estore
}

func recordToCellEntity(_ *sql.Tx, record *sqlstore.DatabaseRecord) (*api.Cell, error) {
	cell := &api.Cell{}
	if id, ok := record.Values["id"].(string); ok {
		cell.ID = id
	} else {
		return nil, fmt.Errorf("invalid cell id reading record")
	}

	if version, ok := record.Values["version"].(int64); ok {
		cell.Version = version
	} else {
		return nil, fmt.Errorf("invalid cell version reading record")
	}

	if state, ok := record.Values["state"].(string); ok {
		cell.State = api.DeploymentState(state)
	} else {
		return nil, fmt.Errorf("invalid cell state reading record")
	}

	if timestamp, ok := record.Values["stateTimestamp"].(time.Time); ok {
		cell.StateTimestamp = timestamp
	} else {
		return nil, fmt.Errorf("invalid cell state timestamp reading record")
	}

	if bytes, ok := record.Values["properties"].([]byte); ok && bytes != nil {
		if err := json.Unmarshal(bytes, &cell.Properties); err != nil {
			return nil, err
		}
	}
	return cell, nil
}

func cellEntityToRecord(cell *api.Cell) (*sqlstore.DatabaseRecord, error) {
	record := &sqlstore.DatabaseRecord{
		Values: make(map[string]any),
	}

	record.Values["id"] = cell.ID
	record.Values["version"] = cell.Version
	record.Values["state"] = cell.State
	record.Values["stateTimestamp"] = cell.StateTimestamp

	if cell.Properties != nil {
		metadataBytes, err := json.Marshal(cell.Properties)
		if err != nil {
			return record, err
		}
		record.Values["properties"] = metadataBytes
	}

	return record, nil
}
