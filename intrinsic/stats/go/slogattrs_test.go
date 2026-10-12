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
	"bytes"
	"context"
	"encoding/json"
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

func TestLabels(t *testing.T) {
	attr := Labels(slog.String("foo", "bar"), slog.Int("baz", 123))
	if attr.Key != CloudLoggingLabelsKey {
		t.Errorf("Labels() key = %q, want %q", attr.Key, CloudLoggingLabelsKey)
	}

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	logger.Info("test", attr)

	var payload map[string]any
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	labels, ok := payload[CloudLoggingLabelsKey].(map[string]any)
	if !ok {
		t.Fatalf("expected map under %q, got %T: %v", CloudLoggingLabelsKey, payload[CloudLoggingLabelsKey], payload)
	}
	if got := labels["foo"]; got != "bar" {
		t.Errorf("labels[foo] = %v, want bar", got)
	}
	if got := labels["baz"]; got != float64(123) {
		t.Errorf("labels[baz] = %v, want 123", got)
	}
}

func TestAuditLog(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	logger.Info("audited event", AuditLog(slog.String("subsystem", "accounts")))

	var payload map[string]any
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	labels, ok := payload[CloudLoggingLabelsKey].(map[string]any)
	if !ok {
		t.Fatalf("expected map under %q, got %T: %v", CloudLoggingLabelsKey, payload[CloudLoggingLabelsKey], payload)
	}
	if got := labels["is_audit_log"]; got != "true" {
		t.Errorf("labels[is_audit_log] = %v, want true", got)
	}
	if got := labels["subsystem"]; got != "accounts" {
		t.Errorf("labels[subsystem] = %v, want accounts", got)
	}
}
