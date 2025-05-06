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

package runtime

import (
	"context"
	"flag"
	"fmt"
	"github.com/metaform/connector-fabric-manager/assembly/routing"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	key  = "httpPort"
	mode = "mode"
)

func LoadLogMonitor(mode system.RuntimeMode) monitor.LogMonitor {
	var config zap.Config
	var options []zap.Option

	switch mode {
	case system.DebugMode, system.DevelopmentMode:
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		config.EncoderConfig.StacktraceKey = "stacktrace"
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	default:
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		// Disable stacktrace in production
		config.EncoderConfig.StacktraceKey = ""
	}

	config.DisableCaller = true

	// Add caller skip for accurate source locations
	options = append(options, zap.AddCallerSkip(1))

	logger, err := config.Build(options...)
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
	}

	return NewSugaredLogMonitor(logger.Sugar())
}

func LoadMode() system.RuntimeMode {
	modeFlag := flag.String(mode, system.ProductionMode, "Runtime mode: development, production, or debug")
	flag.Parse()

	mode, err := system.ParseRuntimeMode(*modeFlag)
	if err != nil {
		panic(fmt.Errorf("error parsing runtime mode: %w", err))
	}
	return mode
}

// SugaredLogMonitor implements LogMonitor by wrapping a zap.SugaredLogger
type SugaredLogMonitor struct {
	logger *zap.SugaredLogger
}

// NewSugaredLogMonitor creates a new LogMonitor that wraps a zap.SugaredLogger
func NewSugaredLogMonitor(logger *zap.SugaredLogger) monitor.LogMonitor {
	return &SugaredLogMonitor{logger: logger}
}

func (s *SugaredLogMonitor) Named(name string) monitor.LogMonitor {
	return &SugaredLogMonitor{logger: s.logger.Named(name)}
}

func (s *SugaredLogMonitor) Severef(message string, args ...any) {
	s.logger.Errorf(message, args...)
}

func (s *SugaredLogMonitor) Warnf(message string, args ...any) {
	s.logger.Warnf(message, args...)
}

func (s *SugaredLogMonitor) Infof(message string, args ...any) {
	s.logger.Infof(message, args...)
}

func (s *SugaredLogMonitor) Debugf(message string, args ...any) {
	s.logger.Debugf(message, args...)
}

func (s *SugaredLogMonitor) Severew(message string, keysValues ...any) {
	s.logger.Errorw(message, keysValues...)
}

func (s *SugaredLogMonitor) Warnw(message string, keysValues ...any) {
	s.logger.Warnw(message, keysValues...)
}

func (s *SugaredLogMonitor) Infow(message string, keysValues ...any) {
	s.logger.Infow(message, keysValues...)
}

func (s *SugaredLogMonitor) Debugw(message string, keysValues ...any) {
	s.logger.Debugw(message, keysValues...)
}

func (s *SugaredLogMonitor) Sync() error {
	return s.logger.Sync()
}

// AssembleLaunch assembles and launches the runtime with the given name and configuration.
// The runtime will be shutdown when the program is terminated.
func AssembleLaunch(assembler *system.ServiceAssembler, name string, vConfig *viper.Viper, logMonitor monitor.LogMonitor) {

	err := assembler.Assemble()
	if err != nil {
		panic(fmt.Errorf("error assembling runtime: %w", err))
	}

	router := assembler.Resolve(routing.RouterKey).(http.Handler)

	port := vConfig.GetInt(key)
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: router,
	}

	// channel for shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logMonitor.Infof("%s starting [%d]", name, port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logMonitor.Severew("failed to start", "error", err)
		}
	}()

	// wait for interrupt signal
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logMonitor.Severew("Error attempting server shutdown", "error", err)
	}

	if err := assembler.Shutdown(); err != nil {
		logMonitor.Severew("Error attempting shutdown", "error", err)
	}

	logMonitor.Infof("%s shutdown", name)
}
