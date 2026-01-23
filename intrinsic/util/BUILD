# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
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
    deps = ["@or_tools//ortools/base:strong_int"],
)

cc_library(
    name = "fixed_vector",
    hdrs = ["fixed_vector.h"],
    deps = [
        "@abseil-cpp//absl/container:inlined_vector",
        "@abseil-cpp//absl/log",
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
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/strings:str_format",
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
        "@abseil-cpp//absl/status",
    ],
)

cc_library(
    name = "string_type",
    hdrs = ["string_type.h"],
    deps = [
        "@abseil-cpp//absl/flags:marshalling",
        "@abseil-cpp//absl/hash",
        "@abseil-cpp//absl/meta:type_traits",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:cord",
    ],
)

cc_library(
    name = "proto_time",
    srcs = ["proto_time.cc"],
    hdrs = ["proto_time.h"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_cc_proto",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "log_lines",
    srcs = ["log_lines.cc"],
    hdrs = ["log_lines.h"],
    deps = [
        "@abseil-cpp//absl/base:log_severity",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/strings",
    ],
)

go_library(
    name = "flag",
    srcs = ["flag.go"],
)
