// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	req = httptest.NewRequest(http.MethodGet, "/foobar", nil)
	req.Header.Set("x-request-key", "123")
	_, err = app.Test(req)
	assert.NoError(t, err)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/requests/123", nil)
	resp, _ = app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/requests/234", nil)
	resp, _ = app.Test(req)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
