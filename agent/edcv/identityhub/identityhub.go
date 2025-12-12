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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/metaform/connector-fabric-manager/common/token"
)

const (
	CreateParticipantURL = "/v1alpha/participants"
)

type ApiClient interface {
	CreateParticipantContext(participantContextID string, did string) (string, error)
}

type IdentityAPIClient struct {
	baseURL       string
	tokenProvider token.TokenProvider
	httpClient    *http.Client
}

func (a IdentityAPIClient) CreateParticipantContext(manifest ParticipantManifest) error {
	accessToken, err := a.tokenProvider.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get API access token: %w", err)
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
		return err
	}

	url := a.baseURL + CreateParticipantURL
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create participant context on IdentityHub: %w", err)
	}
	defer func() {
		// drain and close response body to avoid connection/resource leak
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create participant context on IdentityHub: received status code %d, body: %s", resp.StatusCode, string(body))
	}

	return nil

}
