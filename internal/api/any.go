// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"cosmoparrot/internal/cache"
	"cosmoparrot/internal/config"
	"cosmoparrot/internal/utils"
	"encoding/json"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	go_cache "github.com/patrickmn/go-cache"
)

const maxResponseDelayMs = 60000
const maxResponseSize = 10_000_000
const maxResponseSizePaddingWindowSize = 4096 // How often the padding is shifted before wrapping around
var paddingSource = newRandomString(maxResponseSize + maxResponseSizePaddingWindowSize)
var paddingOffset atomic.Int64

// newRandomString builds a random string
func newRandomString(length int) string {
	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func handleAnyRequest(c *fiber.Ctx) error {
	userAgent := c.Get("User-Agent")
	if c.Path() == "/" && utils.IsBrowser(userAgent) {
		return c.Next()
	}

	var responseBody json.RawMessage

	if body := c.Body(); len(body) > 0 {
		if !json.Valid(body) {
			log.Debug("failed to deserialize request body: invalid JSON")
			return c.SendStatus(fiber.StatusBadRequest)
		}
		responseBody = body
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

	// suppress the echoed request body if the client opted out; the stored
	// copy above intentionally keeps the full body for later inspection.
	if !getMirrorBody(c) {
		reqData.Body = nil
	}

	if delay := getResponseDelay(c); delay > 0 {
		time.Sleep(delay)
	}

	if size := getResponseSize(c); size > 0 {
		offset := int(paddingOffset.Add(1) % maxResponseSizePaddingWindowSize)
		reqData.Padding = paddingSource[offset : offset+size]
	}

	return c.Status(getResponseCode(c)).JSON(reqData)
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

// queryCaseInsensitive returns the value of the query parameter whose name
// matches key case-insensitively, or "" if none is present.
func queryCaseInsensitive(c *fiber.Ctx, key string) string {
	for k, v := range c.Queries() {
		if strings.EqualFold(k, key) {
			return v
		}
	}
	return ""
}

// getMirrorBody reads the optional "mirrorBody" query parameter. Body mirroring
// is enabled by default; it is only disabled when the parameter is explicitly set
// to a false-y value ("false", "0", "f"). Missing or unparsable values keep the
// default (true), so the request body is echoed back in the response body.
func getMirrorBody(c *fiber.Ctx) bool {
	raw := queryCaseInsensitive(c, "mirrorBody")
	if raw == "" {
		return true
	}

	enabled, err := strconv.ParseBool(strings.TrimSpace(raw))
	if err != nil {
		return true
	}

	return enabled
}

func getResponseCode(c *fiber.Ctx) int {
	if rc := queryCaseInsensitive(c, "responseCode"); rc != "" {
		if code, err := strconv.Atoi(strings.TrimSpace(rc)); err == nil && code >= 100 && code <= 599 {
			return code
		}
		// invalid override -> fallthrough to mapping/default
	}

	if code, ok := config.LoadedConfiguration.MethodResponseCodeMap[c.Method()]; ok {
		return code
	}

	return config.LoadedConfiguration.ResponseCode
}

// getResponseDelay reads the optional "responseDelay" query parameter (milliseconds)
// and returns the corresponding time.Duration. Returns 0 for missing, non-integer,
// negative, or out-of-range (>60000 ms) values.
func getResponseDelay(c *fiber.Ctx) time.Duration {
	raw := queryCaseInsensitive(c, "responseDelay")
	if raw == "" {
		return 0
	}

	ms, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || ms < 0 || ms > maxResponseDelayMs {
		return 0
	}

	return time.Duration(ms) * time.Millisecond
}

// getResponseSize reads the optional "responseSize" query parameter (bytes)
// and returns the corresponding integer. Returns 0 for missing, non-integer,
// negative, or out-of-range (> 100 MB) values.
func getResponseSize(c *fiber.Ctx) int {
	raw := queryCaseInsensitive(c, "responseSize")
	if raw == "" {
		return 0
	}

	ms, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || ms < 0 || ms > maxResponseSize {
		return 0
	}

	return ms
}
