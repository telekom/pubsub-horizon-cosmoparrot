// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"bufio"
	"cosmoparrot/internal/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
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

			err := w.WriteByte('.')
			if err != nil {
				log.Error().Err(err).Msg("failed to write to slowloris buffer")
				return
			}

			err = w.Flush()
			if err != nil {
				log.Error().Err(err).Msg("error flushing slowloris buffer")
				return
			}

			time.Sleep(time.Duration(intervalSec) * time.Second)
		}
	})
	return nil
}
