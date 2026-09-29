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

package throttle

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

func TestNewConcurrencyLimiter(t *testing.T) {
	for _, limit := range []int{-5, -1, 0} {
		t.Run(fmt.Sprintf("limit_%d_returns_nil", limit), func(t *testing.T) {
			if cl := NewConcurrencyLimiter(limit); cl != nil {
				t.Errorf("NewConcurrencyLimiter(%d) = %v, want nil", limit, cl)
			}
		})
	}

	t.Run("positive_limit", func(t *testing.T) {
		cl := NewConcurrencyLimiter(4)
		if cl == nil {
			t.Fatal("NewConcurrencyLimiter(4) = nil, want non-nil")
		}
		if got := cap(cl.sem); got != 4 {
			t.Errorf("cap(cl.sem) = %d, want 4", got)
		}
	})
}

func TestConcurrencyLimiter_Do(t *testing.T) {
	t.Run("empty_fns", func(t *testing.T) {
		cl := NewConcurrencyLimiter(1)
		if err := cl.Do(t.Context()); err != nil {
			t.Fatalf("Do() unexpected error: %v", err)
		}
	})

	for _, tc := range []struct {
		name string
		cl   *ConcurrencyLimiter
	}{
		{name: "nil_limiter_runs_concurrently", cl: nil},
		{name: "zero_value_limiter_runs_concurrently", cl: &ConcurrencyLimiter{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			const n = 10
			var count atomic.Int32
			var wg sync.WaitGroup
			wg.Add(n)

			fns := make([]func(context.Context) error, n)
			for i := range fns {
				fns[i] = func(ctx context.Context) error {
					count.Add(1)
					wg.Done()
					wg.Wait()
					return nil
				}
			}
			if err := tc.cl.Do(t.Context(), fns...); err != nil {
				t.Fatalf("Do() error: %v", err)
			}
			if got := count.Load(); got != n {
				t.Errorf("count = %d, want %d", got, n)
			}
		})
	}

	t.Run("propagates_error_and_releases_slot", func(t *testing.T) {
		cl := NewConcurrencyLimiter(1)
		wantErr := errors.New("boom")

		err := cl.Do(t.Context(), func(ctx context.Context) error {
			return wantErr
		})
		if !errors.Is(err, wantErr) {
			t.Fatalf("Do() error = %v, want %v", err, wantErr)
		}

		// Slot must be released after error.
		if err := cl.Do(t.Context(), func(ctx context.Context) error { return nil }); err != nil {
			t.Fatalf("subsequent Do() failed: %v", err)
		}
	})

	t.Run("top_level_strictly_respects_limit", func(t *testing.T) {
		const limit = 3
		cl := NewConcurrencyLimiter(limit)

		var active atomic.Int32
		var maxActive atomic.Int32
		var completed atomic.Int32
		var firstBatch sync.WaitGroup
		firstBatch.Add(limit)

		fns := make([]func(context.Context) error, 15)
		for i := range fns {
			fns[i] = func(ctx context.Context) error {
				cur := active.Add(1)
				for {
					prev := maxActive.Load()
					if cur <= prev || maxActive.CompareAndSwap(prev, cur) {
						break
					}
				}
				if i < limit {
					firstBatch.Done()
					firstBatch.Wait()
				} else {
					time.Sleep(5 * time.Millisecond)
				}
				active.Add(-1)
				completed.Add(1)
				return nil
			}
		}

		if err := cl.Do(t.Context(), fns...); err != nil {
			t.Fatalf("Do() error: %v", err)
		}
		if got := completed.Load(); got != 15 {
			t.Errorf("completed = %d, want 15", got)
		}
		if got := maxActive.Load(); got != limit {
			t.Errorf("maxActive = %d, want %d", got, limit)
		}
	})

	t.Run("nested_Do_borrows_parent_slot_without_deadlocking", func(t *testing.T) {
		const limit = 4
		cl := NewConcurrencyLimiter(limit)

		var activeLeafWork atomic.Int32
		var maxActiveLeafWork atomic.Int32
		var totalChildren atomic.Int32

		// 6 parent tasks, each executing 5 child tasks sharing the same limiter of size 4.
		// Without parent slot borrowing, parents would consume all 4 slots and deadlock
		// waiting on child Do calls. With slot borrowing, max active leaf work never
		// exceeds 4 and never deadlocks.
		parentFns := make([]func(context.Context) error, 6)
		for i := range parentFns {
			parentFns[i] = func(ctx context.Context) error {
				childFns := make([]func(context.Context) error, 5)
				for j := range childFns {
					childFns[j] = func(ctx context.Context) error {
						cur := activeLeafWork.Add(1)
						for {
							prev := maxActiveLeafWork.Load()
							if cur <= prev || maxActiveLeafWork.CompareAndSwap(prev, cur) {
								break
							}
						}
						time.Sleep(10 * time.Millisecond)
						activeLeafWork.Add(-1)
						totalChildren.Add(1)
						return nil
					}
				}
				return cl.Do(ctx, childFns...)
			}
		}

		if err := cl.Do(t.Context(), parentFns...); err != nil {
			t.Fatalf("Do() error: %v", err)
		}
		if got := totalChildren.Load(); got != 30 {
			t.Errorf("totalChildren = %d, want 30", got)
		}
		if got := maxActiveLeafWork.Load(); got > limit {
			t.Errorf("maxActiveLeafWork = %d, want <= %d", got, limit)
		}
	})

	t.Run("single_parent_parallelizes_children_up_to_limit", func(t *testing.T) {
		const limit = 3
		cl := NewConcurrencyLimiter(limit)

		var active atomic.Int32
		var maxActive atomic.Int32
		var wg sync.WaitGroup
		wg.Add(limit)

		// 1 parent holds 1 slot, leaving 2 spare slots. When it dispatches 3 children,
		// 2 acquire the spare slots in cl.sem and 1 borrows the parent's slot via
		// parentSem, so all 3 run simultaneously!
		err := cl.Do(t.Context(), func(ctx context.Context) error {
			childFns := make([]func(context.Context) error, limit)
			for i := range childFns {
				childFns[i] = func(ctx context.Context) error {
					cur := active.Add(1)
					for {
						prev := maxActive.Load()
						if cur <= prev || maxActive.CompareAndSwap(prev, cur) {
							break
						}
					}
					wg.Done()
					wg.Wait()
					active.Add(-1)
					return nil
				}
			}
			return cl.Do(ctx, childFns...)
		})

		if err != nil {
			t.Fatalf("Do() error: %v", err)
		}
		if got := maxActive.Load(); got != limit {
			t.Errorf("maxActive = %d, want %d", got, limit)
		}
	})

	t.Run("multi_level_nesting_with_limit_1_does_not_deadlock", func(t *testing.T) {
		cl := NewConcurrencyLimiter(1)
		wantErr := errors.New("grandchild error")

		err := cl.Do(t.Context(), func(ctx context.Context) error {
			return cl.Do(ctx, func(ctx context.Context) error {
				return cl.Do(ctx, func(ctx context.Context) error {
					return wantErr
				})
			})
		})

		if !errors.Is(err, wantErr) {
			t.Errorf("Do() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("waits_for_inflight_goroutines_on_nested_error", func(t *testing.T) {
		// Limit 2: parent holds 1 slot; 1st child and 2nd child run concurrently using the
		// spare slot in cl.sem and the parent token's sem slot.
		cl := NewConcurrencyLimiter(2)
		wantErr := errors.New("nested child failed")
		firstChildStarted := make(chan struct{})
		var firstChildFinished atomic.Bool

		err := cl.Do(t.Context(), func(ctx context.Context) error {
			return cl.Do(
				ctx,
				func(ctx context.Context) error {
					close(firstChildStarted)
					<-ctx.Done()
					time.Sleep(10 * time.Millisecond)
					firstChildFinished.Store(true)
					return nil
				},
				func(ctx context.Context) error {
					<-firstChildStarted
					return wantErr
				},
			)
		})

		if !errors.Is(err, wantErr) {
			t.Errorf("Do() error = %v, want %v", err, wantErr)
		}
		if !firstChildFinished.Load() {
			t.Error("Do() returned before in-flight child goroutine finished")
		}
	})

	t.Run("independent_limiters_do_not_share_tokens", func(t *testing.T) {
		cl1 := NewConcurrencyLimiter(1)
		cl2 := NewConcurrencyLimiter(1)

		var active2 atomic.Int32
		var maxActive2 atomic.Int32

		// Calling cl2.Do inside cl1.Do should still respect cl2's limit of 1, even if
		// another goroutine already occupies cl2's single slot.
		cl2Occupied := make(chan struct{})
		releaseCl2 := make(chan struct{})
		done := make(chan error, 1)
		go func() {
			done <- cl2.Do(t.Context(), func(ctx context.Context) error {
				active2.Add(1)
				close(cl2Occupied)
				<-releaseCl2
				active2.Add(-1)
				return nil
			})
		}()
		<-cl2Occupied

		err := cl1.Do(t.Context(), func(ctx context.Context) error {
			// ctx has cl1's token, NOT cl2's token. If cl2 mistakenly treated cl1's
			// token as its own, it would borrow cl1's slot immediately while cl2 is
			// still occupied, causing active2 to reach 2.
			go func() {
				time.Sleep(15 * time.Millisecond)
				close(releaseCl2)
			}()
			return cl2.Do(ctx, func(ctx context.Context) error {
				cur := active2.Add(1)
				for {
					prev := maxActive2.Load()
					if cur <= prev || maxActive2.CompareAndSwap(prev, cur) {
						break
					}
				}
				active2.Add(-1)
				return nil
			})
		})
		if err != nil {
			t.Fatalf("cl1.Do() error: %v", err)
		}
		if err := <-done; err != nil {
			t.Fatalf("cl2.Do() error: %v", err)
		}
		if got := maxActive2.Load(); got != 1 {
			t.Errorf("maxActive2 = %d, want 1", got)
		}
	})

	t.Run("top_level_canceled_context", func(t *testing.T) {
		cl := NewConcurrencyLimiter(1)
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		err := cl.Do(ctx, func(ctx context.Context) error {
			t.Error("fn should not run when context is already canceled")
			return nil
		})
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Do() error = %v, want %v", err, context.Canceled)
		}
	})
}

type mockClientStream struct {
	grpc.ClientStream
}

type mockClientConn struct {
	grpc.ClientConnInterface
	invokeCount atomic.Int32
	streamCount atomic.Int32
	invokeErr   error
	streamErr   error
}

func (m *mockClientConn) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	m.invokeCount.Add(1)
	return m.invokeErr
}

