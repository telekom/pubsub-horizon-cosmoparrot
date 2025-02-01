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
	var list []*request
	for _, v := range cache.Current.Items() {
		var requests []*request

		err := json.Unmarshal([]byte(v.Object.(string)), &requests)
		if err != nil {
			log.Errorf("failed to deserialize data, error: %s", err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		list = append(list, requests...)
	}

	if list == nil {
		list = []*request{}
	} else {
		sort.Slice(list, func(i, j int) bool {
			return list[i].Time.After(list[j].Time)
		})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}
func handleGetRequestByKey(c *fiber.Ctx) error {
	if key := c.Params("key"); key != "" {
		log.Debugf("reading from cache with key %s", key)

		if entry, found := cache.Current.Get(key); found {
			var requests []*request

			err := json.Unmarshal([]byte(entry.(string)), &requests)
			if err != nil {
				log.Errorf("failed to deserialize data, error: %s", err.Error())
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			sort.Slice(requests, func(i, j int) bool {
				return requests[i].Time.After(requests[j].Time)
			})

			return c.Status(fiber.StatusOK).JSON(requests)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}
