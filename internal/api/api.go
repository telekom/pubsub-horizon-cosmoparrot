// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"cosmoparrot/internal/config"
	"embed"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"net/http"
)

var app *fiber.App

func init() {
	app = fiber.New()
	app.Use(createNewLogHandler())

}

func AddHandlers(f embed.FS) {
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/requests", handleGetAllRequests)
	v1.Get("/requests/:key", handleGetRequestByKey)

	app.Use(handleAnyRequest)

	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(f),
		PathPrefix: "web",
	}))
}

func Listen() {
	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.LoadedConfiguration.Port)))
}
