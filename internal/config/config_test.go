// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSetDefaults(t *testing.T) {
	// Reset Viper settings before testing
	viper.Reset()
	setDefaults()

	assert.Equal(t, 8080, viper.GetInt("port"))
	assert.Equal(t, 200, viper.GetInt("responseCode"))
	assert.Equal(t, []string{}, viper.GetStringSlice("methodResponseCodeMapping"))
	assert.Equal(t, []string{"x-request-key"}, viper.GetStringSlice("storeKeyRequestHeaders"))
}

func TestEnvironmentVariableOverride(t *testing.T) {
	// Reset Viper to ensure clean state
	viper.Reset()
	setDefaults()

	// Set environment variables
	_ = os.Setenv("COSMOPARROT_PORT", "9090")
	_ = os.Setenv("COSMOPARROT_RESPONSECODE", "500")
	_ = os.Setenv("COSMOPARROT_STOREKEYREQUESTHEADERS", "x-custom-header")

	// Reload configuration
	loadConfiguration()

	// Assert that environment variables override defaults
	assert.Equal(t, 9090, LoadedConfiguration.Port)
	assert.Equal(t, 500, LoadedConfiguration.ResponseCode)
	assert.Equal(t, []string{"x-custom-header"}, LoadedConfiguration.StoreKeyRequestHeaders)

	// Cleanup
	_ = os.Unsetenv("COSMOPARROT_PORT")
	_ = os.Unsetenv("COSMOPARROT_RESPONSECODE")
	_ = os.Unsetenv("COSMOPARROT_STOREKEYREQUESTHEADERS")
}

func TestConfigFileNotFound(t *testing.T) {
	// Reset Viper to prevent contamination
	viper.Reset()
	setDefaults()

	// Run without config file
	loadConfiguration()

	// Ensure defaults are still applied
	assert.Equal(t, 8080, LoadedConfiguration.Port)
	assert.Equal(t, 200, LoadedConfiguration.ResponseCode)
	assert.Equal(t, []string{"x-request-key"}, LoadedConfiguration.StoreKeyRequestHeaders)
}
