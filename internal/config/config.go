// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/viper"
)

var LoadedConfiguration configuration

func init() {
	setDefaults()
	loadConfiguration()
}

type configuration struct {
	LogLevel                        string         `mapstructure:"logLevel"`
	Port                            int            `mapstructure:"port"`
	ResponseCode                    int            `mapstructure:"responseCode"`
	MethodResponseCodeMapping       []string       `mapstructure:"methodResponseCodeMapping"`
	RequestLogging                  bool           `mapstructure:"requestLogging"`
	ReadBufferSize                  int            `mapstructure:"readBufferSize"`
	StoreKeyRequestHeaders          []string       `mapstructure:"storeKeyRequestHeaders"`
	SlowlorisDefaultDurationSeconds int            `mapstructure:"slowlorisDefaultDurationSeconds"`
	SlowlorisDefaultIntervalSeconds int            `mapstructure:"slowlorisDefaultIntervalSeconds"`
	MethodResponseCodeMap           map[string]int `mapstructure:"-"`
}

func setDefaults() {
	viper.SetDefault("logLevel", "info")
	viper.SetDefault("port", 8080)
	viper.SetDefault("responseCode", 200)
	viper.SetDefault("methodResponseCodeMapping", []string{})
	viper.SetDefault("requestLogging", true)
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
			log.Info("configuration not found but environment variables will be taken into account.")
		}
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&LoadedConfiguration); err != nil {
		panic(err)
	}

	LoadedConfiguration.BuildMethodResponseCodeMap()
	log.SetLevel(parseLogLevel(LoadedConfiguration.LogLevel))
}

func (c *configuration) BuildMethodResponseCodeMap() {
	c.MethodResponseCodeMap = make(map[string]int, len(c.MethodResponseCodeMapping))
	for _, m := range c.MethodResponseCodeMapping {
		parts := strings.SplitN(m, ":", 2)
		if len(parts) != 2 {
			continue
		}
		method := strings.ToUpper(strings.TrimSpace(parts[0]))
		code, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			log.Warnf("ignoring invalid method response code mapping: %s", m)
			continue
		}
		c.MethodResponseCodeMap[method] = code
	}
}

func parseLogLevel(lvl string) log.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return log.LevelDebug
	case "info":
		return log.LevelInfo
	case "warn":
		return log.LevelWarn
	case "error":
		return log.LevelError
	default:
		log.Warnf("invalid log-level '%s', falling back to info", lvl)
		return log.LevelInfo
	}
}
