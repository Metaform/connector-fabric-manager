# Storage Extensibility

## Decision

The CFM will support in-memory and [Postgres](https://www.postgresql.org/) as pluggable implementations for all storage
interfaces.

## Rationale

Deployment environments have different storage requirements. For example, unit tests are generally best performed using
an in-memory database for speed and to limit the potential for side effects. Postgres is selected for environments
requiring persistent storage due to its ubiquity and wide support by many cloud platforms.

## Approach

All storage extension points will be defined using interfaces and backed by service assemblies that support in-memory
implementations and Postgres. Runtimes built for specific targets will include one of the service assemblies.