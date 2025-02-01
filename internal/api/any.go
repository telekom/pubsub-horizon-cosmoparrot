// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"cosmoparrot/internal/cache"
	"cosmoparrot/internal/config"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	go_cache "github.com/patrickmn/go-cache"
	"slices"
	"strconv"
	"strings"
	"time"
)

func handleAnyRequest(c *fiber.Ctx) error {
	var responseBody any

	if len(c.Body()) > 0 {
		err := json.Unmarshal(c.Body(), &responseBody)
		if err != nil {
			log.Errorf("failed to deserialize data, error: %s", err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	setResponseHeaders(c)

	reqData := &request{
		Time:    time.Now(),
		Path:    c.Path(),
		Method:  c.Method(),
		Headers: c.GetReqHeaders(),
		Body:    responseBody,
	}

	// write request to store if request key is found
	// in the request headers
	if key := extractStoreKey(c); key != "" {
		log.Debugf("writing to cache with key %s", key)

		var requests []*request
		if entry, found := cache.Current.Get(key); found {
			err := json.Unmarshal([]byte(entry.(string)), &requests)
			if err != nil {
				log.Errorf("failed to deserialize data, error: %s", err.Error())
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		} else {
			requests = []*request{}
		}

		requests = append(requests, reqData)

		jsonData, err := json.Marshal(requests)
		if err != nil {
			log.Errorf("failed to serialize data, error: %s", err.Error())

			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// write to store
		cache.Current.Set(key, string(jsonData), go_cache.DefaultExpiration)
	}

	return c.Status(getResponseCode(c)).JSON(reqData)
}

func handleStoreWrites() {

}

func extractStoreKey(c *fiber.Ctx) string {
	list := config.LoadedConfiguration.StoreKeyRequestHeaders

	var result []string
	for _, v := range list {
		result = append(result, strings.ToLower(v))
	}

	for k, v := range c.GetReqHeaders() {
		if slices.Contains(result, strings.ToLower(k)) {
			if len(v) > 0 {
				return strings.Clone(v[0])
			}
		}
	}

	return ""
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

	// set some default request headers
	c.Set("Current-Control", "max-age=0, must-revalidate")
}

func getResponseCode(c *fiber.Ctx) int {
	mapping := config.LoadedConfiguration.MethodResponseCodeMapping

	for _, m := range mapping {
		s := strings.Split(m, ":")
		if len(s) == 2 {
			if strings.ToUpper(strings.TrimSpace(s[0])) == c.Method() {

				responseCode, err := strconv.Atoi(s[1])
				if err != nil {
					log.Errorf("could not successfully parse method request code mapping configuration. Fallback to request code: %d", config.LoadedConfiguration.ResponseCode)

					return config.LoadedConfiguration.ResponseCode
				}

				return responseCode
			}
		}
	}

	return config.LoadedConfiguration.ResponseCode
}
