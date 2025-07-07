// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetSlowloris(t *testing.T) {
	app := fiber.New()
	app.Get("/slowloris", handleGetSlowloris)

	req := httptest.NewRequest("GET", "/slowloris?duration=9&interval=3", nil)
	resp, err := app.Test(req, 10000)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(body))
}
