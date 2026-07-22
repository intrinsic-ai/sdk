# Copyright 2023 Intrinsic Innovation LLC

# This package contains public protos that are exposed to clients of Intrinsic perception.

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

proto_library(
    name = "vector_proto",
    srcs = ["vector.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

py_proto_library(
    name = "vector_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":vector_proto"],
)

cc_proto_library(
    name = "vector_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":vector_proto"],
)

go_proto_library(
    name = "vector_go_proto",
    importpath = "intrinsic/perception/proto/vector_go_proto",
    protos = [":vector_proto"],
    visibility = ["//intrinsic:public_api_users"],
)
