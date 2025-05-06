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
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/common/system"
	"net/http"
	"time"
)

const (
	RouterKey system.ServiceType = "router:Router"
)

type RouterServiceAssembly struct {
}

func (r *RouterServiceAssembly) ID() string {
	return "common:RouterServiceAssembly"
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
	router := r.setupRouter(ctx.LogMonitor, ctx.Mode)
	ctx.Registry.Register(RouterKey, router)
	return nil
}

func (r *RouterServiceAssembly) Prepare() error {
	return nil
}

func (r *RouterServiceAssembly) Start() error {
	return nil
}

func (r *RouterServiceAssembly) Shutdown() error {
	return nil
}

func (r *RouterServiceAssembly) Finalize() error {
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
