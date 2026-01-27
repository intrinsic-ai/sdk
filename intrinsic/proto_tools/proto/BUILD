# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = [
    "//intrinsic:public_api_users",
    "//intrinsic/proto_tools/registry:__subpackages__",
    "//intrinsic/proto_tools/tools:__subpackages__",
])

proto_library(
    name = "proto_registry_proto",
    srcs = ["proto_registry.proto"],
    deps = ["@com_google_protobuf//:descriptor_proto"],
)

cc_proto_library(
    name = "proto_registry_cc_proto",
    deps = [":proto_registry_proto"],
)

cc_grpc_library(
    name = "proto_registry_cc_grpc_proto",
    srcs = [":proto_registry_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":proto_registry_cc_proto"],
)

py_proto_library(
    name = "proto_registry_py_pb2",
    deps = [":proto_registry_proto"],
)

py_grpc_library(
    name = "proto_registry_py_pb2_grpc",
    srcs = [":proto_registry_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":proto_registry_py_pb2"],
)

go_proto_library(
    name = "proto_registry_go_proto",
    deps = [":proto_registry_proto"],
)

go_grpc_library(
    name = "proto_registry_go_grpc_proto",
    protos = [":proto_registry_proto"],
    deps = ["@io_bazel_rules_go//proto/wkt:descriptor_go_proto"],
)
