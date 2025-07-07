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
	"github.com/metaform/connector-fabric-manager/common/testfixtures"
	"github.com/metaform/connector-fabric-manager/pmanager/natsorchestration"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	testTimeout = 30 * time.Second
	streamName  = "cfm-orchestration"
)

func TestTestAgent_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Set up NATS container
	nt, err := natsorchestration.SetupNatsContainer(ctx, "test-orchestration-bucket")

	require.NoError(t, err)

	defer natsorchestration.TeardownNatsContainer(ctx, nt)

	natsorchestration.SetupTestStream(t, ctx, nt.Client, streamName)

	// Required agent config
	_ = os.Setenv("PM_URI", nt.Uri)
	_ = os.Setenv("PM_BUCKET", "test-agent-bucket")
	_ = os.Setenv("PM_STREAM", streamName)
	_ = os.Setenv("PM_HTTPPORT", strconv.Itoa(testfixtures.GetRandomPort(t)))

	// Create and start the test agent
	shutdownChannel := make(chan struct{})
	go func() {
		Launch(shutdownChannel)
	}()

	// shut agent down
	shutdownChannel <- struct{}{}
}
