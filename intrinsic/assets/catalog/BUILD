# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "assetutils",
    srcs = ["assetutils.go"],
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/kubernetes/workcell_spec/proto:image_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
