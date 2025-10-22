# Copyright 2023 Intrinsic Innovation LLC

# GRPC service definition for the quota management service

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "provisioner_proto",
    srcs = ["provisioner.proto"],
    deps = ["//intrinsic/platform/pubsub/admin_set_grpc/v1:admin_set_proto"],
)

go_proto_library(
    name = "provisioner_go_proto",
    go_deps = [
        "//intrinsic/platform/pubsub/admin_set_grpc/v1:admin_set_go_grpc_proto",
    ],
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":provisioner_proto"],
)

go_grpc_library(
    name = "provisioner_go_grpc_proto",
    srcs = [":provisioner_proto"],
    visibility = ["//intrinsic:__subpackages__"],
    deps = [
        ":provisioner_go_proto",
        "//intrinsic/platform/pubsub/admin_set_grpc/v1:admin_set_go_grpc_proto",
    ],
)

py_proto_library(
    name = "provisioner_py_pb2",
    deps = [":provisioner_proto"],
)

py_grpc_library(
    name = "provisioner_py_pb2_grpc",
    srcs = [":provisioner_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":provisioner_py_pb2"],
)

cc_proto_library(
    name = "provisioner_cc_proto",
    deps = [":provisioner_proto"],
)

cc_grpc_library(
    name = "provisioner_cc_grpc",
    srcs = [":provisioner_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":provisioner_cc_proto"],
)
