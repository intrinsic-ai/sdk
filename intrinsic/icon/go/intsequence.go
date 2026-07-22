// Copyright 2023 Intrinsic Innovation LLC

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
