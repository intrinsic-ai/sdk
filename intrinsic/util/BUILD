# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library", "cc_test")
load("@rules_python//python:defs.bzl", "py_library")
load("//bazel:go_macros.bzl", "go_library")

package(
    default_visibility = ["//visibility:public"],
)

cc_library(
    name = "aggregate_type",
    hdrs = ["aggregate_type.h"],
)

cc_test(
    name = "aggregate_type_test",
    srcs = ["aggregate_type_test.cc"],
    deps = [
        ":aggregate_type",
        "//intrinsic/eigenmath",
        "//intrinsic/util/testing:gtest_wrapper_main",
    ],
)

cc_library(
    name = "atomic_sequence_num",
    hdrs = ["atomic_sequence_num.h"],
    deps = ["//intrinsic/production/external/intops:strong_int"],
)

cc_library(
    name = "fixed_vector",
    hdrs = ["fixed_vector.h"],
    deps = [
        "@com_google_absl//absl/container:inlined_vector",
        "@com_google_absl//absl/log",
    ],
)

cc_library(
    name = "invalid_until_set",
    hdrs = ["invalid_until_set.h"],
)

cc_library(
    name = "page_fault_info",
    srcs = ["page_fault_info.cc"],
    hdrs = ["page_fault_info.h"],
)

py_library(
    name = "decorators",
    srcs = ["decorators.py"],
)

cc_library(
    name = "eigen",
    hdrs = ["eigen.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/eigenmath:rotation_utils",
        "//intrinsic/math:pose3",
        "//intrinsic/util/status:status_macros",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/log:check",
        "@com_google_absl//absl/strings:str_format",
    ],
)

cc_test(
    name = "eigen_test",
    srcs = ["eigen_test.cc"],
    deps = [
        ":eigen",
        "//intrinsic/eigenmath",
        "//intrinsic/eigenmath:rotation_utils",
        "//intrinsic/math:pose3",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@com_google_protobuf//:protobuf_lite",
    ],
)

cc_library(
    name = "memory_lock",
    hdrs = ["memory_lock.h"],
    deps = [
        "@com_google_absl//absl/status",
    ],
)

cc_library(
    name = "string_type",
    hdrs = ["string_type.h"],
    deps = [
        "@com_google_absl//absl/flags:marshalling",
        "@com_google_absl//absl/hash",
        "@com_google_absl//absl/meta:type_traits",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/strings:cord",
    ],
)

cc_library(
    name = "proto_time",
    srcs = ["proto_time.cc"],
    hdrs = ["proto_time.h"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_cc_proto",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/time",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "log_lines",
    srcs = ["log_lines.cc"],
    hdrs = ["log_lines.h"],
    deps = [
        "@com_google_absl//absl/base:log_severity",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/strings",
    ],
)

go_library(
    name = "flag",
    srcs = ["flag.go"],
)
