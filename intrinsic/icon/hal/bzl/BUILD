# Copyright 2023 Intrinsic Innovation LLC

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")

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
