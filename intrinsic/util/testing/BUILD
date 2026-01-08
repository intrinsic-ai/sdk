# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("//bazel:go_macros.bzl", "go_library")

# Helper utilities for handling testing.

go_library(
    name = "testio",
    srcs = ["test_io.go"],
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [
        "//intrinsic/util/path_resolver:pathresolver",
        "//intrinsic/util/proto:protoio",
        "@org_golang_google_protobuf//encoding/prototext:go_default_library",
        "@org_golang_google_protobuf//proto",
    ],
)

cc_library(
    name = "gtest_wrapper",
    testonly = 1,
    hdrs = ["gtest_wrapper.h"],
    visibility = ["//visibility:public"],
    deps = [
        ":status_payload_matchers",
        "@abseil-cpp//absl/status:status_matchers",
        "@com_github_inazarenko_protobuf_matchers//protobuf-matchers",
        "@com_google_benchmark//:benchmark",
        "@com_google_cel_cpp//internal:testing",
        "@com_google_googletest//:gtest",
    ],
)

cc_library(
    name = "gtest_wrapper_main",
    testonly = 1,
    srcs = ["gtest_wrapper_main.cc"],
    visibility = ["//visibility:public"],
    deps = [
        ":gtest_wrapper",
        "@abseil-cpp//absl/flags:parse",
        "@com_google_fuzztest//fuzztest:init_fuzztest",
    ],
)

cc_library(
    name = "status_payload_matchers",
    testonly = 1,
    hdrs = ["status_payload_matchers.h"],
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [
        "//intrinsic/util/proto:type_url",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:cord",
        "@com_google_googleapis//google/rpc:status_cc_proto",
        "@com_google_googletest//:gtest",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "status_payload_matchers_test",
    srcs = ["status_payload_matchers_test.cc"],
    deps = [
        ":status_payload_matchers",
        "//intrinsic/util/proto:type_url",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:cord",
        "@com_google_protobuf//:protobuf",
    ],
)
