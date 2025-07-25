# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = [
    "//visibility:public",
])

cc_library(
    name = "binary_futex",
    srcs = ["binary_futex.cc"],
    hdrs = [
        "binary_futex.h",
    ],
    deps = [
        "//intrinsic/icon/testing:realtime_annotations",
        "//intrinsic/icon/utils:realtime_status",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "lockable_binary_futex",
    hdrs = [
        "lockable_binary_futex.h",
    ],
    deps = [
        ":binary_futex",
        "//intrinsic/icon/testing:realtime_annotations",
        "//intrinsic/icon/utils:realtime_status",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "binary_futex_condition_variable",
    srcs = ["binary_futex_condition_variable.cc"],
    hdrs = [
        "binary_futex_condition_variable.h",
    ],
    deps = [
        ":binary_futex",
        ":lockable_binary_futex",
        "//intrinsic/icon/testing:realtime_annotations",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
    ],
)
