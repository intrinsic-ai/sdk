# Copyright 2023 Intrinsic Innovation LLC

load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")

proto_library(
    name = "motion_target_proto",
    srcs = ["motion_target.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/icon/proto:joint_space_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

cc_proto_library(
    name = "motion_target_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":motion_target_proto"],
)

py_proto_library(
    name = "motion_target_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":motion_target_proto"],
)
