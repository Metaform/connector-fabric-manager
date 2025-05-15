# Testing and Test Frameworks

## Decision

The following test frameworks will be used:

- [Testify](https://github.com/stretchr/testify) for test assertions.
- [Mockery](https://github.com/vektra/mockery/) for generating test mocks.
- [Test Containers](https://github.com/testcontainers/testcontainers-go) when code under test requires infrastructure such as a
  database.

## Rationale

Testing should be standardized throughout the codebase.

## Approach

The following guidelines should be followed:

- All code must be covered by unit tests.
- Do not require Test Containers or external setup for unit tests unless the code under test is tied to a particular
  type of infrastructure. For example, it is OK to test code that specifically persists to Postgres using a Postgres
  container. It is NOT OK to use a Postgres container to test generic code.
- Do not rely on integration tests to verify discrete units of code such as a method or structure.
- Include integration tests to verify multiple units of code work together.
- End-to-end tests should only verify basic operations. They should not be used as a substitute for unit or integration
  tests.
- Use Mockery when writing unit tests to mock dependencies. Do not instantiate implementations of dependent interfaces.  