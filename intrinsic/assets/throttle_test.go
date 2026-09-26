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
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

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
