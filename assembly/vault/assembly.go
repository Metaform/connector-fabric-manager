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

package vault

import (
	"context"
	"fmt"

	"github.com/metaform/connector-fabric-manager/assembly/serviceapi"
	"github.com/metaform/connector-fabric-manager/common/runtime"
	"github.com/metaform/connector-fabric-manager/common/system"
)

const (
	urlKey      = "vaultUrl"
	roleIDKey   = "vaultRoleId"
	secretIDKey = "vaultSecretId"
)

// VaultServiceAssembly defines an assembly that provides a client to Hashicorp Vault.
type VaultServiceAssembly struct {
	system.DefaultServiceAssembly
	client *vaultClient
}

func (v VaultServiceAssembly) Name() string {
	return "Vault"
}

func (v VaultServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{serviceapi.VaultKey}
}

func (v VaultServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{}
}

func (v VaultServiceAssembly) Init(ctx *system.InitContext) error {
	vaultURL := ctx.Config.GetString(urlKey)
	roleID := ctx.Config.GetString(roleIDKey)
	secretID := ctx.Config.GetString(secretIDKey)
	if err := runtime.CheckRequiredParams(urlKey, vaultURL, roleIDKey, roleID, secretIDKey, secretID); err != nil {
		return err
	}
	var err error
	v.client, err = newVaultClient(vaultURL, roleID, secretID, ctx.LogMonitor)
	if err != nil {
		return fmt.Errorf("failed to create Vault client: %w", err)
	}

	err = v.client.init(context.Background())
	if err != nil {
		return fmt.Errorf("failed to initialize Vault client: %w", err)
	}

	ctx.Registry.Register(serviceapi.VaultKey, v.client)

	return nil
}

func (v VaultServiceAssembly) Shutdown() error {
	if v.client != nil {
		err := v.client.Close()
		if err != nil {
			return fmt.Errorf("failed to close Vault client: %w", err)
		}
	}
	return nil
}
