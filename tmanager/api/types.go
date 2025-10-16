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
	"time"

	"github.com/metaform/connector-fabric-manager/common/dmodel"
)

// Entity is the base type for all entities.
type Entity struct {
	ID      string
	Version int64
}

// DeployableEntity is an entity that can be deployed and follows a deployment lifecycle.
// It extends Entity with state tracking capabilities for managing deployment phases.
type DeployableEntity struct {
	Entity
	State          DeploymentState
	StateTimestamp time.Time
}

// Tenant represents an organization. A tenant may have multiple organizational units (e.g., departments), or ParticipantProfiles.
type Tenant struct {
	Entity
	ParticipantProfiles []ParticipantProfile
	Properties          Properties
}

// ParticipantProfile represents a participant in a dataspace. A participant can be an entire organization, in which case
// the tenant will have a single ParticipantProfile. Or, an organization can have a participant per organizational unit
// (e.g., department).
type ParticipantProfile struct {
	Entity
	Identifier        string
	DataSpaceProfiles []DataspaceProfile
	VPAs              []VirtualParticipantAgent
	Properties        Properties
}

// DataspaceProfile represents a specific dataspace, protocol, and policies tuple. For example, The Foo Dataspace that
// runs version 2025-1 with version 2 of its policies schema.
type DataspaceProfile struct {
	Entity
	Artifacts   []string
	Deployments []DataspaceDeployment
	Properties  Properties
}

// VirtualParticipantAgent is a runtime context deployed when a participant profile is provisioned to a cell. A runtime
// context could be a connector, credential service, or another component.
type VirtualParticipantAgent struct {
	DeployableEntity
	Type       dmodel.VPAType
	Cell       Cell
	Properties Properties
}

// DataspaceDeployment is runtime capabilities and configuration deployed when a dataspace profile to a cell.
type DataspaceDeployment struct {
	DeployableEntity
	Cell       Cell
	Properties Properties
}

// Cell is a homogenous deployment zone. A cell could be a Kubernetes cluster or some other infrastructure.
type Cell struct {
	DeployableEntity
	Properties Properties
}

// DeploymentState represents the current state of a deployable entity
type DeploymentState string

const (
	DeploymentStateInitial DeploymentState = "initial"
	DeploymentStatePending DeploymentState = "pending"
	DeploymentStateActive  DeploymentState = "active"
	DeploymentStateLocked  DeploymentState = "locked"
	DeploymentStateOffline DeploymentState = "offline"
	DeploymentStateError   DeploymentState = "error"
)

// String implements the Stringer interface
func (c DeploymentState) String() string {
	return string(c)
}

// IsValid validates the enum value
func (c DeploymentState) IsValid() bool {
	switch c {
	case DeploymentStateInitial, DeploymentStatePending, DeploymentStateActive, DeploymentStateOffline, DeploymentStateError, DeploymentStateLocked:
		return true
	default:
		return false
	}
}

// MarshalJSON implements json.Marshaler
func (c DeploymentState) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(c))
}

// UnmarshalJSON implements json.Unmarshaler
func (c *DeploymentState) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	state := DeploymentState(s)
	if !state.IsValid() {
		return fmt.Errorf("invalid cell state: %s", s)
	}

	*c = state
	return nil
}

// Value implements the driver.Valuer interface for database serialization
func (c DeploymentState) Value() (driver.Value, error) {
	if !c.IsValid() {
		return nil, fmt.Errorf("invalid cell state: %s", c)
	}
	return string(c), nil
}

// Scan implements the sql.Scanner interface for database deserialization
func (c *DeploymentState) Scan(value interface{}) error {
	if value == nil {
		*c = ""
		return nil
	}

	switch v := value.(type) {
	case string:
		*c = DeploymentState(v)
	case []byte:
		*c = DeploymentState(v)
	default:
		return fmt.Errorf("cannot scan %T into DeploymentState", value)
	}

	if !c.IsValid() {
		return fmt.Errorf("invalid cell state: %s", *c)
	}

	return nil
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
