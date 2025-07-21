# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "validate_link_parameters",
    srcs = ["validate_link_parameters.cc"],
    hdrs = ["validate_link_parameters.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/math:inertia_utils",
    ],
)
