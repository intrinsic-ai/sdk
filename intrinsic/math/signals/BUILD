# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "butter_filter2",
    srcs = [
        "butter_filter2.cc",
    ],
    hdrs = [
        "butter_filter2.h",
    ],
    deps = [
        "//intrinsic/icon/utils:log",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
    ],
)
