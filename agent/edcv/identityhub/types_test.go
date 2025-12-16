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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewParticipantManifest_WithDefaults(t *testing.T) {
	manifest := NewParticipantManifest("test-id", "did:web:foo", "http://example.com/credentialservice", "http://example.com/protocol")

	require.Equal(t, manifest.CredentialServiceID, "test-id-credentialservice")
	require.Equal(t, manifest.ProtocolServiceID, "test-id-dsp")
	require.Equal(t, manifest.IsActive, true)
	require.Equal(t, manifest.KeyGeneratorParameters.KeyID, "did:web:foo#"+DefaultKeyID)
	require.Equal(t, manifest.KeyGeneratorParameters.PrivateKeyAlias, "did:web:foo#"+DefaultKeyID)
	require.Equal(t, manifest.VaultConfig.SecretPath, "v1/participants")
	require.Equal(t, manifest.VaultConfig.FolderPath, "test-id/identityhub")
}
