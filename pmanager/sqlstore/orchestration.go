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
	"time"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/sqlstore"
	"github.com/metaform/connector-fabric-manager/common/store"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

func newOrchestrationEntryStore() store.EntityStore[*api.OrchestrationEntry] {
	columnNames := []string{"id", "version", "correlationId", "state", "stateTimestamp", "createdTimestamp", "orchestrationType"}
	builder := sqlstore.NewPostgresJSONBBuilder().WithJSONBFieldTypes(map[string]sqlstore.JSONBFieldType{})

	estore := sqlstore.NewPostgresEntityStore[*api.OrchestrationEntry](
		"orchestration_entries",
		columnNames,
		recordToOrchestrationEntry,
		orchestrationEntryToRecord,
		builder,
	)

	return estore
}

func recordToOrchestrationEntry(tx *sql.Tx, record *sqlstore.DatabaseRecord) (*api.OrchestrationEntry, error) {
	profile := &api.OrchestrationEntry{}
	if id, ok := record.Values["id"].(string); ok {
		profile.ID = id
	} else {
		return nil, fmt.Errorf("invalid orchestration entry id reading record")
	}

	if version, ok := record.Values["version"].(int64); ok {
		profile.Version = version
	} else {
		return nil, fmt.Errorf("invalid orchestration entry version reading record")
	}

	if version, ok := record.Values["correlationId"].(string); ok {
		profile.CorrelationID = version
	} else {
		return nil, fmt.Errorf("invalid orchestration entry correlationId reading record")
	}

	if state, ok := record.Values["state"].(int64); ok {
		profile.State = api.OrchestrationState(state)
	} else {
		return nil, fmt.Errorf("invalid orchestration entry state reading record")
	}

	if timestamp, ok := record.Values["stateTimestamp"].(time.Time); ok {
		profile.StateTimestamp = timestamp
	} else {
		return nil, fmt.Errorf("invalid orchestration entry stateTimestamp reading record")
	}

	if timestamp, ok := record.Values["createdTimestamp"].(time.Time); ok {
		profile.CreatedTimestamp = timestamp
	} else {
		return nil, fmt.Errorf("invalid orchestration entry createdTimestamp reading record")
	}

	if otype, ok := record.Values["orchestrationType"].(string); ok {
		profile.OrchestrationType = model.OrchestrationType(otype)
	} else {
		return nil, fmt.Errorf("invalid orchestration entry type reading record")
	}

	return profile, nil

}

func orchestrationEntryToRecord(profile *api.OrchestrationEntry) (*sqlstore.DatabaseRecord, error) {
	record := &sqlstore.DatabaseRecord{
		Values: make(map[string]any),
	}

	record.Values["id"] = profile.ID
	record.Values["version"] = profile.Version
	record.Values["correlationId"] = profile.CorrelationID
	record.Values["state"] = profile.State
	record.Values["stateTimestamp"] = profile.StateTimestamp
	record.Values["createdTimestamp"] = profile.CreatedTimestamp
	record.Values["orchestrationType"] = profile.OrchestrationType

	return record, nil
}
