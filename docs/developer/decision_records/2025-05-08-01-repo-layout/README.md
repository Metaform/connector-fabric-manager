# Repository Layout

## Decision

The CFM repository layout will mirror its extensibility architecture, with distinct go modules for subsystem runtimes, a
go module for shared service assemblies, and a go module for shared code (e.g., types, structures, and functions).

## Rationale

Adopting a repository layout that mirrors the extensibility architecture will promote code modularity and
maintainability. A consistent layout will also make it easier to understand the codebase.

## Approach

The following approach will be instituted:

- Assemblies for specific subsystem runtimes will be contained in a Go module placed in a top-level directory.
- Shared assemblies will be placed under the top-level `assemblies` directory in a separate module.
- Shared code will be located in a separate module under the `common` directory.