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

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

// DeploymentManifest represents the configuration details for a system deployment. An Orchestration is instantiated
// from the manifest and executed.
//
// The manifest includes a unique identifier, the type of deployment specified by a DeploymentDefinition, and a payload
// of deployment-specific data, which will be passed as input to the Orchestration.
type DeploymentManifest struct {
	ID             string         `json:"id"`
	DeploymentType string         `json:"deploymentType"`
	Payload        map[string]any `json:"payload"`
}

type OrchestrationState uint

const (
	OrchestrationStateStateInitialized OrchestrationState = 0
	OrchestrationStateStateRunning     OrchestrationState = 1
	OrchestrationStateStateCompleted   OrchestrationState = 2
	OrchestrationStateStateErrored     OrchestrationState = 3
)

// Orchestration is a collection of activities that are executed to effect a deployment.
//
// The DeploymentID is a reference to the original DeploymentManifest. As actions are completed, the orchestration
// system will update the Completed map.
type Orchestration struct {
	ID             string             `json:"id"`
	DeploymentID   string             `json:"deploymentId"`
	State          OrchestrationState `json:"state"`
	Steps          []OrchestrationStep
	Inputs         map[string]any
	ProcessingData map[string]any
	Completed      map[string]struct{}
}

// CanProceedToNextActivity returns if the orchestration is able to proceed to the next activity or must wait.
func (o *Orchestration) CanProceedToNextActivity(activityId string, validator func([]string) bool) (bool, error) {
	step, err := o.GetStepForActivity(activityId)
	if err != nil {
		return true, err
	}
	if !step.Parallel {
		return true, nil
	}
	activityIds := make([]string, 0, len(step.Activities))
	for _, activity := range step.Activities {
		activityIds = append(activityIds, activity.ID)
	}
	return validator(activityIds), nil
}

// GetStepForActivity retrieves the orchestration step containing the specified activity ID. Returns an error if not found.
func (o *Orchestration) GetStepForActivity(activityId string) (*OrchestrationStep, error) {
	for _, step := range o.Steps {
		for _, activity := range step.Activities {
			if activity.ID == activityId {
				return &step, nil
			}
		}
	}
	return nil, errors.New("step not found for activity: " + activityId)
}

func (o *Orchestration) GetNextActivities(current string) ([]Activity, bool) {
	reachedCurrent := false
	for _, step := range o.Steps {
		if reachedCurrent {
			if step.Parallel {
				return step.Activities[0 : len(step.Activities)-1], true
			}
			if len(step.Activities) == 0 {
				return []Activity{}, false
			}
			return step.Activities[0:1], false
		}

		for i, activity := range step.Activities {
			if activity.ID == current {
				reachedCurrent = true
				if (i + 1) < len(step.Activities) {
					if step.Parallel {
						continue
					}
					return step.Activities[i+1 : i+2], false
				}
			}
		}
	}
	return []Activity{}, false
}

type OrchestrationStep struct {
	Parallel   bool       `json:"parallel"`
	Activities []Activity `json:"activities"`
}

type Activity struct {
	ID     string         `json:"id"`
	Type   string         `json:"type"`
	Inputs []MappingEntry `json:"inputs"`
}

// ActivityMessage used to enqueue an activity for processing.
type ActivityMessage struct {
	OrchestrationID string   `json:"orchestrationID"`
	Activity        Activity `json:"activity"`
	Parallel        bool     `json:"parallel"`
}

type MappingEntry struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// UnmarshalJSON handles deserializing a MappingEntry from a string to a source/target pair.
func (m *MappingEntry) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		// If successful, use the string as both source and target
		m.Source = s
		m.Target = s
		return nil
	}

	// If string unmarshal fails, try as an object
	var objMap struct {
		Source string `json:"source"`
		Target string `json:"target"`
	}
	if err := json.Unmarshal(data, &objMap); err != nil {
		return fmt.Errorf("failed to unmarshal MappingEntry: %w", err)
	}

	m.Source = objMap.Source
	m.Target = objMap.Target
	return nil
}

type DeploymentDefinition struct {
	Type       string    `json:"type"`
	ApiVersion string    `json:"apiVersion"`
	Resource   Resource  `json:"resource"`
	Versions   []Version `json:"versions"`
}

type Resource struct {
	Group       string `json:"group"`
	Singular    string `json:"singular"`
	Plural      string `json:"plural"`
	Description string `json:"description"`
}

type Version struct {
	Version                 string                  `json:"version"`
	Active                  bool                    `json:"active"`
	Schema                  map[string]any          `json:"schema"`
	OrchestrationDefinition OrchestrationDefinition `json:"orchestration"`
}

type OrchestrationDefinition []OrchestrationStepDefinition

// OrchestrationStepDefinition represents a group of activities that can be executed in parallel or sequentially
type OrchestrationStepDefinition struct {
	Parallel   bool       `json:"parallel"`
	Activities []Activity `json:"activities"`
}

// ActivityDefinition represents a single activity in the orchestration
type ActivityDefinition struct {
	Type         string `json:"type"`
	Provider     string `json:"provider"`
	Description  string `json:"description"`
	InputSchema  string `json:"inputSchema"`
	OutputSchema string `json:"outputSchema"`
}

func ParseDeploymentDefinition(data []byte) (*DeploymentDefinition, error) {
	var definition DeploymentDefinition

	if err := json.Unmarshal(data, &definition); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &definition, nil
}

func InstantiateOrchestration(deploymentID string, definition OrchestrationDefinition, data map[string]any) *Orchestration {
	orchestration := &Orchestration{
		ID:             uuid.New().String(),
		DeploymentID:   deploymentID,
		State:          OrchestrationStateStateInitialized,
		Steps:          make([]OrchestrationStep, len(definition)),
		Inputs:         data,
		ProcessingData: make(map[string]any),
		Completed:      make(map[string]struct{}),
	}

	// Create steps
	for i, stepDef := range definition {
		step := OrchestrationStep{
			Parallel:   stepDef.Parallel,
			Activities: make([]Activity, len(stepDef.Activities)),
		}

		// Create activities
		for j, _ := range stepDef.Activities {
			step.Activities[j] = Activity{
				ID:     uuid.New().String(),
				Inputs: make([]MappingEntry, 0),
			}
		}

		orchestration.Steps[i] = step
	}

	return orchestration
}
