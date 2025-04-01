// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"cosmoparrot/internal/api"
	"embed"
	_ "embed"
)

//go:embed web/*
var webDir embed.FS

func main() {
	api.Listen(webDir)
}
