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

func TestGetMirrorBody_MissingParamDefaultsTrue(t *testing.T) {
	app := fiber.New()
	var result bool
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getMirrorBody(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result)
}

func TestGetMirrorBody_InvalidValueDefaultsTrue(t *testing.T) {
	app := fiber.New()
	var result bool
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getMirrorBody(c)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test?mirrorBody=maybe", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result)
}

func TestGetMirrorBody_ExplicitValues(t *testing.T) {
	app := fiber.New()
	var result bool
	app.Get("/test", func(c *fiber.Ctx) error {
		result = getMirrorBody(c)
		return c.SendStatus(http.StatusOK)
	})

	cases := []struct {
		url  string
		want bool
	}{
		{"/test?mirrorBody=false", false},
		{"/test?mirrorBody=0", false},
		{"/test?mirrorBody=f", false},
		{"/test?mirrorBody=FALSE", false},
		{"/test?mirrorBody=true", true},
		{"/test?mirrorBody=1", true},
		// case-insensitive param name, "-"/"_" variants are NOT recognized
		{"/test?MIRRORBODY=false", false},
		{"/test?mirror_body=false", true},
		{"/test?mirror-body=false", true},
	}

	for _, tc := range cases {
		result = true
		req := httptest.NewRequest(http.MethodGet, tc.url, nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "url: %s", tc.url)
		assert.Equal(t, tc.want, result, "url: %s", tc.url)
	}
}
