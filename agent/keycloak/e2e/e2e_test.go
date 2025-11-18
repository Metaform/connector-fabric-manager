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

package e2e

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/metaform/connector-fabric-manager/agent/keycloak/fixtures"
	"github.com/metaform/connector-fabric-manager/agent/keycloak/launcher"
	"github.com/metaform/connector-fabric-manager/assembly/vault"
	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/common/natsfixtures"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/stretchr/testify/require"
)

const (
	streamName = "cfm-stream"
	cfmBucket  = "cfm-bucket"

	orchestrationID = "1234"
)

func Test_Launch(t *testing.T) {
	ctx := context.Background()

	keycloakURL, token, err := fixtures.SetupKeyCloakContainer(ctx)
	require.NoError(t, err, "Failed to start Keycloak container")

	nt, err := natsfixtures.SetupNatsContainer(ctx, cfmBucket)

	require.NoError(t, err)

	containerResult, err := vault.StartVaultContainer(ctx)
	require.NoError(t, err, "Failed to start Vault container")

	setupResult, err := vault.SetupVault(containerResult.URL, containerResult.Token)
	if err != nil {
		containerResult.Cleanup()
		t.Fatalf("Failed to setup Vault: %v", err)
	}

	_ = os.Setenv("KCAGENT_VAULT_URL", containerResult.URL)
	_ = os.Setenv("KCAGENT_VAULT_ROLEID", setupResult.RoleID)
	_ = os.Setenv("KCAGENT_VAULT_SECRETID", setupResult.SecretID)

	_ = os.Setenv("KCAGENT_URI", nt.Uri)
	_ = os.Setenv("KCAGENT_BUCKET", cfmBucket)
	_ = os.Setenv("KCAGENT_STREAM", streamName)
	_ = os.Setenv("KCAGENT_KEYCLOAK_URL", keycloakURL)
	_ = os.Setenv("KCAGENT_KEYCLOAK_TOKEN", token)
	_ = os.Setenv("KCAGENT_KEYCLOAK_REALM", "edcv")

	shutdownChannel := make(chan struct{})
	go func() {
		launcher.LaunchAndWaitSignal(shutdownChannel)
	}()

	err = createOrchestration(ctx, orchestrationID, nt.Client)
	require.NoError(t, err)

	err = publishActivityMessage(ctx, orchestrationID, nt.Client)
	require.NoError(t, err)

	vaultClient, err := vault.NewVaultClient(containerResult.URL, setupResult.RoleID, setupResult.SecretID)
	require.NoError(t, err)
	defer vaultClient.Close()

	var clientID string
	require.Eventually(t, func() bool {
		oEntry, err := nt.Client.KVStore.Get(ctx, "1234")
		if err != nil {
			return false
		}
		var orchestration api.Orchestration
		err = json.Unmarshal(oEntry.Value(), &orchestration)
		if err != nil {
			return false
		}
		if orchestration.State == api.OrchestrationStateCompleted {
			clientID = orchestration.ProcessingData["clientID"].(string)
			return true
		}
		return false
	}, 10*time.Second, 10*time.Millisecond, "Orchestration did not complete in time")

	require.NotEmpty(t, clientID, "Expected clientID to be set")
	clientSecret, err := vaultClient.ResolveSecret(ctx, clientID)
	require.NoError(t, err, "Failed to resolve secret")
	require.NotEmpty(t, clientSecret, "Expected client secret to be set")

	shutdownChannel <- struct{}{}
}

func createOrchestration(ctx context.Context, id string, client *natsclient.NatsClient) error {
	orchestration := api.Orchestration{
		ID:                id,
		CorrelationID:     "correlation-id",
		State:             0,
		OrchestrationType: "test",
		Steps: []api.OrchestrationStep{
			{
				Activities: []api.Activity{
					{ID: "test-activity", Type: launcher.ActivityType},
				},
			},
		},
		ProcessingData: make(map[string]any),
		OutputData:     make(map[string]any),
		Completed:      make(map[string]struct{}),
	}
	serialized, err := json.Marshal(orchestration)
	if err != nil {
		return err
	}
	_, err = client.KVStore.Update(ctx, "1234", serialized, 0)
	return err
}

func publishActivityMessage(ctx context.Context, id string, client *natsclient.NatsClient) error {
	message := &api.ActivityMessage{
		OrchestrationID: id,
		Activity: api.Activity{
			ID:            "test-activity",
			Type:          launcher.ActivityType,
			Discriminator: "",
			Inputs:        make([]api.MappingEntry, 0),
			DependsOn:     make([]string, 0),
		},
	}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	subject := natsclient.CFMSubjectPrefix + "." + launcher.ActivityType
	_, err = client.JetStream.Publish(ctx, subject, data)
	return nil
}
