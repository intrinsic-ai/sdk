# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "telemetry",
    srcs = ["telemetry.go"],
    deps = [
        "@com_github_golang_glog//:go_default_library",
        "@com_github_pborman_uuid//:go_default_library",
        "@io_opencensus_go//plugin/ocgrpc:go_default_library",
        "@io_opencensus_go//stats/view:go_default_library",
        "@io_opencensus_go//trace:go_default_library",
        "@io_opencensus_go_contrib_exporter_ocagent//:go_default_library",
        "@io_opencensus_go_contrib_exporter_prometheus//:go_default_library",
        "@io_opencensus_go_contrib_exporter_stackdriver//:go_default_library",
    ],
)

go_test(
    name = "telemetry_test",
    srcs = ["telemetry_test.go"],
    library = ":telemetry",
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
    srcs = ["slogattrs.go"],
)

go_test(
    name = "slogattrs_test",
    srcs = ["slogattrs_test.go"],
    library = ":slogattrs",
    deps = ["@com_github_google_go_cmp//cmp:go_default_library"],
)
