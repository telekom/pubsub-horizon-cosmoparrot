package api

import (
	"cosmoparrot/internal/cache"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func handleStoreAccess(c *fiber.Ctx) error {
	storeKey := c.Params("key")

	if storeKey != "" {
		fmt.Printf("getting cache key %s \n", storeKey)
		fmt.Println(cache.Current.Items())

		if entry, found := cache.Current.Get(storeKey); found {
			jsonString := entry.(string)

			var resp *response
			if err := json.Unmarshal([]byte(jsonString), &resp); err == nil {
				return c.Status(fiber.StatusOK).JSON(resp)
			}
		}
	}

	return c.SendStatus(fiber.StatusBadRequest)
}
