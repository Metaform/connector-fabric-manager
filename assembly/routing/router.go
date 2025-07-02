//  Copyright (c) 2025 Metaform Systems, Inc
//
//  This program and the accompanying materials are made available under the
//  terms of the Apache License, Version 2.0 which is available at
//  https://www.apache.org/licenses/LICENSE-2.0
//
//  SPDX-License-Identifier: Apache-2.0
//
//  Contributors:
//       Metaform Systems, Inc. - initial API and implementation
//

package routing

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"time"
)

const (
	RouterKey system.ServiceType = "router:Router"
	key                          = "httpPort"
)

type RouterServiceAssembly struct {
	system.DefaultServiceAssembly
	server     *http.Server
	router     *chi.Mux
	logMonitor monitor.LogMonitor
	config     *viper.Viper
}

func (r *RouterServiceAssembly) Name() string {
	return "Router"
}

func (r *RouterServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{RouterKey}
}

func (r *RouterServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{}
}

func (r *RouterServiceAssembly) Init(ctx *system.InitContext) error {
	r.router = r.setupRouter(ctx.LogMonitor, ctx.Mode)
	ctx.Registry.Register(RouterKey, r.router)
	r.logMonitor = ctx.LogMonitor
	r.config = ctx.Config
	return nil
}

func (r *RouterServiceAssembly) Start(ctx *system.StartContext) error {
	port := r.config.GetInt(key)
	r.server = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: r.router,
	}

	go func() {
		r.logMonitor.Infof("HTTP server listening on [%d]", port)
		if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			r.logMonitor.Severew("failed to start", "error", err)
		}
	}()
	return nil
}

func (r *RouterServiceAssembly) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := r.server.Shutdown(ctx); err != nil {
		r.logMonitor.Severew("Error attempting HTTP server shutdown", "error", err)
	}
	return nil
}

// SetupRouter configures and returns the HTTP router
func (r *RouterServiceAssembly) setupRouter(logMonitor monitor.LogMonitor, mode system.RuntimeMode) *chi.Mux {
	router := chi.NewRouter()

	if mode == system.DebugMode {
		router.Use(createLoggerHandler(logMonitor))
	}
	router.Use(middleware.Recoverer)

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
