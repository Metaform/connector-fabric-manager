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

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/agent/common/identityhub"
	"github.com/metaform/connector-fabric-manager/assembly/serviceapi"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

type OnboardingActivityProcessor struct {
	Monitor           system.LogMonitor
	Vault             serviceapi.VaultClient
	IdentityApiClient identityhub.IdentityAPIClient
}

type credentialRequestData struct {
	//CredentialRequest    identityhub.CredentialRequest `json:"credentialRequest"`
	ParticipantContextID string `json:"clientID.apiAccess"`
}

func (p OnboardingActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {

	var credentialRequest credentialRequestData
	if err := ctx.ReadValues(&credentialRequest); err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("error processing Onboarding activity for orchestration %s: %w", ctx.OID(), err)}
	}

	// todo: have this come in through CFM REST API -> dataspace profile
	holderPid := uuid.New().String()
	cr := identityhub.CredentialRequest{
		IssuerDID: "did:web:issuerservice.edc-v.svc.cluster.local%3A10016:issuer",
		HolderPID: holderPid,
		Credentials: []identityhub.CredentialType{
			{
				Format: "VC1_0_JWT",
				Type:   "MembershipCredential",
				ID:     "membership-credential-def",
			},
		},
	}
	location, err := p.IdentityApiClient.RequestCredentials(credentialRequest.ParticipantContextID, cr)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("error requesting credentials: %w", err)}
	}
	p.Monitor.Infof("Credentials request for participant '%s' submitted successfully, credential is at %s", credentialRequest.ParticipantContextID, location)

	ctx.SetOutputValue("holderPid", holderPid)
	ctx.SetOutputValue("credentialRequest", location)

	return api.ActivityResult{Result: api.ActivityResultComplete}
}
