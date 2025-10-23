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
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/metaform/connector-fabric-manager/common/natstestfixtures"
	"github.com/metaform/connector-fabric-manager/e2e/e2efixtures"
	dnsauncher "github.com/metaform/connector-fabric-manager/mvd/dns/launcher"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	plauncher "github.com/metaform/connector-fabric-manager/pmanager/cmd/server/launcher"
	pv1alpha1 "github.com/metaform/connector-fabric-manager/pmanager/model/v1alpha1"
	tlauncher "github.com/metaform/connector-fabric-manager/tmanager/cmd/server/launcher"
)

const (
	streamName = "cfm-stream"
	cfmBucket  = "cfm-bucket"
)

func LaunchMVD() {
	ctx := context.Background()

	nt, err := natstestfixtures.SetupNatsContainer(ctx, cfmBucket)
	if err != nil {
		panic(err)
	}
	defer natstestfixtures.TeardownNatsContainer(ctx, nt)

	_ = os.Setenv("DNSAGENT_URI", nt.Uri)
	_ = os.Setenv("DNSAGENT_BUCKET", cfmBucket)
	_ = os.Setenv("DNSAGENT_STREAM", streamName)

	_ = os.Setenv("TM_URI", nt.Uri)
	_ = os.Setenv("TM_BUCKET", cfmBucket)
	_ = os.Setenv("TM_STREAM", streamName)

	_ = os.Setenv("PM_URI", nt.Uri)
	_ = os.Setenv("PM_BUCKET", cfmBucket)
	_ = os.Setenv("PM_STREAM", streamName)

	tPort := GetRandomPort()
	_ = os.Setenv("TM_HTTPPORT", strconv.Itoa(tPort))
	pPort := GetRandomPort()
	_ = os.Setenv("PM_HTTPPORT", strconv.Itoa(pPort))

	shutdownChannel := make(chan struct{})
	go func() {
		plauncher.Launch(shutdownChannel)
	}()

	go func() {
		tlauncher.Launch(shutdownChannel)
	}()

	go func() {
		dnsauncher.LaunchAndWaitSignal(shutdownChannel)
	}()

	client := e2efixtures.NewApiClient(fmt.Sprintf("http://localhost:%d", tPort), fmt.Sprintf("http://localhost:%d", pPort))
	// Wait for the pmanager to be ready
	for start := time.Now(); time.Since(start) < 5*time.Second; {
		if err = CreateActivityDefinitions(client); err == nil {
			break
		}
	}

	err = CreateDeploymentDefinition(client)
	if err != nil {
		panic(err)
	}
}

func GetRandomPort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port
}

func CreateActivityDefinitions(apiClient *e2efixtures.ApiClient) error {
	requestBody := api.ActivityDefinition{
		Type:        dnsauncher.ActivityType,
		Description: "Provisions A DNS subdomain and ingress routing",
	}

	return apiClient.PostToPManager("activity-definition", requestBody)
}

func CreateDeploymentDefinition(apiClient *e2efixtures.ApiClient) error {
	requestBody := pv1alpha1.DeploymentDefinition{
		Type: dnsauncher.ActivityType,
		Activities: []pv1alpha1.Activity{
			{
				ID:   "dns-provisioner",
				Type: dnsauncher.ActivityType,
			},
		},
	}
	return apiClient.PostToPManager("deployment-definition", requestBody)
}
