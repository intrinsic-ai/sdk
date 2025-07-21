# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

package(default_visibility = [
    "//visibility:public",
])

cc_library(
    name = "entity_id",
    srcs = ["entity_id.cc"],
    hdrs = ["entity_id.h"],
    deps = [
        "//intrinsic/production/external/intops:strong_int",
        "//intrinsic/util:aggregate_type",
    ],
)

cc_test(
    name = "entity_id_test",
    srcs = ["entity_id_test.cc"],
    deps = [
        ":entity_id",
        "//intrinsic/util:aggregate_type",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/world/hashing",
        "@abseil-cpp//absl/hash:hash_testing",
    ],
)

cc_library(
    name = "labels",
    hdrs = ["labels.h"],
    deps = [
        "//intrinsic/util:string_type",
        "@abseil-cpp//absl/strings",
    ],
)
