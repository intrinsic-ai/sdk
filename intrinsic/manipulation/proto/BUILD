# Copyright 2023 Intrinsic Innovation LLC

# Commonly used protos for manipulation

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "workspace_proto",
    srcs = ["workspace.proto"],
    deps = [
        "//intrinsic/geometry/proto:oriented_bounding_box_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

py_proto_library(
    name = "workspace_py_pb2",
    deps = [":workspace_proto"],
)

cc_proto_library(
    name = "workspace_cc_proto",
    deps = [":workspace_proto"],
)

proto_library(
    name = "gripper_assembly_metadata_proto",
    srcs = ["gripper_assembly_metadata.proto"],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/math/proto:vector3_proto",
        "//intrinsic/scene/proto/v1:scene_object_proto",
    ],
)

cc_proto_library(
    name = "gripper_assembly_metadata_cc_proto",
    deps = [":gripper_assembly_metadata_proto"],
)

py_proto_library(
    name = "gripper_assembly_metadata_py_pb2",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [":gripper_assembly_metadata_proto"],
)

proto_library(
    name = "gripper_exchange_proto",
    srcs = ["gripper_exchange.proto"],
    deps = [
        ":gripper_assembly_metadata_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

py_proto_library(
    name = "gripper_exchange_py_pb2",
    deps = [":gripper_exchange_proto"],
)

cc_proto_library(
    name = "gripper_exchange_cc_proto",
    deps = [":gripper_exchange_proto"],
)
