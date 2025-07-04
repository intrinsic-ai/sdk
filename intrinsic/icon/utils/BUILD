# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library", "cc_test")
load("@rules_pkg//:pkg.bzl", "pkg_tar")
load("//bazel:go_macros.bzl", "go_binary")
load("//intrinsic/icon/testing:cc_test_and_malloc_test.bzl", "cc_test_and_malloc_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "constants",
    hdrs = ["constants.h"],
    deps = [
        "//intrinsic/eigenmath",
    ],
)

cc_library(
    name = "fixed_string",
    hdrs = ["fixed_string.h"],
    deps = [
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/strings",
    ],
)

cc_library(
    name = "fixed_str_cat",
    hdrs = ["fixed_str_cat.h"],
    deps = [
        ":fixed_string",
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/strings",
    ],
)

cc_library(
    name = "realtime_status",
    srcs = ["realtime_status.cc"],
    hdrs = ["realtime_status.h"],
    deps = [
        ":fixed_str_cat",
        ":fixed_string",
        ":realtime_guard",
        "//intrinsic/icon/testing:realtime_annotations",
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/strings",
    ],
)

cc_library(
    name = "realtime_status_or",
    hdrs = ["realtime_status_or.h"],
    deps = [
        ":realtime_status",
        "//intrinsic/icon/testing:realtime_annotations",
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/log:check",
        "@com_google_absl//absl/status",
    ],
)

cc_library(
    name = "realtime_status_matchers",
    testonly = True,
    hdrs = ["realtime_status_matchers.h"],
    deps = [
        ":realtime_status",
        ":realtime_status_macro",
        ":realtime_status_or",
        "//intrinsic/util/testing:gtest_wrapper",
        "@com_google_absl//absl/status",
    ],
)

cc_library(
    name = "realtime_status_macro",
    hdrs = ["realtime_status_macro.h"],
    deps = [
        ":fixed_str_cat",
        ":realtime_status",
        "//intrinsic/util/status:status_macros",
        "@com_google_absl//absl/base:core_headers",
    ],
)

cc_library(
    name = "exit_code",
    hdrs = ["exit_code.h"],
)

go_binary(
    name = "init_hwm",
    srcs = ["init_hwm.go"],
    deps = [
        "//intrinsic/production:intrinsic",
        "@com_github_golang_glog//:go_default_library",
    ],
)

pkg_tar(
    name = "init_hwm_tar",
    srcs = [":init_hwm"],
    include_runfiles = True,
)

cc_library(
    name = "core_time",
    srcs = [
        "clock.cc",
    ],
    hdrs = [
        "clock.h",
    ],
    deps = [
        ":log_internal",
        ":time",
        "//intrinsic/icon/testing:realtime_annotations",
        "@com_google_absl//absl/log:check",
    ],
)

cc_library(
    name = "time",
    srcs = [
        "duration.cc",
    ],
    hdrs = [
        "clock_base.h",
        "duration.h",
    ],
    deps = [
        "@com_google_absl//absl/base:core_headers",
    ],
)

