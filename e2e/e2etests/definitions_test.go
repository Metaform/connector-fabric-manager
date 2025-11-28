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

package e2etests

import (
	"context"
	"fmt"
	"testing"

	"github.com/metaform/connector-fabric-manager/common/natsfixtures"
	"github.com/metaform/connector-fabric-manager/e2e/e2efixtures"
	"github.com/metaform/connector-fabric-manager/pmanager/model/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_VerifyDefinitionOperations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	nt, err := natsfixtures.SetupNatsContainer(ctx, cfmBucket)
	require.NoError(t, err)

	defer natsfixtures.TeardownNatsContainer(ctx, nt)
	defer cleanup()

	client := launchPlatform(t, nt)

	waitTManager(t, client)
	waitPManager(t, client)

	var activityDefinitions []v1alpha1.ActivityDefinition

	err = e2efixtures.CreateTestActivityDefinition(client)
	require.NoError(t, err)

	err = client.GetPManager("activity-definitions", &activityDefinitions)
	require.NoError(t, err)
	assert.Equal(t, 1, len(activityDefinitions))

	err = e2efixtures.CreateTestOrchestrationDefinitions(client)
	require.NoError(t, err)

	var orchestrationDefinitions []v1alpha1.OrchestrationDefinition
	err = client.GetPManager("orchestration-definitions", &orchestrationDefinitions)
	require.NoError(t, err)
	assert.Equal(t, 2, len(orchestrationDefinitions))

	for _, definition := range orchestrationDefinitions {
		err = client.DeleteToPManager(fmt.Sprintf("orchestration-definitions/%s", definition.Type))
		require.NoError(t, err)
	}

	orchestrationDefinitions = nil
	err = client.GetPManager("orchestration-definitions", &orchestrationDefinitions)
	require.NoError(t, err)
	assert.Empty(t, len(orchestrationDefinitions))

	for _, definition := range activityDefinitions {
		err = client.DeleteToPManager(fmt.Sprintf("activity-definitions/%s", definition.Type))
		require.NoError(t, err)
	}
	err = client.GetPManager("activity-definitions", &activityDefinitions)
	require.NoError(t, err)
	assert.Empty(t, len(activityDefinitions))
}
