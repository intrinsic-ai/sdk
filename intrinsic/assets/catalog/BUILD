# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "assetdescriptions",
    srcs = ["assetdescriptions.go"],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_grpc_proto",
    ],
)

go_library(
    name = "assetutils",
    srcs = ["assetutils.go"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_grpc_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/kubernetes/workcell_spec/proto:image_go_proto",
        "@com_github_golang_glog//:go_default_library",
    ],
)
