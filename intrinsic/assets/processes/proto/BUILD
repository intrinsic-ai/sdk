# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//intrinsic:public_api_users"])

proto_library(
    name = "process_asset_proto",
    srcs = ["process_asset.proto"],
    deps = [
        "//intrinsic/assets/proto:metadata_proto",
        "//intrinsic/executive/proto:behavior_tree_proto",
    ],
)

py_proto_library(
    name = "process_asset_py_pb2",
    deps = [":process_asset_proto"],
)

go_proto_library(
    name = "process_asset_go_proto",
    protos = [":process_asset_proto"],
    deps = [
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/executive/proto:behavior_tree_go_proto",
    ],
)

proto_library(
    name = "process_manifest_proto",
    srcs = ["process_manifest.proto"],
    deps = [
        "//intrinsic/assets/proto:asset_tag_proto",
        "//intrinsic/assets/proto:documentation_proto",
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/assets/proto:vendor_proto",
        "//intrinsic/executive/proto:behavior_tree_proto",
    ],
)

py_proto_library(
    name = "process_manifest_py_pb2",
    deps = [":process_manifest_proto"],
)

go_proto_library(
    name = "process_manifest_go_proto",
    protos = [":process_manifest_proto"],
    deps = [
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/executive/proto:behavior_tree_go_proto",
    ],
)
