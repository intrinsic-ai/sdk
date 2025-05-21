# Copyright 2023 Intrinsic Innovation LLC

# Plugins for modeling equipment behaviors in Gazebo simulation.

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@rules_cc//cc:defs.bzl", "cc_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_visibility = ["//visibility:public"],
)

cc_library(
    name = "world_model_config_plugin_constants",
    hdrs = ["world_model_config_plugin_constants.h"],
    deps = [
        "@com_google_absl//absl/strings:string_view",
    ],
)

proto_library(
    name = "sim_hardware_module_config_proto",
    srcs = ["sim_hardware_module_config.proto"],
)

go_proto_library(
    name = "sim_hardware_module_config_go_proto",
    deps = [":sim_hardware_module_config_proto"],
)

cc_proto_library(
    name = "sim_hardware_module_config_cc_proto",
    deps = [":sim_hardware_module_config_proto"],
)
