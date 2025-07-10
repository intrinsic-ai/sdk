// Copyright 2023 Intrinsic Innovation LLC

// Package telemetry sets up OpenCensus tracing and metrics.
package telemetry

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/prometheus"
	"contrib.go.opencensus.io/exporter/stackdriver"
	log "github.com/golang/glog"
	"github.com/pborman/uuid"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"
)

// tracingConfig contains the configuration parameters for tracing.
type tracingConfig struct {
	Enabled         bool
	CloudDeployment bool
	ProjectName     string
	ServiceName     string
	Probability     float64
}

// String returns the string representation of the tracing config.
// The Stringer interface improves debuggability in unit tests.
func (tc *tracingConfig) String() string {
	if tc == nil {
		return "<nil>"
	}
	str := `{
		Enabled: %t
		CloudDeployment: %t
		ProjectName: '%s'
		ServiceName: '%s'
		Probability: %v
	}`
	return fmt.Sprintf(str, tc.Enabled, tc.CloudDeployment, tc.ProjectName,
		tc.ServiceName, tc.Probability)
}

// metricsConfig contains the configuration parameters for metrics.
type metricsConfig struct {
	Enabled     bool
	MetricsPath string
	MetricsPort int64
	Views       []*view.View
}

// String returns the string representation of the tracing config.
// The Stringer interface improves debuggability in unit tests.
func (mc *metricsConfig) String() string {
	if mc == nil {
		return "<nil>"
	}
	str := `{
		Enabled: %t
		MetricsPath: %t
		MetricsPort: %d
	}`
	return fmt.Sprintf(str, mc.Enabled, mc.MetricsPath, mc.MetricsPort)
}

// config contains the configuration parameters for telemetry.
type config struct {
	TracingCfg *tracingConfig
	MetricsCfg *metricsConfig
}

