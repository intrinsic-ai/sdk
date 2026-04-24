# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_visibility = [
        "//sdk:__subpackages__",
        "//visibility:public",
    ],
)

proto_library(
    name = "ati_force_torque_bus_component_config_proto",
    srcs = ["ati_force_torque_bus_component_config.proto"],
)

cc_proto_library(
    name = "ati_force_torque_bus_component_config_cc_proto",
    deps = [":ati_force_torque_bus_component_config_proto"],
)

go_proto_library(
    name = "ati_force_torque_bus_component_config_go_proto",
    importpath = "intrinsic/icon/hal/lib/ati/v1/ati_force_torque_bus_component_config_go_proto",
    protos = [":ati_force_torque_bus_component_config_proto"],
)
