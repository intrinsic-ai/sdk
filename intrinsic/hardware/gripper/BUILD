# Copyright 2023 Intrinsic Innovation LLC

# Gripper hardware protos and headers.

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "gripper_proto",
    srcs = ["gripper.proto"],
    deps = ["//intrinsic/perception/proto/v1:settings_proto"],
)

py_proto_library(
    name = "gripper_py_pb2",
    deps = [":gripper_proto"],
)

cc_proto_library(
    name = "gripper_cc_proto",
    deps = [":gripper_proto"],
)
