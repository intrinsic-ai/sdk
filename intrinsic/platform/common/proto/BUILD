# Copyright 2023 Intrinsic Innovation LLC

# The protobuf definition for serialization of various type objects.

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "test_proto",
    srcs = ["test.proto"],
    visibility = ["//visibility:private"],
    deps = [
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

cc_proto_library(
    name = "test_cc_proto",
    visibility = ["//intrinsic/platform:__subpackages__"],
    deps = ["test_proto"],
)

py_proto_library(
    name = "test_py_pb2",
    visibility = ["//intrinsic/platform:__subpackages__"],
    deps = [":test_proto"],
)
