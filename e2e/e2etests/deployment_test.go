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
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/metaform/connector-fabric-manager/common/dmodel"
	"github.com/metaform/connector-fabric-manager/common/natstestfixtures"
	alauncher "github.com/metaform/connector-fabric-manager/pmanager/agent/testagent/launcher"
	plauncher "github.com/metaform/connector-fabric-manager/pmanager/cmd/server/launcher"
	tlauncher "github.com/metaform/connector-fabric-manager/tmanager/cmd/server/launcher"
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

	_ = os.Setenv("TM_URI", nt.Uri)
	_ = os.Setenv("TM_BUCKET", cfmBucket)
	_ = os.Setenv("TM_STREAM", streamName)

	_ = os.Setenv("PM_URI", nt.Uri)
	_ = os.Setenv("PM_BUCKET", cfmBucket)
	_ = os.Setenv("PM_STREAM", streamName)

	_ = os.Setenv("TESTAGENT_URI", nt.Uri)
	_ = os.Setenv("TESTAGENT_BUCKET", cfmBucket)
	_ = os.Setenv("TESTAGENT_STREAM", streamName)

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

	//m := &dmodel.DeploymentResponse{
	//	ID:             "test-deployment-123",
	//	Success:        true,
	//	ErrorDetail:    "",
	//	ManifestID:     "1234567890",
	//	DeploymentType: dmodel.VpaDeploymentType,
	//	Properties:     make(map[string]any),
	//}

	m := &dmodel.DeploymentManifest{
		ID:             "test-deployment-123",
		DeploymentType: dmodel.VpaDeploymentType,
		Payload:        make(map[string]any),
	}
	ser, err := json.Marshal(m)
	require.NoError(t, err)

	_, err = nt.Client.JetStream.Publish(context.Background(), "event.cfm-deployment", ser)
	//_, err = nt.Client.JetStream.Publish(context.Background(), "event.cfm-deployment-response", ser)

	require.NoError(t, err)

}
