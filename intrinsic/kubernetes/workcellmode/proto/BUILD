# Copyright 2023 Intrinsic Innovation LLC

# GRPC service definition for the workcell mode service.

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "workcellmode_service_proto",
    srcs = ["workcellmode.proto"],
)

go_proto_library(
    name = "workcellmode_service_go_proto",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":workcellmode_service_proto"],
)

go_grpc_library(
    name = "workcellmode_service_go_grpc_proto",
    srcs = [":workcellmode_service_proto"],
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":workcellmode_service_go_proto"],
)

py_proto_library(
    name = "workcellmode_service_py_pb2",
    deps = [":workcellmode_service_proto"],
)

py_grpc_library(
    name = "workcellmodea_service_py_pb2_grpc",
    srcs = [":workcellmode_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":workcellmode_service_py_pb2"],
)
