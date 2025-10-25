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

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/natstestfixtures"
	"github.com/metaform/connector-fabric-manager/e2e/e2efixtures"
	clauncher "github.com/metaform/connector-fabric-manager/mvd/connector/launcher"
	dnslauncher "github.com/metaform/connector-fabric-manager/mvd/dns/launcher"
	oblauncher "github.com/metaform/connector-fabric-manager/mvd/onboard/launcher"
	papi "github.com/metaform/connector-fabric-manager/pmanager/api"
	plauncher "github.com/metaform/connector-fabric-manager/pmanager/cmd/server/launcher"
	pv1alpha1 "github.com/metaform/connector-fabric-manager/pmanager/model/v1alpha1"
	tapi "github.com/metaform/connector-fabric-manager/tmanager/api"
	tlauncher "github.com/metaform/connector-fabric-manager/tmanager/cmd/server/launcher"
	"github.com/metaform/connector-fabric-manager/tmanager/model/v1alpha1"
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

	_ = os.Setenv("CAGENT_URI", nt.Uri)
	_ = os.Setenv("CAGENT_BUCKET", cfmBucket)
	_ = os.Setenv("CAGENT_STREAM", streamName)

	_ = os.Setenv("OBAGENT_URI", nt.Uri)
	_ = os.Setenv("OBAGENT_BUCKET", cfmBucket)
	_ = os.Setenv("OBAGENT_STREAM", streamName)

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
		dnslauncher.LaunchAndWaitSignal(shutdownChannel)
	}()

	go func() {
		clauncher.LaunchAndWaitSignal(shutdownChannel)
	}()

	go func() {
		oblauncher.LaunchAndWaitSignal(shutdownChannel)
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

	cell, err := e2efixtures.CreateCell(client)
	if err != nil {
		panic(err)
	}

	dProfile, err := e2efixtures.CreateDataspaceProfile(client)
	if err != nil {
		panic(err)
	}

	deployment := v1alpha1.NewDataspaceProfileDeployment{
		ProfileID: dProfile.ID,
		CellID:    cell.ID,
	}
	err = e2efixtures.DeployDataspaceProfile(deployment, client)
	if err != nil {
		panic(err)
	}

	newProfile := v1alpha1.NewParticipantProfileDeployment{
		Identifier:    "did:web:foo.com",
		VPAProperties: map[string]map[string]any{string(model.ConnectorType): {"connectorkey": "connectorvalue"}},
	}
	var participantProfile v1alpha1.ParticipantProfile
	err = client.PostToTManagerWithResponse("participants", newProfile, &participantProfile)
	if err != nil {
		panic(err)
	}

	var statusProfile v1alpha1.ParticipantProfile

	// Verify all VPAs are active
	deployCount := 0
	for start := time.Now(); time.Since(start) < 5*time.Second; {
		err = client.GetTManager(fmt.Sprintf("participants/%s", participantProfile.ID), &statusProfile)
		for _, vpa := range statusProfile.VPAs {
			if vpa.State == tapi.DeploymentStateActive.String() {
				deployCount++
			}
		}
		if deployCount == 3 {
			break
		}
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
	err := CreateActivityDefinition(apiClient, dnslauncher.ActivityType, "Provisions a DNS subdomain and ingress routing")
	if err != nil {
		return err
	}
	err = CreateActivityDefinition(apiClient, clauncher.ActivityType, "Provisions Connector VPA")
	if err != nil {
		return err
	}

	return CreateActivityDefinition(apiClient, oblauncher.ActivityType, "Performs onboarding")
}

func CreateActivityDefinition(apiClient *e2efixtures.ApiClient, activityType string, description string) error {
	requestBody := papi.ActivityDefinition{
		Type:        papi.ActivityType(activityType),
		Description: description,
	}
	return apiClient.PostToPManager("activity-definition", requestBody)
}

func CreateDeploymentDefinition(apiClient *e2efixtures.ApiClient) error {
	requestBody := pv1alpha1.DeploymentDefinition{
		Type: model.VPADeploymentType.String(),
		Activities: []pv1alpha1.Activity{
			{
				ID:   "dns-provisioner",
				Type: dnslauncher.ActivityType,
			},
			{
				ID:   "connector-provisioner",
				Type: clauncher.ActivityType,
				DependsOn: []string{
					"dns-provisioner",
				},
			},
			{
				ID:   "onboarder",
				Type: oblauncher.ActivityType,
				DependsOn: []string{
					"connector-provisioner",
				},
			},
		},
	}
	return apiClient.PostToPManager("deployment-definition", requestBody)
}
