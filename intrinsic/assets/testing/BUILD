# Copyright 2023 Intrinsic Innovation LLC

# Testing utilities for assets.

load("@rules_cc//cc:cc_library.bzl", "cc_library")

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
