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

package activity

import (
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

type KeyCloakActivityProcessor struct {
	Monitor system.LogMonitor
}

func (p KeyCloakActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	return api.ActivityResult{Result: api.ActivityResultComplete}
}
