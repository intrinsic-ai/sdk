# Copyright 2023 Intrinsic Innovation LLC

load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "constrained_ik_proto",
    srcs = ["constrained_ik.proto"],
    deps = [
        "//intrinsic/icon/proto:joint_space_proto",
    ],
)

cc_proto_library(
    name = "constrained_ik_cc_proto",
    deps = [":constrained_ik_proto"],
)

py_proto_library(
    name = "constrained_ik_py_pb2",
    deps = [":constrained_ik_proto"],
)
