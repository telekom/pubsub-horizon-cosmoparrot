// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"bytes"
	"cosmoparrot/internal/config"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

func NewApp(f embed.FS) *fiber.App {
	app := fiber.New(fiber.Config{
		StreamRequestBody: true,
		JSONEncoder:       marshalJSONWithoutHTMLEscape,
	})
	var shutdownTracing func()
	apiMiddleware := make([]fiber.Handler, 0)
	if config.LoadedConfiguration.OTelEnabled {
		shutdownTracing = initTracerProvider()
		apiMiddleware = append(apiMiddleware, otelfiber.Middleware())
	}
	app.Use(createNewLogHandler())
	app.Use(healthcheck.New())
	app.Hooks().OnShutdown(func() error {
		if shutdownTracing != nil {
			shutdownTracing()
		}
		return nil
	})

	// Attach tracing only on /api routes to avoid exporting health/static traffic.
	api := app.Group("/api", apiMiddleware...)
	v1 := api.Group("/v1")
	v1.Get("/requests", handleGetAllRequests)
	v1.Get("/requests/:key", handleGetRequestByKey)
	v1.Get("/slowloris", handleGetSlowloris)
	v1.All("/devnull", handleDevNull)

	app.Use(handleAnyRequest)

	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(f),
		PathPrefix: "web",
	}))

	return app
}

func Listen(f embed.FS) {
	app := NewApp(f)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.LoadedConfiguration.Port)))
}

// marshalJSONWithoutHTMLEscape works like json.Marshal, but keeps &, < and >
// as literal characters, so echoed values match the request byte for byte.
func marshalJSONWithoutHTMLEscape(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}
