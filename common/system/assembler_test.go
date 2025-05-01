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

package system

import (
	"errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"testing"
)

// MockServiceAssembly implements ServiceAssembly interface for testing
type MockServiceAssembly struct {
	id        string
	name      string
	provides  []ServiceType
	requires  []ServiceType
	initFunc  func(*ServiceRegistry) error
	destroyed bool
}

func (m *MockServiceAssembly) ID() string              { return m.id }
func (m *MockServiceAssembly) Name() string            { return m.name }
func (m *MockServiceAssembly) Provides() []ServiceType { return m.provides }
func (m *MockServiceAssembly) Requires() []ServiceType { return m.requires }
func (m *MockServiceAssembly) Init(ctx *InitContext) error {
	if m.initFunc != nil {
		return m.initFunc(ctx.Registry)
	}
	return nil
}
func (m *MockServiceAssembly) Destroy(logger *zap.Logger) error {
	m.destroyed = true
	return nil
}

func TestServiceAssembler_Register(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	assembler := NewServiceAssembler(logger, viper.New(), DebugMode)
	mock := &MockServiceAssembly{id: "test", name: "Test Assembly"}

	assembler.Register(mock)

	if len(assembler.assemblies) != 1 {
		t.Errorf("Expected 1 assembly, got %d", len(assembler.assemblies))
	}
	if assembler.assemblies[0] != mock {
		t.Error("Registered assembly does not match mock")
	}
}

func TestServiceAssembler_Assemble_Simple(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	assembler := NewServiceAssembler(logger, viper.New(), DebugMode)
	mock := &MockServiceAssembly{
		id:       "test",
		name:     "Test Assembly",
		provides: []ServiceType{"service1"},
	}

	assembler.Register(mock)

	err := assembler.Assemble()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestServiceAssembler_Assemble_WithDependencies(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	assembler := NewServiceAssembler(logger, viper.New(), DebugMode)

	mock1 := &MockServiceAssembly{
		id:       "test1",
		name:     "Test Assembly 1",
		provides: []ServiceType{"service1"},
	}

	mock2 := &MockServiceAssembly{
		id:       "test2",
		name:     "Test Assembly 2",
		provides: []ServiceType{"service2"},
		requires: []ServiceType{"service1"},
	}

	assembler.Register(mock2)
	assembler.Register(mock1)

	err := assembler.Assemble()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestServiceAssembler_Assemble_MissingDependency(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	assembler := NewServiceAssembler(logger, viper.New(), DebugMode)

	mock := &MockServiceAssembly{
		id:       "test",
		name:     "Test Assembly",
		requires: []ServiceType{"missing-service"},
	}

	assembler.Register(mock)

	err := assembler.Assemble()
	if err == nil {
		t.Error("Expected error for missing dependency, got nil")
	}
}

func TestServiceAssembler_Assemble_CyclicDependency(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	assembler := NewServiceAssembler(logger, viper.New(), DebugMode)

	mock1 := &MockServiceAssembly{
		id:       "test1",
		name:     "Test Assembly 1",
		provides: []ServiceType{"service1"},
		requires: []ServiceType{"service2"},
	}

	mock2 := &MockServiceAssembly{
		id:       "test2",
		name:     "Test Assembly 2",
		provides: []ServiceType{"service2"},
		requires: []ServiceType{"service1"},
	}

	assembler.Register(mock1)
	assembler.Register(mock2)

	err := assembler.Assemble()
	if err == nil {
		t.Error("Expected error for cyclic dependency, got nil")
	}
}

func TestServiceAssembler_Assemble_InitializationError(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	assembler := NewServiceAssembler(logger, viper.New(), DebugMode)

	expectedError := errors.New("initialization failed")
	mock := &MockServiceAssembly{
		id:       "test",
		name:     "Test Assembly",
		provides: []ServiceType{"service1"},
		initFunc: func(*ServiceRegistry) error {
			return expectedError
		},
	}

	assembler.Register(mock)

	err := assembler.Assemble()
	if err == nil {
		t.Error("Expected initialization error, got nil")
	}
}
