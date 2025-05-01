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
	"flag"
	"fmt"
	"github.com/metaform/connector-fabric-manager/common/config"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/tmanager/tmcore"
	"github.com/metaform/connector-fabric-manager/tmanager/tmrouter"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
	mode         = "mode"
)

// The entry point for the Tenant Manager runtime.
func main() {
	mode := loadMode()

	logger := loadLogger(mode)
	//goland:noinspection GoUnhandledErrorResult
	defer logger.Sync()

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

	sLogger := logger.Sugar()

	tManager := tmcore.NewTenantManager(logger, vConfig, mode)
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
		sLogger.Infof("Tenant Manager starting [%d]", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sLogger.Fatal("Tenant Manager failed to start", zap.Error(err))
		}
	}()

	// wait for interrupt signal
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		sLogger.Fatal("Error attempting shutdown", zap.Error(err))
	}

	logger.Info("Tenant Manager shutdown")
}

func loadLogger(mode system.RuntimeMode) *zap.Logger {
	var logger *zap.Logger
	var err error
	if mode == system.DevelopmentMode || mode == system.DebugMode {
		logger, err = zap.NewDevelopment(zap.WithCaller(false))
	} else {
		logger, err = zap.NewProduction(zap.WithCaller(false))
	}
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}
	return logger
}

func loadMode() system.RuntimeMode {
	modeFlag := flag.String(mode, system.ProductionMode, "Runtime mode: development, production, or debug")
	flag.Parse()

	mode, err := system.ParseRuntimeMode(*modeFlag)
	if err != nil {
		panic(fmt.Errorf("error parsing runtime mode: %w", err))
	}
	return mode
}
