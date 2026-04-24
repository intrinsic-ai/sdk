# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = [
    "//intrinsic/icon/hal/lib/process_variable_value:__subpackages__",
    "//intrinsic/icon/hardware_modules:__subpackages__",
    "//sdk:__subpackages__",
])

proto_library(
    name = "process_variable_config_proto",
    srcs = ["process_variable_config.proto"],
)

cc_proto_library(
    name = "process_variable_config_cc_proto",
    deps = [":process_variable_config_proto"],
)

go_proto_library(
    name = "process_variable_config_go_proto",
    importpath = "intrinsic/icon/hal/lib/process_variable_value/v1/process_variable_config_go_proto",
    protos = [":process_variable_config_proto"],
)
