# Copyright 2023 Intrinsic Innovation LLC

# This package is to provide a consistent way throughout Intrinsic components to
# store structured data about errors.

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "error_report_proto",
    srcs = ["error_report.proto"],
    deps = [
        "@com_google_googleapis//google/rpc:status_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "error_report_cc_proto",
    deps = [":error_report_proto"],
)

py_proto_library(
    name = "error_report_py_pb2",
    deps = [":error_report_proto"],
)

go_proto_library(
    name = "error_report_go_proto",
    go_deps = [
        "@org_golang_google_genproto_googleapis_rpc//status",
        "//intrinsic/executive/proto:clips_snapshot_go_proto",
    ],
    deps = [":error_report_proto"],
)
