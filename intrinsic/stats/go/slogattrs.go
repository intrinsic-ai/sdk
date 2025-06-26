// Copyright 2023 Intrinsic Innovation LLC

// Package slogattrs provides a slog.Handler to populate log attributes from the context.
// See ContextHandler for a usage example.
package slogattrs

import (
	"context"
	"log/slog"
	"path/filepath"
	"strings"
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
type ContextHandler struct {
	slog.Handler
}

// Handle populates log attributes from the context.
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		r.AddAttrs(attrs...)
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
	return slog.String("Error", e.Error())
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
