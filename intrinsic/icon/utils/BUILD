# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("@rules_pkg//:pkg.bzl", "pkg_tar")
load("//bazel:go_macros.bzl", "go_binary")
load("//intrinsic/icon/testing:cc_test_and_malloc_test.bzl", "cc_test_and_malloc_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "bitset",
    hdrs = ["bitset.h"],
)

cc_test(
    name = "bitset_test",
    srcs = ["bitset_test.cc"],
    deps = [
        ":bitset",
        "//intrinsic/util/testing:gtest_wrapper_main",
    ],
)

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
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "fixed_str_cat",
    hdrs = ["fixed_str_cat.h"],
    deps = [
        ":fixed_string",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/strings",
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
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "realtime_status_or",
    hdrs = ["realtime_status_or.h"],
    deps = [
        ":realtime_status",
        "//intrinsic/icon/testing:realtime_annotations",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
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
        "@abseil-cpp//absl/status",
    ],
)

cc_library(
    name = "realtime_status_macro",
    hdrs = ["realtime_status_macro.h"],
    deps = [
        ":fixed_str_cat",
        ":realtime_status",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:core_headers",
    ],
)

cc_library(
    name = "exit_code",
    hdrs = ["exit_code.h"],
)

go_binary(
    name = "init_hwm",
    srcs = ["init_hwm.go"],
    importpath = "intrinsic/icon/utils/init_hwm",
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
        "@abseil-cpp//absl/log:check",
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
        "@abseil-cpp//absl/base:core_headers",
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
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
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
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
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
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
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
        "@abseil-cpp//absl/base:log_severity",
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
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "log_sink",
    srcs = ["log_sink.cc"],
    hdrs = ["log_sink.h"],
    deps = [
        "@abseil-cpp//absl/base:log_severity",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/time",
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
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/base:log_severity",
        "@abseil-cpp//absl/debugging:stacktrace",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/types:span",
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
        "@abseil-cpp//absl/base:log_severity",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/strings",
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
        "@abseil-cpp//absl/strings:string_view",
    ],
)

cc_library(
    name = "current_cycle",
    hdrs = ["current_cycle.h"],
)

cc_library(
    name = "async_buffer",
    hdrs = ["async_buffer.h"],
    deps = ["@abseil-cpp//absl/log:check"],
)

cc_test(
    name = "async_buffer_test",
    srcs = ["async_buffer_test.cc"],
    deps = [
        ":async_buffer",
        "//intrinsic/icon/release:source_location",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/util/thread",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "multiple_producer_single_consumer_async_buffer",
    hdrs = ["multiple_producer_single_consumer_async_buffer.h"],
    deps = [
        ":async_buffer",
        "//intrinsic/icon/testing:realtime_annotations",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/synchronization",
    ],
)

cc_test(
    name = "multiple_producer_single_consumer_async_buffer_test",
    srcs = ["multiple_producer_single_consumer_async_buffer_test.cc"],
    deps = [
        ":multiple_producer_single_consumer_async_buffer",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/util/thread",
        "@abseil-cpp//absl/synchronization",
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
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/time",
        "@abseil-cpp//absl/types:span",
        "@com_google_googletest//:gtest_prod",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "inspection_publisher",
    hdrs = ["inspection_publisher.h"],
    deps = [
        "//intrinsic/assets/services/proto/v1:service_inspection_cc_proto",
        "//intrinsic/icon/hal:module_config",
        "//intrinsic/platform/pubsub",
        "//intrinsic/platform/pubsub:publisher",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/time",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "inspection_publisher_test",
    srcs = ["inspection_publisher_test.cc"],
    deps = [
        ":inspection_publisher",
        "//intrinsic/assets/services/proto/v1:service_inspection_cc_proto",
        "//intrinsic/icon/hal:hardware_module_inspection_cc_proto",
        "//intrinsic/platform/pubsub",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
    ],
)
