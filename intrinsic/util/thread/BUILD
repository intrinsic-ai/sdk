# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library", "cc_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "thread_options",
    srcs = ["thread_options.cc"],
    hdrs = ["thread_options.h"],
    deps = ["@com_google_absl//absl/strings"],
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
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/functional:any_invocable",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/synchronization",
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
        "@com_google_absl//absl/synchronization",
        "@com_google_absl//absl/time",
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
    deps = ["@com_google_absl//absl/base:core_headers"],
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
        "@com_google_absl//absl/time",
    ],
)

cc_library(
    name = "util",
    srcs = ["util.cc"],
    hdrs = ["util.h"],
    deps = [
        "//intrinsic/util/status:status_macros",
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/container:flat_hash_set",
        "@com_google_absl//absl/functional:any_invocable",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/synchronization",
        "@com_google_absl//absl/time",
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
        "@com_google_absl//absl/container:flat_hash_set",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/synchronization",
        "@com_google_absl//absl/time",
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
        "@com_google_absl//absl/time",
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
        "@com_google_absl//absl/log:check",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/synchronization",
        "@com_google_absl//absl/time",
    ],
)
