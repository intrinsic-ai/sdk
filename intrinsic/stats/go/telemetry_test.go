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

package telemetry

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"syscall"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	oteltrace "go.opentelemetry.io/otel/trace"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// newRecordingTracerProvider returns a TracerProvider whose spans are always
// sampled, together with the recorder that collects them once they end.
func newRecordingTracerProvider(t *testing.T) (*sdktrace.TracerProvider, *tracetest.SpanRecorder) {
	t.Helper()
	recorder := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(recorder),
	)
	t.Cleanup(func() {
		// t.Context() is canceled just before cleanup functions run, so the
		// shutdown needs a context of its own.
		if err := tp.Shutdown(context.Background()); err != nil {
			t.Errorf("TracerProvider.Shutdown() returned an unexpected error: %v", err)
		}
	})
	return tp, recorder
}

// newRecordingTracer returns a tracer for starting spans and a recorder for collecting them.
func newRecordingTracer(t *testing.T) (oteltrace.Tracer, *tracetest.SpanRecorder) {
	t.Helper()
	tp, recorder := newRecordingTracerProvider(t)
	return tp.Tracer("intrinsic/stats/go/telemetry_test"), recorder
}

// onlyEndedSpan returns the single span captured by the recorder.
func onlyEndedSpan(t *testing.T, recorder *tracetest.SpanRecorder) sdktrace.ReadOnlySpan {
	t.Helper()
	spans := recorder.Ended()
	if len(spans) != 1 {
		t.Fatalf("recorder captured %d ended spans, want 1", len(spans))
	}
	return spans[0]
}

// spanAttributes renders the attributes of a span as a name to value map.
func spanAttributes(s sdktrace.ReadOnlySpan) map[string]string {
	attrs := make(map[string]string)
	for _, kv := range s.Attributes() {
		attrs[string(kv.Key)] = kv.Value.String()
	}
	return attrs
}

// exceptionMessages returns the messages of the exception events recorded on a span.
func exceptionMessages(s sdktrace.ReadOnlySpan) []string {
	var messages []string
	for _, event := range s.Events() {
		if event.Name != "exception" {
			continue
		}
		for _, kv := range event.Attributes {
			if kv.Key == "exception.message" {
				messages = append(messages, kv.Value.String())
			}
		}
	}
	return messages
}

type WasCalledHandler struct {
	called bool
}

func (h *WasCalledHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.called = true
	// check Hijacker
	whi, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, fmt.Sprintf("response writer %T does not implement Hijacker", w), http.StatusInternalServerError)
		return
	}
	// Response recorder does not support Hijacker, but at least we make sure the call chain
	// is working correctly.
	_, _, err := whi.Hijack()
	wantErr := "http.Hijacker interface not supported by *httptest.ResponseRecorder"
	if err.Error() != wantErr {
		http.Error(w, fmt.Sprintf("got error %q, want %q", err.Error(), wantErr), http.StatusInternalServerError)
	}
}

type TraceIDHandlerTest struct {
	desc    string
	r       *http.Request
	wantSet bool
}

