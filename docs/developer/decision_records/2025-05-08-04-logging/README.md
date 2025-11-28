# Logging

## Decision

CFM runtimes will use the [Zap](https://github.com/uber-go/zap) logging framework. However, this logging framework will
be accessed through a `LogMonitor` interface that will help enforce logging conventions and allow users to build
runtimes that use alternative logging systems.

## Rationale

Zap will be the default logging implementation since it is performant and provides convenient features such as
stacktraces. At the same time, we need to ensure consistent logging and allow users to override this default selection.
This will be done by accessing all log operations through an interface.

## Approach

The `LogMonitor` interface will be provided as part of the `InitContext` for `ServiceAssembly`:

```
type LogMonitor interface {
	Named(name string) LogMonitor

	Severef(message string, args ...any)
	Warnf(message string, args ...any)
	Infof(message string, args ...any)
	Debugf(message string, args ...any)

	Severew(message string, keyValues ...any)
	Warnw(message string, keyValues ...any)
	Infow(message string, keyValues ...any)
	Debugw(message string, keyValues ...any)

	Sync() error
}
```

The `...f()` and `...w()` methods will delegate to their counterparts on Zap's `SugaredLogger`.

### Logging Conventions

The following logging conventions should be followed.

#### 1. Named loggers

Each `ServiceAssembly` should create a named logger using `Named()` for services it instantiates.

#### 3. Error

Do not log an error if it is returned to a caller as this will create duplicate log messages. Errors should only be
logged at the top of the call hierarchy.

#### 2. Warn

The WARN level should only be used to indicate problematic or insecure configuration. For example, if a security
feature is disabled for testing.

#### 3. Info

The INFO level should be used in very limited circumstances to indicate important information. For example, that a
runtime is ready to process requests. **Do not log mundane or operational messages that do not require attention.**

#### 4. Debug

Debugging is best done with a debugger, not a message log. The DEBUG level should therefore be used sparingly. **Never
log messages such as "Processing request," "Entering method," or similar events as they serve no useful purpose.**
