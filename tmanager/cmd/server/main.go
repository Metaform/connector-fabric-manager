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

package main

import (
	"context"
	"fmt"
	"github.com/metaform/connector-fabric-manager/common/config"
	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/tmanager/tmcore"
	"github.com/metaform/connector-fabric-manager/tmanager/tmrouter"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	defaultPort  = 8080
	configPrefix = "tmconfig"
	key          = "httpPort"
)

// The entry point for the Tenant Manager runtime.
func main() {
	mode := runtime.LoadMode()

	logMonitor := runtime.LoadLogMonitor(mode)
	//goland:noinspection GoUnhandledErrorResult
	defer logMonitor.Sync()

	vConfig, err := config.LoadConfig(configPrefix)
	if err != nil {
		// ignore not found error, otherwise panic
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("error reading config file: %w", err))
		}
	} else if vConfig == nil {
		panic(fmt.Errorf("error loading config: %w", err))
	}

	//goland:noinspection GoDfaErrorMayBeNotNil
	vConfig.SetDefault(key, defaultPort)

	tManager := tmcore.NewTenantManager(logMonitor, vConfig, mode)
	tManager.ServiceAssembler.Register(tmrouter.RouterServiceAssembly{})
	err = tManager.ServiceAssembler.Assemble()
	if err != nil {
		panic(fmt.Errorf("error assembling runtime: %w", err))
	}

	router, found := tManager.ServiceAssembler.Resolve(tmrouter.RouterKey)
	if !found {
		panic("Router not configured")
	}

	port := vConfig.GetInt(key)
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: router.(http.Handler),
	}

	// channel for shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logMonitor.Infof("Tenant Manager starting [%d]", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logMonitor.Severef("Tenant Manager failed to start", "error", err)
		}
	}()

	// wait for interrupt signal
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logMonitor.Severef("Error attempting shutdown", "error", err)
	}

	logMonitor.Infof("Tenant Manager shutdown")
}
