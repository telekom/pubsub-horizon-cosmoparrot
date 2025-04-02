// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"cosmoparrot/internal/cache"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	c "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupTestApp() *fiber.App {
	app := fiber.New()
	app.Get("/api/v1/requests", handleGetAllRequests)
	app.Get("/api/v1/requests/:key", handleGetRequestByKey)
	return app
}

func TestHandleGetAllRequestsx(t *testing.T) {
	cache.Current = c.New(5*time.Minute, 10*time.Minute)

	requests := []*request{{Time: time.Now().Add(-time.Minute)}, {Time: time.Now()}}
	data, _ := json.Marshal(requests)
	cache.Current.Set("test-key", string(data), c.DefaultExpiration)

	app := setupTestApp()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/requests", nil)

	// Perform the request
	resp, err := app.Test(r, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	assert.Contains(t, string(respBytes), "time")
}

func TestHandleGetRequestByKeyx(t *testing.T) {
	cache.Current = c.New(5*time.Minute, 10*time.Minute)

	requests := []*request{{Time: time.Now()}}
	data, _ := json.Marshal(requests)
	cache.Current.Set("test-key", string(data), c.DefaultExpiration)

	app := setupTestApp()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/requests/test-key", nil)

	// Perform the request
	resp, err := app.Test(r, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	assert.Contains(t, string(respBytes), "time")
}

func TestHandleGetRequestByKey_NotFound(t *testing.T) {
	cache.Current = c.New(5*time.Minute, 10*time.Minute)

	app := setupTestApp()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/requests/non-existent", nil)

	// Perform the request
	resp, err := app.Test(r, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
