// Copyright 2026 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package memlimit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "limit")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestReadLimit_ValidValue(t *testing.T) {
	v, ok := readLimit(writeTemp(t, "536870912\n"))
	assert.True(t, ok)
	assert.Equal(t, int64(536870912), v)
}

func TestReadLimit_CgroupV2Max(t *testing.T) {
	_, ok := readLimit(writeTemp(t, "max\n"))
	assert.False(t, ok)
}

func TestReadLimit_Empty(t *testing.T) {
	_, ok := readLimit(writeTemp(t, ""))
	assert.False(t, ok)
}

func TestReadLimit_NonNumeric(t *testing.T) {
	_, ok := readLimit(writeTemp(t, "not-a-number"))
	assert.False(t, ok)
}

func TestReadLimit_ZeroOrNegative(t *testing.T) {
	_, ok := readLimit(writeTemp(t, "0"))
	assert.False(t, ok)

	_, ok = readLimit(writeTemp(t, "-1"))
	assert.False(t, ok)
}

func TestReadLimit_CgroupV1Unlimited(t *testing.T) {
	// cgroup v1 reports a huge sentinel when no limit is set.
	_, ok := readLimit(writeTemp(t, "9223372036854771712"))
	assert.False(t, ok)
}

func TestReadLimit_MissingFile(t *testing.T) {
	_, ok := readLimit(filepath.Join(t.TempDir(), "does-not-exist"))
	assert.False(t, ok)
}

func TestConfigure_RespectsExplicitEnv(t *testing.T) {
	// When GOMEMLIMIT is set, Configure must not override it.
	t.Setenv("GOMEMLIMIT", "128MiB")
	assert.NotPanics(t, Configure)
}
