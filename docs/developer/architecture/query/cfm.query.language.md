# The CFM Query Language

## Overview

The CFM query language in is a predicate-based system designed to work with both in-memory storage and Postgres
databases. It uses ANTLR (ANother Tool for Language Recognition) for parsing human-readable query syntax into
structured predicate objects that can be evaluated against any data source.

## Architecture

### Core Components

1. **Query String Syntax** - Human-readable query expressions (parsed via ANTLR grammar)
2. **Predicates** - Abstract representation of query conditions
3. **Field Matcher** - Strategy pattern for field value extraction and comparison
4. **Storage Adapters** - Implementations for in-memory and database backends

## Query Syntax

The query language supports both simple and complex query expressions.

### Basic Syntax Structure

`field operator value`

### Comparison Operators

| Operator | Meaning               | Example              |
|----------|-----------------------|----------------------|
| =        | Equal                 | status = 'active'    |
| !=       | Not equal             | status != 'inactive' |
| \>       | Greater than          | count > 10           |
| =        | Greater than or equal | count >= 10          |
| <        | Less than             | count < 100          |
| <=       | Less than or equal    | count <= 100         |

### String Operations

| Operator    | Meaning                        | Example                       |
|-------------|--------------------------------|-------------------------------|
| LIKE        | Pattern matching (database)    | name LIKE 'John%'             |
| CONTAINS    | Substring contains (in-memory) | description CONTAINS 'urgent' |
| STARTS_WITH | Prefix matching                | code STARTS_WITH 'ERR'        |
| ENDS_WITH   | Suffix matching                | filename ENDS_WITH '.log'     |

## Collection Operations

| Operator | Meaning          | Example                               |
|----------|------------------|---------------------------------------|
| IN       | Value in set     | status IN ('active', 'pending')       |
| NOT IN   | Value not in set | status NOT IN ('deleted', 'archived') |

## NULL Checks

| Operator    | Meaning           | Example                 |
|-------------|-------------------|-------------------------|
| IS NULL     | Field is null     | description IS NULL     |
| IS NOT NULL | Field is not null | description IS NOT NULL |

### Compound Predicates

Combine multiple conditions using logical operators:

**AND Operator** - All conditions must be true:

`field1 = 'value1' AND field2 > 100 AND field3 STARTS_WITH 'prefix'`

**OR Operator** - At least one condition must be true:

`status = 'active' OR status = 'pending' OR status = 'processing'`

**Nested Expressions** with parentheses:

`(status = 'active' OR status = 'pending') AND priority > 5`

### Example Queries

Simple equality:

`correlationID = '12345'`

Nested object field access:

`properties.group = 'suppliers'`

Complex filtering:

`(state = 'RUNNING' AND priority >= 1) OR (state = 'PENDING' AND createdTimestamp > '2025-01-01')`

String operations:

`orchestrationType STARTS_WITH 'provisioning' AND state != 'FAILED'`

## Predicate System

### Predicate Interface

All predicates implement the Predicate interface with two methods:

```
type Predicate interface {
    // Matches evaluates the predicate against an object
    Matches(obj any, matcher FieldMatcher) bool
    
    // String returns a human-readable representation
    String() string
}
```

### Atomic Predicates

Single condition predicates consisting of:

- **Field** - The object property to test
- **Operator** - The comparison operation
- **Value** - The comparison value

```
predicate := &AtomicPredicate{
    Field:    "state",
    Operator: "=",
    Value:    "RUNNING",
}
```

### Compound Predicates

Logical combinations of multiple predicates:

```
query.And(
    query.Or(
        query.Eq("state", "RUNNING"),
        query.Eq("state", "PENDING"),
    ),
    query.Gte("priority", 1),
)
```

## Field Matcher Strategy

The `FieldMatcher` interface enables custom logic for different data types.

### Default Field Matcher

The `DefaultFieldMatcher` provides reflection-based matching: Navigates nested object structures using dot notation (
e.g.,
properties.group) Supports standard Go types (strings, numbers, times, booleans) Type-safe comparisons

### Custom Field Matchers

Different storage backends can implement custom matchers for optimization: In-Memory: Direct object navigation
Postgres: Field extraction from JSON/JSONB columns Other DBs: Adapter-specific implementations.

## In-Memory Storage Implementation

### Query Processing Flow

```go
1. Query String (e.g., "status = 'active' AND count > 10")
↓
2. ANTLR Parser → Predicate Object
↓
3. InMemoryEntityStore.FindByPredicate()
↓
4. Iterate entities and evaluate Predicate.Matches()
↓
5. Return matching entities
```

## Postgres Storage Implementation

### Query Translation

```
Predicate Object → SQL WHERE Clause → Postgres Query
```

### Operator Mapping

| Predicate Operator | SQL Translation                                     |
|--------------------|-----------------------------------------------------|
| =                  | column = value                                      |
| !=                 | column != value                                     |
| \>                 | column > value                                      |
| IN                 | column IN (...)                                     |
| LIKE               | column LIKE pattern                                 |
| CONTAINS           | column::text @> value (JSONB) or substring matching |
| IS NULL            | column IS NULL                                      |

### Example Translation

Predicate:

```
query.And(
    query.Eq("state", "ACTIVE"),
    query.Gt("priority", 5),
)
```

SQL WHERE Clause:

`WHERE state = 'ACTIVE' AND priority > 5`

## CFM Entities

The query language can be used to query CFM entities:

- **Tenants and participant profiles**: Query tenant properties, profiles, and configurations
- **Orchestrations**:Filter by state, type, correlation ID, timestamps
- **Definitions**: Search orchestration and activity definitions
- **VPA Metadata**: Find virtual participant agents by type, state, cell assignment 

