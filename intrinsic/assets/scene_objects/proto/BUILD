# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "scene_object_manifest_proto",
    srcs = ["scene_object_manifest.proto"],
    deps = [
        "//intrinsic/assets/proto:asset_tag_proto",
        "//intrinsic/assets/proto:documentation_proto",
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/assets/proto:vendor_proto",
        "//intrinsic/scene/proto/v1:scene_object_config_proto",
        "//intrinsic/scene/proto/v1:scene_object_proto",
        "@com_google_protobuf//:descriptor_proto",
    ],
)

cc_proto_library(
    name = "scene_object_manifest_cc_proto",
    deps = [":scene_object_manifest_proto"],
)

go_proto_library(
    name = "scene_object_manifest_go_proto",
    go_deps = [
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_config_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
    ],
    deps = [":scene_object_manifest_proto"],
)

py_proto_library(
    name = "scene_object_manifest_py_proto",
    deps = [":scene_object_manifest_proto"],
)
