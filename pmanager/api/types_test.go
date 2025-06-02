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
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
	"testing"
)

func TestParseDeploymentDefinition(t *testing.T) {
	result, err := ParseDeploymentDefinition([]byte(testDefinition))
	require.NoError(t, err)
	assert.Equal(t, result.Type, "tenant.example.com")
	assert.Equal(t, result.ApiVersion, "1.0")

	assert.Equal(t, result.Resource.Group, "example.com")
	assert.Equal(t, result.Resource.Singular, "tenant")
	assert.Equal(t, result.Resource.Plural, "tenants")
	assert.Equal(t, result.Resource.Description, "Deploys infrastructure and configuration required to support a tenant")

	assert.Equal(t, len(result.Versions), 1)
	assert.Equal(t, result.Versions[0].Version, "1.0.0")
	assert.Equal(t, result.Versions[0].Active, true)

	orchestrationDefinition := result.Versions[0].OrchestrationDefinition
	assert.Equal(t, len(orchestrationDefinition), 2)
	assert.Equal(t, orchestrationDefinition[0].Parallel, false)
	assert.Equal(t, len(orchestrationDefinition[0].Activities), 5)
	assert.Equal(t, orchestrationDefinition[1].Parallel, true)
	assert.Equal(t, len(orchestrationDefinition[1].Activities), 0)
}

const testDefinition = `{
  "type": "tenant.example.com",
  "apiVersion": "1.0",
  "resource": {
    "group": "example.com",
    "singular": "tenant",
    "plural": "tenants",
    "description": "Deploys infrastructure and configuration required to support a tenant"
  },
  "versions": [
    {
      "version": "1.0.0",
      "active": true,
      "schema": {
        "openAPIV3Schema": {
          "type": "object",
          "properties": {
            "cell": {
              "type": "string"
            },
            "did": {
              "type": "string",
              "format": "uri"
            },
            "baseUrl": {
              "type": "string",
              "format": "uri"
            },
            "dataspaces": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          },
          "required": [
            "cell",
            "did",
            "baseUrl",
            "dataspaces"
          ]
        }
      },
      "orchestration": [
        {
          "parallel": false,
          "activities": [
            {
              "type": "dns.example.com",
              "dataMapping": [
                "cell",
                "baseUrl"
              ]
            },
            {
              "type": "ihtenant.example.com",
              "dataMapping": [
                "cell",
                "did",
                "dataspaces"
              ]
            },
            {
              "type": "edctenant.example.com",
              "dataMapping": [
                "cell",
                "did",
                "dataspaces"
              ]
            },
            {
              "type": "ingress.example.com",
              "dataMapping": [
                "cell",
                "baseUrl",
                "did"
              ]
            },
            {
              "type": "onboard.example.com",
              "dataMapping": [
                "cell",
                "did",
                "dataspaces"
              ]
            }
          ]
        },
        {
          "parallel": true,
          "activities": []
        }
      ]
    }
  ]
}`
