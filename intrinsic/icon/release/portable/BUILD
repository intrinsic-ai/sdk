# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library", "cc_test")

package(default_visibility = [
    "//visibility:public",
])

cc_library(
    name = "init_xfa_absl",
    srcs = ["init_xfa_absl.cc"],
    hdrs = ["init_xfa.h"],
    deps = [
        "//intrinsic/icon/utils:log",
        "@com_google_absl//absl/flags:flag",
        "@com_google_absl//absl/flags:parse",
        "@com_google_absl//absl/flags:usage",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/log:globals",
        "@com_google_absl//absl/log:initialize",
        "@com_google_absl//absl/time",
    ],
)

cc_test(
    name = "init_xfa_absl_test",
    srcs = ["init_xfa_test.cc"],
    deps = [
        ":init_xfa_absl",
        "@com_google_absl//absl/base:log_severity",
        "@com_google_absl//absl/flags:flag",
        "@com_google_absl//absl/log:globals",
        "@com_google_absl//absl/strings",
    ],
)
