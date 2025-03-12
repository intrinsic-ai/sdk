# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library", "cc_test")
load("//bazel:go_macros.bzl", "go_library", "go_test")

# Application, process, product, and resource set describe all user-accessible configuration of a
# solution.

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "environments_cc",
    srcs = ["environments.cc"],
    hdrs = ["environments.h"],
    visibility = [
        "//visibility:public",
    ],
    deps = [
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
    ],
)

cc_test(
    name = "environments_test_cc",
    srcs = ["environments_test.cc"],
    deps = [
        ":environments_cc",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
    ],
)

go_library(
    name = "environments",
    srcs = ["environments.go"],
    visibility = [
        "//visibility:public",
    ],
)

go_test(
    name = "environments_test",
    srcs = ["environments_test.go"],
    library = ":environments",
)
