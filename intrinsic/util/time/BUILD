# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "clock",
    srcs = ["clock.cc"],
    hdrs = ["clock.h"],
    deps = ["@abseil-cpp//absl/time"],
)

cc_library(
    name = "time",
    srcs = [
        "clock_steady.cc",
        "time.cc",
    ],
    hdrs = [
        "clock_steady.h",
        "clock_steady_interface.h",
        "time.h",
    ],
    deps = [
        "@abseil-cpp//absl/time",
        "@com_google_googletest//:gtest_prod",
    ],
)

cc_test(
    name = "time_test",
    size = "small",
    srcs = ["time_test.cc"],
    deps = [
        ":time",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/time",
    ],
)

cc_test(
    name = "clock_steady_test",
    size = "small",
    srcs = ["clock_steady_test.cc"],
    deps = [
        ":time",
        "//intrinsic/util/testing:gtest_wrapper_main",
    ],
)

cc_library(
    name = "timer",
    srcs = [
        "elapsed_timer.cc",
    ],
    hdrs = [
        "elapsed_timer.h",
    ],
    deps = [
        ":time",
        "@abseil-cpp//absl/time",
    ],
)

cc_test(
    name = "elapsed_timer_test",
    size = "small",
    srcs = ["elapsed_timer_test.cc"],
    deps = [
        ":time",
        ":timer",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/util/time/testing:clock_steady_fake",
        "@abseil-cpp//absl/time",
    ],
)
