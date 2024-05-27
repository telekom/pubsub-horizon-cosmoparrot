// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var LoadedConfiguration configuration

func init() {
	setDefaults()
	loadConfiguration()
}

type configuration struct {
	Port                            int      `mapstructure:"port"`
	ResponseCode                    int      `mapstructure:"responseCode"`
	MethodResponseCodeMapping       []string `mapstructure:"methodResponseCodeMapping"`
	ReadBufferSize                  int      `mapstructure:"readBufferSize"`
	StoreKeyRequestHeaders          []string `mapstructure:"storeKeyRequestHeaders"`
	SlowlorisDefaultDurationSeconds int      `mapstructure:"slowlorisDefaultDurationSeconds"`
	SlowlorisDefaultIntervalSeconds int      `mapstructure:"slowlorisDefaultIntervalSeconds"`
}

func setDefaults() {
	viper.SetDefault("port", 8080)
	viper.SetDefault("responseCode", 200)
	viper.SetDefault("methodResponseCodeMapping", []string{})
	viper.SetDefault("readBufferSize", 4096)
	viper.SetDefault("storeKeyRequestHeaders", []string{"x-request-key"})
	viper.SetDefault("slowlorisDefaultDurationSeconds", 15)
	viper.SetDefault("slowlorisDefaultIntervalSeconds", 1)
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
