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
	"flag"
	"fmt"
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/metaform/connector-fabric-manager/common/system"
	"go.uber.org/zap"
)

const mode = "mode"

func LoadLogMonitor(mode system.RuntimeMode) monitor.LogMonitor {
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
