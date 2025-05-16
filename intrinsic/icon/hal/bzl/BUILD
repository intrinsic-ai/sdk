# Copyright 2023 Intrinsic Innovation LLC

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("//bazel:go_macros.bzl", "go_binary")

package(
    default_visibility = ["//visibility:private"],
)

bzl_library(
    name = "hardware_module_binary_bzl",
    srcs = ["hardware_module_binary.bzl"],
)

bzl_library(
    name = "hardware_module_image_bzl",
    srcs = ["hardware_module_image.bzl"],
    deps = [
        ":hardware_module_binary_bzl",
        "//bazel:container_bzl",
    ],
)

bzl_library(
    name = "resources_bzl",
    srcs = ["resources.bzl"],
)

go_binary(
    name = "hal_manifest",
    srcs = ["hal_manifest.go"],
    embedsrcs = [
        "hal_service_manifest.textproto.tmpl",
    ],
    visibility = ["//intrinsic:__subpackages__"],
    deps = [
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/production:intrinsic",
        "//intrinsic/util:flag",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_protobuf//encoding/prototext:go_default_library",
    ],
)
