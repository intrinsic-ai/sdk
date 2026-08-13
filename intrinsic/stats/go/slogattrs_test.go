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

package slogattrs

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type FakeHandler struct {
	slog.Handler
}

func (h *FakeHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := []slog.Attr{}
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})
	want := []slog.Attr{slog.String("testkey", "testvalue")}
	if diff := cmp.Diff(attrs, want); diff != "" {
		return fmt.Errorf("unexpected attrs, diff (-want +got):\n%s", diff)
	}
	return nil
}

func TestSlogAttrs(t *testing.T) {
	ctx := context.Background()
	ctx = Append(ctx, slog.String("testkey", "testvalue"))
	h := ContextHandler{Handler: &FakeHandler{}}
	r := slog.Record{}
	err := h.Handle(ctx, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
