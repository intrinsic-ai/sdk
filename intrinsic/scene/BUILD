# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "constants",
    hdrs = ["constants.h"],
    deps = [
        ":user_data_keys",
        "@abseil-cpp//absl/strings:string_view",
    ],
)

cc_library(
    name = "user_data_keys",
    hdrs = ["user_data_keys.h"],
    visibility = [
        "//intrinsic/scene:__subpackages__",
        "//intrinsic/simulation/world:__subpackages__",
        "//intrinsic/world/conversion/sdf:__subpackages__",
    ],
)
