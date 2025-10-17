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
	"os"
	"testing"
	"time"

	"github.com/metaform/connector-fabric-manager/common/natsclient"
	"github.com/metaform/connector-fabric-manager/common/natstestfixtures"
	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/common/system"

	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsorchestration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testTimeout  = 30 * time.Second
	pollInterval = 100 * time.Millisecond
	streamName   = "cfm-activity"
)

func TestTestAgent_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Set up NATS container
	nt, err := natstestfixtures.SetupNatsContainer(ctx, "cfm-bucket")

	require.NoError(t, err)

	defer natstestfixtures.TeardownNatsContainer(ctx, nt)

	natstestfixtures.SetupTestStream(t, ctx, nt.Client, streamName)

	// Set up an orchestration for the test agent to process
	orchestration := api.Orchestration{
		ID:             "test-agent-orchestration",
		State:          api.OrchestrationStateRunning,
		Completed:      make(map[string]struct{}),
		ProcessingData: make(map[string]any),
		Steps: []api.OrchestrationStep{
			{
				Activities: []api.Activity{
					{ID: "test-activity", Type: "test.activity"},
				},
			},
		},
	}

	// Required agent config
	_ = os.Setenv("TESTAGENT_URI", nt.Uri)
	_ = os.Setenv("TESTAGENT_BUCKET", "cfm-bucket")
	_ = os.Setenv("TESTAGENT_STREAM", streamName)

	// Create and start the test agent
	shutdownChannel := make(chan struct{})
	go func() {
		Launch(shutdownChannel)
	}()

	// Submit orchestration
	adapter := natsclient.NewMsgClient(nt.Client)
	logMonitor := runtime.LoadLogMonitor("test-agent", system.DevelopmentMode)
	orchestrator := natsorchestration.NewNatsDeploymentOrchestrator(adapter, logMonitor)

	err = orchestrator.ExecuteOrchestration(ctx, &orchestration)
	require.NoError(t, err)

	// Wait for the activity to be processed
	assert.Eventually(t, func() bool {
		updatedOrchestration, _, err := natsorchestration.ReadOrchestration(ctx, orchestration.ID, adapter)
		require.NoError(t, err)
		return updatedOrchestration.State == api.OrchestrationStateCompleted
	}, testTimeout, pollInterval, "Activity should be processed")

	// shut agent down
	shutdownChannel <- struct{}{}
}
