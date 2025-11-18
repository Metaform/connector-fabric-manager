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

package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/metaform/connector-fabric-manager/assembly/serviceapi"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	vaultImage           = "hashicorp/vault:1.20"
	vaultPort            = "8200/tcp"
	vaultRootToken       = "myroot"
	vaultRequestTimeout  = 30 * time.Second
	containerStartupTime = 15 * time.Second
)

func NewVaultClient(vaultURL, roleID, secretID string) (serviceapi.VaultClient, error) {
	return createClient(vaultURL, roleID, secretID)
}

type ContainerResult struct {
	URL     string
	Token   string
	Cleanup func()
}

// StartVaultContainer starts a Vault container and returns a ContainerResult containing the URL, root token, and cleanup function
func StartVaultContainer(ctx context.Context) (*ContainerResult, error) {
	req := testcontainers.ContainerRequest{
		Image:        vaultImage,
		ExposedPorts: []string{vaultPort},
		Env: map[string]string{
			"VAULT_DEV_ROOT_TOKEN_ID": vaultRootToken,
		},
		Cmd: []string{"server", "-dev"},
		WaitingFor: wait.ForLog("WARNING! dev mode is enabled!").
			WithStartupTimeout(containerStartupTime),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, vaultPort)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	vaultURL := fmt.Sprintf("http://%s:%s", host, port.Port())

	cleanup := func() {
		_ = container.Terminate(context.Background())
	}

	return &ContainerResult{
		URL:     vaultURL,
		Token:   vaultRootToken,
		Cleanup: cleanup,
	}, nil
}

type VaultSetupResult struct {
	RoleID   string
	SecretID string
}

func SetupVault(vaultURL, rootToken string) (*VaultSetupResult, error) {
	client := &http.Client{Timeout: vaultRequestTimeout}

	// Enable KV v2 secrets engine
	if err := enableSecretsEngine(client, vaultURL, rootToken, "kv-v2", "kv", ""); err != nil {
		return nil, fmt.Errorf("failed to enable KV v2 engine: %w", err)
	}

	// Enable AppRole auth method
	if err := enableAuthMethod(client, vaultURL, rootToken, "approle", "approle", ""); err != nil {
		return nil, fmt.Errorf("failed to enable AppRole auth: %w", err)
	}

	roleResult, err := createAppRole(client, vaultURL, rootToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create AppRole: %w", err)
	}

	return roleResult, nil
}

func setupTestFixtures(ctx context.Context, t *testing.T) (*vaultClient, func()) {
	containerResult, err := StartVaultContainer(ctx)
	require.NoError(t, err, "Failed to start Vault container")

	setupResult, err := SetupVault(containerResult.URL, containerResult.Token)
	if err != nil {
		containerResult.Cleanup()
		t.Fatalf("Failed to setup Vault: %v", err)
	}

	client, err := createClient(containerResult.URL, setupResult.RoleID, setupResult.SecretID)
	if err != nil {
		containerResult.Cleanup()
		t.Fatalf("Failed to create Vault client: %v", err)
	}

	return client, containerResult.Cleanup
}

// enableSecretsEngine enables a secrets engine at a given path
func enableSecretsEngine(client *http.Client, vaultURL, token, path, engineType, description string) error {
	url := fmt.Sprintf("%s/v1/sys/mounts/%s", vaultURL, path)

	data := map[string]any{
		"type":        engineType,
		"description": description,
	}

	_, err := vaultRequest(client, url, token, http.MethodPost, data)
	return err
}

// enableAuthMethod enables an auth method at a given path
func enableAuthMethod(client *http.Client, vaultURL, token, path, methodType, description string) error {
	url := fmt.Sprintf("%s/v1/sys/auth/%s", vaultURL, path)

	data := map[string]any{
		"type":        methodType,
		"description": description,
	}

	_, err := vaultRequest(client, url, token, http.MethodPost, data)
	return err
}

// createAppRole creates an AppRole with a role ID and secret ID and returns a VaultSetupResult
func createAppRole(client *http.Client, vaultURL, token string) (*VaultSetupResult, error) {
	// Create a policy that allows reading from kv-v2
	policyURL := fmt.Sprintf("%s/v1/sys/policies/acl/test-policy", vaultURL)
	policyData := map[string]any{
		"policy": `path "kv-v2/data/*" {capabilities = ["create", "read", "update", "delete", "list"]}`,
	}

	if _, err := vaultRequest(client, policyURL, token, http.MethodPut, policyData); err != nil {
		return nil, fmt.Errorf("failed to create policy: %w", err)
	}

	// Create the AppRole role with the policy
	roleURL := fmt.Sprintf("%s/v1/auth/approle/role/test-role", vaultURL)
	roleData := map[string]any{
		"token_ttl":     "1h",
		"token_max_ttl": "4h",
		"policies":      []string{"test-policy"},
	}

	if _, err := vaultRequest(client, roleURL, token, http.MethodPost, roleData); err != nil {
		return nil, fmt.Errorf("failed to create AppRole role: %w", err)
	}

	// Get the role ID
	roleIDURL := fmt.Sprintf("%s/v1/auth/approle/role/test-role/role-id", vaultURL)
	roleIDResp, err := vaultGetRequest(client, roleIDURL, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get role ID: %w", err)
	}
	defer roleIDResp.Body.Close()

	var roleIDData map[string]any
	if err := json.NewDecoder(roleIDResp.Body).Decode(&roleIDData); err != nil {
		return nil, fmt.Errorf("failed to decode role ID response: %w", err)
	}

	roleID := roleIDData["data"].(map[string]any)["role_id"].(string)

	// Generate a secret ID
	secretIDURL := fmt.Sprintf("%s/v1/auth/approle/role/test-role/secret-id", vaultURL)
	secretIDResp, err := vaultRequest(client, secretIDURL, token, http.MethodPost, map[string]any{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret ID: %w", err)
	}

	var secretIDData map[string]any
	if err := json.NewDecoder(secretIDResp.Body).Decode(&secretIDData); err != nil {
		return nil, fmt.Errorf("failed to decode secret ID response: %w", err)
	}
	_ = secretIDResp.Body.Close()

	secretID := secretIDData["data"].(map[string]any)["secret_id"].(string)

	return &VaultSetupResult{
		RoleID:   roleID,
		SecretID: secretID,
	}, nil
}

// vaultRequest makes an authenticated HTTP request to Vault and returns the response
func vaultRequest(client *http.Client, url, token, method string, data any) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vault-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return nil, fmt.Errorf("Vault request failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// vaultGetRequest makes an authenticated GET request to Vault
func vaultGetRequest(client *http.Client, url, token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Vault-Token", token)

	return client.Do(req)
}

func createClient(vaultURL, roleID string, secretID string) (*vaultClient, error) {
	vaultClient, err := newVaultClient(vaultURL, roleID, secretID, system.NoopMonitor{})
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %w", err)
	}
	err = vaultClient.init(context.Background())
	if err != nil {
		return nil, err
	}
	return vaultClient, nil
}
