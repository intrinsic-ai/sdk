# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = [
    "//intrinsic/icon/control/parts/hal:__subpackages__",
    "//visibility:public",
])

proto_library(
    name = "hal_part_config_proto",
    srcs = ["hal_part_config.proto"],
)

cc_proto_library(
    name = "hal_part_config_cc_proto",
    deps = [":hal_part_config_proto"],
)

go_proto_library(
    name = "hal_part_config_go_proto",
    deps = [":hal_part_config_proto"],
)
