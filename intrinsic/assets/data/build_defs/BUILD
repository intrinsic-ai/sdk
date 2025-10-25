# Copyright 2023 Intrinsic Innovation LLC

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("//bazel:go_macros.bzl", "go_binary", "go_library")

bzl_library(
    name = "data_bzl",
    srcs = ["data.bzl"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/assets/build_defs:asset_bzl",
        "//intrinsic/util/proto/build_defs:descriptor_set",
        "@bazel_skylib//rules:common_settings",
    ],
)

go_library(
    name = "datagen",
    srcs = ["datagen.go"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/assets:bundleio",
        "//intrinsic/assets/data:datamanifest",
        "//intrinsic/assets/data:utils",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/proto/v1:data_manifest_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)

go_binary(
    name = "datagen_main",
    srcs = ["datagen_main.go"],
    visibility = ["//visibility:private"],
    deps = [
        ":datagen",
        "//intrinsic/production:intrinsic",
        "//intrinsic/util:flag",
        "@com_github_golang_glog//:go_default_library",
    ],
)
