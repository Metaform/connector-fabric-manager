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

	"github.com/metaform/connector-fabric-manager/agent/keycloak/activity"
	"github.com/metaform/connector-fabric-manager/assembly/httpclient"
	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsagent"
)

const (
	ActivityType = "keycloak-activity"
	AgentPrefix  = "kcagent"
	urlKey       = "url"
	tokenKey     = "token"
	realmKey     = "realm"
)

func LaunchAndWaitSignal(shutdown <-chan struct{}) {
	config := natsagent.LauncherConfig{
		AgentName:    "KeyCloak Agent",
		ConfigPrefix: AgentPrefix,
		ActivityType: ActivityType,
		AssemblyProvider: func() []system.ServiceAssembly {
			return []system.ServiceAssembly{&httpclient.HttpClientServiceAssembly{}}
		},
		NewProcessor: func(ctx *natsagent.AgentContext) api.ActivityProcessor {
			client := ctx.Registry.Resolve(httpclient.HttpClientKey).(http.Client)
			url := ctx.Config.GetString(urlKey)
			token := ctx.Config.GetString(tokenKey)
			realm := ctx.Config.GetString(realmKey)
			if err := runtime.CheckRequiredParams(urlKey, url, tokenKey, token, realmKey, realm); err != nil {
				panic(err)
			}
			return activity.NewProcessor(&activity.Config{
				KeycloakURL: url,
				Token:       token,
				Realm:       realm,
				HTTPClient:  &client,
				Monitor:     ctx.Monitor,
			})
		},
	}
	natsagent.LaunchAgent(shutdown, config)
}
