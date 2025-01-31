// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"cosmoparrot/internal/cache"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"sort"
)

func handleGetAllRequests(c *fiber.Ctx) error {
	var list []*response
	for _, v := range cache.Current.Items() {
		var resp *response

		err := json.Unmarshal([]byte(v.Object.(string)), &resp)
		if err != nil {
			log.Errorf("failed to deserialize data, error: %s", err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		list = append(list, resp)
	}

	if list == nil {
		list = []*response{}
	} else {
		sort.Slice(list, func(i, j int) bool {
			return list[i].Time.After(list[j].Time)
		})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}
func handleGetRequestByKey(c *fiber.Ctx) error {
	if storeKey := c.Params("key"); storeKey != "" {
		log.Debugf("reading from cache with key %s", storeKey)

		if entry, found := cache.Current.Get(storeKey); found {
			var resp *response

			err := json.Unmarshal([]byte(entry.(string)), &resp)
			if err != nil {
				log.Errorf("failed to deserialize data, error: %s", err.Error())
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			return c.Status(fiber.StatusOK).JSON(resp)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}
