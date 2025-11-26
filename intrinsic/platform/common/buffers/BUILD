# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

# Specialized Buffer Implementations

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "rt_queue",
    hdrs = ["rt_queue.h"],
    deps = [
        ":rt_queue_buffer",
        "//intrinsic/icon/utils:realtime_guard",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/functional:function_ref",
        "@abseil-cpp//absl/types:optional",
    ],
)

cc_library(
    name = "realtime_write_queue",
    hdrs = ["realtime_write_queue.h"],
    deps = [
        ":rt_queue_buffer",
        "//intrinsic/icon/interprocess:binary_futex",
        "//intrinsic/icon/utils:realtime_status",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "rt_queue_buffer",
    hdrs = ["rt_queue_buffer.h"],
    visibility = ["//visibility:private"],
    deps = [
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/functional:function_ref",
        "@abseil-cpp//absl/log:check",
    ],
)

cc_library(
    name = "rt_promise",
    hdrs = ["rt_promise.h"],
    deps = [
        ":rt_queue",
        "//intrinsic/icon/interprocess:binary_futex",
        "//intrinsic/icon/interprocess:lockable_binary_futex",
        "//intrinsic/icon/testing:realtime_annotations",
        "//intrinsic/icon/utils:log",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/time",
    ],
)

cc_test(
    name = "rt_queue_buffer_test",
    size = "small",
    srcs = ["rt_queue_buffer_test.cc"],
    deps = [
        ":rt_queue_buffer",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/util/thread",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "rt_queue_multi_writer",
    hdrs = ["rt_queue_multi_writer.h"],
    deps = [
        ":rt_queue",
        "//intrinsic/icon/utils:realtime_guard",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/synchronization",
    ],
)

cc_test(
    name = "rt_queue_multi_writer_test",
    srcs = ["rt_queue_multi_writer_test.cc"],
    deps = [
        ":rt_queue",
        ":rt_queue_multi_writer",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/util/thread",
        "@abseil-cpp//absl/status",
    ],
)
