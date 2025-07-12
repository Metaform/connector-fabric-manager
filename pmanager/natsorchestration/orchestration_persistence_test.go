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

package natsorchestration_test

import (
	"context"
	"fmt"
	"github.com/metaform/connector-fabric-manager/pmanager/natsorchestration"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	persistenceTimeout = 10 * time.Second
	pollInterval       = 10 * time.Millisecond
	streamName         = "cfm-activity"
)

func Test_ValuePersistence(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), persistenceTimeout)
	defer cancel()

	nt, err := natsorchestration.SetupNatsContainer(ctx, "cfm-activity-context-bucket")
	require.NoError(t, err)

	defer natsorchestration.TeardownNatsContainer(ctx, nt)

	stream := natsorchestration.SetupTestStream(t, ctx, nt.Client, streamName)
	natsorchestration.SetupTestConsumer(t, ctx, stream, "test.context.persistence")

	var wg sync.WaitGroup
	wg.Add(1)

	processor := &ValueSettingProcessor{
		onProcess: func(activityCtx api.ActivityContext) {
			defer wg.Done()
			activityCtx.SetValue("string_key", "test_value")
			activityCtx.SetValue("int_key", 42)
			activityCtx.SetValue("bool_key", true)
			activityCtx.SetValue("map_key", map[string]interface{}{
				"nested": "value",
				"count":  123,
			})
		},
	}

	orchestration := createTestOrchestration("test-context-persistence", "test.context.persistence")
	adapter := natsorchestration.NatsClientAdapter{Client: nt.Client}

	orchestrator := &natsorchestration.NatsDeploymentOrchestrator{
		Client:  adapter,
		Monitor: monitor.NoopMonitor{},
	}

	err = orchestrator.ExecuteOrchestration(ctx, &orchestration)
	require.NoError(t, err)

	executor := &natsorchestration.NatsActivityExecutor{
		Client:            adapter,
		StreamName:        "cfm-activity",
		ActivityType:      "test.context.persistence",
		ActivityProcessor: processor,
		Monitor:           monitor.NoopMonitor{},
	}

	err = executor.Execute(ctx)
	require.NoError(t, err)

	// Wait for activity to complete
	wg.Wait()

	// Verify values were persisted
	require.Eventually(t, func() bool {
		updatedOrchestration, _, err := natsorchestration.ReadOrchestration(ctx, orchestration.ID, adapter)
		if err != nil {
			return false
		}

		// Check if all expected values are present
		if updatedOrchestration.ProcessingData["string_key"] == nil {
			return false
		}

		assert.Equal(t, "test_value", updatedOrchestration.ProcessingData["string_key"])
		assert.Equal(t, float64(42), updatedOrchestration.ProcessingData["int_key"]) // JSON unmarshalling converts numbers to float64
		assert.Equal(t, true, updatedOrchestration.ProcessingData["bool_key"])

		mapValue, ok := updatedOrchestration.ProcessingData["map_key"].(map[string]interface{})
		require.True(t, ok, "map_key should be a map[string]interface{}")
		assert.Equal(t, "value", mapValue["nested"])
		assert.Equal(t, float64(123), mapValue["count"])

		return true
	}, persistenceTimeout, pollInterval, "Values should be persisted")
}

