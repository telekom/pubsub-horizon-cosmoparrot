// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/gofiber/fiber/v2"
)

func handleDevNull(c *fiber.Ctx) error {
	// With StreamRequestBody enabled, the body is never read into memory
	// because we never call c.Body(). No deserialization, no caching, no logging.
	return c.SendStatus(getResponseCode(c))
}
