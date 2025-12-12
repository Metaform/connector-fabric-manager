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
	"github.com/metaform/connector-fabric-manager/agent/edcv/identityhub"
	"github.com/metaform/connector-fabric-manager/assembly/serviceapi"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/common/token"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

type EDCVActivityProcessor struct {
	VaultClient       serviceapi.VaultClient
	HTTPClient        *http.Client
	Monitor           system.LogMonitor
	TokenProvider     token.TokenProvider
	IdentityAPIClient identityhub.IdentityAPIClient
}

type EDCVData struct {
	ParticipantID       string `json:"cfm.participant.id" validate:"required"`
	VaultAccessClientID string `json:"clientID.vaultAccess" validate:"required"`
	ApiAccessClientID   string `json:"clientID.apiAccess" validate:"required"`
	// PublicURL the public URL which is used for resolving Web DIDs. If not specified, must contain the IdentityHub's public endpoint.
	PublicURL string `json:"publicURL" validate:"required"`
	// CredentialServiceURL the URL of the credential service, i.e., the query and storage endpoints of IdentityHub
	CredentialServiceURL string `json:"cfm.participant.credentialservice" validate:"required"`
	// ProtocolServiceURL the URL of the protocol service, i.e., the DSP protocol endpoint of the control plane
	ProtocolServiceURL string `json:"cfm.participant.protocolservice" validate:"required"`
}

func NewProcessor(config *Config) *EDCVActivityProcessor {
	return &EDCVActivityProcessor{
		VaultClient:       config.VaultClient,
		HTTPClient:        config.Client,
		Monitor:           config.LogMonitor,
		IdentityAPIClient: config.IdentityAPIClient,
	}
}

type Config struct {
	serviceapi.VaultClient
	*http.Client
	system.LogMonitor
	identityhub.IdentityAPIClient
}

func (p EDCVActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	var data EDCVData
	err := ctx.ReadValues(&data)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("error processing EDC-V activity for orchestration %s: %w", ctx.OID(), err)}
	}

	participantContextId := createParticipantContextID()

	// create participant-context in IdentityHub
	did, err := p.extractWebDid(data.PublicURL)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("cannot convert URL to did:web: %w", err)}
	}
	if err := p.createParticipantInIdentityHub(participantContextId, did, data.CredentialServiceURL, data.ProtocolServiceURL); err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("cannot create participant in identity hub: %w", err)}
	}

	_, err = p.VaultClient.ResolveSecret(ctx.Context(), data.VaultAccessClientID)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("error retrieving client secret for orchestration %s: %w", ctx.OID(), err)}
	}
	p.Monitor.Infof("EDCV activity for participant '%s' (client ID = %s) completed successfully", data.ParticipantID, data.VaultAccessClientID)
	return api.ActivityResult{Result: api.ActivityResultComplete}
}

func (p EDCVActivityProcessor) createParticipantInIdentityHub(participantContextID string, did string, credentialServiceURL string, protocolURL string) error {

	manifest := identityhub.NewParticipantManifest(participantContextID, did, credentialServiceURL, protocolURL)
	if err := p.IdentityAPIClient.CreateParticipantContext(manifest); err != nil {
		return fmt.Errorf("error creating participant context in identity hub: %w", err)
	}

	return nil
}

func (p EDCVActivityProcessor) extractWebDid(url string) (string, error) {

	did := strings.Replace(url, "https", "http", -1)
	did = strings.Replace(did, "http://", "", -1)
	did = strings.Replace(did, ":", "%3A", 8)
	did = strings.ReplaceAll(did, "/", ":")
	did = "did:web:" + did

	return did, nil
}

func (p EDCVActivityProcessor) createParticipantInControlPlane() error {
	return nil
}

func createParticipantContextID() string {
	return uuid.New().String()
}
