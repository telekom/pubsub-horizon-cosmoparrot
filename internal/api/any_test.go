// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"bytes"
	"cosmoparrot/internal/cache"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
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
	// the request body is mirrored verbatim into the response body
	assert.Equal(t, map[string]interface{}{"message": "test"}, responseData["body"])

	var cachedRequests []*request
	jsonStr, _ := cache.Current.Get("test-key")
	err = json.Unmarshal([]byte(jsonStr.(string)), &cachedRequests)

	assert.NoError(t, err)
	assert.Equal(t, cachedRequests[0].Path, "/test")
	assert.Equal(t, cachedRequests[0].Method, "POST")
	assert.JSONEq(t, `{"message": "test"}`, string(cachedRequests[0].Body))
}

func TestHandleAnyRequest_EchoesRawQueryAndRequestURI(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		query      string
		requestURI string
	}{
		{
			name:       "no query",
			target:     "/apis/x/y",
			query:      "",
			requestURI: "/apis/x/y",
		},
		{
			name:       "empty query",
			target:     "/apis/x/y?",
			query:      "",
			requestURI: "/apis/x/y?",
		},
		{
			name:       "several parameters keep their order",
			target:     "/apis/x/y?zeta=1&alpha=2&mu=3",
			query:      "zeta=1&alpha=2&mu=3",
			requestURI: "/apis/x/y?zeta=1&alpha=2&mu=3",
		},
		{
			name:       "percent-encoded characters stay encoded",
			target:     "/apis/a%2Fb/c%20d?b=%2F&c=a%20b",
			query:      "b=%2F&c=a%20b",
			requestURI: "/apis/a%2Fb/c%20d?b=%2F&c=a%20b",
		},
		{
			name:       "repeated parameter",
			target:     "/apis/x/y?limit=1&b=%2F&limit=2",
			query:      "limit=1&b=%2F&limit=2",
			requestURI: "/apis/x/y?limit=1&b=%2F&limit=2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(handleAnyRequest)

			storeKey := "raw-query-" + tt.name
			cache.Current.Delete(storeKey)
			t.Cleanup(func() { cache.Current.Delete(storeKey) })

			r := httptest.NewRequest(http.MethodGet, tt.target, nil)
			r.Header.Set("X-Request-Key", storeKey)

			resp, err := app.Test(r, -1)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var responseData map[string]any
			assert.NoError(t, json.NewDecoder(resp.Body).Decode(&responseData))

			assert.Equal(t, tt.requestURI, responseData["requestUri"])
			if tt.query == "" {
				assert.NotContains(t, responseData, "query")
			} else {
				assert.Equal(t, tt.query, responseData["query"])
			}

			var cachedRequests []*request
			entry, found := cache.Current.Get(storeKey)
			assert.True(t, found)
			assert.NoError(t, json.Unmarshal([]byte(entry.(string)), &cachedRequests))
			assert.Len(t, cachedRequests, 1)
			assert.Equal(t, tt.query, cachedRequests[0].Query)
			assert.Equal(t, tt.requestURI, cachedRequests[0].RequestURI)
		})
	}
}

func TestHandleAnyRequest_MirrorBodyDisabled(t *testing.T) {
	app := fiber.New()
	app.Use(handleAnyRequest)

	requestBody := []byte(`{"message": "test"}`)
	r := httptest.NewRequest("POST", "/test?mirrorBody=false", bytes.NewReader(requestBody))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Request-Key", "no-mirror-key")

	resp, err := app.Test(r, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	assert.NoError(t, err)
	assert.Equal(t, "/test", responseData["path"])
	assert.Equal(t, "POST", responseData["method"])
	// body must NOT be mirrored back in the response
	_, hasBody := responseData["body"]
	assert.False(t, hasBody, "response body should be omitted when mirrorBody=false")

	// the request is recorded without its body when mirroring is disabled
	var cachedRequests []*request
	jsonStr, _ := cache.Current.Get("no-mirror-key")
	err = json.Unmarshal([]byte(jsonStr.(string)), &cachedRequests)
	assert.NoError(t, err)
	assert.Equal(t, "/test", cachedRequests[0].Path)
	assert.Nil(t, cachedRequests[0].Body)
}

func TestHandleAnyRequest_MalformedBody(t *testing.T) {
	app := fiber.New()
	app.Use(handleAnyRequest)

	r := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte(`{not valid json`)))
	r.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(r, -1)
	assert.NoError(t, err)
	// A malformed client body is a client error, not a server fault.
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHandleAnyRequest_DrainsStreamedBodyWhenNotMirroring(t *testing.T) {
	app := fiber.New(fiber.Config{StreamRequestBody: true})
	app.Use(handleAnyRequest)

	// A large body forces fasthttp to expose it as a stream; with
	// mirrorBody=false the handler must drain it so the request still
	// completes cleanly instead of leaving the body unread.
	largeBody := bytes.Repeat([]byte("a"), 5*1024*1024)
	r := httptest.NewRequest("POST", "/test?mirrorBody=false", bytes.NewReader(largeBody))
	r.Header.Set("Content-Type", "application/octet-stream")

	resp, err := app.Test(r, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	assert.NoError(t, err)
	assert.Equal(t, "/test", responseData["path"])
	_, hasBody := responseData["body"]
	assert.False(t, hasBody)
}

func TestHandleAnyRequest_WithResponseDelay(t *testing.T) {
	app := fiber.New()
	app.Use(handleAnyRequest)

	r := httptest.NewRequest("GET", "/test?responseDelay=200", nil)

	start := time.Now()
	resp, err := app.Test(r, -1)
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(150), "response should be delayed by at least ~200ms")

	var responseData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	assert.NoError(t, err)
	assert.Equal(t, "/test", responseData["path"])
	assert.Equal(t, "GET", responseData["method"])
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
