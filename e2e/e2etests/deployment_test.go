// Copyright (c) 2025 Metaform Systems, Inc
//
// This program and the accompanying materials are made available under the
// terms of the Apache License, Version 2.0 which is available at
// https://www.apache.org/licenses/LICENSE-2.0
//
// SPDX-License-Identifier: Apache-2.0
//
// Contributors:
//
//	Metaform Systems, Inc. - initial API and implementation
package e2etests

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/metaform/connector-fabric-manager/common/natstestfixtures"
	"github.com/metaform/connector-fabric-manager/common/testfixtures"
	"github.com/metaform/connector-fabric-manager/e2e/e2efixtures"
	alauncher "github.com/metaform/connector-fabric-manager/pmanager/agent/testagent/launcher"
	plauncher "github.com/metaform/connector-fabric-manager/pmanager/cmd/server/launcher"
	tlauncher "github.com/metaform/connector-fabric-manager/tmanager/cmd/server/launcher"
	"github.com/metaform/connector-fabric-manager/tmanager/model/v1alpha1"
	"github.com/stretchr/testify/require"
)

const (
	testTimeout = 30 * time.Second
	streamName  = "cfm-stream"
	cfmBucket   = "cfm-bucket"
)

func Test_VerifyE2E(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	nt, err := natstestfixtures.SetupNatsContainer(ctx, cfmBucket)

	require.NoError(t, err)

	defer natstestfixtures.TeardownNatsContainer(ctx, nt)
	defer cleanup()

	_ = os.Setenv("TM_URI", nt.Uri)
	_ = os.Setenv("TM_BUCKET", cfmBucket)
	_ = os.Setenv("TM_STREAM", streamName)

	_ = os.Setenv("PM_URI", nt.Uri)
	_ = os.Setenv("PM_BUCKET", cfmBucket)
	_ = os.Setenv("PM_STREAM", streamName)

	_ = os.Setenv("TESTAGENT_URI", nt.Uri)
	_ = os.Setenv("TESTAGENT_BUCKET", cfmBucket)
	_ = os.Setenv("TESTAGENT_STREAM", streamName)

	tPort := testfixtures.GetRandomPort(t)
	_ = os.Setenv("TM_HTTPPORT", strconv.Itoa(tPort))
	pPort := testfixtures.GetRandomPort(t)
	_ = os.Setenv("PM_HTTPPORT", strconv.Itoa(pPort))

	shutdownChannel := make(chan struct{})
	go func() {
		plauncher.Launch(shutdownChannel)
	}()

	go func() {
		tlauncher.Launch(shutdownChannel)
	}()

	go func() {
		alauncher.Launch(shutdownChannel)
	}()

	client := e2efixtures.NewApiClient(fmt.Sprintf("http://localhost:%d", tPort), fmt.Sprintf("http://localhost:%d", pPort))
	// Wait for the pmanager to be ready
	for start := time.Now(); time.Since(start) < 5*time.Second; {
		if err = e2efixtures.CreateTestActivityDefinition(client); err == nil {
			break
		}
	}
	require.NoError(t, err)

	err = e2efixtures.CreateTestDeploymentDefinition(client)
	require.NoError(t, err)

	cell, err := e2efixtures.CreateCell(client)
	require.NoError(t, err)

	profile, err := e2efixtures.CreateDataspaceProfile(client)
	require.NoError(t, err)

	deployment := v1alpha1.NewDataspaceProfileDeployment{
		ProfileID: profile.ID,
		CellID:    cell.ID,
	}
	err = e2efixtures.DeployDataspaceProfile(deployment, client)
	require.NoError(t, err)

	fmt.Println(cell)
	fmt.Println(profile)
	err = client.PostToTManager("participant/"+uuid.New().String(), "{}")

	require.NoError(t, err)
}

func cleanup() {
	_ = os.Unsetenv("TM_URI")
	_ = os.Unsetenv("TM_BUCKET")
	_ = os.Unsetenv("TM_STREAM")

	_ = os.Unsetenv("PM_URI")
	_ = os.Unsetenv("PM_BUCKET")
	_ = os.Unsetenv("PM_STREAM")

	_ = os.Unsetenv("TESTAGENT_URI")
	_ = os.Unsetenv("TESTAGENT_BUCKET")
	_ = os.Unsetenv("TESTAGENT_STREAM")

	_ = os.Unsetenv("TM_HTTPPORT")
	_ = os.Unsetenv("PM_HTTPPORT")
}
