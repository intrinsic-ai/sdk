# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "event_proto",
    srcs = ["event.proto"],
    deps = [
        "@com_google_protobuf//:struct_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

py_proto_library(
    name = "event_py_pb2",
    deps = ["//intrinsic/performance/analysis/proto:event_proto"],
)

cc_proto_library(
    name = "event_cc_proto",
    deps = ["//intrinsic/performance/analysis/proto:event_proto"],
)

proto_library(
    name = "performance_metrics_proto",
    srcs = ["performance_metrics.proto"],
    deps = ["@com_google_protobuf//:struct_proto"],
)

py_proto_library(
    name = "performance_metrics_py_pb2",
    deps = [":performance_metrics_proto"],
)

cc_proto_library(
    name = "performance_metrics_cc_proto",
    visibility = ["//intrinsic:__subpackages__"],
    deps = ["//intrinsic/performance/analysis/proto:performance_metrics_proto"],
)

go_proto_library(
    name = "performance_metrics_go_proto",
    deps = ["//intrinsic/performance/analysis/proto:performance_metrics_proto"],
)
