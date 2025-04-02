// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	go_cache "github.com/patrickmn/go-cache"
	"time"
)

var Current *go_cache.Cache

func init() {
	Current = go_cache.New(1*time.Hour, 10*time.Minute)
}
