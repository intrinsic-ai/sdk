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
    name = "sdo_value_config_proto",
    srcs = ["sdo_value_config.proto"],
    deps = [
        "//intrinsic/icon/hal/lib/ecat/v1:value_parsing_proto",
        "@com_google_protobuf//:duration_proto",
    ],
)

cc_proto_library(
    name = "sdo_value_config_cc_proto",
    deps = [":sdo_value_config_proto"],
)

go_proto_library(
    name = "sdo_value_config_go_proto",
    importpath = "intrinsic/icon/hal/lib/sdo_value/v1/sdo_value_config_go_proto",
    protos = [":sdo_value_config_proto"],
    deps = ["//intrinsic/icon/hal/lib/ecat/v1:value_parsing_go_proto"],
)
