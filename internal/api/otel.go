// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"cosmoparrot/internal/config"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func initTracerProvider() func() {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize OTel OTLP trace exporter; tracing disabled")
		return func() {}
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(config.LoadedConfiguration.OTelServiceName),
		),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize OTel resource; tracing disabled")
		return func() {}
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	log.Info().Msg("OpenTelemetry tracing enabled")

	return func() {
		if shutdownErr := tracerProvider.Shutdown(ctx); shutdownErr != nil {
			log.Error().Err(shutdownErr).Msg("Failed to shutdown OTel tracer provider")
		}
	}
}
