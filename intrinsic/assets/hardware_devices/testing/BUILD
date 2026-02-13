# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "utils",
    testonly = True,
    srcs = ["utils.go"],
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/catalog/proto/v1:asset_catalog_go_proto",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/testing:utils",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/assets/proto/v1:asset_graph_go_proto",
        "//intrinsic/assets/proto/v1:reference_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/assets/scene_objects/testing:utils",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/assets/services/testing:utils",
        "//intrinsic/util/proto:names",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
