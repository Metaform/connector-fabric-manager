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
	"time"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/agent/common/identityhub"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

type OnboardingActivityProcessor struct {
	Monitor           system.LogMonitor
	IdentityApiClient identityhub.IdentityAPIClient
}

type credentialRequestData struct {
	//CredentialRequest    identityhub.CredentialRequest `json:"credentialRequest"`
	ParticipantContextID string `json:"clientID.apiAccess" validate:"required"`
	HolderPID            string `json:"holderPid"`
	CredentialRequestURL string `json:"credentialRequest"`
}

func (p OnboardingActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {

	var credentialRequest credentialRequestData
	if err := ctx.ReadValues(&credentialRequest); err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("error processing Onboarding activity for orchestration %s: %w", ctx.OID(), err)}
	}

	// no credential request was made yet -> make one
	if "" == credentialRequest.HolderPID {
		return p.processNewRequest(ctx, credentialRequest)
	} else { // holderPID exists, check the status of the issuance
		return p.processExistingRequest(ctx, credentialRequest)
	}

}

func (p OnboardingActivityProcessor) processExistingRequest(ctx api.ActivityContext, credentialRequest credentialRequestData) api.ActivityResult {
	state, err := p.IdentityApiClient.GetCredentialRequestState(credentialRequest.ParticipantContextID, credentialRequest.HolderPID)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("error getting credential request state: %w", err)}
	}
	p.Monitor.Infof("Credential request for participant '%s' is in state '%d'", credentialRequest.ParticipantContextID, state)

	switch state {

	case identityhub.CredentialRequestStateCreated:
		return api.ActivityResult{Result: api.ActivityResultSchedule, WaitOnReschedule: time.Duration(5) * time.Second}
	case identityhub.CredentialRequestStateIssued:
		ctx.SetOutputValue("holderPid", credentialRequest.HolderPID)
		ctx.SetOutputValue("credentialRequest", credentialRequest.CredentialRequestURL)
		ctx.SetOutputValue("participantContextId", credentialRequest.ParticipantContextID)
		return api.ActivityResult{Result: api.ActivityResultComplete}
	case identityhub.CredentialRequestStateRejected:
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("credential request for participant '%s' was rejected", credentialRequest.ParticipantContextID)}

	default:
		return api.ActivityResult{Result: api.ActivityResultRetryError, Error: fmt.Errorf("unexpected credential request state '%d'", state)}
	}
}

func (p OnboardingActivityProcessor) processNewRequest(ctx api.ActivityContext, credentialRequest credentialRequestData) api.ActivityResult {
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
	// make credential request
	location, err := p.IdentityApiClient.RequestCredentials(credentialRequest.ParticipantContextID, cr)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("error requesting credentials: %w", err)}
	}
	p.Monitor.Infof("Credentials request for participant '%s' submitted successfully, credential is at %s", credentialRequest.ParticipantContextID, location)
	ctx.SetValue("participantContextId", credentialRequest.ParticipantContextID)
	ctx.SetValue("holderPid", holderPid)
	ctx.SetValue("credentialRequest", location)

	return api.ActivityResult{
		Result:           api.ActivityResultSchedule,
		WaitOnReschedule: time.Duration(5) * time.Second,
		Error:            nil,
	}
}
