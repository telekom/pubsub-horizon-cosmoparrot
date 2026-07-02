// Copyright 2026 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetResponseSize_ValidValue(t *testing.T) {
	app := fiber.New()
	var result int
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseSize(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseSize=500", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 500, result)
}

func TestGetResponseSize_MissingParam(t *testing.T) {
	app := fiber.New()
	var result int
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseSize(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 0, result)
}

func TestGetResponseSize_NonInteger(t *testing.T) {
	app := fiber.New()
	var result int
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseSize(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseSize=abc", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 0, result)
}

func TestGetResponseSize_NegativeValue(t *testing.T) {
	app := fiber.New()
	var result int
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseSize(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseSize=-100", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 0, result)
}

func TestGetResponseSize_ExceedsMax(t *testing.T) {
	app := fiber.New()
	var result int
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseSize(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseSize=10000001", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 0, result)
}

func TestGetResponseSize_BoundaryMax(t *testing.T) {
	app := fiber.New()
	var result int
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseSize(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseSize=10000000", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, maxResponseSize, result)
}

func TestGetResponseSize_Zero(t *testing.T) {
	app := fiber.New()
	var result int
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseSize(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?responseSize=0", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 0, result)
}

// Case sensitivity should be ignored for query parameters
// "-" and "_" should NOT be recognized
func TestGetResponseSize_CaseInsensitiveVariants(t *testing.T) {
	app := fiber.New()
	var result int
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getResponseSize(c)
		return c.SendStatus(http.StatusOK)
	})

	cases := []struct {
		url  string
		want int
	}{
		{"/test?response_size=500", 0},
		{"/test?Response-Size=500", 0},
		// exact match should work
		{"/test?reSpoNsesIze=500", 500},
		{"/test?responsesize=500", 500},
		{"/test?RESPONSESIZE=500", 500},
		{"/test?responseSize=500", 500},
	}

	for _, tc := range cases {
		result = 0
		req := httptest.NewRequest(http.MethodGet, tc.url, nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "url: %s", tc.url)
		assert.Equal(t, tc.want, result, "url: %s", tc.url)
	}
}