// newConfig creates a new config.
func newConfig(opts ...ConfigOption) config {
	cfg := config{}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

// getTracingConfig returns the tracing config.
func (c *config) getTracingConfig() *tracingConfig {
	if c.TracingCfg == nil {
		c.TracingCfg = &tracingConfig{}
		c.TracingCfg.Enabled = false
		c.TracingCfg.CloudDeployment = false
		c.TracingCfg.Probability = 1.0
	}
	return c.TracingCfg
}

// getMetricsConfig returns the metrics config.
func (c *config) getMetricsConfig() *metricsConfig {
	if c.MetricsCfg == nil {
		c.MetricsCfg = &metricsConfig{}
		c.MetricsCfg.Enabled = false
		c.MetricsCfg.MetricsPath = "/metrics"
		c.MetricsCfg.Views = ocgrpc.DefaultServerViews
	}
	return c.MetricsCfg
}

// ConfigOption is a function that can be used to modify the config.
type ConfigOption func(c *config)

// WithCloudTracing enables cloud tracing. For on-prem workloads, please use
// EnableTracing() or WithTracing() instead. The intended use for this function
// are cases where tracing can be turned on and off based on a parameter, e.g.
// one specified as a command line parameter.
func WithCloudTracing(enabled bool) ConfigOption {
	return func(c *config) {
		cfg := c.getTracingConfig()
		cfg.Enabled = enabled
		cfg.CloudDeployment = true
	}
}

// EnableCloudTracing enables cloud tracing.
func EnableCloudTracing() ConfigOption {
	return WithCloudTracing(true)
}

// WithTracing enables tracing for on-prem (non-cloud) workloads. For cloud
// hosted workloads, please use EnableCloudTracing() or WithCloudTracing()
// instead. The intended use for this function are cases where tracing can be
// turned on and off based on a parameter, e.g. one specified as a command line
// parameter.
func WithTracing(enabled bool) ConfigOption {
	return func(c *config) {
		cfg := c.getTracingConfig()
		cfg.Enabled = enabled
		cfg.CloudDeployment = false
	}
}

// EnableTracing enables tracing.
func EnableTracing() ConfigOption {
	return WithTracing(true)
}

// WithProjectName sets the project name for cloud tracing.
func WithProjectName(pn string) ConfigOption {
	return func(c *config) {
		cfg := c.getTracingConfig()
		cfg.ProjectName = pn
	}
}

// WithServiceName sets the service name for cloud tracing.
func WithServiceName(sn string) ConfigOption {
	return func(c *config) {
		cfg := c.getTracingConfig()
		cfg.ServiceName = sn
	}
}

// WithProbability sets the probability for tracing.
func WithProbability(p float64) ConfigOption {
	return func(c *config) {
		cfg := c.getTracingConfig()
		cfg.Probability = p
	}
}

// WithMetrics enables metrics.
func WithMetrics(enabled bool, port int64) ConfigOption {
	return func(c *config) {
		cfg := c.getMetricsConfig()
		cfg.Enabled = enabled
		cfg.MetricsPort = port
	}
}

// EnableMetrics enables and serves metrics at "0.0.0.0:<port>/metrics".
func EnableMetrics(port int64) ConfigOption {
	return WithMetrics(true, port)
}

// WithViews sets the views for metrics.
func WithViews(Views []*view.View) ConfigOption {
	return func(c *config) {
		cfg := c.getMetricsConfig()
		cfg.Views = Views
	}
}

// Telemetry is an object to configure tracing and metrics export in services.
type Telemetry struct {
	exporter      *trace.Exporter
	metricsServer *http.Server
}

// Initialize creates a new telemetry instance.
func Initialize(opts ...ConfigOption) Telemetry {
	t := Telemetry{}
	t.init(newConfig(opts...))
	return t
}

func (t *Telemetry) createCloudTraceExporter(project string, serviceName string) (trace.Exporter, error) {
	if project == "" {
		return nil, fmt.Errorf("project must be specified")
	}
	defaultAttributes := map[string]any{
		"service.name": serviceName,
	}
	if h, err := os.Hostname(); err == nil { // no error
		defaultAttributes["hostname"] = h
	}
	log.Info("Creating Stackdriver exporter for tracing.")
	return stackdriver.NewExporter(stackdriver.Options{
		ProjectID:              project,
		DefaultTraceAttributes: defaultAttributes,
	})
}

func (t *Telemetry) createTraceExporter(project string, serviceName string, cloudTracing bool) (trace.Exporter, error) {
	if cloudTracing {
		return t.createCloudTraceExporter(project, serviceName)
	}

	log.Info("Creating OpenCensus Agent exporter for tracing.")
	return ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithAddress("oc-agent.app-intrinsic-base.svc.cluster.local:55678"),
		ocagent.WithServiceName(serviceName),
	)
}

// init initializes the telemetry instance.
func (t *Telemetry) init(cfg config) {
	if cfg.TracingCfg != nil {
		t.enableTracing(*cfg.TracingCfg)
	}
	if cfg.MetricsCfg != nil {
		t.enableMetrics(*cfg.MetricsCfg)
	}
}

// Shutdown gracefully shuts down the current telemetry instance without
// interrupting active metrics connections.
//
// Example usage:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	if err := tele.Shutdown(ctx); err != nil {
//		log.Fatalf("Telemetry shutdown error: %v", err)
//	}
//	log.Println("Graceful telemetry shutdown complete.")
func (t *Telemetry) Shutdown(ctx context.Context) error {
	if t.exporter != nil {
		trace.UnregisterExporter(*t.exporter)
	}
	if t.metricsServer != nil {
		return t.metricsServer.Shutdown(ctx)
	}
	return nil
}

// enableTracing enables tracing for the current telemetry instance.
func (t *Telemetry) enableTracing(c tracingConfig) {
	if !c.Enabled {
		return
	}
	exp, err := t.createTraceExporter(c.ProjectName, c.ServiceName, c.CloudDeployment)
	if err != nil {
		log.Warningf("Tracing is disabled! Tracing setup failed: %v", err)
	}
	t.exporter = &exp
	trace.RegisterExporter(exp)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(c.Probability)})
}

