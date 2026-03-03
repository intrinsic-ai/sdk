# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

package(default_visibility = ["//intrinsic/proto_tools/registry:__subpackages__"])

proto_library(
    name = "descriptor_test_message_proto",
    testonly = True,
    srcs = ["descriptor_test_message.proto"],
    deps = ["@com_google_protobuf//:timestamp_proto"],
)

proto_source_code_info_transitive_descriptor_set(
    name = "descriptor_test_message_proto_descriptor_set",
    testonly = True,
    deps = [":descriptor_test_message_proto"],
)
