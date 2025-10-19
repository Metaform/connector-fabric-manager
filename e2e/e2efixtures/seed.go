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

package e2efixtures

import (
	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

func CreateTestActivityDefinition(apiClient *ApiClient) error {
	requestBody := api.ActivityDefinition{
		Type:        "test.activity",
		Description: "Performs a test activity",
	}

	return apiClient.PostToPManager("activity-definition", requestBody)
}

func CreateTestDeploymentDefinition(apiClient *ApiClient) error {
	requestBody := api.DeploymentDefinition{
		Type:       model.VpaDeploymentType,
		ApiVersion: "v1",
		Resource: api.Resource{
			Group:       "deployments.example.com",
			Singular:    "TestDeployment",
			Plural:      "TestDeployments",
			Description: "Test deployment",
		},
		Versions: []api.Version{
			{
				Version: "1.0.0",
				Active:  true,
				Activities: []api.Activity{
					{
						ID:   "activity1",
						Type: "test-activity",
					},
				},
			},
		},
	}

	return apiClient.PostToPManager("deployment-definition", requestBody)
}
