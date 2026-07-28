# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_license//rules:license.bzl", "license")

package(default_visibility = ["//visibility:public"])

exports_files(["LICENSE"])

license(
    name = "license",
    package_name = "or-tools",
    license_kinds = [
        "@rules_license//licenses/spdx:Apache-2.0",
    ],
    license_text = "LICENSE",
    package_url = "https://github.com/google/or-tools/blob/v9.4/ortools/base/thorough_hash.h",
)

cc_library(
    name = "thorough_hash",
    hdrs = [":thorough_hash.h"],
)