func mustCreateNewRequestWithSpan(ctx context.Context, t *testing.T, method, url string, body io.Reader) *http.Request {
	t.Helper()
	tracer, _ := newRecordingTracer(t)
	spanCtx, span := tracer.Start(ctx, "test")
	t.Cleanup(func() { span.End() })
	r, err := http.NewRequestWithContext(spanCtx, method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestTraceIDHandler(t *testing.T) {
	tests := []*TraceIDHandlerTest{
		{
			desc:    "no header set",
			r:       httptest.NewRequest(http.MethodGet, "/", nil),
			wantSet: false,
		},
		{
			desc:    "header set",
			r:       mustCreateNewRequestWithSpan(context.Background(), t, "GET", "/", nil),
			wantSet: true,
		},
	}
	for _, test := range tests {
		runTraceIDHandlerTest(t, test)
	}
}

func runTraceIDHandlerTest(t *testing.T, test *TraceIDHandlerTest) {
	rr := httptest.NewRecorder()
	wantCalled := &WasCalledHandler{}
	h := TraceIDHandler(wantCalled)
	h.ServeHTTP(rr, test.r)
	// the handler should never block so next must be called always
	if !wantCalled.called {
		t.Errorf("wanted wrapped handler called but did not happen")
	}
	if rr.Code != http.StatusOK {
		if body, err := io.ReadAll(rr.Body); err != nil {
			t.Error("failed to read body")
		} else {
			t.Errorf("got code %v, want %v, body: %q", rr.Code, http.StatusOK, body)
		}
	}
	if got := (rr.Header().Get("X-Intrinsic-TraceID") != ""); got != test.wantSet {
		t.Errorf("got header %v, want %v", got, test.wantSet)
	}
}

func TestHijacker(t *testing.T) {
	tid := &TraceIDWriter{}
	if _, ok := any(tid).(http.Hijacker); !ok {
		t.Fatalf("%T does not implement http.Hijacker", tid)
	}
}

func TestFlusher(t *testing.T) {
	tid := &TraceIDWriter{}
	if _, ok := any(tid).(http.Flusher); !ok {
		t.Fatalf("%T does not implement http.Flusher", tid)
	}
}

func TestEnableCloudTracing(t *testing.T) {
	got := config{}
	EnableCloudTracing()(&got)
	want := config{}
	WithCloudTracing(true)(&want)
	if !cmp.Equal(got, want) {
		t.Errorf("got config %+v, want %+v", got, want)
	}
}

func TestWithCloudTracingEnabled(t *testing.T) {
	got := config{}
	WithCloudTracing(true)(&got)
	want := config{
		TracingCfg: &tracingConfig{
			Enabled:         true,
			CloudDeployment: true,
			Probability:     1,
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("got config %+v, want %+v", got, want)
	}
}

func TestWithCloudTracingDisabled(t *testing.T) {
	got := config{}
	WithCloudTracing(false)(&got)
	want := config{
		TracingCfg: &tracingConfig{
			Enabled:         false,
			CloudDeployment: true,
			Probability:     1,
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("got config %+v, want %+v", got, want)
	}
}

func TestEnableTracing(t *testing.T) {
	got := config{}
	EnableTracing()(&got)
	want := config{}
	WithTracing(true)(&want)
	if !cmp.Equal(got, want) {
		t.Errorf("got config %+v, want %+v", got, want)
	}
}

func TestWithTracingEnabled(t *testing.T) {
	got := config{}
	WithTracing(true)(&got)
	want := config{
		TracingCfg: &tracingConfig{
			Enabled:         true,
			CloudDeployment: false,
			Probability:     1,
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("got config %+v, want %+v", got, want)
	}
}

func TestWithTracingDisabled(t *testing.T) {
	got := config{}
	WithTracing(false)(&got)
	want := config{
		TracingCfg: &tracingConfig{
			Enabled:         false,
			CloudDeployment: false,
			Probability:     1,
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("got config %+v, want %+v", got, want)
	}
}

func TestWithProjectName(t *testing.T) {
	pn := "a project name"
	got := config{}
	WithProjectName(pn)(&got)
	if got, want := got.TracingCfg.ProjectName, pn; got != want {
		t.Errorf("got config %+v, want %+v", got, want)
	}
}

func TestWithServiceName(t *testing.T) {
	sn := "a service name"
	got := config{}
	WithServiceName(sn)(&got)
	if got, want := got.TracingCfg.ServiceName, sn; got != want {
		t.Errorf("got config %+v, want %+v", got, want)
	}
}

func TestWithProbability(t *testing.T) {
	p := 0.25
	got := config{}
	WithProbability(p)(&got)
	if got, want := got.TracingCfg.Probability, p; got != want {
		t.Errorf("got config %+v, want %+v", got, want)
	}
}

func TestWithMetricsEnabled(t *testing.T) {
	got := config{}
	var port int64 = 9101
	WithMetrics(true, port)(&got)
	// Manual testing since cmp.Equals fatal crashed.
	if got, want := got.MetricsCfg.Enabled, true; got != want {
		t.Errorf("got config %+v, want %+v", got, want)
	}
	if got, want := got.MetricsCfg.MetricsPath, "/metrics"; got != want {
		t.Errorf("got config %+v, want %+v", got, want)
	}
	if got, want := got.MetricsCfg.MetricsPort, port; got != want {
		t.Errorf("got config %+v, want %+v", got, want)
	}
}

func TestWithMetricsDisabled(t *testing.T) {
	got := config{}
	var port int64 = 9101
	WithMetrics(false, port)(&got)
	// Manual testing since cmp.Equals fatal crashed.
	if got, want := got.MetricsCfg.Enabled, false; got != want {
		t.Errorf("got config %+v, want %+v", got, want)
	}
	if got, want := got.MetricsCfg.MetricsPath, "/metrics"; got != want {
		t.Errorf("got config %+v, want %+v", got, want)
	}
	if got, want := got.MetricsCfg.MetricsPort, port; got != want {
		t.Errorf("got config %+v, want %+v", got, want)
	}
}

func TestEnableMetrics(t *testing.T) {
	got := config{}
	var port int64 = 9101
	EnableMetrics(port)(&got)
	// Manual testing since cmp.Equals fatal crashed.
	if got, want := got.MetricsCfg.Enabled, true; got != want {
		t.Errorf("got config %+v, want %+v", got, want)
	}
	if got, want := got.MetricsCfg.MetricsPath, "/metrics"; got != want {
		t.Errorf("got config %+v, want %+v", got, want)
	}
	if got, want := got.MetricsCfg.MetricsPort, port; got != want {
		t.Errorf("got config %+v, want %+v", got, want)
	}
}

func TestAddSpanAttributes(t *testing.T) {
	tracer, recorder := newRecordingTracer(t)

	ctx, span := tracer.Start(t.Context(), "test-span")
	AddSpanAttributes(ctx, attribute.String("key1", "value1"), attribute.Bool("key2", true))
	span.End()

	want := map[string]string{
		"key1": "value1",
		"key2": "true",
	}
	got := spanAttributes(onlyEndedSpan(t, recorder))
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("AddSpanAttributes() recorded unexpected attributes (-want +got):\n%s", diff)
	}
}

func TestAddSpanAttributesWithoutSpan(t *testing.T) {
	// There is no span in the context, so this must not panic.
	AddSpanAttributes(t.Context(), attribute.String("key1", "value1"))
}

func TestSetSpanError(t *testing.T) {
	tracer, recorder := newRecordingTracer(t)

	_, span := tracer.Start(t.Context(), "test-span")
	SetSpanError(span, grpccodes.NotFound, "loading item", errors.New("no such key"))
	span.End()

	ended := onlyEndedSpan(t, recorder)
	if got, want := ended.Status().Code, codes.Error; got != want {
		t.Errorf("SetSpanError() set status code %v, want %v", got, want)
	}
	if got, want := ended.Status().Description, "loading item: no such key"; got != want {
		t.Errorf("SetSpanError() set status description %q, want %q", got, want)
	}
	want := map[string]string{"error.type": "NotFound"}
	if diff := cmp.Diff(want, spanAttributes(ended)); diff != "" {
		t.Errorf("SetSpanError() recorded unexpected attributes (-want +got):\n%s", diff)
	}
	wantExceptions := []string{"loading item: no such key"}
	if diff := cmp.Diff(wantExceptions, exceptionMessages(ended)); diff != "" {
		t.Errorf("SetSpanError() recorded unexpected exception events (-want +got):\n%s", diff)
	}
}

func TestSetSpanErrorf(t *testing.T) {
	tracer, recorder := newRecordingTracer(t)

	_, span := tracer.Start(t.Context(), "test-span")
	SetSpanErrorf(span, grpccodes.InvalidArgument, "value %d is out of range", 42)
	span.End()

	ended := onlyEndedSpan(t, recorder)
	if got, want := ended.Status().Code, codes.Error; got != want {
		t.Errorf("SetSpanErrorf() set status code %v, want %v", got, want)
	}
	if got, want := ended.Status().Description, "value 42 is out of range"; got != want {
		t.Errorf("SetSpanErrorf() set status description %q, want %q", got, want)
	}
	want := map[string]string{"error.type": "InvalidArgument"}
	if diff := cmp.Diff(want, spanAttributes(ended)); diff != "" {
		t.Errorf("SetSpanErrorf() recorded unexpected attributes (-want +got):\n%s", diff)
	}
	wantExceptions := []string{"value 42 is out of range"}
	if diff := cmp.Diff(wantExceptions, exceptionMessages(ended)); diff != "" {
		t.Errorf("SetSpanErrorf() recorded unexpected exception events (-want +got):\n%s", diff)
	}
}

func TestSpanStatusWithError(t *testing.T) {
	tests := []struct {
		desc            string
		err             error
		wantCode        codes.Code
		wantDescription string
		wantAttributes  map[string]string
		wantExceptions  []string
	}{
		{
			desc:            "grpc_status_error",
			err:             status.Error(grpccodes.PermissionDenied, "not allowed"),
			wantCode:        codes.Error,
			wantDescription: "not allowed",
			wantAttributes:  map[string]string{"error.type": "PermissionDenied"},
			wantExceptions:  []string{"not allowed"},
		},
		{
			desc:            "plain_error_is_reported_as_unknown",
			err:             errors.New("something broke"),
			wantCode:        codes.Error,
			wantDescription: "something broke",
			wantAttributes:  map[string]string{"error.type": "Unknown"},
			wantExceptions:  []string{"something broke"},
		},
		{
			desc:            "nil_error_leaves_the_span_untouched",
			err:             nil,
			wantCode:        codes.Unset,
			wantDescription: "",
			wantAttributes:  map[string]string{},
			wantExceptions:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			tracer, recorder := newRecordingTracer(t)

			_, span := tracer.Start(t.Context(), "test-span")
			if got := SpanStatusWithError(span, test.err); got != test.err {
				t.Errorf("SpanStatusWithError() returned error %v, want %v", got, test.err)
			}
			span.End()

			ended := onlyEndedSpan(t, recorder)
			if got, want := ended.Status().Code, test.wantCode; got != want {
				t.Errorf("SpanStatusWithError() set status code %v, want %v", got, want)
			}
			if got, want := ended.Status().Description, test.wantDescription; got != want {
				t.Errorf("SpanStatusWithError() set status description %q, want %q", got, want)
			}
			if diff := cmp.Diff(test.wantAttributes, spanAttributes(ended)); diff != "" {
				t.Errorf("SpanStatusWithError() recorded unexpected attributes (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(test.wantExceptions, exceptionMessages(ended)); diff != "" {
				t.Errorf("SpanStatusWithError() recorded unexpected exception events (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOtelHTTPHelpersPropagateTraceContext(t *testing.T) {
	serverTP, serverRecorder := newRecordingTracerProvider(t)
	clientTP, clientRecorder := newRecordingTracerProvider(t)

	mux := http.NewServeMux()
	mux.HandleFunc("/items/list", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	srv := httptest.NewServer(NewOtelHTTPHandler(mux, otelhttp.WithTracerProvider(serverTP)))
	defer srv.Close()

	client := &http.Client{
		Transport: NewOtelHTTPTransport(nil, otelhttp.WithTracerProvider(clientTP)),
	}
	resp, err := client.Get(srv.URL + "/items/list")
	if err != nil {
		t.Fatalf("client.Get() returned an unexpected error: %v", err)
	}
	if _, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("reading the response body returned an unexpected error: %v", err)
	}
	resp.Body.Close()

	serverSpan := onlyEndedSpan(t, serverRecorder)
	clientSpan := onlyEndedSpan(t, clientRecorder)

	if serverSpan.Parent().SpanID() != clientSpan.SpanContext().SpanID() {
		t.Errorf("server span parent is %v, want the client span %v", serverSpan.Parent().SpanID(), clientSpan.SpanContext().SpanID())
	}
	if serverSpan.SpanContext().TraceID() != clientSpan.SpanContext().TraceID() {
		t.Errorf("server span trace is %v, want the client trace %v", serverSpan.SpanContext().TraceID(), clientSpan.SpanContext().TraceID())
	}
	if serverSpan.Name() != "/items/list" {
		t.Errorf("server span is named %q, want %q", serverSpan.Name(), "/items/list")
	}
	if clientSpan.Name() != "/items/list" {
		t.Errorf("client span is named %q, want %q", clientSpan.Name(), "/items/list")
	}
}

func TestTelemetryMetrics(t *testing.T) {
	prevMP := otel.GetMeterProvider()
	opencensusMeasure := stats.Int64("opencensus_metric", "", "")
	opencensusView := &view.View{
		Measure:     opencensusMeasure,
		Name:        opencensusMeasure.Name(),
		Description: opencensusMeasure.Description(),
		Aggregation: view.Count(),
	}

	meter := otel.Meter("testing_otel_meter")
	otelMetric, err := meter.Int64Counter("opentelemetry_metric")
	if err != nil {
		t.Fatalf("meter.Int64Counter(\"opentelemetry_metric\") returned an unexpected error: %v", err)
	}

	tele := Initialize(
		EnableMetrics(9101),
		WithViews([]*view.View{opencensusView}),
	)

	t.Cleanup(func() {
		view.Unregister(opencensusView)
		otel.SetMeterProvider(prevMP)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tele.Shutdown(ctx); err != nil {
			t.Errorf("Telemetry.Shutdown() returned an unexpected error: %v", err)
		}
	})

	otelMetric.Add(t.Context(), 1)
	stats.Record(t.Context(), opencensusMeasure.M(1))

	// Flush OpenCensus's asynchronous stats worker before scraping /metrics.
	if _, err := view.RetrieveData(opencensusView.Name); err != nil {
		t.Fatalf("view.RetrieveData() returned an unexpected error: %v", err)
	}

	var prometheusResp string
	retries := 0
	for {
		if resp, err := http.Get("http://localhost:9101/metrics"); err == nil {
			// No error; break out of retry loop.
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("reading metrics body: %v", err)
			}
			prometheusResp = string(body)
			break
		} else if !errors.Is(err, syscall.ECONNREFUSED) {
			t.Fatalf("error making http request: %s", err)
		}
		// Server is not up yet.
		if retries > 30 {
			t.Fatal("too many retries waiting for server to come up")
		}
		retries++
		time.Sleep(time.Duration(retries) * time.Millisecond)
	}

	// Both metrics must be exported with value 1.
	for _, re := range []*regexp.Regexp{
		regexp.MustCompile(`(?m)^opencensus_metric\{.*\} 1$`),
		regexp.MustCompile(`(?m)^opentelemetry_metric\{.*\} 1$`),
	} {
		if !re.MatchString(prometheusResp) {
			t.Errorf("/metrics output does not match %q:\n%s", re.String(), prometheusResp)
		}
	}
}
