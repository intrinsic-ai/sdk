# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

go_library(
    name = "utils",
    testonly = True,
    srcs = ["utils.go"],
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/assets/processes/proto:process_asset_go_proto",
        "//intrinsic/assets/processes/proto:process_manifest_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/executive/proto:behavior_tree_go_proto",
        "//intrinsic/skills/proto:skills_go_proto",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_protobuf//proto",
    ],
)
