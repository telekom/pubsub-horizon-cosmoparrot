// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/gofiber/fiber/v2"
)

func handleDevNull(c *fiber.Ctx) error {
	// Body is already read by Fiber into c.Body(); we simply ignore it.
	// No deserialization, no caching, no logging — just return the status code.
	return c.SendStatus(getResponseCode(c))
}
