# Configuration

## Decision

The [Viper library](https://github.com/spf13/viper) will be used for loading configuration.

## Rationale

While Viper brings in a number of dependencies, it fits our requirements for configuration:

- The ability to resolve configuration from multiple sources, including environment variables and files
- The ability to override configuration set in files with environment variables
- Support for multiple source formats such as JSON, YAML, and TOML
- Support for configuration hierarchies (e.g., foo.bar.baz) that map to multiple source formats

## Approach

Viper will be configured to support environment variable overrides and loaded during runtime boostrap. The Viper
configuration instance will be contained in the `InitContext`:

```
type InitContext struct {
    Config     *viper.Viper
    ...
}
```

While this will tie extensions to Viper, it does not make sense to provide a wrapper that will simply serve as a
passthrough. 