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

package e2etests

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/metaform/connector-fabric-manager/common/natstestfixtures"
	"github.com/metaform/connector-fabric-manager/common/testfixtures"
	"github.com/metaform/connector-fabric-manager/e2e/e2efixtures"
	alauncher "github.com/metaform/connector-fabric-manager/e2e/testagent/launcher"
	plauncher "github.com/metaform/connector-fabric-manager/pmanager/cmd/server/launcher"
	tlauncher "github.com/metaform/connector-fabric-manager/tmanager/cmd/server/launcher"
)

const (
	testTimeout = 30 * time.Second
	streamName  = "cfm-stream"
	cfmBucket   = "cfm-bucket"
)

func launchPlatform(t *testing.T, nt *natstestfixtures.NatsTestContainer) *e2efixtures.ApiClient {
	_ = os.Setenv("TM_URI", nt.Uri)
	_ = os.Setenv("TM_BUCKET", cfmBucket)
	_ = os.Setenv("TM_STREAM", streamName)

	_ = os.Setenv("PM_URI", nt.Uri)
	_ = os.Setenv("PM_BUCKET", cfmBucket)
	_ = os.Setenv("PM_STREAM", streamName)

	_ = os.Setenv("TESTAGENT_URI", nt.Uri)
	_ = os.Setenv("TESTAGENT_BUCKET", cfmBucket)
	_ = os.Setenv("TESTAGENT_STREAM", streamName)

	tPort := testfixtures.GetRandomPort(t)
	_ = os.Setenv("TM_HTTPPORT", strconv.Itoa(tPort))
	pPort := testfixtures.GetRandomPort(t)
	_ = os.Setenv("PM_HTTPPORT", strconv.Itoa(pPort))

	shutdownChannel := make(chan struct{})
	go func() {
		plauncher.Launch(shutdownChannel)
	}()

	go func() {
		tlauncher.Launch(shutdownChannel)
	}()

	go func() {
		alauncher.Launch(shutdownChannel)
	}()
	return e2efixtures.NewApiClient(fmt.Sprintf("http://localhost:%d", tPort), fmt.Sprintf("http://localhost:%d", pPort))
}

func cleanup() {
	_ = os.Unsetenv("TM_URI")
	_ = os.Unsetenv("TM_BUCKET")
	_ = os.Unsetenv("TM_STREAM")

	_ = os.Unsetenv("PM_URI")
	_ = os.Unsetenv("PM_BUCKET")
	_ = os.Unsetenv("PM_STREAM")

	_ = os.Unsetenv("TESTAGENT_URI")
	_ = os.Unsetenv("TESTAGENT_BUCKET")
	_ = os.Unsetenv("TESTAGENT_STREAM")

	_ = os.Unsetenv("TM_HTTPPORT")
	_ = os.Unsetenv("PM_HTTPPORT")
}
