# Modularity

## Decision

All CFM runtime components will be composed of independent, reusable modules. These modules will define extensibility
points that allow users to extend runtimes with custom functionality. For example, a runtime may be configured with
different storage backends based on its target deployment environment. In order not to clash with Go `modules`, the unit
of modularity in CFM is termed an `assembly`.

## Rationale

CFM assemblies will address the following requirements:

- Efficiently compose runtimes from a set of independent, reusable units of functionality.
- Be able to extend a runtime with custom functionality.
- Have a standard way to bootstrap runtimes.

## Approach

A `ServiceAssembler` will be responsible for instantiating a set of `ServiceAssembly` instances. Each runtime will
configure the set of `ServiceAssembly` types that must be installed.

A `ServiceAssembly` may _provide_ services that can be used by other assemblies and _require_ services from other
assemblies. During bootstrapping, the `ServiceAssembler` will order all `ServiceAssembly` instances according to their
dependencies and invoke a set of lifecycle callbacks. During initialization, the `ServiceAssembler` will call the
`ServiceAssmbly.Init` method, passing an `InitContext`. The `InitContext` contains a `ServiceRegistry` that can
used to register and resolve services.

```
type ServiceAssembly interface {
	Name() string
	Provides() []ServiceType
	Requires() []ServiceType
	RequiresOptional() []ServiceType
	Init(*InitContext) error
	Prepare() error
	Start() error
	Finalize() error
	Shutdown() error
}
```

All shared services must be associated with a `ServiceType` key. Services are registered and resolved using a
`ServiceType` key with the `ServiceRegistry` obtained from the`InitContext`. A `ServiceAssembly` declares the services
it provides and requires using `ServiceType` keys.

Runtimes boostrap using a `ServiceAssembler` configured with a set of `ServiceAssembly` instances:

```
assembler := system.NewServiceAssembler(...)
assembler.Register(&httpclient.HttpClientServiceAssembly{})
assembler.Register(&routing.RouterServiceAssembly{})
assembler.Register(&tmhandler.HandlerServiceAssembly{})
err := assembler.Assemble()
if err != nil {
    panic(fmt.Errorf("error assembling runtime: %w", err))
}
```

Note that the `ServiceAssembly` is a modularity system, not a dependency injection framework. While valuable in some
contexts, a DI framework adds complexity not needed for the CFM.

### Lifecycles

Lifecycle methods are defined and invoked as follows:

- **Init(*InitContext)**: All provided services must be registered. Required services may be resolved but not yet
  invoked.
- **Prepare()**: Additional setup may be performed. For example, a registry of extension points may need to be.
  initialized after those extension points have been registered with it during the `Init` phase.
- **Start()**: Services such as transport listeners should be readied for processing.
- **Finalize()**: Resource cleanup should be performed, such as canceling ongoing processing.
- **Shutdown()**: All services and subsystems should be disposed.

Note that the lifecycle callback sequence is single-threaded in order to keep the boostrap simple (no need to deal with
concurrent access) and efficient.  