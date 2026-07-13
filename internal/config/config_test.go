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
	assert.Equal(t, true, viper.GetBool("requestLogging"))
	assert.Equal(t, []string{"x-request-key"}, viper.GetStringSlice("storeKeyRequestHeaders"))
	assert.Equal(t, false, viper.GetBool("otelEnabled"))
	assert.Equal(t, "cosmoparrot", viper.GetString("otelServiceName"))
}

func TestRequestLoggingEnvironmentOverride(t *testing.T) {
	// Reset Viper to ensure clean state
	viper.Reset()
	setDefaults()

	os.Setenv("COSMOPARROT_REQUESTLOGGING", "false")

	loadConfiguration()

	assert.Equal(t, false, LoadedConfiguration.RequestLogging)

	// Cleanup
	os.Unsetenv("COSMOPARROT_REQUESTLOGGING")
}

func TestEnvironmentVariableOverride(t *testing.T) {
	// Reset Viper to ensure clean state
	viper.Reset()
	setDefaults()

	// Set environment variables
	os.Setenv("COSMOPARROT_PORT", "9090")
	os.Setenv("COSMOPARROT_RESPONSECODE", "500")
	os.Setenv("COSMOPARROT_STOREKEYREQUESTHEADERS", "x-custom-header")
	os.Setenv("COSMOPARROT_OTELENABLED", "true")
	os.Setenv("COSMOPARROT_OTELSERVICENAME", "cosmoparrot-tests")

	// Reload configuration
	loadConfiguration()

	// Assert that environment variables override defaults
	assert.Equal(t, 9090, LoadedConfiguration.Port)
	assert.Equal(t, 500, LoadedConfiguration.ResponseCode)
	assert.Equal(t, []string{"x-custom-header"}, LoadedConfiguration.StoreKeyRequestHeaders)
	assert.Equal(t, true, LoadedConfiguration.OTelEnabled)
	assert.Equal(t, "cosmoparrot-tests", LoadedConfiguration.OTelServiceName)

	// Cleanup
	os.Unsetenv("COSMOPARROT_PORT")
	os.Unsetenv("COSMOPARROT_RESPONSECODE")
	os.Unsetenv("COSMOPARROT_STOREKEYREQUESTHEADERS")
	os.Unsetenv("COSMOPARROT_OTELENABLED")
	os.Unsetenv("COSMOPARROT_OTELSERVICENAME")
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