func Test_ValuePersistenceOnRetry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), persistenceTimeout)
	defer cancel()

	nt, err := natsorchestration.SetupNatsContainer(ctx, "cfm-activity-retry-bucket")
	require.NoError(t, err)

	defer natsorchestration.TeardownNatsContainer(ctx, nt)

	stream := natsorchestration.SetupTestStream(t, ctx, nt.Client, streamName)
	natsorchestration.SetupTestConsumer(t, ctx, stream, "test.retry.persistence")

	var wg sync.WaitGroup
	var callCount int32

	// Create processor that sets values and fails on first call, succeeds on second
	processor := &RetryWithValueProcessor{
		onProcess: func(activityCtx api.ActivityContext) api.ActivityResult {
			currentCall := atomic.AddInt32(&callCount, 1)

			if currentCall == 1 {
				// First call: set values and return retry error
				activityCtx.SetValue("retry_count", int(currentCall))
				activityCtx.SetValue("first_attempt_data", "initial_value")
				return api.ActivityResult{
					Result: api.ActivityResultRetryError,
					Error:  fmt.Errorf("simulated retry error"),
				}
			}

			// Second call: verify values from first call are available and set additional values
			activityCtx.SetValue("retry_count", int(currentCall))
			activityCtx.SetValue("second_attempt_data", "retry_value")
			wg.Done()
			return api.ActivityResult{
				Result: api.ActivityResultComplete,
			}
		},
	}

	orchestration := createTestOrchestration("test-retry-persistence", "test.retry.persistence")
	adapter := natsorchestration.NatsClientAdapter{Client: nt.Client}

	orchestrator := &natsorchestration.NatsDeploymentOrchestrator{
		Client:  adapter,
		Monitor: monitor.NoopMonitor{},
	}

	err = orchestrator.ExecuteOrchestration(ctx, &orchestration)
	require.NoError(t, err)

	executor := &natsorchestration.NatsActivityExecutor{
		Client:            adapter,
		StreamName:        "cfm-activity",
		ActivityType:      "test.retry.persistence",
		ActivityProcessor: processor,
		Monitor:           monitor.NoopMonitor{},
	}

	err = executor.Execute(ctx)
	require.NoError(t, err)

	wg.Add(1)
	wg.Wait()

	// Verify the processor was called twice
	assert.Equal(t, int32(2), atomic.LoadInt32(&callCount), "Processor should have been called twice")

	// Verify values were persisted
	require.Eventually(t, func() bool {
		finalOrchestration, _, err := natsorchestration.ReadOrchestration(ctx, orchestration.ID, adapter)
		if err != nil || finalOrchestration.ProcessingData["retry_count"] == nil {
			return false
		}

		retryCount, ok := finalOrchestration.ProcessingData["retry_count"].(float64)
		if !ok || retryCount < 2 {
			return false
		}

		// Verify values from both attempts are present
		assert.Equal(t, float64(2), finalOrchestration.ProcessingData["retry_count"])
		assert.Equal(t, "initial_value", finalOrchestration.ProcessingData["first_attempt_data"])
		assert.Equal(t, "retry_value", finalOrchestration.ProcessingData["second_attempt_data"])
		return true
	}, persistenceTimeout, pollInterval, "Retry values should be persisted")
}

