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

package fixtures

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	formContentType = "application/x-www-form-urlencoded"
	jsonContentType = "application/json"
)

func SetupKeyCloakContainer(ctx context.Context) (string, string, error) {

	req := testcontainers.ContainerRequest{
		Image:        "quay.io/keycloak/keycloak:latest",
		ExposedPorts: []string{"8080/tcp"},
		Env: map[string]string{
			"KEYCLOAK_ADMIN":          "admin",
			"KEYCLOAK_ADMIN_PASSWORD": "admin",
			"KC_HTTP_ENABLED":         "true",
			"KC_HEALTH_ENABLED":       "true",
		},
		Cmd: []string{
			"start-dev",
		},
		WaitingFor: wait.ForLog("Profile dev activated").
			WithStartupTimeout(10 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to start container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "8080/tcp")
	if err != nil {
		return "", "", fmt.Errorf("failed to get container port: %w", err)
	}

	keycloakURL := fmt.Sprintf("http://%s:%s", host, port.Port())

	token, err := getAdminToken(keycloakURL, "admin", "admin")
	if err != nil {
		return "", "", fmt.Errorf("failed to get admin token: %w", err)
	}

	realmName := "edcv"
	err = createRealm(keycloakURL, token, realmName)
	if err != nil {
		return "", "", fmt.Errorf("failed to create EDC-V realm: %w", err)

	}

	// Create user
	err = createUser(keycloakURL, token, realmName, "testuser", "testpassword", "test@example.com")
	if err != nil {
		return "", "", fmt.Errorf("failed to create user: %w", err)
	}

	// Create client scope
	err = createClientScope(keycloakURL, token, "issuer-admin-api", "full access to the Issuer Admin API")
	if err != nil {
		return "", "", fmt.Errorf("failed to create client scope: %w", err)
	}

	// Create provisioner client
	err = createProvisionerClient(keycloakURL, token, realmName)
	if err != nil {
		return "", "", fmt.Errorf("failed to create provisioner client: %w", err)
	}
	return keycloakURL, token, nil
}

// getAdminToken obtains an access token from the master realm
func getAdminToken(keycloakURL, username, password string) (string, error) {
	tokenURL := fmt.Sprintf("%s/realms/master/protocol/openid-connect/token", keycloakURL)

	data := map[string]string{
		"grant_type": "password",
		"client_id":  "admin-cli",
		"username":   username,
		"password":   password,
	}

	formData := bytes.NewBufferString("")
	for k, v := range data {
		formData.WriteString(k + "=" + v + "&")
	}

	resp, err := http.Post(tokenURL, formContentType, formData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get admin token: status %d, body: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access_token not found in response")
	}

	return token, nil
}

// createRealm creates a new realm with DCR policies configured
// createRealm creates a new realm
func createRealm(keycloakURL, token, realmName string) error {
	realmURL := fmt.Sprintf("%s/admin/realms", keycloakURL)

	realmData := map[string]interface{}{
		"realm":   realmName,
		"enabled": true,
	}

	return AdminRequest(realmURL, token, "POST", realmData)
}

// createUser creates a new user in a realm
func createUser(keycloakURL, token, realmName, username, password, email string) error {
	userURL := fmt.Sprintf("%s/admin/realms/%s/users", keycloakURL, realmName)

	userData := map[string]interface{}{
		"username": username,
		"enabled":  true,
		"email":    email,
		"credentials": []map[string]interface{}{
			{
				"type":      "password",
				"value":     password,
				"temporary": false,
			},
		},
	}

	return AdminRequest(userURL, token, "POST", userData)
}

// AdminRequest makes an authenticated request to the Keycloak admin API
func AdminRequest(url, token, method string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", jsonContentType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// createClientScope creates a new client scope in a realm
func createClientScope(keycloakURL, token, scopeName, description string) error {
	scopeURL := fmt.Sprintf("%s/admin/realms/master/client-scopes", keycloakURL)

	scopeData := map[string]interface{}{
		"name":        scopeName,
		"description": description,
		"protocol":    "openid-connect",
		"attributes": map[string]interface{}{
			"include.in.token.scope":    "true",
			"display.on.consent.screen": "true",
		},
	}

	return AdminRequest(scopeURL, token, "POST", scopeData)
}

// createProvisionerClient creates the provisioner service client with role mapper
func createProvisionerClient(keycloakURL, token, realmName string) error {
	clientURL := fmt.Sprintf("%s/admin/realms/%s/clients", keycloakURL, realmName)

	protocolMapper := map[string]interface{}{
		"name":            "role",
		"protocol":        "openid-connect",
		"protocolMapper":  "oidc-hardcoded-claim-mapper",
		"consentRequired": false,
		"config": map[string]interface{}{
			"claim.name":           "role",
			"claim.value":          "provisioner",
			"jsonType.label":       "String",
			"access.token.claim":   "true",
			"id.token.claim":       "true",
			"userinfo.token.claim": "true",
		},
	}

	clientData := map[string]interface{}{
		"clientId":                  "edcv-provisioner",
		"name":                      "Provisioner User",
		"description":               "Can create and delete tenants",
		"enabled":                   true,
		"protocol":                  "openid-connect",
		"publicClient":              false,
		"serviceAccountsEnabled":    true,
		"secret":                    "provisioner-secret",
		"standardFlowEnabled":       false,
		"directAccessGrantsEnabled": false,
		"fullScopeAllowed":          true,
		"protocolMappers":           []map[string]interface{}{protocolMapper},
		"defaultClientScopes":       []string{"issuer-admin-api"},
	}

	return AdminRequest(clientURL, token, "POST", clientData)
}