// enableMetrics enables metrics for the current telemetry instance.
func (t *Telemetry) enableMetrics(c metricsConfig) {
	if !c.Enabled {
		return
	}
	if err := view.Register(c.Views...); err != nil {
		log.Warningf("Failed to register views: %v", err)
	}
	pe, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		log.Errorf("Metrics is disabled! Metrics setup failed: %v", err)
	}
	view.RegisterExporter(pe)

	t.metricsServer = &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%v", c.MetricsPort),
	}
	http.Handle(c.MetricsPath, pe)

	go func() {
		if err := t.metricsServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("Metrics is disabled! Failed to start HTTP server: %v", err)
		}
		log.Info("Stopped serving new metrics connections.")
	}()
}

// Tracing session identifier
const trsid = "trsid"

// TraceOnCookie samples the request if a custom tracing session cookie is present.
func TraceOnCookie(req *http.Request) trace.StartOptions {
	hdr, err := req.Cookie(trsid)
	if err != nil || !validateTRSID.MatchString(hdr.Value) {
		return trace.StartOptions{}
	}
	// If we have a valid cookie set, we trace always, independently on the user defined
	// probability. I.e. even in a project where we might by default sample with a probability
	// of 20%, we can set to always sample using the cookie.
	return trace.StartOptions{Sampler: trace.AlwaysSample()}
}

// Validates the tracing session identifier
var validateTRSID = regexp.MustCompile(`^[A-Za-z0-9\-_]+$`)

// TraceIDWriter makes sure only one trace identifier header is written.
// Otherwise our proxy setup in the portal would also add the
// upstream header after we set the header in the first handler.
// This has to be done only once in WriteHeader or Write.
type TraceIDWriter struct {
	http.ResponseWriter
	traceID string
}

// IDHeader is the trace identifier header
const IDHeader = "X-Intrinsic-TraceID"

// WriteHeader writes the headers but makes sure the trace identifier is only present once.
func (w *TraceIDWriter) WriteHeader(status int) {
	if w.traceID != "" {
		w.ResponseWriter.Header().Set(IDHeader, w.traceID)
		w.traceID = ""
	}
	w.ResponseWriter.WriteHeader(status)
}

// Write writes the bytes to the transport and also makes sure only one trace identifier is present.
func (w *TraceIDWriter) Write(b []byte) (int, error) {
	if w.traceID != "" {
		w.ResponseWriter.Header().Set(IDHeader, w.traceID)
		w.traceID = ""
	}
	return w.ResponseWriter.Write(b)
}

// Flush triggers Flush in the underlying response writer.
func (w *TraceIDWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack allows hijacking of the underlying response writer.
func (w *TraceIDWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("http.Hijacker interface not supported by %T", w.ResponseWriter)
	}
	return h.Hijack()
}

// DefaultHandler returns a handler which adds a trace ID, a TRS ID (if there exists a
// cookie on the request) and the 'service.name' attribute.
func DefaultHandler(serviceName string) http.Handler {
	return ServiceNameHandler(serviceName, AddSpanTRSIDHandler(TraceIDHandler(http.DefaultServeMux)))
}

// ServiceNameHandler returns a handler which annotates the current span with the given service
// name.
func ServiceNameHandler(serviceName string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if span := trace.FromContext(r.Context()); span != nil {
			span.AddAttributes(trace.StringAttribute("service.name", serviceName))
		}
		h.ServeHTTP(w, r)
	})
}

// TraceID returns the current trace identifier if tracing is active.
func TraceID(ctx context.Context) string {
	if span := trace.FromContext(ctx); span != nil && span.SpanContext().IsSampled() {
		return span.SpanContext().TraceID.String()
	}
	return ""
}

// TraceIDHandler adds the current trace identifier to the response.
// No header is added if no tracing context exists.
func TraceIDHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tid := ""
		if tid := TraceID(r.Context()); tid != "" {
			// Set tracing header in addition to using TraceIDWriter.
			// TraceIDWriter does not cover cases where implicit writes are performed.
			w.Header().Set(IDHeader, tid)
		}
		h.ServeHTTP(&TraceIDWriter{w, tid}, r)
	})
}

