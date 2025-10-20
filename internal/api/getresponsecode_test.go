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

// New test: case-insensitive and variant recognition for responseCode
func TestGetResponseCode_CaseInsensitiveVariants(t *testing.T) {
	config.LoadedConfiguration.ResponseCode = 200
	config.LoadedConfiguration.MethodResponseCodeMapping = []string{"GET:202"}

	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		code := getResponseCode(c)
		return c.SendStatus(code)
	})

	// Non-exact variants should NOT be recognized anymore and therefore fall back to mapping (202).
	cases := []struct {
		url  string
		want int
	}{
		{"/test?response_code=205", 202},
		{"/test?RESPONSECODE=206", 202},
		{"/test?Response-Code=207", 202},
		// exact match should still work
		{"/test?responseCode=205", 205},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, tc.url, nil)
		resp, _ := app.Test(req)
		assert.Equal(t, tc.want, resp.StatusCode, "url: %s", tc.url)
	}
}
