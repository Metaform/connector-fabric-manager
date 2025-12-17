/*
 *  Copyright (c) 2025 Metaform Systems, Inc.
 *
 *  This program and the accompanying materials are made available under the
 *  terms of the Apache License, Version 2.0 which is available at
 *  https://www.apache.org/licenses/LICENSE-2.0
 *
 *  SPDX-License-Identifier: Apache-2.0
 *
 *  Contributors:
 *       Metaform Systems, Inc. - initial API and implementation
 *
 */

package identityhub

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/metaform/connector-fabric-manager/common/token"
)

const (
	CreateParticipantURL = "/v1alpha/participants"
)

type IdentityAPIClient interface {
	CreateParticipantContext(manifest ParticipantManifest) (*CreateParticipantContextResponse, error)
	RequestCredentials(participantContextID string, credentialRequest CredentialRequest) (string, error)
}

type HttpIdentityAPIClient struct {
	BaseURL       string
	TokenProvider token.TokenProvider
	HttpClient    *http.Client
}

func (a HttpIdentityAPIClient) RequestCredentials(participantContextID string, credentialRequest CredentialRequest) (string, error) {
	accessToken, err := a.TokenProvider.GetToken() // this should be the participant context's access token!
	if err != nil {
		return "", fmt.Errorf("failed to get API access token: %w", err)
	}

	payload, err := json.Marshal(credentialRequest)
	if err != nil {
		return "", err
	}

	b64 := base64.RawURLEncoding.EncodeToString([]byte(participantContextID))
	url := fmt.Sprintf("%s/v1alpha/participants/%s/credentials/request", a.BaseURL, b64)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := a.HttpClient.Do(req)
	defer a.closeResponse(resp)

	if err != nil {
		return "", fmt.Errorf("failed to create credentials request for %s: %w", participantContextID, err)
	}

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to request credentials for participant context on IdentityHub: received status code %d, body: %s", resp.StatusCode, string(body))
	}

	location := resp.Header.Get("Location")

	return location, nil
}

func (a HttpIdentityAPIClient) CreateParticipantContext(manifest ParticipantManifest) (*CreateParticipantContextResponse, error) {
	accessToken, err := a.TokenProvider.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get API access token: %w", err)
	}

	data := map[string]any{
		"roles": []string{"participant"},
		"serviceEndpoints": []map[string]any{
			{
				"type":            "CredentialService",
				"id":              manifest.CredentialServiceID,
				"serviceEndpoint": manifest.CredentialServiceURL,
			},
			{
				"type":            "ProtocolEndpoint",
				"id":              manifest.ProtocolServiceID,
				"serviceEndpoint": manifest.ProtocolServiceURL,
			},
		},
		"active":               manifest.IsActive,
		"participantContextId": manifest.ParticipantContextID,
		"did":                  manifest.DID,
		"key": map[string]any{
			"keyId":           manifest.KeyGeneratorParameters.KeyID,
			"privateKeyAlias": manifest.KeyGeneratorParameters.PrivateKeyAlias,
			"keyGeneratorParams": map[string]string{
				"algorithm": manifest.KeyGeneratorParameters.KeyAlgorithm,
				"curve":     manifest.KeyGeneratorParameters.Curve,
			},
		},
		"additionalProperties": map[string]any{
			"edc.vault.hashicorp.config": map[string]any{
				"credentials": map[string]string{
					"clientId":     manifest.VaultCredentials.ClientID,
					"clientSecret": manifest.VaultCredentials.ClientSecret,
					"tokenUrl":     manifest.VaultCredentials.TokenURL,
				},
				"config": map[string]string{
					"secretPath": manifest.VaultConfig.SecretPath,
					"folderPath": manifest.VaultConfig.FolderPath,
					"vaultUrl":   manifest.VaultConfig.VaultURL,
				},
			},
		},
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	url := a.BaseURL + CreateParticipantURL
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := a.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create participant context on IdentityHub: %w", err)
	}
	defer a.closeResponse(resp)

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to create participant context on IdentityHub: received status code %d, body: %s", resp.StatusCode, string(body))
	}

	createResponse := &CreateParticipantContextResponse{}

	if err := json.Unmarshal(body, createResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal participant context creation response: %w", err)
	}

	return createResponse, nil
}

func (a HttpIdentityAPIClient) closeResponse(resp *http.Response) {
	func() {
		// drain and close response body to avoid connection/resource leak
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
}
