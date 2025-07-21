# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

# Logging utilities.
package(default_visibility = ["//visibility:public"])

cc_library(
    name = "downsampler",
    srcs = ["downsampler.cc"],
    hdrs = ["downsampler.h"],
    deps = [
        "//intrinsic/logging/proto:downsampler_cc_proto",
        "//intrinsic/logging/proto:log_item_cc_proto",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/time",
    ],
)

cc_test(
    name = "downsampler_test",
    srcs = ["downsampler_test.cc"],
    deps = [
        ":downsampler",
        "//intrinsic/logging/proto:downsampler_cc_proto",
        "//intrinsic/logging/proto:log_item_cc_proto",
        "//intrinsic/logging/proto:logger_service_cc_proto",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "proto_conversion",
    srcs = ["proto_conversion.cc"],
    hdrs = ["proto_conversion.h"],
    deps = [
        ":downsampler",
        "//intrinsic/logging/proto:downsampler_cc_proto",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/time",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "proto_conversion_test",
    srcs = ["proto_conversion_test.cc"],
    deps = [
        ":downsampler",
        ":proto_conversion",
        "//intrinsic/logging/proto:downsampler_cc_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/time",
    ],
)
