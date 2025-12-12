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
	"strings"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/assembly/serviceapi"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/common/token"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

const (
	jsonContentType   = "application/json"
	contentTypeHeader = "Content-Type"
)

type EDCVActivityProcessor struct {
	VaultClient     serviceapi.VaultClient
	HTTPClient      *http.Client
	Monitor         system.LogMonitor
	identityHubURL  string
	controlPlaneURL string
	TokenProvider   token.TokenProvider
}

func NewProcessor(config *Config) *EDCVActivityProcessor {
	return &EDCVActivityProcessor{
		VaultClient:     config.VaultClient,
		HTTPClient:      config.HTTPClient,
		Monitor:         config.Monitor,
		TokenProvider:   config.TokenProvider,
		identityHubURL:  config.IdentityHubBaseURL,
		controlPlaneURL: config.ControlPlaneBaseURL,
	}
}

type Config struct {
	VaultClient         serviceapi.VaultClient
	HTTPClient          *http.Client
	Monitor             system.LogMonitor
	TokenProvider       token.TokenProvider
	IdentityHubBaseURL  string
	ControlPlaneBaseURL string
}

func (p EDCVActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	var data EDCVData
	err := ctx.ReadValues(&data)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("error processing EDC-V activity for orchestration %s: %w", ctx.OID(), err)}
	}

	participantContextId := createParticipantContextID()
	// create participant-context in IdentityHub

	did, err := p.extractWebDid(data.PublicURL, participantContextId)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("cannot convert URL to did:web: %w", err)}
	}
	if err := p.createParticipantInIdentityHub(participantContextId, did); err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("cannot create participant in identity hub: %w", err)}
	}

	_, err = p.VaultClient.ResolveSecret(ctx.Context(), data.VaultAccessClientID)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("error retrieving client secret for orchestration %s: %w", ctx.OID(), err)}
	}
	p.Monitor.Infof("EDCV activity for participant '%s' (client ID = %s) completed successfully", data.ParticipantID, data.VaultAccessClientID)
	return api.ActivityResult{Result: api.ActivityResultComplete}
}

func (p EDCVActivityProcessor) createParticipantInIdentityHub(participantContextId string, participantContextDID string) error {
	jwt, err := p.TokenProvider.GetToken()
	if err != nil {
		return fmt.Errorf("error getting provisioner token for participant context '%s': %w", participantContextId, err)
	}
	p.Monitor.Infof(jwt)
	return nil
}

func (p EDCVActivityProcessor) extractWebDid(url string, participantContextId string) (string, error) {
	if p.identityHubURL == "" {
		return "", fmt.Errorf("IdentityHub base URL must not be empty")
	}
	if url == "" {
		url = p.identityHubURL + "/" + participantContextId
	}

	did := strings.Replace(url, "https", "http", -1)
	did = strings.Replace(did, "http://", "", -1)
	did = strings.Replace(did, ":", "%3A", 8)
	did = strings.ReplaceAll(did, "/", ":")
	did = "did:web:" + did

	return did, nil
}

type EDCVData struct {
	ParticipantID       string `json:"cfm.participant.id" validate:"required"`
	VaultAccessClientID string `json:"clientID.vaultAccess" validate:"required"`
	ApiAccessClientID   string `json:"clientID.apiAccess" validate:"required"`
	// PublicURL the public URL which is used for resolving Web DIDs. Optional.
	PublicURL string `json:"publicURL"`
}

func createParticipantInControlPlane() error {
	return nil
}

func createParticipantContextID() string {
	return uuid.New().String()
}
