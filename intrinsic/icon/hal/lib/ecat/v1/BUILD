# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

proto_library(
    name = "value_parsing_proto",
    srcs = ["value_parsing.proto"],
)

cc_proto_library(
    name = "value_parsing_cc_proto",
    deps = [":value_parsing_proto"],
)

go_proto_library(
    name = "value_parsing_go_proto",
    importpath = "intrinsic/icon/hal/lib/ecat/v1/value_parsing_go_proto",
    protos = [":value_parsing_proto"],
)
