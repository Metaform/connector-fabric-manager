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
	"testing"
	"time"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/natsfixtures"
	"github.com/metaform/connector-fabric-manager/e2e/e2efixtures"
	"github.com/metaform/connector-fabric-manager/tmanager/model/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_VerifyTenantQueries(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	nt, err := natsfixtures.SetupNatsContainer(ctx, cfmBucket)
	require.NoError(t, err)

	defer natsfixtures.TeardownNatsContainer(ctx, nt)
	defer cleanup()

	client := launchPlatform(t, nt)

	// Wait for the tmanager to be ready
	for start := time.Now(); time.Since(start) < 5*time.Second; {
		if _, err = e2efixtures.CreateTenant(client, map[string]any{"group": "suppliers"}); err == nil {
			break
		}
	}
	require.NoError(t, err)

	_, err = e2efixtures.CreateTenant(client, map[string]any{"group": "manufacturers"})
	require.NoError(t, err)

	var result []v1alpha1.Tenant
	err = client.PostToTManagerWithResponse("tenants/query", model.Query{Predicate: "properties.group = 'suppliers'"}, &result)
	require.NoError(t, err)
	assert.Equal(t, 1, len(result))

	err = client.PostToTManagerWithResponse("tenants/query", model.Query{Predicate: "properties.group = 'suppliers' OR properties.group = 'manufacturers'"}, &result)
	require.NoError(t, err)
	assert.Equal(t, 2, len(result))

}
