# HTTP Routing

## Decision

CFM runtimes will use the [Chi router](https://github.com/go-chi/chi) for HTTP endpoints.

## Rationale

Chi is a popular Go router.

## Approach

A `ServiceAssembly` will be provided that registers a root Chi router. Other assemblies can depend on this assembly to
provide routes by registering handlers, middleware, and child routers in a modular way. For example:

```
func (h *HandlerServiceAssembly) Init(context *system.InitContext) error {
	router := context.Registry.Resolve(routing.RouterKey).(chi.Router)
	router.Use(middleware.Recoverer)

	router.Get("/endpoint", func(w http.ResponseWriter, r *http.Request) {
		response := response{Message: "OK"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}
```
