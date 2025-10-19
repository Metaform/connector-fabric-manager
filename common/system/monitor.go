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

// LogMonitor is a sink for sending log messages to a destination where they can be monitored.
type LogMonitor interface {
	Named(name string) LogMonitor

	Severef(message string, args ...any)
	Warnf(message string, args ...any)
	Infof(message string, args ...any)
	Debugf(message string, args ...any)

	Severew(message string, keyValues ...any)
	Warnw(message string, keyValues ...any)
	Infow(message string, keyValues ...any)
	Debugw(message string, keyValues ...any)

	Sync() error
}

type NoopMonitor struct{}

func (n NoopMonitor) Named(name string) LogMonitor {
	return n
}

func (n NoopMonitor) Severef(message string, args ...any) {
}

func (n NoopMonitor) Warnf(message string, args ...any) {
}

func (n NoopMonitor) Infof(message string, args ...any) {
}

func (n NoopMonitor) Debugf(message string, args ...any) {
}

func (n NoopMonitor) Severew(message string, keyValues ...any) {
}

func (n NoopMonitor) Warnw(message string, keyValues ...any) {
}

func (n NoopMonitor) Infow(message string, keyValues ...any) {
}

func (n NoopMonitor) Debugw(message string, keyValues ...any) {
}

func (n NoopMonitor) Sync() error {
	return nil
}
