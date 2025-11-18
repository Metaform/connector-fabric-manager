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
	"github.com/metaform/connector-fabric-manager/common/model"
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
	_, found := ctx.Value(model.ParticipantIdentifier)
	if !found {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("missing participant identifier")}
	}
	clientID, found := ctx.Value(clientIDKey)
	if !found {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("missing clientID")}
	}
	_, err := p.VaultClient.ResolveSecret(ctx.Context(), clientID.(string))
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("failed to resolve secret for client ID %s: %w", clientID, err)}
	}
	return api.ActivityResult{Result: api.ActivityResultComplete}
}
