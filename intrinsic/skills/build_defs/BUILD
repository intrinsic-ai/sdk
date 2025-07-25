# Copyright 2023 Intrinsic Innovation LLC

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("//bazel:go_macros.bzl", "go_binary", "go_library")

package(default_visibility = ["//visibility:public"])

bzl_library(
    name = "manifest_bzl",
    srcs = ["manifest.bzl"],
    deps = [
        "//intrinsic/util/proto/build_defs:descriptor_set",
    ],
)

bzl_library(
    name = "skill_bzl",
    srcs = ["skill.bzl"],
    deps = [
        ":manifest_bzl",
        "//bazel:cc_oci_image_bzl",
        "//bazel:container_bzl",
        "//bazel:python_oci_image_bzl",
        "@bazel_skylib//lib:paths",
        "@rules_cc//cc:core_rules",
    ],
)

go_binary(
    name = "skillmanifestgen",
    srcs = ["skillmanifestgen.go"],
    visibility = ["//visibility:public"],
    deps = [
        "//intrinsic/production:intrinsic",
        "//intrinsic/skills/internal:skillmanifest",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/proto:protoio",
        "//intrinsic/util/proto:registryutil",
        "@com_github_golang_glog//:go_default_library",
    ],
)

go_binary(
    name = "skillbundlegen",
    srcs = ["skillbundlegen.go"],
    visibility = ["//visibility:public"],
    deps = [
        "//intrinsic/assets:bundleio",
        "//intrinsic/production:intrinsic",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/proto:protoio",
        "@com_github_golang_glog//:go_default_library",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
    ],
)

go_library(
    name = "skillserviceconfiggen",
    srcs = ["skillserviceconfiggen.go"],
    visibility = ["//intrinsic/tools/inbuild:__subpackages__"],
    deps = [
        "//intrinsic/assets:idutils",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/skills/proto:skill_service_config_go_proto",
        "//intrinsic/skills/proto:skills_go_proto",
        "//intrinsic/util/proto:sourcecodeinfoview",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_protobuf//proto",
    ],
)

go_binary(
    name = "skillserviceconfiggen_main",
    srcs = ["skillserviceconfiggen_main.go"],
    visibility = ["//visibility:public"],
    deps = [
        ":skillserviceconfiggen",
        "//intrinsic/production:intrinsic",
        "//intrinsic/skills/proto:skill_manifest_go_proto",
        "//intrinsic/util/proto:protoio",
        "@com_github_golang_glog//:go_default_library",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
    ],
)
