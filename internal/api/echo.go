package api

import (
	"cosmoparrot/internal/cache"
	"cosmoparrot/internal/config"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	go_cache "github.com/patrickmn/go-cache"
	"slices"
	"strconv"
	"strings"
)

func createNewEchoHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var responseBody any

		if len(c.Body()) > 0 {
			err := json.Unmarshal(c.Body(), &responseBody)
			if err != nil {
				return err
			}
		}

		var storeKey string
		storeKeyRequestHeaders := getStoreKeyRequestHeaders()

		for k, v := range c.GetReqHeaders() {
			if slices.Contains(storeKeyRequestHeaders, strings.ToLower(k)) {
				if len(v) > 0 {
					storeKey = v[0]
					break
				}
			}
		}

		setResponseHeaders(c)

		responseCode := decideResponseCode(c)

		resp := &response{
			Path:    c.Path(),
			Method:  c.Method(),
			Headers: c.GetReqHeaders(),
			Body:    responseBody,
		}

		jsonData, _ := json.Marshal(resp)

		if storeKey != "" {
			fmt.Printf("writing to cache key %s \n", storeKey)

			cache.Current.Set("123", string(jsonData), go_cache.DefaultExpiration)

			fmt.Println(cache.Current.Items())
		}

		return c.Status(responseCode).JSON(resp)
	}
}

func getStoreKeyRequestHeaders() []string {
	list := config.LoadedConfiguration.StoreKeyRequestHeaders

	var result []string
	for _, v := range list {
		result = append(result, strings.ToLower(v))
	}

	return result
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

	// set some default response headers
	c.Set("Current-Control", "max-age=0, must-revalidate")
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
