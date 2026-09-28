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
	"io"
	"net/http"
	"os"
	"time"

	"intrinsic/stats/go/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	metadataServerBaseURL = "http://metadata.google.internal/computeMetadata/v1/"
	projectIDPath         = "project/project-id"
	metadataFlavorHeader  = "Metadata-Flavor"
	googleMetadataFlavor  = "Google"
)

func getMetadata(ctx context.Context, path string) (string, error) {
	ctx, span := telemetry.StartSpan(ctx, "metadata.getMetadata", trace.WithAttributes(
		attribute.String("path", path),
	))
	defer span.End()

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

func spanStatusFromResponse(span trace.Span, resp *http.Response) {
	if span != nil && span.IsRecording() && resp.StatusCode >= 300 {
		span.AddEvent("http: "+resp.Status, trace.WithAttributes(
			attribute.String("method", resp.Request.Method),
			attribute.String("url", resp.Request.URL.String()),
			attribute.Int64("code", int64(resp.StatusCode)),
			attribute.String("status", resp.Status),
		))
		span.SetStatus(codes.Error, resp.Status)
	}
}

func spanSetErrorStatus(span trace.Span, err error) error {
	if span != nil && span.IsRecording() && err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}
