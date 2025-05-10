# Core Subsystems

## Decision

CFM will consist of three runtime types: the Tenant Manager, the Provision Manager, and Provisioners. These runtimes may
be hosted in separate OS processes or collocated in a single OS process. Each runtime type is responsible for performing
a set of discrete aspects:

- **Tenant Manager**: Manages tenants and deployments associated with those tenants
- **Provision Manager**: Orchestrates processes that provision resources. These resources may be associated with a
  tenant or system infrastructure.
- **Provisioner**: Provisions a resource type.

## Rationale

Tenant management and provisioning need to be scaled independently, so it must be possible to deploy them as separate
processes. In addition, it must be possible to add provisioners without requiring a redeployment of the Provision
Manager.

## Approach

Each runtime will be composed of a set of assemblies. The runtimes must be deployable:

- As a standalone process
- As a single process
- To a Kubernetes cluster
- In a unit test

Note that the CSM may introduce additional auxiliary services in the future
