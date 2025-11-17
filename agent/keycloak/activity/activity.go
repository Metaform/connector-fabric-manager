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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	HTTPClient  *http.Client
}

type KeyCloakActivityProcessor struct {
	keycloakURL string
	token       string
	realm       string
	monitor     system.LogMonitor
	HTTPClient  *http.Client
}

func NewProcessor(config *Config) *KeyCloakActivityProcessor {
	return &KeyCloakActivityProcessor{
		keycloakURL: config.KeycloakURL,
		token:       config.Token,
		realm:       config.Realm,
		monitor:     config.Monitor,
		HTTPClient:  config.HTTPClient}
}

func (p KeyCloakActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	return api.ActivityResult{Result: api.ActivityResultComplete}
}

// CreateClientWithSecret creates a confidential client with the specified secret
func (p KeyCloakActivityProcessor) createClient(name string, clientID string, clientSecret string) error {
	clientURL := fmt.Sprintf(clientUrl, p.keycloakURL, p.realm)

	clientData := map[string]any{
		"name":     name,
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
	resp, err := p.HTTPClient.Do(req)
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
