// Copyright (c) 2025 Metaform Systems, Inc
//
// This program and the accompanying materials are made available under the
// terms of the Apache License, Version 2.0 which is available at
// https://www.apache.org/licenses/LICENSE-2.0
//
// SPDX-License-Identifier: Apache-2.0
//
// Contributors:
//
//	Metaform Systems, Inc. - initial API and implementation

package httpclient

import (
	"context"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/metaform/connector-fabric-manager/common/system"
	"net/http"
	"time"
)

const (
	HttpClientKey         system.ServiceType = "httpclient:HttpClient"
	ConfigKeyRetryMax     string             = "httpclient.retrymax"
	DefaultRetryMax       int                = 5
	ConfigKeyRetryWaitMin string             = "httpclient.retrywaitmin"
	DefaultRetryWaitMin   int                = 1
	ConfigKeyRetryWaitMax string             = "httpclient.retrywaitmax"
	DefaultRetryWaitMax   int                = 5
)

type HttpClientServiceAssembly struct {
	system.DefaultServiceAssembly
}

func (h HttpClientServiceAssembly) Name() string {
	return "HttpClient"
}

func (h HttpClientServiceAssembly) Provides() []system.ServiceType {
	return []system.ServiceType{HttpClientKey}
}

func (h HttpClientServiceAssembly) Requires() []system.ServiceType {
	return []system.ServiceType{}
}

func (h HttpClientServiceAssembly) Init(context *system.InitContext) error {
	retryClient := retryablehttp.NewClient()

	retryClient.RetryMax = context.GetConfigIntOrDefault(ConfigKeyRetryMax, DefaultRetryMax)
	retryClient.RetryWaitMin = time.Duration(context.GetConfigIntOrDefault(ConfigKeyRetryWaitMin, DefaultRetryWaitMin)) * time.Second
	retryClient.RetryWaitMax = time.Duration(context.GetConfigIntOrDefault(ConfigKeyRetryWaitMax, DefaultRetryWaitMax)) * time.Second
	retryClient.CheckRetry = customCheckRetry

	standardClient := retryClient.StandardClient()
	context.Registry.Register(HttpClientKey, *standardClient)

	return nil
}

// customCheckRetry determines whether a request should be retried
// It will not retry on 4xx client errors (bad requests)
func customCheckRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// Don't retry if the context is canceled or timed out
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	// Don't retry on 4xx client errors (bad requests)
	if resp != nil && resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return false, nil
	}

	// Use the default retry logic for other cases
	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}
