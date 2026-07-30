# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

# Simulation tools built on top of Gazebo simulation engine.
package(default_visibility = ["//visibility:public"])

cc_library(
    name = "type_conversion",
    srcs = ["type_conversion.cc"],
    hdrs = ["type_conversion.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/math:pose3",
        "//intrinsic/util/status:status_builder",
        "@abseil-cpp//absl/status:statusor",
        "@gz-math",
    ],
)
