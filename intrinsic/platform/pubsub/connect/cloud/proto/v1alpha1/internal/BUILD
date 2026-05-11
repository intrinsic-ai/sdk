# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "line_orchestration_config_proto",
    srcs = ["line_orchestration_config.proto"],
)

go_proto_library(
    name = "line_orchestration_config_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/platform/pubsub/connect/cloud/proto/v1alpha1/internal/line_orchestration_config_go_proto",
    protos = [":line_orchestration_config_proto"],
)

cc_proto_library(
    name = "line_orchestration_config_cc_proto",
    deps = [":line_orchestration_config_proto"],
)

py_proto_library(
    name = "line_orchestration_config_py_pb2",
    deps = [":line_orchestration_config_proto"],
)
