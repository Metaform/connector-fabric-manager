# Registries

## Decision

Registries will be used for extension points that have multiple implementations. The `ServiceAssembler` system
and `ServiceRegistry` will therefore not require direct support for service multiplicities.

## Rationale

Registries are a well-established way to support extensibility without complicating the core modularity system.

## Approach

A registry can be created by a `ServiceAssembly` and registered with the `ServiceRegistry`. The registry can be declared
as a dependency and resolved by other `ServiceAssembly` implementations. During the `Init` phase, those
implementations can register providers with the resolved registry:

```
func (a *SampleAssembly) Init(context *system.InitContext) error {
	testProvider := ...
	registry := context.Registry.Resolve(registry.RegistryKey).(registry.SampleRegistry)
    registry.RegisterProvider(testProvider)
	return nil
}
```

At runtime, registries can be used to dispatch to the correct provider.

Note that providers should be registered during the `Init` phase so that the registry can be further initialized if
needed during the `Prepare` phase.