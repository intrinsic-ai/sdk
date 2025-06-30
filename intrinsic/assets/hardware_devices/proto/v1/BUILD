# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

proto_library(
    name = "hardware_device_manifest_proto",
    srcs = ["hardware_device_manifest.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/data/proto/v1:data_asset_proto",
        "//intrinsic/assets/proto:asset_tag_proto",
        "//intrinsic/assets/proto:documentation_proto",
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/assets/proto:vendor_proto",
        "//intrinsic/assets/proto/v1:asset_graph_proto",
        "//intrinsic/assets/proto/v1:reference_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_proto",
        "//intrinsic/assets/services/proto:service_manifest_proto",
    ],
)

go_proto_library(
    name = "hardware_device_manifest_go_proto",
    go_deps = [
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/assets/proto/v1:asset_graph_go_proto",
        "//intrinsic/assets/proto/v1:reference_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":hardware_device_manifest_proto"],
)

py_proto_library(
    name = "hardware_device_manifest_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":hardware_device_manifest_proto"],
)

cc_proto_library(
    name = "hardware_device_manifest_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":hardware_device_manifest_proto"],
)
