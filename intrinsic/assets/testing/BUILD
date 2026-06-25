# Copyright 2023 Intrinsic Innovation LLC

# Testing utilities for assets.

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("//bazel:go_macros.bzl", "go_library")

package(
    default_testonly = True,
    default_visibility = ["//visibility:public"],
)

cc_library(
    name = "id_test_utils",
    srcs = ["id_test_utils.cc"],
    hdrs = ["id_test_utils.h"],
    deps = ["@abseil-cpp//absl/strings"],
)

go_library(
    name = "assetbuilder",
    testonly = True,
    srcs = ["assetbuilder.go"],
    importpath = "intrinsic/assets/testing/assetbuilder",
    visibility = [
        "//intrinsic:intrinsic_internal",
    ],
    deps = [
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/processes/proto:process_asset_go_proto",
        "//intrinsic/assets/proto/v1:asset_go_proto",
        "//intrinsic/assets/proto/v1:processed_asset_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/skills/proto:processed_skill_manifest_go_proto",
    ],
)
