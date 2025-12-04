# Copyright 2023 Intrinsic Innovation LLC

# Provides the bindings and other utility with externally built zenoh pubsub

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

exports_files(
    [
        "peer_config.json",
        "zenoh_handle.h",
    ],
    visibility = [
        "//intrinsic:intrinsic_internal",
    ],
)

cc_library(
    name = "zenoh_handle",
    srcs = ["zenoh_handle.cc"],
    hdrs = ["zenoh_handle.h"],
    data = [
        "//intrinsic/insrc/middleware:libimw_zenoh.so.1",
    ],
    linkopts = ["-ldl"],  # for dlfcn.h, dlsym in zenoh_handle.cc
    visibility = ["//intrinsic/platform/pubsub:__subpackages__"],
    deps = [
        ":zenoh_helpers",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
    alwayslink = True,
)

cc_test(
    name = "zenoh_handle_test",
    size = "small",
    srcs = ["zenoh_handle_test.cc"],
    deps = [
        ":zenoh_handle",
        "//intrinsic/util/testing:gtest_wrapper_main",
    ],
)

cc_library(
    name = "zenoh_helpers",
    srcs = ["zenoh_helpers.cc"],
    hdrs = ["zenoh_helpers.h"],
    visibility = [
        "//intrinsic/platform/pubsub:__subpackages__",
    ],
    deps = [
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@or_tools//ortools/base:path",
        "@rules_cc//cc/runfiles",
    ],
)

cc_library(
    name = "zenoh_config",
    srcs = ["zenoh_config.cc"],
    hdrs = ["zenoh_config.h"],
    data = [
        ":peer_config.json",
    ],
    visibility = [
        "//intrinsic/platform/pubsub:__subpackages__",
    ],
    deps = [
        ":zenoh_helpers",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log",
    ],
)

cc_test(
    name = "zenoh_config_test",
    size = "small",
    srcs = ["zenoh_config_test.cc"],
    deps = [
        ":zenoh_config",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/flags:flag",
    ],
)
