# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library", "cc_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "thread_options",
    srcs = ["thread_options.cc"],
    hdrs = ["thread_options.h"],
    deps = ["@abseil-cpp//absl/strings"],
)

cc_library(
    name = "thread_utils",
    srcs = [
        "thread_utils.cc",
    ],
    hdrs = [
        "rt_thread.h",
        "thread_utils.h",
    ],
    copts = ["-D_GNU_SOURCE"],
    deps = [
        ":stop_token",
        ":sysinfo",
        ":thread",
        ":thread_options",
        "//intrinsic/icon/testing:realtime_annotations",
        "//intrinsic/icon/utils:log",  # buildcleaner: keep
        "//intrinsic/icon/utils:realtime_guard",
        "//intrinsic/icon/utils:realtime_stack_trace",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/functional:any_invocable",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
    ],
)

cc_test(
    name = "thread_utils_test",
    srcs = ["thread_utils_test.cc"],
    deps = [
        ":stop_token",
        ":thread",
        ":thread_options",
        ":thread_utils",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "thread",
    srcs = ["thread.cc"],
    hdrs = [
        "thread.h",
    ],
    deps = [
        ":stop_token",
        "//intrinsic/icon/utils:realtime_guard",
        "//intrinsic/util/status:status_macros",
    ],
)

cc_library(
    name = "stop_token",
    hdrs = ["stop_token.h"],
    deps = ["@abseil-cpp//absl/base:core_headers"],
)

cc_test(
    name = "stop_token_test",
    srcs = ["stop_token_test.cc"],
    deps = [
        ":stop_token",
        "//intrinsic/util/testing:gtest_wrapper_main",
    ],
)

cc_library(
    name = "sysinfo",
    srcs = ["sysinfo.cc"],
    hdrs = ["sysinfo.h"],
    deps = [
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "util",
    srcs = ["util.cc"],
    hdrs = ["util.h"],
    deps = [
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/functional:any_invocable",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@com_googlesource_code_re2//:re2",
    ],
)

cc_test(
    name = "util_test",
    srcs = ["util_test.cc"],
    deps = [
        ":thread",
        ":util",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@or_tools//ortools/base:file",
        "@or_tools//ortools/base:path",
    ],
)

cc_library(
    name = "lockstep",
    srcs = ["lockstep.cc"],
    hdrs = ["lockstep.h"],
    deps = [
        "//intrinsic/icon/interprocess:binary_futex",
        "//intrinsic/icon/utils:log",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "@abseil-cpp//absl/time",
    ],
)

cc_test(
    name = "lockstep_test",
    srcs = ["lockstep_test.cc"],
    deps = [
        ":lockstep",
        ":thread",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_matchers",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
    ],
)
