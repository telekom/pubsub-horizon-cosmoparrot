// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"bytes"
	"cosmoparrot/internal/cache"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleAnyRequest(t *testing.T) {
	app := fiber.New()
	app.Use(handleAnyRequest)

	requestBody := []byte(`{"message": "test"}`)
	r := httptest.NewRequest("POST", "/test", bytes.NewReader(requestBody))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Request-Key", "test-key")

	resp, err := app.Test(r, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseData)

	assert.NoError(t, err)
	assert.Equal(t, "/test", responseData["path"])
	assert.Equal(t, "POST", responseData["method"])

	var cachedRequests []*request
	jsonStr, _ := cache.Current.Get("test-key")
	err = json.Unmarshal([]byte(jsonStr.(string)), &cachedRequests)

	assert.NoError(t, err)
	assert.Equal(t, cachedRequests[0].Path, "/test")
	assert.Equal(t, cachedRequests[0].Method, "POST")
	assert.Equal(t, cachedRequests[0].Body.(map[string]interface{})["message"], "test")
}

func TestExtractStoreKey(t *testing.T) {
	app := fiber.New()
	app.Post("/test", func(c *fiber.Ctx) error {
		key := extractStoreKey(c)
		return c.SendString(key)
	})

	r := httptest.NewRequest("POST", "/test", nil)
	r.Header.Set("X-Request-Key", "stored-key")

	resp, err := app.Test(r, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBytes, err := io.ReadAll(resp.Body)

	assert.NoError(t, err)
	assert.Equal(t, "stored-key", string(respBytes))
}
