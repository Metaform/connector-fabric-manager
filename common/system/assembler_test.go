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
	"github.com/metaform/connector-fabric-manager/common/monitor"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestServiceAssembler(t *testing.T) {
	t.Run("Register", func(t *testing.T) {
		logMonitor := monitor.NoopMonitor{}
		assembler := NewServiceAssembler(logMonitor, viper.New(), DebugMode)
		mock := &MockServiceAssembly{name: "Test Assembly"}

		assembler.Register(mock)

		assert.Len(t, assembler.assemblies, 1)
		assert.Equal(t, mock, assembler.assemblies[0])
	})

	t.Run("Assemble", func(t *testing.T) {
		t.Run("Simple", func(t *testing.T) {
			logMonitor := monitor.NoopMonitor{}
			assembler := NewServiceAssembler(logMonitor, viper.New(), DebugMode)
			mock := &MockServiceAssembly{
				name:     "Test Assembly",
				provides: []ServiceType{"service1"},
			}

			assembler.Register(mock)

			err := assembler.Assemble()
			assert.NoError(t, err)
		})

		t.Run("WithDependencies", func(t *testing.T) {
			logMonitor := monitor.NoopMonitor{}
			assembler := NewServiceAssembler(logMonitor, viper.New(), DebugMode)

			mock1 := &MockServiceAssembly{
				name:     "Test Assembly 1",
				provides: []ServiceType{"service1"},
			}

			mock2 := &MockServiceAssembly{
				name:     "Test Assembly 2",
				provides: []ServiceType{"service2"},
				requires: []ServiceType{"service1"},
			}

			assembler.Register(mock2)
			assembler.Register(mock1)

			err := assembler.Assemble()
			assert.NoError(t, err)
		})

		t.Run("MissingDependency", func(t *testing.T) {
			logMonitor := monitor.NoopMonitor{}
			assembler := NewServiceAssembler(logMonitor, viper.New(), DebugMode)

			mock := &MockServiceAssembly{
				name:     "Test Assembly",
				requires: []ServiceType{"missing-service"},
			}

			assembler.Register(mock)

			err := assembler.Assemble()
			assert.Error(t, err)
		})

		t.Run("CyclicDependency", func(t *testing.T) {
			logMonitor := monitor.NoopMonitor{}
			assembler := NewServiceAssembler(logMonitor, viper.New(), DebugMode)

			mock1 := &MockServiceAssembly{
				name:     "Test Assembly 1",
				provides: []ServiceType{"service1"},
				requires: []ServiceType{"service2"},
			}

			mock2 := &MockServiceAssembly{
				name:     "Test Assembly 2",
				provides: []ServiceType{"service2"},
				requires: []ServiceType{"service1"},
			}

			assembler.Register(mock1)
			assembler.Register(mock2)

			err := assembler.Assemble()
			assert.Error(t, err)
		})

		t.Run("InitializationError", func(t *testing.T) {
			logMonitor := monitor.NoopMonitor{}
			assembler := NewServiceAssembler(logMonitor, viper.New(), DebugMode)

			expectedError := errors.New("initialization failed")
			mock := &MockServiceAssembly{
				name:     "Test Assembly",
				provides: []ServiceType{"service1"},
				initFunc: func(*ServiceRegistry) error {
					return expectedError
				},
			}

			assembler.Register(mock)

			err := assembler.Assemble()
			assert.Error(t, err)
		})
	})

	t.Run("LifecycleMethods", func(t *testing.T) {
		logMonitor := monitor.NoopMonitor{}
		assembler := NewServiceAssembler(logMonitor, viper.New(), DebugMode)

		preparedCh := make(chan bool, 1)
		startedCh := make(chan bool, 1)
		shutdownCh := make(chan bool, 1)

		mock := &MockServiceAssembly{
			name:     "Test Assembly",
			provides: []ServiceType{"service1"},
			prepareFunc: func() error {
				preparedCh <- true
				return nil
			},
			startFunc: func() error {
				startedCh <- true
				return nil
			},
			shutdownFunc: func() error {
				shutdownCh <- true
				return nil
			},
		}

		assembler.Register(mock)

		err := assembler.Assemble()
		require.NoError(t, err)

		select {
		case <-preparedCh:
			// Success
		case <-time.After(time.Second):
			assert.Fail(t, "Prepare method was not called")
		}

		select {
		case <-startedCh:
			// Success
		case <-time.After(time.Second):
			assert.Fail(t, "Start method was not called")
		}

		err = assembler.Shutdown()
		require.NoError(t, err)

		select {
		case <-shutdownCh:
			// Success
		case <-time.After(time.Second):
			assert.Fail(t, "Shutdown method was not called")
		}
	})
}

// MockServiceAssembly implements ServiceAssembly interface for testing
type MockServiceAssembly struct {
	name         string
	provides     []ServiceType
	requires     []ServiceType
	initFunc     func(*ServiceRegistry) error
	prepareFunc  func() error
	startFunc    func() error
	shutdownFunc func() error
	destroyed    bool
}

func (m *MockServiceAssembly) Name() string            { return m.name }
func (m *MockServiceAssembly) Provides() []ServiceType { return m.provides }
func (m *MockServiceAssembly) Requires() []ServiceType { return m.requires }
func (m *MockServiceAssembly) Init(ctx *InitContext) error {
	if m.initFunc != nil {
		return m.initFunc(ctx.Registry)
	}
	return nil
}
func (m *MockServiceAssembly) Finalize() error {
	m.destroyed = true
	return nil
}

func (m *MockServiceAssembly) Prepare() error {
	if m.prepareFunc != nil {
		return m.prepareFunc()
	}
	return nil
}

func (m *MockServiceAssembly) Start() error {
	if m.startFunc != nil {
		return m.startFunc()
	}
	return nil
}

func (m *MockServiceAssembly) Shutdown() error {
	if m.shutdownFunc != nil {
		return m.shutdownFunc()
	}
	return nil
}
