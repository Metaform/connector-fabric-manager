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

package model

type Tenant struct {
	ID                  string
	ParticipantContexts []ParticipantContext
}

type ParticipantContext struct {
	DID         string
	DataSpaceId string
}

type Dataspace struct {
	ID string
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
