package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func createNewLogHandler() fiber.Handler {
	return logger.New(logger.Config{
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
	})
}
