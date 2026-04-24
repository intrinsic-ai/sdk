# Copyright 2023 Intrinsic Innovation LLC

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "realtime_annotations",
    hdrs = ["realtime_annotations.h"],
    visibility = [
        "//intrinsic:intrinsic_internal",
        "//sdk:__subpackages__",
    ],
)

bzl_library(
    name = "cc_test_and_malloc_test_bzl",
    srcs = ["cc_test_and_malloc_test.bzl"],
    visibility = [
        "//sdk:__subpackages__",
        "//visibility:public",
    ],
    deps = [
        "@bazel_skylib//lib:new_sets",
        "@rules_cc//cc:core_rules",
    ],
)
