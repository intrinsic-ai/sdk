# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "telemetry",
    srcs = [
        "otel_telemetry.go",
        "telemetry.go",
    ],
    importpath = "intrinsic/stats/go/telemetry",
    deps = [
        "@com_github_golang_glog//:go_default_library",
        "@com_github_googlecloudplatform_opentelemetry_operations_go_exporter_trace//:go_default_library",
        "@com_github_pborman_uuid//:go_default_library",
        "@io_opencensus_go//plugin/ocgrpc:go_default_library",
        "@io_opencensus_go//stats:go_default_library",
        "@io_opencensus_go//stats/view:go_default_library",
        "@io_opencensus_go//tag:go_default_library",
        "@io_opencensus_go//trace:go_default_library",
        "@io_opencensus_go_contrib_exporter_ocagent//:go_default_library",
        "@io_opencensus_go_contrib_exporter_prometheus//:go_default_library",
        "@io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc//:otelgrpc",
        "@io_opentelemetry_go_contrib_propagators_b3//:b3",
        "@io_opentelemetry_go_contrib_propagators_opencensus//:opencensus",
        "@io_opentelemetry_go_otel//:otel",
        "@io_opentelemetry_go_otel//attribute",
        "@io_opentelemetry_go_otel//propagation",
        "@io_opentelemetry_go_otel//semconv/v1.34.0:v1_34_0",
        "@io_opentelemetry_go_otel_bridge_opencensus//:opencensus",
        "@io_opentelemetry_go_otel_sdk//resource",
        "@io_opentelemetry_go_otel_sdk//trace",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//stats:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_test(
    name = "telemetry_test",
    srcs = ["telemetry_test.go"],
    embed = [":telemetry"],
    importpath = "intrinsic/stats/go/telemetry_test",
    tags = [
        "requires-net:loopback",
    ],
    deps = [
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@io_opencensus_go//plugin/ocgrpc:go_default_library",
        "@io_opencensus_go//trace:go_default_library",
    ],
)

go_library(
    name = "slogattrs",
    srcs = [
        "metadata.go",
        "slogattrs.go",
        "slogattrs_metrics.go",
    ],
    importpath = "intrinsic/stats/go/slogattrs",
    deps = [
        "@io_opencensus_go//stats:go_default_library",
        "@io_opencensus_go//stats/view:go_default_library",
        "@io_opencensus_go//tag:go_default_library",
        "@io_opencensus_go//trace:go_default_library",
    ],
)

go_test(
    name = "slogattrs_test",
    srcs = ["slogattrs_test.go"],
    embed = [":slogattrs"],
    importpath = "intrinsic/stats/go/slogattrs_test",
    deps = ["@com_github_google_go_cmp//cmp:go_default_library"],
)
