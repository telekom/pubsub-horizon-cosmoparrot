// Copyright 2024 Deutsche Telekom IT GmbH
//
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

func TestCacheInitialization(t *testing.T) {
	// Ensure the cache is initialized
	assert.NotNil(t, Current, "Cache should not be nil")
}

func TestCacheSetAndGet(t *testing.T) {
	// Reset cache for a clean test
	Current = cache.New(1*time.Hour, 10*time.Minute)

	key := "testKey"
	value := "testValue"

	// Store value in cache
	Current.Set(key, value, cache.DefaultExpiration)

	// Retrieve from cache
	storedValue, found := Current.Get(key)
	assert.True(t, found, "Key should exist in cache")
	assert.Equal(t, value, storedValue, "Stored value should match expected value")
}

func TestCacheExpiration(t *testing.T) {
	// Reset cache with short expiration
	Current = cache.New(50*time.Millisecond, 10*time.Millisecond)

	key := "tempKey"
	value := "tempValue"
	Current.Set(key, value, 50*time.Millisecond)

	// Ensure key exists initially
	storedValue, found := Current.Get(key)
	assert.True(t, found, "Key should be found before expiration")
	assert.Equal(t, value, storedValue, "Stored value should match before expiration")

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Ensure key is expired
	_, found = Current.Get(key)
	assert.False(t, found, "Key should be expired and not found")
}

func TestCacheDelete(t *testing.T) {
	// Reset cache for a clean test
	Current = cache.New(1*time.Hour, 10*time.Minute)

	key := "deleteKey"
	value := "toBeDeleted"

	Current.Set(key, value, cache.DefaultExpiration)

	// Ensure key exists before deletion
	_, found := Current.Get(key)
	assert.True(t, found, "Key should exist in cache before deletion")

	// Delete key
	Current.Delete(key)

	// Ensure key no longer exists
	_, found = Current.Get(key)
	assert.False(t, found, "Key should not exist in cache after deletion")
}
