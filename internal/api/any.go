// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"cosmoparrot/internal/cache"
	"cosmoparrot/internal/config"
	"cosmoparrot/internal/utils"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	go_cache "github.com/patrickmn/go-cache"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

var storeWritesOnce sync.Once

// StartBackgroundWorkers starts package background workers (idempotent).
// Call this from your application startup when you want background workers
// to run (e.g. before Listen). Tests that only construct the app can skip
// calling this to avoid background activity.
func StartBackgroundWorkers() {
	handleStoreWrites()
}

func handleAnyRequest(c *fiber.Ctx) error {
	userAgent := c.Get("User-Agent")
	if c.Path() == "/" && utils.IsBrowser(userAgent) {
		return c.Next()
	}

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
	storeWritesOnce.Do(func() {
		// Start a small background worker that can later be extended to flush
		// pending store writes. Keep it non-blocking and test-friendly.
		go func() {
			log.Debug("store write worker started")
			ticker := time.NewTicker(1 * time.Minute)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					// Placeholder: insert batching/flush logic here in the future.
					log.Debug("store write worker tick")
				}
			}
		}()
	})
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
	// Parse query string in a case-insensitive way to find any parameter that
	// corresponds to responseCode (e.g. responseCode, responsecode, response_code)
	qs := string(c.Request().URI().QueryString())
	if qs != "" {
		vals, err := url.ParseQuery(qs)
		if err == nil {
			for k, arr := range vals {
				if len(arr) == 0 {
					continue
				}
				lk := strings.ToLower(k)
				// normalize common variants by removing non-alphanumeric characters
				norm := strings.ReplaceAll(lk, "_", "")
				norm = strings.ReplaceAll(norm, "-", "")
				if norm == "responsecode" {
					val := arr[0]
					responseCode, err := strconv.Atoi(val)
					if err != nil {
						log.Errorf("invalid responseCode query parameter '%s': %v. Falling back to configured mapping/default: %d", val, err, config.LoadedConfiguration.ResponseCode)
						break
					}

					if responseCode < 100 || responseCode > 599 {
						log.Errorf("responseCode query parameter out of range (%d). Falling back to configured mapping/default: %d", responseCode, config.LoadedConfiguration.ResponseCode)
						break
					}

					return responseCode
				}
			}
		}
	}

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
