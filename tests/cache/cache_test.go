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

package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/dominikhei/serverless-statistics/internal/cache"
)

func TestCacheSetGetHas(t *testing.T) {
	c := cache.NewCache()

	key := cache.CacheKey{
		FunctionName: "myFunc",
		Qualifier:    "v1",
		Start:        time.Unix(1000, 0),
		End:          time.Unix(2000, 0),
	}

	if c.Has(key) {
		t.Error("expected key to not exist")
	}

	count, ok := c.Get(key)
	if ok {
		t.Errorf("expected Get to return false for non-existing key, got count %d", count)
	}

	c.Set(key, 42)

	if !c.Has(key) {
		t.Error("expected key to exist")
	}

	count, ok = c.Get(key)
	if !ok {
		t.Error("expected Get to return true for existing key")
	}
	if count != 42 {
		t.Errorf("expected count 42, got %d", count)
	}
}

func TestCacheConcurrency(t *testing.T) {
	c := cache.NewCache()
	key := cache.CacheKey{
		FunctionName: "func",
		Qualifier:    "q",
		Start:        time.Unix(0, 0),
		End:          time.Unix(1, 0),
	}

	const n = 1000
	var wg sync.WaitGroup
	wg.Add(n * 2)

	for i := 0; i < n; i++ {
		go func(val int) {
			defer wg.Done()
			c.Set(key, val)
		}(i)

		go func() {
			defer wg.Done()
			_, _ = c.Get(key)
		}()
	}

	wg.Wait()
}
