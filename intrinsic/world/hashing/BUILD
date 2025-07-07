# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library", "cc_test")

cc_library(
    name = "hashing",
    hdrs = ["hashing.h"],
    visibility = [
        "//intrinsic:__subpackages__",
        "//intrinsic/motion_planner/viewer:__pkg__",
    ],
    deps = [
        "//intrinsic/util:string_type",
        "//intrinsic/world:entity_id",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/container:flat_hash_set",
    ],
)

cc_test(
    name = "hashing_test",
    srcs = ["hashing_test.cc"],
    deps = [
        ":hashing",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/world:entity_id",
        "//intrinsic/world:labels",
    ],
)