func (m *mockClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	m.streamCount.Add(1)
	if m.streamErr != nil {
		return nil, m.streamErr
	}
	return &mockClientStream{}, nil
}

func TestConnectionWithRateLimit_Construction(t *testing.T) {
	mock := &mockClientConn{}

	t.Run("nil connection", func(t *testing.T) {
		conn := ConnectionWithRateLimit(nil, 10, 1)
		if conn != nil {
			t.Errorf("ConnectionWithRateLimit(nil, 10, 1) = %v, want nil", conn)
		}
	})

	for _, limit := range []float64{-1, 0} {
		t.Run(fmt.Sprintf("limit %v returns original connection", limit), func(t *testing.T) {
			conn := ConnectionWithRateLimit(mock, limit, 1)
			if conn != mock {
				t.Errorf("ConnectionWithRateLimit(mock, %v, 1) = %v, want %v", limit, conn, mock)
			}
		})
	}

	t.Run("positive limit returns wrapped connection", func(t *testing.T) {
		conn := ConnectionWithRateLimit(mock, 10, 1)
		if conn == mock {
			t.Errorf("ConnectionWithRateLimit(mock, 10, 1) returned unwrapped mock, want rateLimitedConn")
		}
	})
}

func TestConnectionWithRateLimiter_Construction(t *testing.T) {
	mock := &mockClientConn{}
	limiter := rate.NewLimiter(rate.Limit(10), 1)

	t.Run("nil connection", func(t *testing.T) {
		conn := ConnectionWithRateLimiter(nil, limiter)
		if conn != nil {
			t.Errorf("ConnectionWithRateLimiter(nil, limiter) = %v, want nil", conn)
		}
	})

	t.Run("nil limiter returns original connection", func(t *testing.T) {
		conn := ConnectionWithRateLimiter(mock, nil)
		if conn != mock {
			t.Errorf("ConnectionWithRateLimiter(mock, nil) = %v, want %v", conn, mock)
		}
	})

	t.Run("valid limiter returns wrapped connection", func(t *testing.T) {
		conn := ConnectionWithRateLimiter(mock, limiter)
		if conn == mock {
			t.Errorf("ConnectionWithRateLimiter(mock, limiter) returned unwrapped mock, want rateLimitedConn")
		}
	})
}

