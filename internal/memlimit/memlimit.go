// Copyright 2026 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

// Package memlimit derives the Go soft memory limit (GOMEMLIMIT) from the
// container's cgroup memory limit.
package memlimit

import (
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"
)

const (
	cgroupV2Max = "/sys/fs/cgroup/memory.max"
	cgroupV1Max = "/sys/fs/cgroup/memory/memory.limit_in_bytes"
)

const reserveRatio = 0.9

const cgroupV1Unlimited = int64(1) << 62

// Configure sets GOMEMLIMIT from the container's cgroup memory limit unless it
// has already been provided explicitly via the environment.
func Configure() {
	if _, ok := os.LookupEnv("GOMEMLIMIT"); ok {
		return
	}

	limit, ok := detectCgroupLimit()
	if !ok {
		return
	}

	soft := int64(float64(limit) * reserveRatio)
	debug.SetMemoryLimit(soft)
	log.Infof("GOMEMLIMIT set to %d bytes (%.0f%% of cgroup limit %d bytes)", soft, reserveRatio*100, limit)
}

func detectCgroupLimit() (int64, bool) {
	for _, path := range []string{cgroupV2Max, cgroupV1Max} {
		if v, ok := readLimit(path); ok {
			return v, true
		}
	}
	return 0, false
}

func readLimit(path string) (int64, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, false
	}

	s := strings.TrimSpace(string(data))
	if s == "" || s == "max" {
		return 0, false
	}

	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || v <= 0 || v >= cgroupV1Unlimited {
		return 0, false
	}

	return v, true
}
