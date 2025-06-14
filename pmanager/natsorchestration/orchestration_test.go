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

package natsorchestration

import (
	"context"
	"fmt"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

const (
	parallelActivityWindow = 50 * time.Millisecond
	processTimeout         = 5 * time.Second
)

func TestExecuteOrchestration_NoSteps(t *testing.T) {
	orchestration := api.Orchestration{
		ID:        "test",
		Steps:     []api.OrchestrationStep{},
		Completed: make(map[string]struct{}),
	}

	ctx, cancel := context.WithTimeout(context.Background(), processTimeout)
	defer cancel()

	nt, err := setupNatsContainer(ctx, "cfm-durable-activity-bucket")
	require.NoError(t, err)
	//goland:noinspection GoUnhandledErrorResult
	defer teardownNatsContainer(ctx, nt)

	adapter := natsClientAdapter{client: nt.client}
	orchestrator := NatsDeploymentOrchestrator{client: adapter}
	err = orchestrator.ExecuteOrchestration(ctx, orchestration)
	require.Error(t, err)

}

func TestExecuteOrchestration(t *testing.T) {
	type activityExecution struct {
		id        string
		startTime time.Time
		endTime   time.Time
	}

	tests := []struct {
		name          string
		orchestration api.Orchestration
		validateFn    func(t *testing.T, executions []activityExecution)
	}{
		{
			name: "4 parallel activities in one step",
			orchestration: api.Orchestration{
				ID: "O1",
				Steps: []api.OrchestrationStep{
					{
						Activities: []api.Activity{
							{ID: "A1", Type: "test.activity"},
							{ID: "A2", Type: "test.activity"},
							{ID: "A3", Type: "test.activity"},
							{ID: "A4", Type: "test.activity"},
						},
					},
				},
				Completed: make(map[string]struct{}),
			},
			validateFn: func(t *testing.T, executions []activityExecution) {
				require.Len(t, executions, 4, "Should have 4 activities")

				// Verify all activities started within a small time window (50ms)
				var startTimes []time.Time
				for _, e := range executions {
					startTimes = append(startTimes, e.startTime)
				}

				firstStart := startTimes[0]
				for _, start := range startTimes[1:] {
					timeDiff := start.Sub(firstStart)
					assert.Less(t, timeDiff, parallelActivityWindow, "Parallel activities should start almost simultaneously")
				}

				// Verify all activities completed
				expectedIDs := map[string]bool{"A1": false, "A2": false, "A3": false, "A4": false}
				for _, e := range executions {
					expectedIDs[e.id] = true
				}
				for id, completed := range expectedIDs {
					assert.True(t, completed, "Activity %s should have completed", id)
				}
			},
		},
		{
			name: "2 steps with 2 parallel activities each",
			orchestration: api.Orchestration{
				ID: "O2",
				Steps: []api.OrchestrationStep{
					{
						Activities: []api.Activity{
							{ID: "A1", Type: "test.activity"},
							{ID: "A2", Type: "test.activity"},
						},
					},
					{
						Activities: []api.Activity{
							{ID: "A3", Type: "test.activity"},
							{ID: "A4", Type: "test.activity"},
						},
					},
				},
				Completed: make(map[string]struct{}),
			},
			validateFn: func(t *testing.T, executions []activityExecution) {
				require.Len(t, executions, 4, "Should have 4 activities")

				// Group activities by step
				step1Acts := make([]activityExecution, 0)
				step2Acts := make([]activityExecution, 0)

				for _, e := range executions {
					if e.id == "A1" || e.id == "A2" {
						step1Acts = append(step1Acts, e)
					} else if e.id == "A3" || e.id == "A4" {
						step2Acts = append(step2Acts, e)
					}
				}

				require.Len(t, step1Acts, 2, "Should have 2 activities in step 1")
				require.Len(t, step2Acts, 2, "Should have 2 activities in step 2")

				// Verify step 1 activities started in parallel
				timeDiff := step1Acts[1].startTime.Sub(step1Acts[0].startTime)
				if timeDiff < 0 {
					timeDiff = -timeDiff
				}
				assert.Less(t, timeDiff, parallelActivityWindow, "Step 1 activities should start almost simultaneously")

				// Verify step 2 activities started in parallel
				timeDiff = step2Acts[1].startTime.Sub(step2Acts[0].startTime)
				if timeDiff < 0 {
					timeDiff = -timeDiff
				}
				assert.Less(t, timeDiff, parallelActivityWindow, "Step 2 activities should start almost simultaneously")

				// Find when step 1 completed (latest end time)
				step1EndTime := step1Acts[0].endTime
				if step1Acts[1].endTime.After(step1EndTime) {
					step1EndTime = step1Acts[1].endTime
				}

				// Verify step 2 started after step 1 completed
				for _, step2Act := range step2Acts {
					assert.True(t, step2Act.startTime.After(step1EndTime) || step2Act.startTime.Equal(step1EndTime),
						"Step 2 activity %s should start after step 1 completes", step2Act.id)
				}

				// Verify all expected activities completed
				expectedIDs := map[string]bool{"A1": false, "A2": false, "A3": false, "A4": false}
				for _, e := range executions {
					expectedIDs[e.id] = true
				}
				for id, completed := range expectedIDs {
					assert.True(t, completed, "Activity %s should have completed", id)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), processTimeout)
			defer cancel()

			nt, err := setupNatsContainer(ctx, "cfm-durable-activity-bucket")
			require.NoError(t, err)
			defer teardownNatsContainer(ctx, nt)

			setupConsumer(t, ctx, nt)

			var executions []activityExecution
			executionsMutex := &sync.Mutex{}

			// Count total expected activities
			expectedActivities := 0
			for _, step := range tt.orchestration.Steps {
				expectedActivities += len(step.Activities)
			}

			// Setup synchronization
			resultCh := make(chan struct{})
			var wg sync.WaitGroup
			wg.Add(expectedActivities)

			processor := TestActivityProcessor{
				onProcess: func(id string) {
					execution := activityExecution{
						id:        id,
						startTime: time.Now(),
					}
					time.Sleep(10 * time.Millisecond) // Simulate work
					execution.endTime = time.Now()

					executionsMutex.Lock()
					executions = append(executions, execution)
					executionsMutex.Unlock()

					wg.Done()
				},
			}

			// Start a goroutine to wait for all activities
			go func() {
				wg.Wait()
				close(resultCh)
			}()

			noOpMonitor := monitor.NoopMonitor{}
			adapter := natsClientAdapter{client: nt.client}

			// Create executors
			executors := make([]NatsActivityExecutor, 4)
			for i := range executors {
				executors[i] = NatsActivityExecutor{
					id:                fmt.Sprintf("Executor%d", i+1),
					client:            adapter,
					activityName:      "test.activity",
					activityProcessor: processor,
					monitor:           noOpMonitor,
				}
				executors[i].Execute(ctx)
			}

			orchestrator := NatsDeploymentOrchestrator{client: adapter}
			err = orchestrator.ExecuteOrchestration(ctx, tt.orchestration)
			require.NoError(t, err)

			// Wait for completion or timeout
			select {
			case <-resultCh:
				// All activities completed successfully
			case <-ctx.Done():
				t.Fatalf("Test timed out waiting for activities to complete: %v", ctx.Err())
			}

			// Run validation
			tt.validateFn(t, executions)
		})
	}
}

