# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")

package(
    default_visibility = [
        "//intrinsic/robot_definitions:__subpackages__",
        "//visibility:public",
    ],
)

proto_library(
    name = "digital_input_output_proto",
    srcs = ["digital_input_output.proto"],
)

cc_proto_library(
    name = "digital_input_output_cc_proto",
    deps = [":digital_input_output_proto"],
)
