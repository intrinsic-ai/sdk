// Copyright 2023 Intrinsic Innovation LLC

package slogattrs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"go.opencensus.io/trace"
)

const (
	metadataServerBaseURL = "http://metadata.google.internal/computeMetadata/v1/"
	projectIDPath         = "project/project-id"
	metadataFlavorHeader  = "Metadata-Flavor"
	googleMetadataFlavor  = "Google"
)

func getMetadata(ctx context.Context, path string) (string, error) {
	ctx, span := trace.StartSpan(ctx, "metadata.getMetadata")
	defer span.End()
	if span.IsRecordingEvents() {
		span.AddAttributes(trace.StringAttribute("path", path))
	}

	client := &http.Client{
		Timeout: 5 * time.Second, // Set a timeout for the request
	}

	req, err := http.NewRequestWithContext(ctx, "GET", metadataServerBaseURL+path, nil)
	if err != nil {
		return "", spanSetErrorStatus(span, fmt.Errorf("failed to create request: %w", err))
	}

	// Essential header for accessing the metadata server
	req.Header.Set(metadataFlavorHeader, googleMetadataFlavor)

	resp, err := client.Do(req)
	if err != nil {
		return "", spanSetErrorStatus(span, fmt.Errorf("failed to make HTTP request to metadata server: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		spanStatusFromResponse(span, resp)
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("received non-OK status from metadata server: %s, body: %s", resp.Status, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", spanSetErrorStatus(span, fmt.Errorf("failed to read response body: %w", err))
	}

	return string(bodyBytes), nil
}

// CloudProjectName returns a GCP Project name associated with this execution.
// It obtains this information by contacting the metadata server.
func CloudProjectName(ctx context.Context) (string, error) {
	// Allows to read project from GCP defaults if found.
	if value, ok := os.LookupEnv("GOOGLE_CLOUD_PROJECT"); ok && value != "" {
		return value, nil
	}
	return getMetadata(ctx, projectIDPath)
}

func spanStatusFromResponse(span *trace.Span, resp *http.Response) {
	if span != nil && span.IsRecordingEvents() && resp.StatusCode >= 300 {
		span.Annotatef([]trace.Attribute{
			trace.StringAttribute("method", resp.Request.Method),
			trace.StringAttribute("url", resp.Request.URL.String()),
			trace.Int64Attribute("code", int64(resp.StatusCode)),
			trace.StringAttribute("status", resp.Status),
		}, "http: %s", resp.Status)
		span.SetStatus(trace.Status{
			Code:    int32(resp.StatusCode),
			Message: resp.Status,
		})
	}
}

func spanSetErrorStatus(span *trace.Span, err error) error {
	if span != nil && span.IsRecordingEvents() && err != nil {
		span.SetStatus(trace.Status{
			Code:    13,
			Message: err.Error(),
		})
	}
	return err
}
