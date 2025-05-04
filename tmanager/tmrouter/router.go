package tmrouter

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/common/system"
	"net/http"
	"time"
)

const (
	RouterKey system.ServiceType = "tmrouter:Router"
)

type response struct {
	Message string `json:"message"`
}

type RouterServiceAssembly struct {
}

func (r RouterServiceAssembly) ID() string {
	return "tmanager:RouterServiceAssembly"
}

func (r RouterServiceAssembly) Name() string {
	return "Router"
}

func (r RouterServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{RouterKey}
}

func (r RouterServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{}
}

func (r RouterServiceAssembly) Init(ctx *system.InitContext) error {
	router := r.setupRouter(ctx.LogMonitor, ctx.Mode)
	ctx.Registry.Register(RouterKey, router)
	return nil
}

func (r RouterServiceAssembly) Destroy(logMonitor monitor.LogMonitor) error {
	return nil
}

// SetupRouter configures and returns the HTTP router
func (r RouterServiceAssembly) setupRouter(logMonitor monitor.LogMonitor, mode system.RuntimeMode) *chi.Mux {
	router := chi.NewRouter()

	if mode == system.DebugMode {
		router.Use(createLoggerHandler(logMonitor))
	}
	router.Use(middleware.Recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response := response{Message: "OK"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response := response{Message: "OK"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	return router
}

func createLoggerHandler(logMonitor monitor.LogMonitor) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				logMonitor.Debugw("http",
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.Status(),
					"bytes", ww.BytesWritten(),
					"duration", time.Since(start),
					"reqId", middleware.GetReqID(r.Context()),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