func Test_ValuePersistenceMultipleActivities(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), persistenceTimeout)
	defer cancel()

	nt, err := natsorchestration.SetupNatsContainer(ctx, "cfm-multi-activity-bucket")
	require.NoError(t, err)

	defer natsorchestration.TeardownNatsContainer(ctx, nt)

	stream := natsorchestration.SetupTestStream(t, ctx, nt.Client, streamName)
	natsorchestration.SetupTestConsumer(t, ctx, stream, "test.multi.persistence")

	var wg sync.WaitGroup
	wg.Add(2) // Two activities

	counter := &atomicCounter{}

	// Create processor that sets unique values per activity
	processor := &MultiActivityValueProcessor{
		onProcess: func(activityCtx api.ActivityContext) {
			defer wg.Done()
			activityID := activityCtx.ID()
			activityCtx.SetValue(fmt.Sprintf("%s_key", activityID), fmt.Sprintf("value_from_%s", activityID))
			activityCtx.SetValue("shared_counter", counter.IncrementAndGet())
		},
	}

	orchestration := api.Orchestration{
		ID:             "test-multi-activity-persistence",
		State:          api.OrchestrationStateRunning,
		Completed:      make(map[string]struct{}),
		ProcessingData: make(map[string]any),
		Steps: []api.OrchestrationStep{
			{
				Activities: []api.Activity{
					{ID: "A1", Type: "test.multi.persistence"},
					{ID: "A2", Type: "test.multi.persistence"},
				},
			},
		},
	}

	adapter := natsorchestration.NatsClientAdapter{Client: nt.Client}

	orchestrator := &natsorchestration.NatsDeploymentOrchestrator{
		Client:  adapter,
		Monitor: monitor.NoopMonitor{},
	}

	err = orchestrator.ExecuteOrchestration(ctx, &orchestration)
	require.NoError(t, err)

	// Create multiple executors
	for i := 0; i < 2; i++ {
		executor := &natsorchestration.NatsActivityExecutor{
			Client:            adapter,
			StreamName:        "cfm-activity",
			ActivityType:      "test.multi.persistence",
			ActivityProcessor: processor,
			Monitor:           monitor.NoopMonitor{},
		}
		err = executor.Execute(ctx)
		require.NoError(t, err)
	}

	wg.Wait()

	// Verify values were persisted
	require.Eventually(t, func() bool {
		finalOrchestration, _, err := natsorchestration.ReadOrchestration(ctx, orchestration.ID, adapter)
		if err != nil {
			return false
		}

		if finalOrchestration.ProcessingData["A1_key"] == nil || finalOrchestration.ProcessingData["A2_key"] == nil {
			return false
		}

		// Verify values from both activities are present
		assert.Equal(t, "value_from_A1", finalOrchestration.ProcessingData["A1_key"])
		assert.Equal(t, "value_from_A2", finalOrchestration.ProcessingData["A2_key"])

		// Verify shared counter was handled properly
		_, exists := finalOrchestration.ProcessingData["shared_counter"]
		assert.True(t, exists, "shared_counter should exist")
		return true
	}, persistenceTimeout, pollInterval, "Multi-activity values should be persisted")
}

func Test_ValuePersistenceOnWait(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), persistenceTimeout)
	defer cancel()

	nt, err := natsorchestration.SetupNatsContainer(ctx, "cfm-wait-activity-bucket")
	require.NoError(t, err)

	defer natsorchestration.TeardownNatsContainer(ctx, nt)

	stream := natsorchestration.SetupTestStream(t, ctx, nt.Client, streamName)
	natsorchestration.SetupTestConsumer(t, ctx, stream, "test.wait.persistence")

	var wg sync.WaitGroup
	wg.Add(1)

	// Create processor that sets values and returns wait
	processor := &WaitWithValueProcessor{
		onProcess: func(activityCtx api.ActivityContext) {
			defer wg.Done()
			activityCtx.SetValue("wait_state", "waiting")
			activityCtx.SetValue("wait_timestamp", time.Now().Unix())
		},
	}

	orchestration := createTestOrchestration("test-wait-persistence", "test.wait.persistence")
	adapter := natsorchestration.NatsClientAdapter{Client: nt.Client}

	orchestrator := &natsorchestration.NatsDeploymentOrchestrator{
		Client:  adapter,
		Monitor: monitor.NoopMonitor{},
	}

	err = orchestrator.ExecuteOrchestration(ctx, &orchestration)
	require.NoError(t, err)

	executor := &natsorchestration.NatsActivityExecutor{
		Client:            adapter,
		StreamName:        "cfm-activity",
		ActivityType:      "test.wait.persistence",
		ActivityProcessor: processor,
		Monitor:           monitor.NoopMonitor{},
	}

	err = executor.Execute(ctx)
	require.NoError(t, err)

	wg.Wait()

	// Verify values were persisted
	require.Eventually(t, func() bool {
		waitOrchestration, _, err := natsorchestration.ReadOrchestration(ctx, orchestration.ID, adapter)
		if err != nil || waitOrchestration.ProcessingData["wait_state"] == nil {
			return false
		}

		// Verify values were persisted during wait
		assert.Equal(t, "waiting", waitOrchestration.ProcessingData["wait_state"])
		assert.NotNil(t, waitOrchestration.ProcessingData["wait_timestamp"])
		return true
	}, persistenceTimeout, pollInterval, "Wait values should be persisted")
}

