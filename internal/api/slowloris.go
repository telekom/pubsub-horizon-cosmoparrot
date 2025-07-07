// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"bufio"
	"cosmoparrot/internal/config"
	"github.com/gofiber/fiber/v2"
	"time"
)

func handleGetSlowloris(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/plain")

	durationSec := config.LoadedConfiguration.SlowlorisDefaultDurationSeconds
	if d := ctx.Query("duration"); d != "" {
		if parsed, err := time.ParseDuration(d + "s"); err == nil && parsed > 0 {
			durationSec = int(parsed.Seconds())
		}
	}

	// add a query parameter to set the time between single dots
	intervalSec := config.LoadedConfiguration.SlowlorisDefaultIntervalSeconds
	if t := ctx.Query("interval"); t != "" {
		if parsed, err := time.ParseDuration(t + "s"); err == nil && parsed > 0 {
			intervalSec = int(parsed.Seconds())
		}
	}

	ctx.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		start := time.Now()
		for i := 0; ; i++ {
			if time.Since(start) > time.Duration(durationSec)*time.Second {
				break
			}
			w.WriteByte('.')
			w.Flush()
			time.Sleep(time.Duration(intervalSec) * time.Second)
		}
	})
	return nil
}