// AddSpanTRSIDHandler adds the current tracing session identifier as span attribute.
func AddSpanTRSIDHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hdr, err := r.Cookie(trsid); err == nil { // no error
			if !validateTRSID.MatchString(hdr.Value) {
				log.Warningf("tracing session identifier invalid, ignoring")
			} else {
				if span := trace.FromContext(r.Context()); span != nil {
					span.AddAttributes(trace.StringAttribute("trsid", hdr.Value))
				}
			}
		}
		h.ServeHTTP(w, r)
	})
}

// TracingEndpointsHandler provides endpoints for turning tracing for specific requests on / off.
type TracingEndpointsHandler struct {
	base string
}

// RegisterTracingEndpoints registers the endpoints under a specific base path.
func RegisterTracingEndpoints(mux *http.ServeMux, base string) (*TracingEndpointsHandler, error) {
	h := &TracingEndpointsHandler{base: strings.TrimSuffix(base, "/")}
	mux.HandleFunc(base+"/", h.tracingHandler)
	return h, nil
}

func genCookie(value string, delete bool) *http.Cookie {
	c := http.Cookie{
		Name:     trsid,
		Value:    value,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
	if delete {
		c.MaxAge = -1
		c.Expires = time.Unix(1, 0)
	} else {
		c.Expires = time.Now().Add(30 * time.Minute)
	}
	return &c
}

var uuidNew = uuid.New // for unit testing

// Handles all available tracing operations
// The endpoints accept GET requests so that a new tracing session can be created via golinks.
// Expecting these paths:
// - [GET] [base]/enable: Create a new tracing session, even if one exists
// - [GET] [base]/trsid: Return the current tracing session identifier or NotFound
// - [GET] [base]/disable: Delete the current tracing session
func (h *TracingEndpointsHandler) tracingHandler(w http.ResponseWriter, r *http.Request) {
	_, span := trace.StartSpan(r.Context(), "telemetry.tracingHandler")
	defer span.End()

	// format: [base]/[op]
	op := strings.TrimPrefix(r.URL.Path, h.base+"/")
	switch op {
	case "enable":
		t := uuidNew()[0:8]
		c := genCookie(t, false)
		http.SetCookie(w, c)
		w.Write([]byte(fmt.Sprintf("created tracing session %s for you", t)))
	case "trsid":
		if t, err := r.Cookie(trsid); err == nil { // no error, cookie found
			w.Write([]byte(t.Value))
		} else { // error or cookie not found
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("no tracing session cookie found"))
		}
	case "disable":
		c := genCookie("", true)
		http.SetCookie(w, c)
		w.Write([]byte("disabled your tracing session"))
	default:
		http.Error(w, "invalid tracing operation", http.StatusBadRequest)
	}
}

// SetError sets the error status code and message for the given span.
// The helper avoids the string-formatting operation for the error if the span is not recorded.
func SetError(span *trace.Span, statusCode int, message string, err error) {
	if !span.IsRecordingEvents() {
		return
	}
	span.SetStatus(trace.Status{Code: int32(statusCode), Message: fmt.Sprintf("%s: %v", message, err)})
}

// SetErrorf sets the error status code and message for the given span.
// The helper avoids the string-formatting operation for the error if the span is not recorded.
func SetErrorf(span *trace.Span, statusCode int, format string, a ...any) {
	if !span.IsRecordingEvents() {
		return
	}
	span.SetStatus(trace.Status{Code: int32(statusCode), Message: fmt.Sprintf(format, a...)})
}

// StatusWithError takes span and error and treats error as grpc status
// to set status on the span. Returns error for easy daisy-chaining.
func StatusWithError(span *trace.Span, err error) error {
	if err != nil {
		errStat, _ := status.FromError(err)
		span.SetStatus(trace.Status{
			Code:    int32(errStat.Code()),
			Message: errStat.Message(),
		})
	}
	return err
}