func TestConnectionWithRateLimit_Invoke(t *testing.T) {
	t.Run("success and rate limiting", func(t *testing.T) {
		mock := &mockClientConn{}
		// Rate of 20/sec (50ms per token), burst 1.
		conn := ConnectionWithRateLimit(mock, 20, 1)

		start := time.Now()
		for range 3 {
			if err := conn.Invoke(t.Context(), "/test/Method", nil, nil); err != nil {
				t.Fatalf("Invoke() unexpected error: %v", err)
			}
		}
		elapsed := time.Since(start)

		if got := mock.invokeCount.Load(); got != 3 {
			t.Errorf("invokeCount = %d, want 3", got)
		}
		// 1st call is immediate (burst 1), 2nd at ~50ms, 3rd at ~100ms.
		if elapsed < 80*time.Millisecond {
			t.Errorf("3 calls completed in %v, expected at least ~100ms for rate limit 20/s burst 1", elapsed)
		}
	})

	t.Run("propagates underlying error", func(t *testing.T) {
		wantErr := errors.New("rpc failed")
		mock := &mockClientConn{invokeErr: wantErr}
		conn := ConnectionWithRateLimit(mock, 100, 1)

		if err := conn.Invoke(t.Context(), "/test/Method", nil, nil); !errors.Is(err, wantErr) {
			t.Errorf("Invoke() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		mock := &mockClientConn{}
		conn := ConnectionWithRateLimit(mock, 1, 1)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		if err := conn.Invoke(ctx, "/test/Method", nil, nil); err == nil {
			t.Errorf("Invoke() with cancelled context succeeded, want error")
		}
		if got := mock.invokeCount.Load(); got != 0 {
			t.Errorf("invokeCount = %d, want 0", got)
		}
	})
}

func TestConnectionWithRateLimit_NewStream(t *testing.T) {
	t.Run("success and rate limiting", func(t *testing.T) {
		mock := &mockClientConn{}
		// Rate of 20/sec (50ms per token), burst 1.
		conn := ConnectionWithRateLimit(mock, 20, 1)

		start := time.Now()
		for range 3 {
			stream, err := conn.NewStream(t.Context(), &grpc.StreamDesc{}, "/test/Stream")
			if err != nil {
				t.Fatalf("NewStream() unexpected error: %v", err)
			}
			if stream == nil {
				t.Fatalf("NewStream() returned nil stream")
			}
		}
		elapsed := time.Since(start)

		if got := mock.streamCount.Load(); got != 3 {
			t.Errorf("streamCount = %d, want 3", got)
		}
		if elapsed < 80*time.Millisecond {
			t.Errorf("3 calls completed in %v, expected at least ~100ms for rate limit 20/s burst 1", elapsed)
		}
	})

	t.Run("propagates underlying error", func(t *testing.T) {
		wantErr := errors.New("stream failed")
		mock := &mockClientConn{streamErr: wantErr}
		conn := ConnectionWithRateLimit(mock, 100, 1)

		if _, err := conn.NewStream(t.Context(), &grpc.StreamDesc{}, "/test/Stream"); !errors.Is(err, wantErr) {
			t.Errorf("NewStream() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		mock := &mockClientConn{}
		conn := ConnectionWithRateLimit(mock, 1, 1)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		if _, err := conn.NewStream(ctx, &grpc.StreamDesc{}, "/test/Stream"); err == nil {
			t.Errorf("NewStream() with cancelled context succeeded, want error")
		}
		if got := mock.streamCount.Load(); got != 0 {
			t.Errorf("streamCount = %d, want 0", got)
		}
	})
}
