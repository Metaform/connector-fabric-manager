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

package natsorchestration_test_test

import (
	"context"
	"fmt"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/pmanager/natsorchestration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

// TestExecuteOrchestration_ParallelActivitiesOneFailsFirst verifies orchestration with parallel activities where one fails first.
// The orchestration should be in the errored state, i.e. the successful process should not change from an errored state.
func TestExecuteOrchestration_ParallelActivitiesOneFailsFirst(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nt, err := natsorchestration.SetupNatsContainer(ctx, "cfm-durable-activity-bucket")
	require.NoError(t, err)

	defer natsorchestration.TeardownNatsContainer(ctx, nt)

	stream := natsorchestration.SetupTestStream(t, ctx, nt.Client, "cfm-activity")
	natsorchestration.SetupTestConsumer(t, ctx, stream, "test.fail.activity")
	natsorchestration.SetupTestConsumer(t, ctx, stream, "test.succeed.activity")

	// Create an orchestration with two parallel activities
	orchestration := api.Orchestration{
		ID:    "test-parallel-fail-succeed",
		State: api.OrchestrationStateRunning,
		Steps: []api.OrchestrationStep{
			{
				Activities: []api.Activity{
					{ID: "A1", Type: "test.fail.activity"},
					{ID: "A2", Type: "test.succeed.activity"},
				},
			},
		},
		Completed: make(map[string]struct{}),
	}

	adapter := natsorchestration.NatsClientAdapter{Client: nt.Client}

	// WaitGroup to coordinate activity execution order
	var activityWg sync.WaitGroup
	activityWg.Add(1) // Only the failing activity signals completion

	// WaitGroup to wait for both activities to complete
	var verificationWg sync.WaitGroup
	verificationWg.Add(2) // Both activities

	// Failing activity processor
	failProcessor := TestActivityProcessor{
		onProcess: func(id string) {
			activityWg.Done() // Signal that failing activity completed
			verificationWg.Done()
		},
	}

	// Succeeding activity processor
	succeedProcessor := TestActivityProcessor{
		onProcess: func(id string) {
			activityWg.Wait() // Wait for the failing activity to complete first
			verificationWg.Done()
		},
	}

	noOpMonitor := monitor.NoopMonitor{}

	// Create executor for failing activity
	failExecutor := natsorchestration.NatsActivityExecutor{
		Client:       adapter,
		StreamName:   "cfm-activity",
		ActivityType: "test.fail.activity",
		ActivityProcessor: FailingActivityProcessor{
			testProcessor: failProcessor,
		},
		Monitor: noOpMonitor,
	}

	// Create executor for succeeding activity
	succeedExecutor := natsorchestration.NatsActivityExecutor{
		Client:            adapter,
		StreamName:        "cfm-activity",
		ActivityType:      "test.succeed.activity",
		ActivityProcessor: succeedProcessor,
		Monitor:           noOpMonitor,
	}

	// Start both executors
	err = failExecutor.Execute(ctx)
	require.NoError(t, err)

	err = succeedExecutor.Execute(ctx)
	require.NoError(t, err)

	// Start orchestration
	orchestrator := natsorchestration.NatsDeploymentOrchestrator{Client: adapter}
	err = orchestrator.ExecuteOrchestration(ctx, &orchestration)
	require.NoError(t, err)

	// Wait for both activities to complete
	verificationWg.Wait()

	// Verify orchestration is in an error state
	var finalOrchestration api.Orchestration
	timeout := time.After(3 * time.Second)
outerLoop:
	for {
		select {
		case <-timeout:
			t.Fatalf("Timeout waiting for activity A2 to complete after 3 seconds")
		default:
			finalOrchestration, _, err = natsorchestration.ReadOrchestration(ctx, orchestration.ID, adapter)
			require.NoError(t, err)
			if _, found := finalOrchestration.Completed["A2"]; found {
				break outerLoop
			}
		}
	}

	assert.Equal(t, api.OrchestrationStateErrored, finalOrchestration.State,
		"Orchestration should be in error state after activity failure")
}

// FailingActivityProcessor wraps a TestActivityProcessor and always returns an error
type FailingActivityProcessor struct {
	testProcessor TestActivityProcessor
}

func (f FailingActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	// Execute the test processor logic first (for timing coordination)
	if f.testProcessor.onProcess != nil {
		f.testProcessor.onProcess(ctx.ID())
	}

	// Always return error result
	return api.ActivityResult{
		Result: api.ActivityResultFatalError,
		Error:  fmt.Errorf("simulated activity failure for %s", ctx.ID()),
	}
}
