# Copyright 2023 Intrinsic Innovation LLC

# Package containing test utilities for gRPC.

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:private"])

proto_library(
    name = "ping_proto",
    srcs = ["ping.proto"],
    deps = [
        "@com_google_protobuf//:empty_proto",
    ],
)

cc_proto_library(
    name = "ping_cc_proto",
    deps = [":ping_proto"],
)

cc_grpc_library(
    name = "ping_cc_grpc_proto",
    srcs = [":ping_proto"],
    generate_mocks = True,
    grpc_only = True,
    visibility = ["//intrinsic/util/grpc:__subpackages__"],
    deps = [":ping_cc_proto"],
)

go_proto_library(
    name = "ping_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/util/grpc/testing/ping_go_proto",
    protos = [":ping_proto"],
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = ["@org_golang_google_protobuf//types/known/emptypb"],
)

py_proto_library(
    name = "ping_proto_py_pb2",
    deps = [":ping_proto"],
)

py_grpc_library(
    name = "ping_proto_py_pb2_grpc",
    srcs = [":ping_proto"],
    grpc_library = requirement("grpcio"),
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [":ping_proto_py_pb2"],
)
