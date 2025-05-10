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

package durableprovision

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/microsoft/durabletask-go/backend"
	"github.com/microsoft/durabletask-go/backend/sqlite"
	"github.com/microsoft/durabletask-go/task"
	"log"
)

const (
	ConfigKeyBackend string = "provisionmanager.backend"
	inMemoryDb       string = ""
)

type DurableProvisionManagerServiceAssembly struct {
}

func (d *DurableProvisionManagerServiceAssembly) Name() string {
	return "Durable Provision Manager"
}

func (d *DurableProvisionManagerServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{}
}

func (d *DurableProvisionManagerServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{}
}

func (d *DurableProvisionManagerServiceAssembly) Init(ctx *system.InitContext) error {
	main(ctx.LogMonitor)
	return nil
}

func (d *DurableProvisionManagerServiceAssembly) Prepare() error {
	return nil
}

func (d *DurableProvisionManagerServiceAssembly) Start() error {
	return nil
}

func (d *DurableProvisionManagerServiceAssembly) Finalize() error {
	return nil
}

func (d *DurableProvisionManagerServiceAssembly) Shutdown() error {
	return nil
}

type logWrapper struct {
	LogMonitor monitor.LogMonitor
}

func (l logWrapper) Debug(v ...any) {
	l.LogMonitor.Debugw("", v)
}

func (l logWrapper) Debugf(format string, v ...any) {
	l.LogMonitor.Debugf(format, v)
}

func (l logWrapper) Info(v ...any) {
	// downgrade log verbosity
	l.LogMonitor.Debugw("", v...)
}

func (l logWrapper) Infof(format string, v ...any) {
	// downgrade log verbosity
	l.LogMonitor.Debugf(format, v...)
}

func (l logWrapper) Warn(v ...any) {
	l.LogMonitor.Warnw("", v...)
}

func (l logWrapper) Warnf(format string, v ...any) {
	l.LogMonitor.Warnf(format, v...)
}

func (l logWrapper) Error(v ...any) {
	l.LogMonitor.Severew("", v...)
}

func (l logWrapper) Errorf(format string, v ...any) {
	l.LogMonitor.Severef(format, v...)
}

func main(logMonitor monitor.LogMonitor) {
	// Create a new task registry and add the orchestrator and activities
	registry := task.NewTaskRegistry()
	registry.AddOrchestrator(ActivitySequenceOrchestrator)
	registry.AddActivity(SayHelloActivity)

	// Init the client
	ctx := context.Background()
	client, worker, err := Init(ctx, registry, logMonitor)
	if err != nil {
		log.Fatalf("Failed to initialize the client: %v", err)
	}
	fmt.Printf("", worker)
	// XCV defer worker.Shutdown(ctx)

	id, err := client.ScheduleNewOrchestration(ctx, ActivitySequenceOrchestrator)
	if err != nil {
		log.Fatalf("Failed to schedule new orchestration: %v", err)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				log.Printf("Waiting for orchestration to complete...")
				// Wait for the orchestration to complete
				metadata, err := client.WaitForOrchestrationCompletion(ctx, id)
				if err != nil {
					log.Fatalf("Failed to wait for orchestration to complete: %v", err)
				}

				// Print the results
				metadataEnc, err := json.MarshalIndent(metadata, "", "  ")
				if err != nil {
					log.Fatalf("Failed to encode result to JSON: %v", err)
				}
				log.Printf("Orchestration completed: %v", string(metadataEnc))
				return
			}

		}
	}()
	//// Wait for the orchestration to complete
	//metadata, err := client.WaitForOrchestrationCompletion(ctx, id)
	//if err != nil {
	//	log.Fatalf("Failed to wait for orchestration to complete: %v", err)
	//}
	//
	//// Print the results
	//metadataEnc, err := json.MarshalIndent(metadata, "", "  ")
	//if err != nil {
	//	log.Fatalf("Failed to encode result to JSON: %v", err)
	//}
	//log.Printf("Orchestration completed: %v", string(metadataEnc))
}

// Init creates and initializes an in-memory client and worker pair.
func Init(ctx context.Context, r *task.TaskRegistry, logMonitor monitor.LogMonitor) (backend.TaskHubClient, backend.TaskHubWorker, error) {
	logger := logWrapper{logMonitor}

	executor := task.NewTaskExecutor(r)

	be := sqlite.NewSqliteBackend(sqlite.NewSqliteOptions(inMemoryDb), logger)
	orchestrationWorker := backend.NewOrchestrationWorker(be, executor, logger)
	activityWorker := backend.NewActivityTaskWorker(be, executor, logger)
	taskHubWorker := backend.NewTaskHubWorker(be, orchestrationWorker, activityWorker, logger)

	// Start the worker
	err := taskHubWorker.Start(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Get the client to the backend
	taskHubClient := backend.NewTaskHubClient(be)

	return taskHubClient, taskHubWorker, nil

}

// ActivitySequenceOrchestrator makes three activity calls in sequence and results the results
// as an array.
func ActivitySequenceOrchestrator(ctx *task.OrchestrationContext) (any, error) {
	var helloTokyo string
	if err := ctx.CallActivity(SayHelloActivity, task.WithActivityInput("Tokyo")).Await(&helloTokyo); err != nil {
		return nil, err
	}
	fmt.Println("\nInvoking London XXXXXXXXXX" + helloTokyo)
	var helloLondon string
	if err := ctx.CallActivity(SayHelloActivity, task.WithActivityInput("London")).Await(&helloLondon); err != nil {
		return nil, err
	}
	var helloSeattle string
	fmt.Println("\nInvoking Seattle XXXXXXXXXX" + helloSeattle)
	if err := ctx.CallActivity(SayHelloActivity, task.WithActivityInput("Seattle")).Await(&helloSeattle); err != nil {
		return nil, err
	}
	return []string{helloTokyo, helloLondon, helloSeattle}, nil
}

// SayHelloActivity can be called by an orchestrator function and will return a friendly greeting.
func SayHelloActivity(ctx task.ActivityContext) (any, error) {
	var input string
	if err := ctx.GetInput(&input); err != nil {
		return "", err
	}
	fmt.Println("Hello from activity")
	return fmt.Sprintf("Hello, %s!", input), nil
}
