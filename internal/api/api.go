// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"cosmoparrot/internal/config"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var app *fiber.App

func init() {
	app = fiber.New()
	app.Use(createNewLogHandler())

	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/requests", handleGetAllRequests)
	v1.Get("/requests/:key", handleGetRequestByKey)

	app.Use(handleAnyRequest)
}

func Listen() {
	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.LoadedConfiguration.Port)))
}
