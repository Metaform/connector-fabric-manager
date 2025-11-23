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

package v1alpha1

import (
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

func ToAPIActivityDefinition(definition *ActivityDefinition) *api.ActivityDefinition {
	if definition == nil {
		return &api.ActivityDefinition{}
	}
	return &api.ActivityDefinition{
		Type:         api.ActivityType(definition.Type),
		Description:  definition.Description,
		InputSchema:  definition.InputSchema,
		OutputSchema: definition.OutputSchema,
	}
}

func ToAPIOrchestrationDefinition(definition *OrchestrationDefinition) *api.OrchestrationDefinition {
	if definition == nil {
		return &api.OrchestrationDefinition{}
	}
	apiActivities := make([]api.Activity, len(definition.Activities))
	for i, activity := range definition.Activities {
		apiActivities[i] = api.Activity{
			ID:            activity.ID,
			Type:          api.ActivityType(activity.Type),
			Discriminator: api.Discriminator(activity.Discriminator),
			Inputs:        ToAPIMappingEntries(activity.Inputs),
			DependsOn:     activity.DependsOn,
		}
	}

	return &api.OrchestrationDefinition{
		Type:        model.OrchestrationType(definition.Type),
		Description: definition.Description,
		Active:      true, // Default to active as the model doesn't have this field
		Schema:      definition.Schema,
		Activities:  apiActivities,
	}
}

func ToAPIMappingEntries(entries []MappingEntry) []api.MappingEntry {
	apiEntries := make([]api.MappingEntry, len(entries))
	for i, entry := range entries {
		apiEntries[i] = api.MappingEntry{
			Source: entry.Source,
			Target: entry.Target,
		}
	}
	return apiEntries
}

func ToOrchestrationEntry(entry *api.OrchestrationEntry) OrchestrationEntry {
	return OrchestrationEntry{
		ID:                entry.ID,
		CorrelationID:     entry.CorrelationID,
		State:             int(entry.State),
		StateTimestamp:    entry.StateTimestamp,
		CreatedTimestamp:  entry.CreatedTimestamp,
		OrchestrationType: entry.OrchestrationType,
	}
}

func ToOrchestration(orchestration *api.Orchestration) Orchestration {
	return Orchestration{
		ID:                orchestration.ID,
		CorrelationID:     orchestration.CorrelationID,
		State:             int(orchestration.State),
		StateTimestamp:    orchestration.StateTimestamp,
		CreatedTimestamp:  orchestration.CreatedTimestamp,
		OrchestrationType: orchestration.OrchestrationType,
		ProcessingData:    orchestration.ProcessingData,
		Steps:             toSteps(orchestration.Steps),
		OutputData:        orchestration.OutputData,
		Completed:         orchestration.Completed,
	}
}

func toSteps(steps []api.OrchestrationStep) []OrchestrationStep {
	result := make([]OrchestrationStep, len(steps))
	for i, step := range steps {
		result[i] = OrchestrationStep{
			Activities: toActivities(step.Activities),
		}
	}
	return result
}

func toActivities(activities []api.Activity) []Activity {
	result := make([]Activity, len(activities))
	for i, activity := range activities {
		result[i] = Activity{
			ID:            activity.ID,
			Type:          string(activity.Type),
			Discriminator: string(activity.Discriminator),
			Inputs:        toMappingEntries(activity.Inputs),
			DependsOn:     activity.DependsOn,
		}
	}
	return result
}

func toMappingEntries(entries []api.MappingEntry) []MappingEntry {
	result := make([]MappingEntry, len(entries))
	for i, entry := range entries {
		result[i] = MappingEntry{
			Source: entry.Source,
			Target: entry.Target,
		}
	}
	return result
}