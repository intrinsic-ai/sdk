// Copyright 2023 Intrinsic Innovation LLC

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
