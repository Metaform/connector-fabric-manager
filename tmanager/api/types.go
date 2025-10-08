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

package api

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Tenant struct {
	ID                  string
	ParticipantProfiles []ParticipantProfile
	Properties          Properties
}

type ParticipantProfile struct {
	Identifier       string
	DataSpaceProfile DataspaceProfile
	VPAs             []VirtualParticipantAgent
	Properties       Properties
}

type DataspaceProfile struct {
	ID         string
	Properties Properties
}

type VirtualParticipantAgent struct {
	ID         string
	Cell       Cell
	Properties Properties
}

type Cell struct {
	ID         string
	Properties Properties
}

type User struct {
	Roles []Role
}

type Role struct {
	Rights []Right
}

type Right interface {
	GetDescription() string
}

// Properties are extensible key-value pairs
type Properties map[string]any

// Value implements the driver.Valuer interface for database serialization
func (p *Properties) Value() (driver.Value, error) {
	if p == nil || *p == nil || len(*p) == 0 {
		return nil, nil
	}
	return json.Marshal(*p)
}

// Scan implements the sql.Scanner interface for database deserialization
func (p *Properties) Scan(value any) error {
	if value == nil {
		*p = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into Properties", value)
	}

	if len(bytes) == 0 {
		*p = make(Properties)
		return nil
	}

	return json.Unmarshal(bytes, p)
}

// Helper methods for common operations
func (p *Properties) Get(key string) (any, bool) {
	if p == nil || *p == nil {
		return nil, false
	}
	value, exists := (*p)[key]
	return value, exists
}

func (p *Properties) GetString(key string) (string, bool) {
	if value, exists := p.Get(key); exists {
		if str, ok := value.(string); ok {
			return str, true
		}
	}
	return "", false
}

func (p *Properties) GetInt(key string) (int, bool) {
	if value, exists := p.Get(key); exists {
		switch v := value.(type) {
		case int:
			return v, true
		case float64:
			return int(v), true
		}
	}
	return 0, false
}

func (p *Properties) Set(key string, value any) {
	if *p == nil {
		*p = make(Properties)
	}
	(*p)[key] = value
}
