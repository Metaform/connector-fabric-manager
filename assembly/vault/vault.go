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
	"strings"
	"time"

	hvault "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/metaform/connector-fabric-manager/common/system"
)

const (
	contentKey = "content"
)

type VaultOptions func(*vaultClient)

func (f VaultOptions) apply(vc *vaultClient) {
	f(vc)
}

// WithMountPath sets the mount path
func WithMountPath(path string) VaultOptions {
	return func(vc *vaultClient) {
		vc.mountPath = path
	}
}

// vaultClient implements a client to Hashicorp Vault supporting token renewal.
type vaultClient struct {
	vaultURL    string
	roleID      string
	secretID    string
	mountPath   string
	monitor     system.LogMonitor
	client      *hvault.Client
	stopCh      chan struct{}
	lastCreated time.Time // When the token was last renewed; will be the zero value if the token has never been renewed or there was an error.
	lastRenew   time.Time // When the token was last renewed; will be the zero value if the token has never been renewed or there was an error.
}

func newVaultClient(vaultURL string, roleID string, secretID string, monitor system.LogMonitor, opts ...VaultOptions) (*vaultClient, error) {
	client := &vaultClient{
		vaultURL: vaultURL,
		roleID:   roleID,
		secretID: secretID,
		monitor:  monitor,
		stopCh:   make(chan struct{}),
	}
	for _, opt := range opts {
		opt.apply(client)
	}
	return client, nil
}

func (v *vaultClient) ResolveSecret(ctx context.Context, path string) (string, error) {
	secret, err := v.client.Secrets.KvV2Read(
		context.Background(),
		path,
		v.getOptions()...,
	)
	if err != nil {
		return "", fmt.Errorf("unable to read secret: %w", err)
	}
	if value, ok := secret.Data.Data["content"].(string); ok {
		return value, nil
	}
	return "", fmt.Errorf("content field not found or not a string")

}

func (v *vaultClient) StoreSecret(ctx context.Context, path string, value string) error {
	_, err := v.client.Secrets.KvV2Write(
		ctx,
		path,
		schema.KvV2WriteRequest{
			Data: map[string]any{
				contentKey: value,
			},
		},
		v.getOptions()...,
	)
	if err != nil {
		return fmt.Errorf("unable to write secret to path %s: %w", path, err)
	}

	return nil
}

func (v *vaultClient) DeleteSecret(ctx context.Context, path string) error {
	_, err := v.client.Secrets.KvV2Delete(
		context.Background(),
		path,
		v.getOptions()...,
	)
	if err != nil {
		return fmt.Errorf("unable to delete secret at path %s: %w", path, err)
	}
	return nil
}

func (v *vaultClient) init(ctx context.Context) error {
	var err error
	v.client, err = hvault.New(
		hvault.WithAddress(v.vaultURL),
		hvault.WithRequestTimeout(10*time.Second), // TODO configure
	)
	if err != nil {
		return fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	// Authenticate using AppRole
	resp, err := v.createToken(ctx)
	if err != nil {
		return fmt.Errorf("unable to initialize Vault client: %w", err)
	}
	go v.renewTokenPeriodically(time.Duration(resp.Auth.LeaseDuration) * time.Second)
	return nil
}

func (v *vaultClient) createToken(ctx context.Context) (*hvault.Response[map[string]any], error) {
	appRoleResp, err := v.client.Auth.AppRoleLogin(
		ctx,
		schema.AppRoleLoginRequest{
			RoleId:   v.roleID,
			SecretId: v.secretID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate with AppRole: %w", err)
	}

	// Set the token obtained from AppRole login
	err = v.client.SetToken(appRoleResp.Auth.ClientToken)
	v.lastCreated = time.Now()
	return appRoleResp, err
}

// leaseDuration specifies the token lease duration and supports any time.Duration unit (milliseconds, seconds, minutes, etc.)
func (v *vaultClient) renewTokenPeriodically(leaseDuration time.Duration) {
	// Renew at 80% of lease duration
	renewInterval := time.Duration(float64(leaseDuration) * 0.8)

	ticker := time.NewTicker(renewInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := v.client.Auth.TokenRenewSelf(context.Background(), schema.TokenRenewSelfRequest{
				Increment: fmt.Sprintf("%ds", int(leaseDuration.Seconds())),
			})
			if err != nil {
				if strings.Contains(err.Error(), "invalid token") {
					// Token cannot be renewed further because it has expired so create a new one
					_, err2 := v.createToken(context.Background())
					if err2 != nil {
						v.monitor.Severef("Error creating token after expiration: %v. Will attempt renewal at next interval", err2)
						continue
					}
				}
				v.lastRenew = time.Time{}
				v.monitor.Severef("Error renewing token: %v. Will attempt renewal at next interval", err)
				continue
			}
			v.lastRenew = time.Now()
		case <-v.stopCh:
			return
		}
	}
}

// Close gracefully shuts down the vaultClient and stops token renewal
func (v *vaultClient) Close() error {
	close(v.stopCh)
	return nil
}

func (v *vaultClient) getOptions() []hvault.RequestOption {
	var opts []hvault.RequestOption
	if v.mountPath != "" {
		opts = append(opts, hvault.WithMountPath(v.mountPath))
	}
	return opts
}
