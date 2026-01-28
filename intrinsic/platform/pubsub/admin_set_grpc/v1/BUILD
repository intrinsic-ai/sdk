# Copyright 2023 Intrinsic Innovation LLC

# GRPC service definition for the copying keys across namespaces.

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "admin_set_proto",
    srcs = ["admin_set.proto"],
    deps = ["@com_google_protobuf//:any_proto"],
)

go_proto_library(
    name = "admin_set_go_proto",
    protos = [":admin_set_proto"],
    visibility = ["//intrinsic:intrinsic_internal"],
)

go_grpc_library(
    name = "admin_set_go_grpc_proto",
    protos = [":admin_set_proto"],
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [
        ":admin_set_go_proto",
    ],
)

py_proto_library(
    name = "admin_set_py_pb2",
    deps = [":admin_set_proto"],
)

py_grpc_library(
    name = "admin_set_py_pb2_grpc",
    srcs = [":admin_set_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":admin_set_py_pb2"],
)

cc_proto_library(
    name = "admin_set_cc_proto",
    deps = [":admin_set_proto"],
)

cc_grpc_library(
    name = "admin_set_cc_grpc",
    srcs = [":admin_set_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":admin_set_cc_proto"],
)
