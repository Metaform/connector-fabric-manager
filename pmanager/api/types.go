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

type DeploymentManifest struct {
	Type    string                 `json:"type"`
	ID      string                 `json:"id"`
	Payload map[string]interface{} `json:"payload"`
}

type Orchestration struct {
	ID             string `json:"id"`
	Steps          []OrchestrationStep
	Data           map[string]any
	ProcessingData map[string]any
	Completed      map[string]struct{}
}

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
	ID         string     `json:"id"`
	Parallel   bool       `json:"parallel"`
	Activities []Activity `json:"activities"`
}

type Activity struct {
	ID string `json:"id"`
	ActivityDefinition
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
	OrchestrationDefinition OrchestrationDefinition `json:"orchestration"`
}

type OrchestrationDefinition []OrchestrationStepDefinition

// OrchestrationStepDefinition represents a group of activities that can be executed in parallel or sequentially
type OrchestrationStepDefinition struct {
	Parallel   bool                 `json:"parallel"`
	Activities []ActivityDefinition `json:"activities"`
}

// ActivityDefinition represents a single activity in the orchestration
type ActivityDefinition struct {
	Type        string   `json:"type"`
	DataMapping []string `json:"dataMapping"`
}

func ParseDeploymentDefinition(data []byte) (*DeploymentDefinition, error) {
	var definition DeploymentDefinition

	if err := json.Unmarshal(data, &definition); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &definition, nil
}

func InstantiateOrchestration(definition OrchestrationDefinition, data map[string]any) *Orchestration {
	orchestration := &Orchestration{
		ID:             uuid.New().String(),
		Steps:          make([]OrchestrationStep, len(definition)),
		Data:           data,
		ProcessingData: make(map[string]any),
		Completed:      make(map[string]struct{}),
	}

	// Create steps
	for i, stepDef := range definition {
		step := OrchestrationStep{
			ID:         uuid.New().String(),
			Parallel:   stepDef.Parallel,
			Activities: make([]Activity, len(stepDef.Activities)),
		}

		// Create activities
		for j, actDef := range stepDef.Activities {
			step.Activities[j] = Activity{
				ID:                 uuid.New().String(),
				ActivityDefinition: actDef,
			}
		}

		orchestration.Steps[i] = step
	}

	return orchestration
}
