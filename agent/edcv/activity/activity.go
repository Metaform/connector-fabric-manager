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
	"fmt"

	"github.com/metaform/connector-fabric-manager/common/model"
	"github.com/metaform/connector-fabric-manager/common/system"
	"github.com/metaform/connector-fabric-manager/pmanager/api"
)

type EDCVActivityProcessor struct {
	monitor system.LogMonitor
}

func (t EDCVActivityProcessor) Process(ctx api.ActivityContext) api.ActivityResult {
	_, found := ctx.InputData().Get(model.ParticipantIdentifier)
	if !found {
		return api.ActivityResult{Result: api.ActivityResultFatalError, Error: fmt.Errorf("missing participant identifier")}
	}
	return api.ActivityResult{Result: api.ActivityResultComplete}
}
