// Copyright 2025 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"cosmoparrot/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetResponseCode_QueryParamPrecedence(t *testing.T) {
	// ensure default
	config.LoadedConfiguration.ResponseCode = 201
	config.LoadedConfiguration.MethodResponseCodeMapping = []string{"GET:202", "POST:203"}
	config.LoadedConfiguration.BuildMethodResponseCodeMap()

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		code := getResponseCode(c)
		return c.SendStatus(code)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseCode=204", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestGetResponseCode_MappingFallback(t *testing.T) {
	config.LoadedConfiguration.ResponseCode = 200
	config.LoadedConfiguration.MethodResponseCodeMapping = []string{"GET:202", "POST:203"}
	config.LoadedConfiguration.BuildMethodResponseCodeMap()

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		code := getResponseCode(c)
		return c.SendStatus(code)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 202, resp.StatusCode)
}

func TestGetResponseCode_InvalidQueryParamFallsback(t *testing.T) {
	config.LoadedConfiguration.ResponseCode = 299
	config.LoadedConfiguration.MethodResponseCodeMapping = []string{"GET:202"}
	config.LoadedConfiguration.BuildMethodResponseCodeMap()

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		code := getResponseCode(c)
		return c.SendStatus(code)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseCode=notanint", nil)
	resp, _ := app.Test(req)
	// should fallback to mapping (202)
	assert.Equal(t, 202, resp.StatusCode)
}

func TestGetResponseCode_OutOfRangeQueryParamFallsback(t *testing.T) {
	config.LoadedConfiguration.ResponseCode = 299
	config.LoadedConfiguration.MethodResponseCodeMapping = []string{"GET:202"}
	config.LoadedConfiguration.BuildMethodResponseCodeMap()

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		code := getResponseCode(c)
		return c.SendStatus(code)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseCode=700", nil)
	resp, _ := app.Test(req)
	// should fallback to mapping (202)
	assert.Equal(t, 202, resp.StatusCode)
}

// Case sensitivity should be ignored for query parameters
// "-" and "_" should NOT be recognized
func TestGetResponseCode_CaseInsensitiveVariants(t *testing.T) {
	config.LoadedConfiguration.ResponseCode = 200
	config.LoadedConfiguration.MethodResponseCodeMapping = []string{"GET:202"}
	config.LoadedConfiguration.BuildMethodResponseCodeMap()

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		code := getResponseCode(c)
		return c.SendStatus(code)
	})

	cases := []struct {
		url  string
		want int
	}{
		{"/test?response_code=205", 202},
		{"/test?Response-Code=206", 202},
		// exact match should still work
		{"/test?responsecode=207", 207},
		{"/test?rEspoNsecOde=208", 208},
		{"/test?RESPONSECODE=209", 209},
		{"/test?responseCode=210", 210},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, tc.url, nil)
		resp, _ := app.Test(req)
		assert.Equal(t, tc.want, resp.StatusCode, "url: %s", tc.url)
	}
}
