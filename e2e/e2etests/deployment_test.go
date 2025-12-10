// Copyright (c) 2025 Metaform Systems, Inc
//
// This program and the accompanying materials are made available under the
// terms of the Apache License, Version 2.0 which is available at
// https://www.apache.org/licenses/LICENSE-2.0
//
// SPDX-License-Identifier: Apache-2.0
//
// Contributors:
//
//	Metaform Systems, Inc. - initial API and implementation
package e2etests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/natsfixtures"
	"github.com/metaform/connector-fabric-manager/common/sqlstore"
	"github.com/metaform/connector-fabric-manager/e2e/e2efixtures"
	papi "github.com/metaform/connector-fabric-manager/pmanager/api"
	"github.com/metaform/connector-fabric-manager/tmanager/api"
	"github.com/metaform/connector-fabric-manager/tmanager/model/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_VerifyE2E(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	nt, err := natsfixtures.SetupNatsContainer(ctx, cfmBucket)

	require.NoError(t, err)

	defer natsfixtures.TeardownNatsContainer(ctx, nt)
	defer cleanup()

	pg, dsn, err := sqlstore.SetupTestContainer(t)
	require.NoError(t, err)
	defer pg.Terminate(context.Background())

	client := launchPlatform(t, nt.URI, dsn)

	waitPManager(t, client)

	err = e2efixtures.CreateTestActivityDefinition(client)
	require.NoError(t, err)

	err = e2efixtures.CreateTestOrchestrationDefinitions(client)
	require.NoError(t, err)

	cell, err := e2efixtures.CreateCell(client)
	require.NoError(t, err)

	var cells []api.Cell
	err = client.GetTManager("cells", &cells)
	require.NoError(t, err)
	assert.Equal(t, 1, len(cells))

	dProfile, err := e2efixtures.CreateDataspaceProfile(client)
	require.NoError(t, err)

	var dprofiles []api.DataspaceProfile
	err = client.GetTManager("dataspace-profiles", &dprofiles)
	require.NoError(t, err)
	assert.Equal(t, 1, len(dprofiles))

	deployment := v1alpha1.NewDataspaceProfileDeployment{
		ProfileID: dProfile.ID,
		CellID:    cell.ID,
	}
	err = e2efixtures.DeployDataspaceProfile(deployment, client)
	require.NoError(t, err)

	var dprofile api.DataspaceProfile
	err = client.GetTManager(fmt.Sprintf("dataspace-profiles/%s", dprofiles[0].ID), &dprofile)
	require.NoError(t, err)
	require.NotNil(t, dprofile)

	tenant, err := e2efixtures.CreateTenant(client, map[string]any{})
	require.NoError(t, err)

	newProfile := v1alpha1.NewParticipantProfileDeployment{
		Identifier:    "did:web:foo.com",
		VPAProperties: map[string]map[string]any{string(model.ConnectorType): {"connectorkey": "connectorvalue"}},
	}
	var participantProfile v1alpha1.ParticipantProfile
	err = client.PostToTManagerWithResponse(fmt.Sprintf("tenants/%s/participant-profiles", tenant.ID), newProfile, &participantProfile)
	require.NoError(t, err)

	var statusProfile v1alpha1.ParticipantProfile

	// Verify all VPAs are active
	deployCount := 0
	for start := time.Now(); time.Since(start) < 5*time.Second; {
		err = client.GetTManager(fmt.Sprintf("tenants/%s/participant-profiles/%s", tenant.ID, participantProfile.ID), &statusProfile)
		require.NoError(t, err)
		for _, vpa := range statusProfile.VPAs {
			if vpa.State == api.DeploymentStateActive.String() {
				deployCount++
			}
		}
		if deployCount == 3 {
			break
		}
	}
	require.Equal(t, 3, deployCount, "Expected 3 deployments to be active")

	var profiles []api.ParticipantProfile
	err = client.PostToTManagerWithResponse("participant-profiles/query", model.Query{Predicate: "vpas.type='cfm.connector' AND vpas.properties.connectorkey='connectorvalue'"}, &profiles)
	require.NoError(t, err)
	assert.Equal(t, 1, len(profiles), "Expected 1 profile to be found")

	// Verify round-tripping of VPA properties - these are supplied during profile creation and are added to the VPA
	//
	// Check for VPA that contains a key with "cfm.connector" value and verify it has "connectorkey"
	var connectorVPA *v1alpha1.VirtualParticipantAgent
	for _, vpa := range statusProfile.VPAs {
		if vpa.Type == model.ConnectorType {
			connectorVPA = &vpa
			break
		}
	}

	require.NotNil(t, connectorVPA, "Expected to find a VPA with cfm.connector type")
	require.NotNil(t, connectorVPA.Properties, "Connector VPA properties should not be nil")
	require.Contains(t, connectorVPA.Properties, "connectorkey", "Connector VPA should contain 'connectorkey' property")

	var orchestrations []papi.OrchestrationEntry
	err = client.PostToPManagerWithResponse("orchestrations/query", model.Query{Predicate: fmt.Sprintf("correlationId = '%s'", statusProfile.ID)}, &orchestrations)
	require.NoError(t, err)
	require.Equal(t, 1, len(orchestrations), "Expected 1 orchestration to be created")
	assert.Equal(t, papi.OrchestrationStateCompleted, papi.OrchestrationState(orchestrations[0].State))

	// Dispose VPAs
	err = client.DeleteToTManager(fmt.Sprintf("tenants/%s/participant-profiles/%s", tenant.ID, participantProfile.ID))
	require.NoError(t, err)

	disposeCount := 0
	for start := time.Now(); time.Since(start) < 5*time.Second; {
		err = client.GetTManager(fmt.Sprintf("tenants/%s/participant-profiles/%s", tenant.ID, participantProfile.ID), &statusProfile)
		require.NoError(t, err)
		for _, vpa := range statusProfile.VPAs {
			if vpa.State == api.DeploymentStateDisposed.String() {
				disposeCount++
			}
		}
		if disposeCount == 3 {
			break
		}
	}
	require.Equal(t, 3, disposeCount, "Expected 3 deployments to be disposed")
}
