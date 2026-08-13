// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package intsequence can be used to produce sequences of positive integers.
// It is safe to use concurrently.
package intsequence

import (
	"sync"
)

// Generator produces a sequence of integers, in counting order, starting with
// 1.
type Generator struct {
	counter int64
	mu      sync.Mutex
}

// Next increments the internal counter and returns its new value. It is safe
// to use concurrently. This will rollover to math.MinInt64 if the internal
// counter overflows.
func (g *Generator) Next() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.counter++
	return g.counter
}