func TestNatsDeploymentOrchestrator_GetOrchestration_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), persistenceTimeout)
	defer cancel()

	nt, err := natsorchestration.SetupNatsContainer(ctx, "cfm-get-orchestration-bucket")
	require.NoError(t, err)
	defer natsorchestration.TeardownNatsContainer(ctx, nt)

	stream := natsorchestration.SetupTestStream(t, ctx, nt.Client, streamName)
	natsorchestration.SetupTestConsumer(t, ctx, stream, "test.activity")

	adapter := natsorchestration.NatsClientAdapter{Client: nt.Client}
	orchestrator := &natsorchestration.NatsDeploymentOrchestrator{
		Client:  adapter,
		Monitor: monitor.NoopMonitor{},
	}

	// Create and execute an orchestration
	orchestration := &api.Orchestration{
		ID:             "test-get-orchestration-success",
		State:          api.OrchestrationStateRunning,
		Completed:      make(map[string]struct{}),
		ProcessingData: make(map[string]any),
		Steps: []api.OrchestrationStep{
			{
				Activities: []api.Activity{
					{ID: "A1", Type: "test.activity"},
				},
			},
		},
	}

	err = orchestrator.ExecuteOrchestration(ctx, orchestration)
	require.NoError(t, err)

	// Test GetOrchestration
	result, err := orchestrator.GetOrchestration(ctx, orchestration.ID)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, orchestration.ID, result.ID)
}

func TestNatsDeploymentOrchestrator_GetOrchestration_NotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), persistenceTimeout)
	defer cancel()

	nt, err := natsorchestration.SetupNatsContainer(ctx, "cfm-get-orchestration-notfound-bucket")
	require.NoError(t, err)
	defer natsorchestration.TeardownNatsContainer(ctx, nt)

	adapter := natsorchestration.NatsClientAdapter{Client: nt.Client}
	orchestrator := &natsorchestration.NatsDeploymentOrchestrator{
		Client:  adapter,
		Monitor: monitor.NoopMonitor{},
	}

	// Test GetOrchestration for non-existent orchestration
	result, err := orchestrator.GetOrchestration(ctx, "non-existent-orchestration")
	require.NoError(t, err)
	assert.Nil(t, result)
}

// Helper function to create test orchestration
func createTestOrchestration(id, activityType string) api.Orchestration {
	return api.Orchestration{
		ID:             id,
		State:          api.OrchestrationStateRunning,
		Completed:      make(map[string]struct{}),
		ProcessingData: make(map[string]any),
		Steps: []api.OrchestrationStep{
			{
				Activities: []api.Activity{
					{ID: "A1", Type: activityType},
				},
			},
		},
	}
}

// Test processors

type ValueSettingProcessor struct {
	onProcess func(api.ActivityContext)
}

func (p *ValueSettingProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	if p.onProcess != nil {
		p.onProcess(ctx)
	}
	return api.ActivityResult{Result: api.ActivityResultComplete}
}

type RetryWithValueProcessor struct {
	onProcess func(api.ActivityContext) api.ActivityResult
}

func (p *RetryWithValueProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	if p.onProcess != nil {
		return p.onProcess(ctx)
	}
	return api.ActivityResult{Result: api.ActivityResultComplete}
}

type MultiActivityValueProcessor struct {
	onProcess func(api.ActivityContext)
}

func (p *MultiActivityValueProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	if p.onProcess != nil {
		p.onProcess(ctx)
	}
	return api.ActivityResult{Result: api.ActivityResultComplete}
}

type WaitWithValueProcessor struct {
	onProcess func(api.ActivityContext)
}

func (p *WaitWithValueProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	if p.onProcess != nil {
		p.onProcess(ctx)
	}
	return api.ActivityResult{Result: api.ActivityResultWait}
}

// Thread-safe atomic counter
type atomicCounter struct {
	count int64
}

func (c *atomicCounter) IncrementAndGet() int {
	return int(atomic.AddInt64(&c.count, 1))
}
