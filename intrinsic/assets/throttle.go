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

// Package throttle provides primitives to limit concurrency and rates.
package throttle

import (
	"context"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

// rateLimitedConn wraps a grpc.ClientConnInterface to rate-limit outbound RPC calls.
type rateLimitedConn struct {
	grpc.ClientConnInterface
	limiter *rate.Limiter
}

// ConnectionWithRateLimit wraps conn to enforce rate limits on all Invoke and NewStream calls.
//
// limit specifies the maximum request rate per second. If limit <= 0, conn is returned directly.
// burst specifies the maximum burst size permitted by the rate limiter.
//
// If conn is nil, ConnectionWithRateLimit returns nil.
func ConnectionWithRateLimit(conn grpc.ClientConnInterface, limit float64, burst int) grpc.ClientConnInterface {
	if conn == nil || limit <= 0 {
		return conn
	}

	return ConnectionWithRateLimiter(conn, rate.NewLimiter(rate.Limit(limit), burst))
}

// ConnectionWithRateLimiter wraps conn with limiter to enforce rate limits on all Invoke and
// NewStream calls.
//
// If conn or limiter is nil, conn is returned directly.
func ConnectionWithRateLimiter(conn grpc.ClientConnInterface, limiter *rate.Limiter) grpc.ClientConnInterface {
	if conn == nil || limiter == nil {
		return conn
	}

	return &rateLimitedConn{
		ClientConnInterface: conn,
		limiter:             limiter,
	}
}

func (c *rateLimitedConn) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	if err := c.limiter.Wait(ctx); err != nil {
		return err
	}
	return c.ClientConnInterface.Invoke(ctx, method, args, reply, opts...)
}

func (c *rateLimitedConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	return c.ClientConnInterface.NewStream(ctx, desc, method, opts...)
}
