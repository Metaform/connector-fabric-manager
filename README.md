![](./docs/logo/cfm.small.logo.svg)

Connector Fabric Manager (CFM) provides infrastructure for managing virtualized dataspace connectors. CFM consists of
the following components:

- **Tenant Manager** - Tenant management and deployment
- **Provision Manager** - Resource provisioning and orchestration

## Prerequisites

The following are required to build and run the system:

- Go
- A Docker-compatible CLI

## Setting up the Workspace

The CFM is a multimodule project and requires the creation of a workspace file:

```bash
go work init ./assembly ./common ./pmanager ./tmanager ./e2e ./mvd ./agent/edcv ./agent/keycloak ./agent/registration
```
