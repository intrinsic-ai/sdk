# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_license//rules:license.bzl", "license")

package(default_visibility = ["//visibility:public"])

exports_files(["LICENSE"])

license(
    name = "license",
    package_name = "jwt_verify_lib",
    license_kinds = [
        "@rules_license//licenses/spdx:Apache-2.0",
    ],
    license_text = "LICENSE",
    package_url = "https://github.com/google/jwt_verify_lib/tree/master/simple_lru_cache",
)

cc_library(
    name = "simple_lru_cache",
    hdrs = [
        "simple_lru_cache.h",
        "simple_lru_cache_inl.h",
    ],
    deps = [
        "@abseil-cpp//absl/container:flat_hash_map",
    ],
)
