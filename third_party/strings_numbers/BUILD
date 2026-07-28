# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("@rules_license//rules:license.bzl", "license")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

# Copied from https://github.com/google/sling/blob/master/LICENSE
exports_files(["LICENSE"])

license(
    name = "license",
    package_name = "sling",
    license_kinds = [
        "@rules_license//licenses/spdx:Apache-2.0",
    ],
    license_text = "LICENSE",
    package_url = "https://github.com/google/sling",
)

# Forked from https://github.com/google/sling/blob/master/sling/string/numbers.h
cc_library(
    name = "sort",
    srcs = ["sort.cc"],
    hdrs = ["sort.h"],
    applicable_licenses = [":license"],
    deps = ["@abseil-cpp//absl/strings"],
)

cc_test(
    name = "sort_test",
    srcs = ["sort_test.cc"],
    deps = [
        ":sort",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/container:btree",
    ],
)