func TestActivityProcessor_ScheduleThenContinue(t *testing.T) {
	// Setup NATS test environment
	ctx, cancel := context.WithTimeout(context.Background(), processTimeout)
	defer cancel()

	nt, err := setupNatsContainer(ctx, "cfm-durable-activity-bucket")
	require.NoError(t, err)
	//goland:noinspection ALL
	defer teardownNatsContainer(ctx, nt)

	setupConsumer(t, ctx, nt)

	// Create a processor that returns schedule first, then continue
	var wg sync.WaitGroup
	processor := &ScheduleThenContinueProcessor{
		callCount: 0,
		wg:        &wg,
	}

	// Create orchestration with single activity
	orchestration := api.Orchestration{
		ID:        "test-schedule-continue-1",
		State:     api.OrchestrationStateRunning,
		Completed: make(map[string]struct{}),
		Steps: []api.OrchestrationStep{
			{
				Activities: []api.Activity{
					{ID: "A1", Type: "test.schedule.continue"},
				},
			},
		},
	}

	adapter := natsClientAdapter{client: nt.client}

	// Create and start the orchestrator
	orchestrator := &NatsDeploymentOrchestrator{
		client:  adapter,
		monitor: monitor.NoopMonitor{},
	}

	err = orchestrator.ExecuteOrchestration(context.Background(), orchestration)
	require.NoError(t, err)

	// Create an activity executor with our test processor
	executor := &NatsActivityExecutor{
		id:                "test-executor-schedule-continue",
		client:            adapter,
		activityName:      "test.schedule.continue",
		activityProcessor: processor,
		monitor:           monitor.NoopMonitor{},
	}

	// Start executor
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = executor.Execute(ctx)
	require.NoError(t, err)

	// Wait for both calls to complete
	wg.Add(2) // Expecting 2 calls: schedule then continue
	wg.Wait()

	// Verify processor was called twice
	assert.Equal(t, 2, processor.callCount, "Processor should have been called twice")
}

// ScheduleThenContinueProcessor implements ActivityProcessor
// Returns ActivityResultSchedule on first call, ActivityResultContinue on subsequent calls
type ScheduleThenContinueProcessor struct {
	callCount int
	wg        *sync.WaitGroup
}

func (p *ScheduleThenContinueProcessor) Process(_ api.ActivityContext) api.ActivityResult {
	p.callCount++
	defer p.wg.Done() // Signal completion of this call

	if p.callCount == 1 {
		// First call: return schedule result with 1 second delay for faster testing
		return api.ActivityResult{
			Result:     api.ActivityResultSchedule,
			WaitMillis: 100 * time.Millisecond,
			Error:      nil,
		}
	}

	// Subsequent calls: return continue result
	return api.ActivityResult{
		Result:     api.ActivityResultContinue,
		WaitMillis: 0,
		Error:      nil,
	}
}

func setupConsumer(t *testing.T, ctx context.Context, nt *natsTestContainer) {
	cfg := jetstream.StreamConfig{
		Name:      "cfm-activity",
		Retention: jetstream.WorkQueuePolicy,
		Subjects:  []string{"event.*"},
	}

	stream, err := nt.client.jetStream.CreateOrUpdateStream(ctx, cfg)
	require.NoError(t, err)

	_, err = stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "cfm-durable-activity",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	require.NoError(t, err)
}

// TestActivityProcessor with timing information
type TestActivityProcessor struct {
	onProcess func(id string)
}

func (t TestActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	ctx.Value("key")

	if t.onProcess != nil {
		t.onProcess(ctx.ID())
	}
	return api.ActivityResult{Result: api.ActivityResultContinue}
}
