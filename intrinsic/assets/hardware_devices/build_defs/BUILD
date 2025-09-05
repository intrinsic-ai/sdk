# Copyright 2023 Intrinsic Innovation LLC

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("//bazel:go_macros.bzl", "go_binary", "go_library")

bzl_library(
    name = "hardware_device_bzl",
    srcs = ["hardware_device.bzl"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/assets/build_defs:asset_bzl",
    ],
)

go_library(
    name = "hardwaredevicegen",
    srcs = ["hardwaredevicegen.go"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/assets:bundleio",
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/build_defs:asset_go_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/proto/v1:reference_go_proto",
        "//intrinsic/util/proto:protoio",
    ],
)

go_binary(
    name = "hardwaredevicegen_main",
    srcs = ["hardwaredevicegen_main.go"],
    visibility = ["//visibility:private"],
    deps = [
        ":hardwaredevicegen",
        "//intrinsic/production:intrinsic",
        "//intrinsic/util:flag",
        "@com_github_golang_glog//:go_default_library",
    ],
)
