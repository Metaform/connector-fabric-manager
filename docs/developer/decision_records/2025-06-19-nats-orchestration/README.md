# NATS Orchestration

## Decision

[NATS](https://nats.io/) will be used as the default orchestration messaging system.

## Rationale

NATS provides is a lightweight, reliable messaging system that supports all of the qualities of service required by the
Deployment Orchestrator. It was chosen over other messaging systems because of its easy configuration, compact runtime
footprint, and support for a wide-range of capabilities, including durable subscriptions, highly performant key-value
storage, and clustering.

## Approach

The Deployment Orchestrator will provide interface definitions that will be used to add NATs and other messaging
systems support.    