# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("@rules_python//python:defs.bzl", "py_library", "py_test")
load("//bazel:go_macros.bzl", "go_library", "go_test")

# Application, process, product, and resource set describe all user-accessible configuration of a
# solution.

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "environments_cc",
    srcs = ["environments.cc"],
    hdrs = ["environments.h"],
    visibility = [
        "//visibility:public",
    ],
    deps = [
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_test(
    name = "environments_test_cc",
    srcs = ["environments_test.cc"],
    deps = [
        ":environments_cc",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

py_library(
    name = "environments_py",
    srcs = ["environments.py"],
    visibility = [
        "//visibility:public",
    ],
)

py_test(
    name = "environments_test_py",
    srcs = ["environments_test.py"],
    main = "environments_test.py",
    deps = [
        ":environments_py",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

go_library(
    name = "environments",
    srcs = ["environments.go"],
    visibility = [
        "//visibility:public",
    ],
)

go_test(
    name = "environments_test",
    srcs = ["environments_test.go"],
    library = ":environments",
)
