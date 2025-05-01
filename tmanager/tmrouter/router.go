package tmrouter

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/metaform/connector-fabric-manager/common/system"
	"go.uber.org/zap"
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
	router := r.setupRouter(ctx.Logger, ctx.Mode)
	ctx.Registry.Register(RouterKey, router)
	return nil
}

func (r RouterServiceAssembly) Destroy(logger *zap.Logger) error {
	return nil
}

// SetupRouter configures and returns the HTTP router
func (r RouterServiceAssembly) setupRouter(logger *zap.Logger, mode system.RuntimeMode) *chi.Mux {
	router := chi.NewRouter()

	if mode == system.DebugMode {
		router.Use(createLoggerHandler(logger))
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

func createLoggerHandler(log *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				log.Debug("http",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Int("status", ww.Status()),
					zap.Int("bytes", ww.BytesWritten()),
					zap.Duration("duration", time.Since(start)),
					zap.String("reqId", middleware.GetReqID(r.Context())),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
