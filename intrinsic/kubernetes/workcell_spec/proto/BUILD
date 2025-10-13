# Copyright 2023 Intrinsic Innovation LLC

# gRPC service definition for the workcell-spec-service.

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "image_proto",
    srcs = ["image.proto"],
)

go_proto_library(
    name = "image_go_proto",
    deps = [":image_proto"],
)

py_proto_library(
    name = "image_py_pb2",
    deps = [":image_proto"],
)

proto_library(
    name = "deploy_proto",
    srcs = ["deploy.proto"],
    deps = ["//intrinsic/config/proto:operation_mode_proto"],
)

go_proto_library(
    name = "deploy_go_proto",
    go_deps = ["//intrinsic/config/proto:operation_mode_go_proto"],
    deps = [":deploy_proto"],
)

go_grpc_library(
    name = "deploy_go_grpc_proto",
    srcs = [":deploy_proto"],
    deps = ["//intrinsic/config/proto:operation_mode_go_proto"],
)
