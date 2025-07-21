# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("//intrinsic/platform:flatbuffers.bzl", "cc_flatbuffers_library", "flatbuffers_library")

package(default_visibility = ["//visibility:public"])

package_group(
    name = "icon",
    packages = [
        "//intrinsic/icon/...",
    ],
)

package_group(
    name = "safety",
    packages = [
        "//intrinsic/icon/control/safety/...",
    ],
)

flatbuffers_library(
    name = "safety_messages_fbs",
    srcs = [
        "safety_messages.fbs",
    ],
    deps = [
        "//intrinsic/icon/control/safety/extern:safety_status_fbs",
    ],
)

cc_flatbuffers_library(
    name = "safety_messages_fbs_cc",
    deps = [":safety_messages_fbs"],
)

cc_library(
    name = "safety_messages_fbs_utils",
    srcs = ["safety_messages_utils.cc"],
    hdrs = ["safety_messages_utils.h"],
    deps = [
        ":safety_messages_fbs_cc",
        "//intrinsic/icon/control/safety/extern:safety_status_fbs_cc",
        "@abseil-cpp//absl/strings",
        "@com_github_google_flatbuffers//:runtime_cc",
    ],
)

cc_test(
    name = "safety_messages_fbs_utils_test",
    srcs = ["safety_messages_utils_test.cc"],
    deps = [
        ":safety_messages_fbs_cc",
        ":safety_messages_fbs_utils",
        "//intrinsic/icon/control/safety/extern:safety_status_fbs_cc",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@com_github_google_flatbuffers//:flatbuffers",
    ],
)
