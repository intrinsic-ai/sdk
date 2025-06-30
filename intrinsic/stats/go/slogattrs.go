// Copyright 2023 Intrinsic Innovation LLC

// Package slogattrs provides a slog.Handler to populate log attributes from the context.
// See NewHandler for a usage example.
package slogattrs

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"go.opencensus.io/trace"
)

type slogAttrsCtxKey string

const (
	slogFields slogAttrsCtxKey = "slog_attrs"
)

// ContextHandler wraps a slog.Handler to populate log attributes from the context.
// Example:
//
//	  // in your main function:
//	  h := &slogattrs.ContextHandler{Handler: slog.NewTextHandler(os.Stdout, nil)}
//	  logger := slog.New(h)
//	  slog.SetDefault(logger)
//
//		// in your application code:
//	  ctx = slogattrs.Append(ctx, slog.String("trace_id", traceID))
//	  slog.InfoContext(ctx, "Hello World!")
//
// Note: use NewHandler instead of this.
type ContextHandler struct {
	slog.Handler
	ProjectName string
}

// NewHandler wraps a slog.Handler to populate log attributes from the context and
// to attach OT trace/span information if present and span is recording events.
// No need to add trace information to context, it will be added automatically
// if present.
//
// Usage:
//
//	h := slogattrs.NewHandler("my_gcp_project", slog.NewTextHandler(os.Stdout, nil))
//	logger := slog.New(h)
//	slog.SetDefault(logger)
func NewHandler(projectName string, handler slog.Handler) slog.Handler {
	return &ContextHandler{
		Handler:     handler,
		ProjectName: projectName,
	}
}

// SetDefaultLogger reconfigures slog default logger to use ContextHandler
func SetDefaultLogger(projectName string, options *slog.HandlerOptions) {
	handler := NewHandler(projectName, slog.NewTextHandler(os.Stdout, options))
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// Handle populates log attributes from the context.
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		r.AddAttrs(attrs...)
	}
	span := trace.FromContext(ctx)
	if span != nil && span.IsRecordingEvents() {
		// We are going to attach trace information IFF span is recording events.
		spanContext := span.SpanContext()
		r.Add(
			// See: https://cloud.google.com/logging/docs/structured-logging#special-payload-fields
			slog.String("logging.googleapis.com/trace", fmt.Sprintf("projects/%s/traces/%s", h.ProjectName, spanContext.TraceID)),
			slog.String("logging.googleapis.com/spanId", spanContext.SpanID.String()),
			slog.Bool("logging.googleapis.com/traceSampled", span.IsRecordingEvents()), // always true in this context
		)
	}
	return h.Handler.Handle(ctx, r)
}

// Append adds a slog attribute to the current context.
func Append(ctx context.Context, attrs ...slog.Attr) context.Context {
	if len(attrs) == 0 {
		return ctx
	}
	v, _ := ctx.Value(slogFields).([]slog.Attr)
	return context.WithValue(ctx, slogFields, append(v, attrs...))
}

// Err adds an error as slog attribute.
func Err(e error) slog.Attr {
	return slog.String("error", e.Error())
}

// ReplaceAttr replaces the key of a slog.Attr with the corresponding key for log explorer.
// Also shortens the source location and function.
func ReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.MessageKey:
		a.Key = "message"
	case slog.TimeKey:
		a.Key = "timestamp"
	case slog.LevelKey:
		a.Key = "severity"
	case slog.SourceKey:
		source, ok := a.Value.Any().(*slog.Source)
		if ok && source != nil {
			source.File = filepath.Base(source.File)
			f := strings.Split(source.Function, "/")
			source.Function = f[len(f)-1]
		}
	}
	return a
}