cc_library(
    name = "metrics_logger",
    srcs = ["metrics_logger.cc"],
    hdrs = ["metrics_logger.h"],
    deps = [
        ":realtime_metrics",
        "//intrinsic/icon/testing:realtime_annotations",
        "//intrinsic/logging:data_logger_client",
        "//intrinsic/logging/proto:log_item_cc_proto",
        "//intrinsic/performance/analysis/proto:performance_metrics_cc_proto",
        "//intrinsic/platform/common/buffers:realtime_write_queue",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/thread",
        "//intrinsic/util/thread:thread_options",
        "//intrinsic/util/thread:thread_utils",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/status",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "realtime_log_sink",
    srcs = ["realtime_log_sink.cc"],
    hdrs = ["realtime_log_sink.h"],
    deps = [
        ":log_sink",
        ":realtime_guard",
        "//intrinsic/icon/interprocess:binary_futex",
        "//intrinsic/platform/common/buffers:realtime_write_queue",
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/container:flat_hash_map",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/synchronization",
        "@com_google_absl//absl/time",
    ],
)

cc_library(
    name = "log_internal",
    srcs = ["log_internal.cc"],
    hdrs = ["log_internal.h"],
    visibility = ["//visibility:private"],
    deps = [
        ":fixed_string",
        ":log_sink",
        "//intrinsic/icon/release:source_location",
        "//intrinsic/icon/testing:realtime_annotations",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/time",
    ],
)

cc_library(
    name = "log",
    srcs = ["log.cc"],
    hdrs = ["log.h"],
    deps = [
        ":log_internal",
        ":log_sink",
        ":realtime_log_sink",
        "//intrinsic/icon/release:source_location",
        "@com_google_absl//absl/base:log_severity",
    ],
)

cc_test_and_malloc_test(
    name = "log_test",
    srcs = ["log_test.cc"],
    deps = [
        ":log",
        ":log_internal",
        "//intrinsic/icon/release:source_location",
        "//intrinsic/util/thread",
        "@com_google_absl//absl/strings",
    ],
)

cc_library(
    name = "log_sink",
    srcs = ["log_sink.cc"],
    hdrs = ["log_sink.h"],
    deps = [
        "@com_google_absl//absl/base:log_severity",
        "@com_google_absl//absl/log:check",
        "@com_google_absl//absl/time",
    ],
)

cc_library(
    name = "realtime_stack_trace",
    srcs = [
        "realtime_stack_trace.cc",
    ],
    hdrs = [
        "log.h",  # Avoid circular dep realtime_stack_trace <-> realtime_log_sink.
        "realtime_stack_trace.h",
    ],
    linkopts = [
        "-ldl",  # for dlfcn.h, dladdr in realtime_stack_trace.cc
    ],
    deps = [
        ":fixed_str_cat",
        ":fixed_string",
        ":log_internal",
        ":log_sink",
        "//intrinsic/icon/release:source_location",
        "//intrinsic/icon/testing:realtime_annotations",
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/base:log_severity",
        "@com_google_absl//absl/debugging:stacktrace",
        "@com_google_absl//absl/log:check",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/synchronization",
        "@com_google_absl//absl/types:span",
    ],
)

cc_library(
    name = "realtime_guard",
    srcs = [
        "log.h",  # Avoid circular dep realtime_guard <-> realtime_log_sink.
        "realtime_guard.cc",
    ],
    hdrs = [
        "realtime_guard.h",
    ],
    linkopts = [
        "-ldl",  # for dlfcn.h, dladdr in realtime_guard.cc
    ],
    deps = [
        ":log_internal",
        ":log_sink",
        ":realtime_stack_trace",
        "//intrinsic/icon/release:source_location",
        "@com_google_absl//absl/base:log_severity",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/log:check",
        "@com_google_absl//absl/strings",
    ],
)

cc_test_and_malloc_test(
    name = "realtime_guard_test",
    srcs = ["realtime_guard_test.cc"],
    deps = [
        ":realtime_guard",
        ":realtime_stack_trace",
    ],
)

cc_library(
    name = "shutdown_signals",
    srcs = ["shutdown_signals.cc"],
    hdrs = ["shutdown_signals.h"],
    deps = [
        "@com_google_absl//absl/strings:string_view",
    ],
)

cc_library(
    name = "current_cycle",
    hdrs = ["current_cycle.h"],
)

cc_library(
    name = "async_buffer",
    hdrs = ["async_buffer.h"],
    deps = ["@com_google_absl//absl/log:check"],
)

cc_test(
    name = "async_buffer_test",
    srcs = ["async_buffer_test.cc"],
    deps = [
        ":async_buffer",
        "//intrinsic/icon/release:source_location",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/util/thread",
        "@com_google_absl//absl/strings",
    ],
)

cc_library(
    name = "multiple_producer_single_consumer_async_buffer",
    hdrs = ["multiple_producer_single_consumer_async_buffer.h"],
    deps = [
        ":async_buffer",
        "//intrinsic/icon/testing:realtime_annotations",
        "@com_google_absl//absl/base:core_headers",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/synchronization",
    ],
)

cc_test(
    name = "multiple_producer_single_consumer_async_buffer_test",
    srcs = ["multiple_producer_single_consumer_async_buffer_test.cc"],
    deps = [
        ":multiple_producer_single_consumer_async_buffer",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/util/thread",
        "@com_google_absl//absl/synchronization",
    ],
)

cc_library(
    name = "async_request",
    hdrs = ["async_request.h"],
    deps = [
        ":realtime_status",
        "//intrinsic/icon/testing:realtime_annotations",
        "//intrinsic/platform/common/buffers:rt_promise",
    ],
)

cc_library(
    name = "realtime_metrics",
    srcs = ["realtime_metrics.cc"],
    hdrs = ["realtime_metrics.h"],
    deps = [
        ":log",
        ":realtime_status",
        ":realtime_status_macro",
        "//intrinsic/icon/testing:realtime_annotations",
        "//intrinsic/performance/analysis/proto:performance_metrics_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/strings:str_format",
        "@com_google_absl//absl/time",
        "@com_google_absl//absl/types:span",
        "@com_google_googletest//:gtest_prod",
        "@com_google_protobuf//:protobuf",
    ],
)
