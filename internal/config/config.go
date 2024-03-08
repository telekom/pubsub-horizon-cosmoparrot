// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
)

var LoadedConfiguration configuration

func init() {
	setDefaults()
	loadConfiguration()
}

type configuration struct {
	Port         int `mapstructure:"port"`
	ResponseCode int `mapstructure:"responseCode"`
}

func setDefaults() {
	viper.SetDefault("port", 8080)
	viper.SetDefault("responseCode", 200)
}

func loadConfiguration() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("COSMOPARROT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Info().Msg("Configuration not found but environment variables will be taken into account.")
		}
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&LoadedConfiguration); err != nil {
		panic(err)
	}
}
