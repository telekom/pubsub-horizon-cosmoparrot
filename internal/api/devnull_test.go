// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHandleDevNull_POST(t *testing.T) {
	app := fiber.New()
	app.All("/api/v1/devnull", handleDevNull)

	body := []byte(`{"message": "should be discarded"}`)
	r := httptest.NewRequest("POST", "/api/v1/devnull", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(r, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHandleDevNull_GET(t *testing.T) {
	app := fiber.New()
	app.All("/api/v1/devnull", handleDevNull)

	r := httptest.NewRequest("GET", "/api/v1/devnull", nil)

	resp, err := app.Test(r, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHandleDevNull_CustomResponseCode(t *testing.T) {
	app := fiber.New()
	app.All("/api/v1/devnull", handleDevNull)

	r := httptest.NewRequest("POST", "/api/v1/devnull?responseCode=201", bytes.NewReader([]byte(`{}`)))
	r.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(r, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestHandleDevNull_LargeBody(t *testing.T) {
	app := fiber.New()
	app.All("/api/v1/devnull", handleDevNull)

	largeBody := make([]byte, 1024*1024) // 1 MB
	for i := range largeBody {
		largeBody[i] = 'x'
	}

	r := httptest.NewRequest("PUT", "/api/v1/devnull", bytes.NewReader(largeBody))

	resp, err := app.Test(r, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
