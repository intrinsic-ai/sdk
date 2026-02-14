# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "hardwaredevicebundle",
    srcs = ["hardwaredevicebundle.go"],
    importpath = "intrinsic/assets/hardware_devices/hardwaredevicebundle",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        ":hardwaredevicevalidate",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:imageutils",
        "//intrinsic/assets:ioutils",
        "//intrinsic/assets/data:databundle",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/services:servicebundle",
        "//intrinsic/util/archive:tartooling",
        "@com_github_google_safearchive//tar",
    ],
)

go_library(
    name = "hardwaredevicevalidate",
    srcs = ["hardwaredevicevalidate.go"],
    importpath = "intrinsic/assets/hardware_devices/hardwaredevicevalidate",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets:metadatautils",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto/v1:asset_graph_go_proto",
        "//intrinsic/assets/proto/v1:reference_go_proto",
    ],
)
