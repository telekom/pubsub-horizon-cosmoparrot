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
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/stretchr/testify/assert"
)

type byteSliceWriter struct {
	b []byte
}

func (w *byteSliceWriter) Write(p []byte) (int, error) {
	w.b = append(w.b, p...)
	return len(p), nil
}

func TestLogFormat(t *testing.T) {
	var w byteSliceWriter

	c := getLoggerConfig()

	// For this test we disable color log output
	c.DisableColors = true

	// Here we capture everything that is logged via our fiber middleware
	c.Done = func(c *fiber.Ctx, logString []byte) {
		writtenBytes, err := w.Write(logString)
		assert.NoError(t, err)
		assert.Equal(t, len(logString), writtenBytes)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Use the custom logger as middleware
	app.Use(logger.New(c))

	// Define a test route
	app.Post("/test", func(c *fiber.Ctx) error {
		// Simulate request body handling
		return c.SendString("Hello, Fiber!")
	})

	// Create a test request
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(`{"key": "value"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Custom-Header", "TestValue")

	// Perform the request
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get the captured log output
	logOutput := string(w.b)

	// Validate log contains expected data
	assert.Contains(t, logOutput, "→ Request received:")
	assert.Contains(t, logOutput, "POST /test")
	assert.Contains(t, logOutput, "200")
	assert.Contains(t, logOutput, `"key": "value"`)
	assert.Contains(t, logOutput, "X-Custom-Header: [TestValue]")
}
