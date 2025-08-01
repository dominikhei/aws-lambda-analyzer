// Copyright 2025 dominikhei
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Cache is a simple cache, which is used to store whether functions have been invoked yet.
// This reduces the amount of API calls and cost in the metrics functions.
// The cache gets deleted as soon as the process is killed.
package cache

import (
	"fmt"
	"sync"
	"time"
)

// CacheKey contains the identifiers of a lambda function and invocation interval.
type CacheKey struct {
	FunctionName string
	Qualifier    string
	Start        time.Time
	End          time.Time
}

// This computes a string out of CacheKey
func (k CacheKey) String() string {
	return fmt.Sprintf("%s|%s|%d|%d", k.FunctionName, k.Qualifier, k.Start.Unix(), k.End.Unix())
}

// The actual cache implementation, thread safety is guaranteed via a mutex.
// A CacheKey is only stored in it, if it has been invoked.
type Cache struct {
	mu    sync.RWMutex
	store map[string]int // map from key string to invocation count
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]int),
	}
}

// Has returns true if the key exists in the cache.
func (c *Cache) Has(key CacheKey) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.store[key.String()]
	return ok
}

// Set stores the invocation count for the given key.
func (c *Cache) Set(key CacheKey, count int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key.String()] = count
}

// Get returns the invocation count for the key and a bool indicating if it was found.
func (c *Cache) Get(key CacheKey) (int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	count, ok := c.store[key.String()]
	return count, ok
}
