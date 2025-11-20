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

package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/pmanager/model/v1alpha1"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)

const docsDir = "docs"

func main() {
	r := spec.NewRouter(
		option.WithTitle("Provision Manager API"),
		option.WithVersion("0.0.1"),
		option.WithDescription("API for managing Orchestrations, Orchestration Definitions, and Activity Definitions"),
		option.WithServer("http://localhost:8080", option.ServerDescription("Development server")),
	)

	generateOrchestrationEndpoints(r)
	generateOrchestrationDefinitionEndpoints(r)
	generateActivityDefinitionEndpoints(r)

	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		if err := os.Mkdir(docsDir, 0755); err != nil {
			panic(err)
		}
	}

	if err := r.WriteSchemaTo(filepath.Join(docsDir, "openapi.json")); err != nil {
		panic(err)
	}
}

func generateOrchestrationEndpoints(r spec.Generator) {
	orchestration := r.Group("/api/v1alpha1/orchestration")

	orchestration.Post("",
		option.Summary("Execute an Orchestration"),
		option.Description("Execute an Orchestration"),
		option.Request(model.OrchestrationManifest{}),
		option.Response(http.StatusAccepted, nil),
	)
}

func generateActivityDefinitionEndpoints(r spec.Generator) {
	activity := r.Group("/api/v1alpha1/activity-definitions")

	activity.Post("",
		option.Summary("Create an Activity Definition"),
		option.Description("Create a new Activity Definition"),
		option.Request(v1alpha1.ActivityDefinition{}),
		option.Response(http.StatusCreated, nil),
	)
}

func generateOrchestrationDefinitionEndpoints(r spec.Generator) {
	orchestration := r.Group("/api/v1alpha1/orchestration-definitions")

	orchestration.Post("",
		option.Summary("Create an Orchestration Definition"),
		option.Description("Create a new Orchestration Definition"),
		option.Request(v1alpha1.OrchestrationDefinition{}),
		option.Response(http.StatusCreated, nil),
	)
}
