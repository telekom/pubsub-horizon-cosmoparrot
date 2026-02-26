// Copyright 2026 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetResponseDelay_ValidValue(t *testing.T) {
	app := fiber.New()
	var result time.Duration
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseDelay(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseDelay=500", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 500*time.Millisecond, result)
}

func TestGetResponseDelay_MissingParam(t *testing.T) {
	app := fiber.New()
	var result time.Duration
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseDelay(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, time.Duration(0), result)
}

func TestGetResponseDelay_NonInteger(t *testing.T) {
	app := fiber.New()
	var result time.Duration
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseDelay(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseDelay=abc", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, time.Duration(0), result)
}

func TestGetResponseDelay_NegativeValue(t *testing.T) {
	app := fiber.New()
	var result time.Duration
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseDelay(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseDelay=-100", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, time.Duration(0), result)
}

func TestGetResponseDelay_ExceedsMax(t *testing.T) {
	app := fiber.New()
	var result time.Duration
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseDelay(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseDelay=60001", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, time.Duration(0), result)
}

func TestGetResponseDelay_BoundaryMax(t *testing.T) {
	app := fiber.New()
	var result time.Duration
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseDelay(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseDelay=60000", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 60*time.Second, result)
}

func TestGetResponseDelay_Zero(t *testing.T) {
	app := fiber.New()
	var result time.Duration
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseDelay(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseDelay=0", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, time.Duration(0), result)
}

// Non-exact param names should NOT be recognized
func TestGetResponseDelay_CaseInsensitiveVariants(t *testing.T) {
	app := fiber.New()
	var result time.Duration
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseDelay(c)
		return c.SendStatus(http.StatusOK)
	})

	cases := []struct {
		url  string
		want time.Duration
	}{
		{"/test?response_delay=500", 0},
		{"/test?RESPONSEDELAY=500", 0},
		{"/test?Response-Delay=500", 0},
		// exact match should work
		{"/test?responseDelay=500", 500 * time.Millisecond},
	}

	for _, tc := range cases {
		result = 0
		req := httptest.NewRequest(http.MethodGet, tc.url, nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "url: %s", tc.url)
		assert.Equal(t, tc.want, result, "url: %s", tc.url)
	}
}
