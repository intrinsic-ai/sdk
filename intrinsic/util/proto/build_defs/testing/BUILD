# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_python//python:defs.bzl", "py_test")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

package(
    default_testonly = True,
    default_visibility = ["//visibility:public"],
)

proto_library(
    name = "test_message_dep",
    srcs = ["test_message_dep.proto"],
)

cc_proto_library(
    name = "test_message_dep_cc_proto",
    deps = [":test_message_dep"],
)

proto_library(
    name = "test_message_proto",
    srcs = ["test_message.proto"],
    deps = [
        ":test_message_dep",
        "@com_google_protobuf//:duration_proto",
    ],
)

py_proto_library(
    name = "test_message_py_pb2",
    deps = [":test_message_proto"],
)

proto_source_code_info_transitive_descriptor_set(
    name = "test_message_proto_descriptors",
    deps = [":test_message_proto"],
)

py_test(
    name = "source_code_info_present_test",
    srcs = ["source_code_info_present_test.py"],
    data = [
        ":test_message_proto_descriptors",
    ],
    deps = [
        "@com_google_absl_py//absl/flags",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)
