# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

package(default_visibility = [
    "//visibility:public",
])

cc_library(
    name = "init_xfa_absl",
    srcs = ["init_xfa_absl.cc"],
    hdrs = [
        "init_intrinsic.h",
        "init_xfa.h",
    ],
    deps = [
        "//intrinsic/icon/utils:log",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/flags:parse",
        "@abseil-cpp//absl/flags:usage",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:globals",
        "@abseil-cpp//absl/log:initialize",
        "@abseil-cpp//absl/time",
    ],
)

# Migration aliases for insrc SoT files.
genrule(
    name = "init_intrinsic_h",
    srcs = [":init_xfa.h"],
    outs = ["init_intrinsic.h"],
    cmd = "cp $< $@",
)

alias(
    name = "init_intrinsic",
    actual = ":init_xfa_absl",
)

cc_test(
    name = "init_xfa_absl_test",
    srcs = ["init_xfa_test.cc"],
    deps = [
        ":init_xfa_absl",
        "@abseil-cpp//absl/base:log_severity",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log:globals",
        "@abseil-cpp//absl/strings",
    ],
)
