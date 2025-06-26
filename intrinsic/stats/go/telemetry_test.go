// Copyright 2023 Intrinsic Innovation LLC

package telemetry

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
)

type ValidSessionIdentifierTest struct {
	value string
	want  bool
}

func TestValidSessionIdentifier(t *testing.T) {
	tests := []ValidSessionIdentifierTest{
		ValidSessionIdentifierTest{
			value: "",
			want:  false,
		},
		ValidSessionIdentifierTest{
			value: "\\.12323",
			want:  false,
		},
		ValidSessionIdentifierTest{
			value: "test",
			want:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			r := validateTRSID.MatchString(test.value)
			if r != test.want {
				t.Errorf("identifier %s: got %v, want %v", test.value, r, test.want)
			}
		})
	}
}

type TraceOnCookieTest struct {
	desc        string
	cookie      *http.Cookie
	wantSampler bool
}

func TestTraceOnCookie(t *testing.T) {
	tests := []TraceOnCookieTest{
		TraceOnCookieTest{
			desc:        "empty request",
			cookie:      nil,
			wantSampler: false,
		},
		TraceOnCookieTest{
			desc:        "with cookie",
			cookie:      &http.Cookie{Name: "trsid", Value: "asd"},
			wantSampler: true,
		},
		TraceOnCookieTest{
			desc:        "with empty cookie",
			cookie:      &http.Cookie{Name: "trsid", Value: ""},
			wantSampler: false,
		},
		TraceOnCookieTest{
			desc:        "with invalid cookie",
			cookie:      &http.Cookie{Name: "trsid", Value: "asd//asd"},
			wantSampler: false,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			runTraceOnCookieTest(t, &test)
		})
	}
}

func runTraceOnCookieTest(t *testing.T, test *TraceOnCookieTest) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if test.cookie != nil {
		r.AddCookie(test.cookie)
	}
	if err != nil {
		t.Fatal(err)
	}

	if gotSampler := (TraceOnCookie(r).Sampler != nil); gotSampler != test.wantSampler {
		t.Errorf("got sampler %v, want sampler %v", gotSampler, test.wantSampler)
	}
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
	spanCtx, _ := trace.StartSpanWithRemoteParent(ctx, "test",
		trace.SpanContext{TraceOptions: 0x1}, // trace span
	)
	r, err := http.NewRequestWithContext(spanCtx, method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestTraceIDHandler(t *testing.T) {
	tests := []*TraceIDHandlerTest{
		&TraceIDHandlerTest{
			desc:    "no header set",
			r:       httptest.NewRequest(http.MethodGet, "/", nil),
			wantSet: false,
		},
		&TraceIDHandlerTest{
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

type AddSpanTRSIDHandlerTest struct {
	desc   string
	cookie *http.Cookie
	trsid  string
}

func TestAddSpanTRSIDHandler(t *testing.T) {
	tests := []AddSpanTRSIDHandlerTest{
		AddSpanTRSIDHandlerTest{
			desc:   "with cookie",
			cookie: &http.Cookie{Name: "X-Intrinsic-Tracing-Session", Value: "asd"},
			trsid:  "asd",
		},
		AddSpanTRSIDHandlerTest{
			desc:   "without cookie",
			cookie: nil,
		},
		AddSpanTRSIDHandlerTest{
			desc:   "wit invalid cookie",
			cookie: &http.Cookie{Name: "X-Intrinsic-Tracing-Session", Value: "asdasd//asd"},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) { runAddSpanTRSIDHandlerTest(t, &test) })
	}
}

func runAddSpanTRSIDHandlerTest(t *testing.T, test *AddSpanTRSIDHandlerTest) {
	// setup
	rr := httptest.NewRecorder()
	wantCalled := &WasCalledHandler{}
	h := AddSpanTRSIDHandler(wantCalled)
	ctx, _ := trace.StartSpan(context.Background(), "")
	req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	if test.cookie != nil {
		req.AddCookie(test.cookie)
	}
	// request
	h.ServeHTTP(rr, req)
	// check
	if !wantCalled.called {
		t.Error("wanted wrapped handler called but did not happen")
	}
}

func mustRequestWithCookie(t *testing.T, method string, path string, body io.Reader, cookie *http.Cookie) *http.Request {
	r := httptest.NewRequest(method, path, body)
	if cookie != nil {
		r.AddCookie(cookie)
	}
	return r
}

type TracingEndpointsTest struct {
	desc           string
	r              *http.Request
	wantSetCookies int
	wantBody       string
	wantCode       int
}

func TestTracingEndpoints(t *testing.T) {
	// static uuid
	oldUUIDNew := uuidNew
	uuidNew = func() string { return "test1234" }
	t.Cleanup(func() { uuidNew = oldUUIDNew })
	tests := []*TracingEndpointsTest{
		&TracingEndpointsTest{
			desc:     "noop",
			r:        httptest.NewRequest(http.MethodGet, "/tracing", nil),
			wantCode: http.StatusBadRequest,
			wantBody: "invalid tracing operation\n",
		},
		&TracingEndpointsTest{
			desc:     "trsid with no session",
			r:        httptest.NewRequest(http.MethodGet, "/tracing/trsid", nil),
			wantBody: "no tracing session cookie found",
			wantCode: http.StatusNotFound,
		},
		&TracingEndpointsTest{
			desc: "trsid with session",
			r: mustRequestWithCookie(t, "GET", "/tracing/trsid", nil,
				&http.Cookie{Name: "trsid", Value: "testtrsid"}),
			wantCode: http.StatusOK,
			wantBody: "testtrsid",
		},
		&TracingEndpointsTest{
			desc:           "enable",
			r:              httptest.NewRequest(http.MethodGet, "/tracing/enable", nil),
			wantCode:       http.StatusOK,
			wantSetCookies: 1,
			wantBody:       "created tracing session test1234 for you",
		},
		&TracingEndpointsTest{
			desc:           "disable",
			r:              httptest.NewRequest(http.MethodGet, "/tracing/disable", nil),
			wantCode:       http.StatusOK,
			wantSetCookies: 1,
			wantBody:       "disabled your tracing session",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) { runTracingEndpointsTest(t, test) })
	}
}

func runTracingEndpointsTest(t *testing.T, test *TracingEndpointsTest) {
	// setup
	rr := httptest.NewRecorder()
	h, err := RegisterTracingEndpoints(http.NewServeMux(), "/tracing")
	if err != nil {
		t.Fatal(err)
	}
	// request
	h.tracingHandler(rr, test.r)
	// check
	if rr.Code != test.wantCode {
		t.Errorf("got code %v, want %v", rr.Code, test.wantCode)
	}
	gotBody := string(rr.Body.Bytes())
	if gotBody != test.wantBody {
		t.Errorf("got body %q, want %q", gotBody, test.wantBody)
	}
	setCookie := rr.Result().Header["Set-Cookie"]
	if len(setCookie) != test.wantSetCookies {
		t.Errorf("got cookies %d, want %d", len(setCookie), test.wantSetCookies)
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
			Probability:     1},
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
			Probability:     1},
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
			Probability:     1},
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
			Probability:     1},
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

func TestTelemetryMetrics(t *testing.T) {
	tele := Initialize(
		EnableMetrics(9101),
		WithViews(ocgrpc.DefaultServerViews),
	)

	retries := 0
	for {
		if _, err := http.Get("http://localhost:9101/metrics"); err == nil {
			// No error; break out of retry loop.
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := tele.Shutdown(ctx); err != nil {
		t.Fatalf("Telemetry shutdown error: %v", err)
	}
}
