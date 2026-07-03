# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "clock",
    srcs = ["clock.cc"],
    hdrs = ["clock.h"],
    deps = ["@abseil-cpp//absl/time"],
)

cc_library(
    name = "deadline_timeout",
    srcs = ["deadline_timeout.cc"],
    hdrs = ["deadline_timeout.h"],
    deps = ["@abseil-cpp//absl/time"],
)
