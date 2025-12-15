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

package launcher

import (
	"net/http"

	"github.com/metaform/connector-fabric-manager/agent/edcv/activity"
	"github.com/metaform/connector-fabric-manager/agent/edcv/identityhub"
	"github.com/metaform/connector-fabric-manager/assembly/httpclient"
	"github.com/metaform/connector-fabric-manager/assembly/serviceapi"
	"github.com/metaform/connector-fabric-manager/assembly/vault"
	"github.com/metaform/connector-fabric-manager/common/oauth2"
	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsagent"
)

const (
	urlKey             = "vault.url" // duplicate of common/vault/assembly.go
	ActivityType       = "edcv-activity"
	clientIDKey        = "keycloak.clientID"
	clientSecretKey    = "keycloak.clientSecret"
	tokenURLKey        = "keycloak.tokenurl"
	identityHubURLKey  = "identityhub.url"
	controlPlaneURLKey = "controlplane.url"
)

func LaunchAndWaitSignal(shutdown <-chan struct{}) {
	config := natsagent.LauncherConfig{
		AgentName:    "EDC-V Agent",
		ConfigPrefix: "edcvagent",
		ActivityType: ActivityType,
		AssemblyProvider: func() []system.ServiceAssembly {
			return []system.ServiceAssembly{
				&httpclient.HttpClientServiceAssembly{},
				&vault.VaultServiceAssembly{},
			}
		},
		NewProcessor: func(ctx *natsagent.AgentContext) api.ActivityProcessor {
			httpClient := ctx.Registry.Resolve(serviceapi.HttpClientKey).(http.Client)
			vaultClient := ctx.Registry.Resolve(serviceapi.VaultKey).(serviceapi.VaultClient)
			clientID := ctx.Config.GetString(clientIDKey)
			clientSecret := ctx.Config.GetString(clientSecretKey)
			tokenURL := ctx.Config.GetString(tokenURLKey) // this may be nil or "" if the in-mem vault is used
			ihURL := ctx.Config.GetString(identityHubURLKey)
			cpURL := ctx.Config.GetString(controlPlaneURLKey)
			vaultURL := ctx.Config.GetString(urlKey) // this may be nil or "" if the in-mem vault is used

			if err := runtime.CheckRequiredParams(clientIDKey, clientID, clientSecretKey, clientSecret, identityHubURLKey, ihURL, controlPlaneURLKey, cpURL); err != nil {
				panic(err)
			}

			return activity.NewProcessor(&activity.Config{
				VaultClient: vaultClient,
				Client:      &httpClient,
				LogMonitor:  ctx.Monitor,
				IdentityAPIClient: identityhub.HttpIdentityAPIClient{
					BaseURL: ihURL,
					TokenProvider: oauth2.NewTokenProvider(
						oauth2.Oauth2Params{
							ClientID:     clientID,
							ClientSecret: clientSecret,
							TokenURL:     tokenURL,
							GrantType:    oauth2.ClientCredentials,
						}, &httpClient),
					HttpClient: &httpClient,
				},
				TokenURL: tokenURL,
				VaultURL: vaultURL,
			})
		},
	}
	natsagent.LaunchAgent(shutdown, config)
}
