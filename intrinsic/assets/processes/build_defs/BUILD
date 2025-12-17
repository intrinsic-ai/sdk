# Copyright 2023 Intrinsic Innovation LLC

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("//bazel:go_macros.bzl", "go_binary", "go_library")

package(default_visibility = ["//visibility:private"])

bzl_library(
    name = "processes_bzl",
    srcs = ["processes.bzl"],
    deps = [
        "//intrinsic/assets/build_defs:asset_bzl",
        "//intrinsic/util/proto/build_defs:descriptor_set",
    ],
)

go_library(
    name = "processgen",
    srcs = ["processgen.go"],
    deps = [
        "//intrinsic/assets:bundleio",
        "//intrinsic/assets/processes/proto:process_manifest_go_proto",
        "//intrinsic/executive/proto:behavior_tree_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
    ],
)

go_binary(
    name = "processgen_main",
    srcs = ["processgen_main.go"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":processgen",
        "//intrinsic/production:intrinsic",
        "//intrinsic/util:flag",
        "@com_github_golang_glog//:go_default_library",
    ],
)
