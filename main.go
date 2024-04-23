// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"cosmoparrot/internal/config"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"strconv"
	"strings"
)

type response struct {
	Path    string              `json:"path"`
	Method  string              `json:"method"`
	Headers map[string][]string `json:"headers,omitempty"`
	Body    any                 `json:"body,omitempty"`
}

func setResponseHeaders(c *fiber.Ctx) {
	for name, values := range c.GetReqHeaders() {
		prefix := "x-parrot-"
		if strings.HasPrefix(strings.ToLower(name), prefix) {
			_, newName, _ := strings.Cut(strings.ToLower(name), prefix)
			for _, value := range values {
				c.Set(newName, value)
			}
		}
	}
}

func main() {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format:   "${green}→ Request received:\n${reset}${time} | ${status} - ${method} ${path}\n${green}→ Request headers:${magenta}\n${custom_tag}${green}→ Request body:${cyan}\n${body}${reset}\n",
		TimeZone: "UTC",
		CustomTags: map[string]logger.LogFunc{
			"custom_tag": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				var str string
				for k, v := range c.GetReqHeaders() {
					str += fmt.Sprintf("%v: %v\n", k, v)
				}

				return output.WriteString(str)
			},
		},
	}))

	app.Use(func(c *fiber.Ctx) error {
		var responseBody any

		if len(c.Body()) > 0 {
			return json.Unmarshal(c.Body(), &responseBody)
		}

		setResponseHeaders(c)

		responseCode := decideResponseCode(c)
		return c.Status(responseCode).JSON(&response{
			Path:    c.Path(),
			Method:  c.Method(),
			Headers: c.GetReqHeaders(),
			Body:    responseBody,
		})
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.LoadedConfiguration.Port)))
}

func decideResponseCode(c *fiber.Ctx) int {
	mapping := config.LoadedConfiguration.MethodResponseCodeMapping

	for _, m := range mapping {
		s := strings.Split(m, ":")
		if len(s) == 2 {
			if strings.ToUpper(strings.TrimSpace(s[0])) == c.Method() {

				responseCode, err := strconv.Atoi(s[1])
				if err != nil {
					log.Errorf("could not successfully parse method response code mapping configuration. Fallback to response code: %d", config.LoadedConfiguration.ResponseCode)

					return config.LoadedConfiguration.ResponseCode
				}

				return responseCode
			}
		}
	}

	return config.LoadedConfiguration.ResponseCode
}
