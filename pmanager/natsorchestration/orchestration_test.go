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
						Parallel: true,
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
			name: "4 sequential activities in one step",
			orchestration: api.Orchestration{
				ID: "O2",
				Steps: []api.OrchestrationStep{
					{
						Parallel: false,
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

				// Verify activities executed in sequence
				expectedOrder := []string{"A1", "A2", "A3", "A4"}
				for i := 0; i < len(executions)-1; i++ {
					assert.Equal(t, expectedOrder[i], executions[i].id, "Activities should execute in order")
					assert.True(t, executions[i+1].startTime.After(executions[i].endTime),
						"Activity %s should start after activity %s ends", executions[i+1].id, executions[i].id)
				}
			},
		},
		{
			name: "2 parallel followed by 2 sequential activities",
			orchestration: api.Orchestration{
				ID: "O3",
				Steps: []api.OrchestrationStep{
					{
						Parallel: true,
						Activities: []api.Activity{
							{ID: "A1", Type: "test.activity"},
							{ID: "A2", Type: "test.activity"},
						},
					},
					{
						Parallel: false,
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

				// Find parallel activities (A1 and A2)
				var parallelActs, sequentialActs []activityExecution
				for _, e := range executions {
					if e.id == "A1" || e.id == "A2" {
						parallelActs = append(parallelActs, e)
					} else {
						sequentialActs = append(sequentialActs, e)
					}
				}

				// Verify parallel activities started together
				require.Len(t, parallelActs, 2, "Should have 2 parallel activities")
				timeDiff := parallelActs[1].startTime.Sub(parallelActs[0].startTime)
				assert.Less(t, timeDiff, parallelActivityWindow, "Parallel activities should start almost simultaneously")

				// Find the last completed parallel activity
				lastParallelEnd := parallelActs[0].endTime
				if parallelActs[1].endTime.After(lastParallelEnd) {
					lastParallelEnd = parallelActs[1].endTime
				}

				// Verify sequential activities started after parallel ones and in order
				require.Len(t, sequentialActs, 2, "Should have 2 sequential activities")
				assert.True(t, sequentialActs[0].startTime.After(lastParallelEnd),
					"First sequential activity should start after parallel activities complete")
				assert.True(t, sequentialActs[1].startTime.After(sequentialActs[0].endTime),
					"Second sequential activity should start after first sequential activity completes")
			},
		},
		{
			name: "No steps",
			orchestration: api.Orchestration{
				ID:        "O4",
				Steps:     []api.OrchestrationStep{},
				Completed: make(map[string]struct{}),
			},
			validateFn: func(t *testing.T, executions []activityExecution) {
				assert.Len(t, executions, 0, "Should have no executions")
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
					time.Sleep(100 * time.Millisecond) // Simulate work
					execution.endTime = time.Now()

					executionsMutex.Lock()
					executions = append(executions, execution)
					executionsMutex.Unlock()

					wg.Done()
				},
			}

			// Start goroutine to wait for all activities
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

// TestActivityProcessor with timing information
type TestActivityProcessor struct {
	onProcess func(id string)
}

func (t TestActivityProcessor) Process(ctx api.ActivityContext) (bool, error) {
	if t.onProcess != nil {
		t.onProcess(ctx.ID())
	}
	return true, nil
}
