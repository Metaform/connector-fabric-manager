# HTTP Client

## Decision

All code that makes HTTP(S) invocations must use a provided HTTP client service. This service will use
the [Retriable HTTP client](https://github.com/hashicorp/go-retryablehttp) library.

## Rationale

All HTTP interactions initiated by CFM runbtimes should be managed with support for such features as automatic retry.

## Approach

A `ServiceAssembly` will be provided that registers an HTTP client service. The following values will be configurable:

- Retry max
- Retry min
- Retry wait max
- Retry wait min

The client can be accessed by other assemblies as follows:

```
client := assembler.Resolve(httpclient.HttpClientKey).(http.Client)
```
