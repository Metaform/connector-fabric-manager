# Repository Layout

## Decision

The CFM will keep its existing repository layout but move to a single Go module.

## Rationale

Multi-module repositories are difficult to maintian and release and don't offer CFM much benefit given existing
components are versioned together.

## Approach

The following approach will be instituted:

- The existing folder and package structure will be maintained.
- All module files except the root will be removed.
