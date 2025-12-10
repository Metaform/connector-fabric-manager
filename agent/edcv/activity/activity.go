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

package activity

import (
	"fmt"
	"net/http"

	"github.com/metaform/connector-fabric-manager/assembly/serviceapi"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

const (
	clientIDKey = "clientID"
)

type EDCVActivityProcessor struct {
	VaultClient serviceapi.VaultClient
	HTTPClient  *http.Client
	Monitor     system.LogMonitor
}

func (p EDCVActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	var data EDCVData
	err := ctx.ReadValues(&data)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("error processing EDC-V activity for orchestration %s: %w", ctx.OID(), err)}
	}
	_, err = p.VaultClient.ResolveSecret(ctx.Context(), data.ClientID)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("error retrieving client secret for orchestration %s: %w", ctx.OID(), err)}
	}
	p.Monitor.Infof("EDCV activity for participant '%s' (client ID = %s) completed successfully", data.ParticipantID, data.ClientID)
	return api.ActivityResult{Result: api.ActivityResultComplete}
}

type EDCVData struct {
	ParticipantID string `json:"cfm.participant.id" validate:"required"`
	ClientID      string `json:"clientID" validate:"required"`
}
