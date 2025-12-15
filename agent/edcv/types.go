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

package edcv

// VaultConfig defines configuration for accessing a vault, including its URL, secret path, and folder path.
type VaultConfig struct {
	// VaultURL the base URL of the vault
	VaultURL string `json:"vaultUrl"`
	// SecretPath the path of the mount point of the secret engine
	SecretPath string `json:"secretPath"`
	// FolderPath the path of the folder within the secret engine where the participant's manifest will be stored. Note that this will be prefixed with the
	// participant context ID
	FolderPath string `json:"folderPath"`
}

// VaultCredentials defines the credentials which are needed to get a JWT which is used to access a vault
type VaultCredentials struct {
	// ClientID client ID of the service account of the IdP which is configured in Vault
	ClientID string `json:"clientId"`
	// ClientSecret secret of the service account of the IdP, which is configured in Vault
	ClientSecret string `json:"clientSecret"`
	// TokenURL URL of the token endpoint of the IdP which is configured in Vault
	TokenURL string `json:"tokenUrl"`
}
