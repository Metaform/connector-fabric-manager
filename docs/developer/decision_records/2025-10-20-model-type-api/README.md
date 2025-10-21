# Model, Type, and API Packages

## Decision

The following convention is adopted for package naming related to data types:

- `model` - Data transfer types (DTOs) used for message exchange and invocations involving remote systems (e.g., API
  clients)
- `types` - General types such as errors that are used throughout the system.
- `api` - Types and interfaces used to interface between extensible system components. These types are not used directly
  by external systems.

## Rationale

This approach provides a consistent naming convention, promotes discoverability, and enables decoupling of boundary
interfaces so that external APIs can be versioned independently of the internal implementation.

## Approach

The `model` packages may directly include unversioned DTOs. Versioned DTOs are defined in packages under a `model`
directory in the form `v[major][minor]`, e.g., `model/v1` and `model/v1beta1`. Transformers are responsible for
converting `model` types to internal types and vice versa.

Note that the `types` plural form is used to avoid a clash with the `type` keyword.