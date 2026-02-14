# Copyright 2023 Intrinsic Innovation LLC

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_binary", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

bzl_library(
    name = "asset_bzl",
    srcs = ["asset.bzl"],
)

proto_library(
    name = "asset_proto",
    srcs = ["asset.proto"],
    deps = [
        "//intrinsic/assets/proto:asset_type_proto",
        "//intrinsic/assets/proto:id_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:descriptor_proto",
    ],
)

py_proto_library(
    name = "asset_py_pb2",
    deps = [":asset_proto"],
)

go_proto_library(
    name = "asset_go_proto",
    importpath = "intrinsic/assets/build_defs/asset_go_proto",
    protos = [":asset_proto"],
    deps = [
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
    ],
)

go_binary(
    name = "assetlocalinfogen",
    srcs = ["assetlocalinfogen.go"],
    importpath = "intrinsic/assets/build_defs/assetlocalinfogen",
    deps = [
        ":asset_go_proto",
        "//intrinsic/assets:typeutils",
        "//intrinsic/assets/data/proto/v1:data_manifest_go_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/processes/proto:process_manifest_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/production:intrinsic",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util:flag",
        "//intrinsic/util/proto:descriptor",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "@com_github_golang_glog//:go_default_library",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
    ],
)

go_binary(
    name = "assetcatalogrefinfogen",
    srcs = ["assetcatalogrefinfogen.go"],
    importpath = "intrinsic/assets/build_defs/assetcatalogrefinfogen",
    deps = [
        ":asset_go_proto",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:typeutils",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/production:intrinsic",
        "//intrinsic/util:flag",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "@com_github_golang_glog//:go_default_library",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
    ],
)

go_binary(
    name = "assetinstancegen",
    srcs = ["assetinstancegen.go"],
    importpath = "intrinsic/assets/build_defs/assetinstancegen",
    deps = [
        ":asset_go_proto",
        "//intrinsic/assets:idutils",
        "//intrinsic/production:intrinsic",
        "//intrinsic/util/proto:protoio",
        "@com_github_golang_glog//:go_default_library",
    ],
)
