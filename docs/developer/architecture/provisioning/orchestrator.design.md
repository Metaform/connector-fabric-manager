# The Deployment Orchestrator

|                |              |
|----------------|--------------|
| **Status**     | In design    |
| **Subsystem**  | Provisioning |
| **Committers** | @jimmarino   |

## Overview

The `DeploymentOrchestrator` is responsible for executing and managing the deployment of resources to a target. A
_**resource**_ can be anything from a compute cluster to tenant configuration at the application layer. Since resources
often depend on one another, a _**deployment**_ is defined as a collection of activities termed an _**orchestration**_.
Each activity collection is executed as a sequence or in parallel.

Consider a tenant deployment that involves the creation of a Web DID using a domain supplied by the tenant owner. The
deployment involves the following activities:

- Input domain name, tenant metadata, and target cell
- Stage and apply the tenant configuration to the application in the target cell
- Apply ingress routing configuration for the tenant domain

The orchestrator is responsible for executing these activities in the correct order, maintaining state, and ensuring
reliable processing. Because activities may have high latency and processing must scale-out, the orchestrator is
designed as a stateful message-based system.

## Stateful Messaging

Activities are executed on worker nodes that dequeue messages. During execution, activities have access to a shared
persistent context managed by the orchestration framework. Activities must be idempotent, that is, they must complete
with the same result if executed multiple times without side effects. For example, if an activity is invoked twice due
to a failure, it must ensure duplicate resources are not created and the same shared stated is applied to the context.

Activities may be implemented using a variety of programming languages and technologies, for example, a custom Go
service or Terraform script. The `DeploymentOrchestrator` delegates to a _**provider**_ for an activity type that is an
extensibility point for the system.

### Messaging Implementation

The messaging implementation will be pluggable. The initial system will be based
on [NATS Jetstream](https://docs.nats.io/nats-concepts/jetstream). A design goal is to allow the use of other
technologies such as [Temporal](https://github.com/temporalio).

### Kubernetes Integration

The Deployment Orchestrator will be deployable as a standalone application or to a Kubernetes cluster. While it is
possible to implement the Orchestration Resource Model described above
as [Custom Resource Definitions](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/),
doing so will add additional complexity (the need to implement Kubernetes operators) and tie the solution to Kubernetes.

## Resource Model: Deployments, Activities, and Definitions

The `DeploymentOrchestrator` is built on a resource model consisting of two types: a `DeploymentDefinition` and an
`ActivityDefinition`. A `DeploymentDefinition` contains a collection of `ActivityDefinitions` that define the
orchestration for a deployment. The following is an example of a `DeploymentDefinition`:

```json
{
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
        "openAPIV3Schema": {}
      },
      "orchestration": [
        {
          "activities": [
            {
              "id": "activity1",
              "type": "activity1.example.com",
              "inputs": [
                "cell",
                "baseUrl"
              ]
            },
            {
              "id": "activity2",
              "type": "activity2.example.com",
              "dependsOn": "activity1",
              "inputs": [
                {
                  "source": "activity1.resource",
                  "target": "resourceId"
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

The `DeploymentDefinition` [JSON Schema](./deployment.definition.schema.json) defines the following properties and
types:

- `type`: The definition type used when creating a corresponding resource.
- `apiVersion`: The API version the definition applies to.
- `resource`: The resource metadata, which is used to create API endpoints for deployment resources of the definition
  type.
- `versions`: Contains one or more versions of the deployment definition.

Each version defines the following properties:

- `version`: The version of the deployment definition.
- `active`: Indicates whether the version is active, i.e., deployments at that version level can be created and run.
- `schema`: The schema for input properties when creating a deployment resource of the definition type. Currently,
  `openAPIV3Schema` is the only supported schema type. Activities may reference input properties.
- `orchestration`: Defines the sequence of activities that are executed to deploy the resource.

The orchestration is a collection of activities. Activities may form a Directed Acyclic Graph (DAG) by declaring
dependencies using the `dependsOn` property. At deployment time, the activities will be ordered using a topological sort
and grouping activities into tiers of parallel execution steps based on their dependencies.

An activity has the following properties:

- `id`: The activity identifier.
- `type`: The activity type.
- `dependsOn`: An array of activity ids the activity depends on.
- `inputs`: An array of input properties. The input properties may include references to properties contained in the
  deployment input data or references to output data properties from a previous activity. References to activity output
  data are prefixed with the activity identifier followed by a '.'. Activity output data is defined in the activity
  definition described below. An input property may be specified using a string or an object containing `source` and
  `target` properties if a mapping is required.

An `ActivityDefinition` defines a work item reliably executed by a worker. For example:

```json
{
  "type": "activity1.example.com",
  "provider": "provisioner.example.com",
  "description": "Provisions a resource for a tenant",
  "inputSchema": {
    "openAPIV3Schema": {}
  },
  "outputSchema": {
    "openAPIV3Schema": {}
  }
}
```

The `ActivityDefinition` [JSON Schema](./activity-definition.schema.json) specifies the following properties:

- `type`: The activity type used as a reference.
- `provider`: The provisioner that executes the activity. A provisioner could be a service, Terraform script, or
  other technology.
- `description`: A description of the activity.
- `inputSchema`: The schema for input properties when creating a deployment resource of the definition type. Currently,
  `openAPIV3Schema` is the only supported schema type.
- `outputSchema`: The schema for output properties when creating a deployment resource of the definition type.
  Currently, `openAPIV3Schema` is the only supported schema type.

## Activity Executors

When an orchestration is executed, the Deployment Orchestrator reliably enqueues activity messages which will be
dequeued and processed by an associated activity executor. The executor delegates to an `ActivityProcessor` to process
the message. The Deployment Orchestrator is responsible for handling system reliability, context persistence, recovery,
and activity coordination.

An `ActivityProcessor` is an extensibility points for integrating technologies such as Terraform or custom operations
code into the deployment process. For example, a Terraform processor would gather input data associated with the
orchestration and pass it to a Terraform script for execution. The `ActivityProcess` interface is defined as follows:

```go
package api

type ActivityProcessor interface {
	Process(activityContext ActivityContext) ActivityResult
}

type ActivityResultType int

type ActivityResult struct {
	Result     ActivityResultType
	WaitMillis time.Duration
	Error      error
}

const (
	ActivityResultWait       = 0
	ActivityResultComplete   = 1
	ActivityResultSchedule   = 2
	ActivityResultRetryError = -1
	ActivityResultFatalError = -2
)

```

The `ActivityResult` indicates the following actions to be taken:

- **ActivityResultWait** - The message is acknowledged and the activity must be marked for completion by an external
  process. This is useful for activity types that asynchronously execute a callback on completion.
- **ActivityResultComplete** - The activity is marked as completed and the message is acknowledged.
- **ActivityResultSchedule** - Schedules the message for redelivery as defined by `WaitMillis`. This can be used to
  implement a completion polling mechanism.
- **ActivityResultRetryError** - A recoverable error was raised andtThe message is negatively acknowledged so that it
  can be redelivered.
- **ActivityResultRetryError** - A fatal error was raised, the orchestration is put into the error state, and the
  messsage is acknowledged so it will not be redelivered.

### Activity Processors

The following providers will be created:

- Terraform/Open Tofu
- HTTP endpoint
- Custom Go providers for handling required application-level resources


