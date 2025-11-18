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
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/assembly/serviceapi"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

const (
	jsonContentType   = "application/json"
	contentTypeHeader = "Content-Type"
	authHeader        = "Authorization"
	clientUrl         = "%s/admin/realms/%s/clients"
)

type Config struct {
	KeycloakURL string
	Token       string
	Realm       string
	Monitor     system.LogMonitor
	VaultClient serviceapi.VaultClient
	HTTPClient  *http.Client
}

// KeyCloakActivityProcessor creates a confidential client in Keycloak and stores the client secret in Vault for use by
// other processors. The client ID is returned as a value in the context.
type KeyCloakActivityProcessor struct {
	keycloakURL string
	token       string
	realm       string
	monitor     system.LogMonitor
	httpClient  *http.Client
	vaultClient serviceapi.VaultClient
}

// NewProcessor creates a new KeyCloakActivityProcessor instance
func NewProcessor(config *Config) *KeyCloakActivityProcessor {
	return &KeyCloakActivityProcessor{
		keycloakURL: config.KeycloakURL,
		token:       config.Token,
		realm:       config.Realm,
		monitor:     config.Monitor,
		httpClient:  config.HTTPClient,
		vaultClient: config.VaultClient,
	}
}

func (p KeyCloakActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	if ctx.Discriminator() == api.DisposeDiscriminator {
		panic("Not yet implemented")
	} else {
		return p.provisionConfidentialClient(ctx)
	}
}

// provisionConfidentialClient creates a confidential client in Keycloak and stores the client secret in Vault for use by
// other processors. The client ID is returned as a value in the context.
// TODO support idempotent provisioning
func (p KeyCloakActivityProcessor) provisionConfidentialClient(ctx api.ActivityContext) api.ActivityResult {
	clientID := generateClientID()
	clientSecret, err := generateClientSecret()
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: err}
	}
	err = p.createClient(clientID, clientSecret)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: err}
	}
	err = p.vaultClient.StoreSecret(ctx.Context(), clientID, clientSecret)
	if err != nil {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: err}
	}
	ctx.SetValue("clientID", clientID)
	return api.ActivityResult{Result: api.ActivityResultComplete}
}

// CreateClientWithSecret creates a confidential client with the specified secret
func (p KeyCloakActivityProcessor) createClient(clientID string, clientSecret string) error {
	clientURL := fmt.Sprintf(clientUrl, p.keycloakURL, p.realm)

	clientData := map[string]any{
		"clientId": clientID,
		"secret":   clientSecret,
	}

	jsonData, err := json.Marshal(clientData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, clientURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating client request: %w", err)
	}

	req.Header.Set(contentTypeHeader, jsonContentType)
	req.Header.Set(authHeader, fmt.Sprintf("Bearer %s", p.token))
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("create client request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create client operation failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// generateClientSecret generates a random secret using encoding.
func generateClientSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// generateClientID generates a unique client ID that complies with Keycloak and typical Vault requirements
func generateClientID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
