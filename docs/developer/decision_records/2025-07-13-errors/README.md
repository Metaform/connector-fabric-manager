# Raising and Handling Errors

## Decision

When possible, code should distinguish between the following error types: fatal, recoverable, and client errors. Use
errors in the common `types` packages to distinguish these different error conditions. When it is not possible or
needed to distinguish error types, normal error handling applies.

## Rationale

Reliable systems need to signal to clients when retry attempts should be made, which is achieved by using error types.

## Approach

### Error Types

Use one of the following types:

```
type RecoverableError interface {
    error
    IsRecoverable() bool
}

type ClientError interface {
    error
    IsClientError() bool
}

type FatalError interface {
    error
    IsFatal() bool
}
```

Functions to instantiate concrete errors will be provided including variants for wrapped types. Client code can then use
specific error functions to discriminate error types:

```
if err != nil {
    switch {
        case model.IsClientError(err):
        //...
        case model.IsRecoverable(err):
        //...
        case model.IsFatal(err):
        //...
        default:
        //...
    }
}
```

### Error Correlation

Error responses sent to external clients should include a correlation id. This can be done in the following way:

```go
id := uuid.New().String()
a.logMonitor.Infof("Recoverable error encountered [%s]: %w ", id, err)
http.Error(w, fmt.Sprintf("Recoverable error encountered [%s], id"), http.StatusServiceUnavailable)
```

### Sending Recoverable Errors to HTTP Clients

When returning a recoverable error to an HTTP client, use `http.StatusServiceUnavailable`, HTTP code 503.