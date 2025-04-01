// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetAllRequests(t *testing.T) {
	app := NewApp(embed.FS{}) // Empty embed.FS for testing

	req := httptest.NewRequest(http.MethodGet, "/api/v1/requests", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHandleGetRequestByKey(t *testing.T) {
	app := NewApp(embed.FS{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/requests/123", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	req = httptest.NewRequest(http.MethodGet, "/foobar", nil)
	req.Header.Set("x-request-key", "123")
	resp, _ = app.Test(req)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/requests/123", nil)
	resp, _ = app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/requests/234", nil)
	resp, _ = app.Test(req)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
