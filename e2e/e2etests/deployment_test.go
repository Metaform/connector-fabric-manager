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
	"os"
	"testing"
	"time"

	"github.com/metaform/connector-fabric-manager/common/natstestfixtures"
	"github.com/stretchr/testify/require"
)

const (
	testTimeout = 30 * time.Second
	streamName  = "cfm-deployment"
)

func Test_VerifyE2E(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Set up NATS container
	nt, err := natstestfixtures.SetupNatsContainer(ctx, "test-agent-bucket")

	require.NoError(t, err)

	defer natstestfixtures.TeardownNatsContainer(ctx, nt)

	_ = os.Setenv("TM_URI", nt.Uri)
	_ = os.Setenv("TM_BUCKET", "cfm-bucket")
	_ = os.Setenv("TM_STREAM", "cfm-stream")

	_ = os.Setenv("PM_URI", nt.Uri)
	_ = os.Setenv("PM_BUCKET", "cfm-bucket")
	_ = os.Setenv("PM_STREAM", "cfm-stream")

	//shutdownChannel := make(chan struct{})
	//go func() {
	//	planucher.Launch(shutdownChannel)
	//}()
	//
	//go func() {
	//	tlauncher.Launch(shutdownChannel)
	//}()

}
