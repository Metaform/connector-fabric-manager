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

package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

// LoadConfigOrPanic initializes a Config instance with the specified configuration name.
// Configuration will be read from a file (if it exists) and can be overridden using environment variables.
func LoadConfigOrPanic(name string) *viper.Viper {
	v := viper.New()
	v.SetConfigName(name)
	v.AddConfigPath("/etc/appname/")
	v.AddConfigPath("$HOME/.appname")
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetEnvPrefix(name)
	err := v.ReadInConfig()

	if err != nil {
		// ignore not found error, otherwise panic
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("error reading config file: %w", err))
		}
	}
	return v
}
