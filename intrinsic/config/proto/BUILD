# Copyright 2023 Intrinsic Innovation LLC

# These protos describe all user-accessible configuration of an Intrinsic application.
#
# See go/intrinsic-workcell-ppr for more details.

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "operation_mode_proto",
    srcs = ["operation_mode.proto"],
)

cc_proto_library(
    name = "operation_mode_cc_proto",
    deps = [":operation_mode_proto"],
)

go_proto_library(
    name = "operation_mode_go_proto",
    deps = [":operation_mode_proto"],
)

py_proto_library(
    name = "operation_mode_py_pb2",
    deps = [":operation_mode_proto"],
)
