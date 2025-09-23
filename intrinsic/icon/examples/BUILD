# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@rules_cc//cc:cc_binary.bzl", "cc_binary")
load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "adio_lib",
    srcs = ["adio_lib.cc"],
    hdrs = ["adio_lib.h"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/icon/actions:adio_cc_proto",
        "//intrinsic/icon/actions:adio_info",
        "//intrinsic/icon/cc_client:client",
        "//intrinsic/icon/cc_client:condition",
        "//intrinsic/icon/cc_client:session",
        "//intrinsic/icon/common:id_types",
        "//intrinsic/icon/proto:io_block_cc_proto",
        "//intrinsic/icon/proto:part_status_cc_proto",
        "//intrinsic/util/grpc:channel_interface",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_binary(
    name = "adio",
    srcs = ["adio.cc"],
    deps = [
        ":adio_lib",
        "//intrinsic/icon/release/portable:init_xfa_absl",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
    ],
)

proto_library(
    name = "joint_move_positions_proto",
    srcs = ["joint_move_positions.proto"],
)

cc_proto_library(
    name = "joint_move_positions_cc_proto",
    deps = [":joint_move_positions_proto"],
)

cc_library(
    name = "joint_move_lib",
    srcs = ["joint_move_lib.cc"],
    hdrs = ["joint_move_lib.h"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/icon/actions:point_to_point_move_info",
        "//intrinsic/icon/cc_client:client",
        "//intrinsic/icon/cc_client:client_utils",
        "//intrinsic/icon/cc_client:condition",
        "//intrinsic/icon/cc_client:session",
        "//intrinsic/icon/common:builtins",
        "//intrinsic/icon/common:id_types",
        "//intrinsic/kinematics/types:joint_limits",
        "//intrinsic/util/grpc:channel_interface",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
    ],
)

cc_binary(
    name = "joint_move",
    srcs = ["joint_move.cc"],
    deps = [
        ":joint_move_lib",
        "//intrinsic/icon/release/portable:init_xfa_absl",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "joint_then_cart_move_lib",
    srcs = ["joint_then_cart_move_lib.cc"],
    hdrs = ["joint_then_cart_move_lib.h"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/icon/actions:cartesian_jogging_info",
        "//intrinsic/icon/actions:point_to_point_move_info",
        "//intrinsic/icon/cc_client:client",
        "//intrinsic/icon/cc_client:client_utils",
        "//intrinsic/icon/cc_client:condition",
        "//intrinsic/icon/cc_client:robot_config",
        "//intrinsic/icon/cc_client:session",
        "//intrinsic/icon/common:builtins",
        "//intrinsic/icon/common:id_types",
        "//intrinsic/icon/proto:cart_space_cc_proto",
        "//intrinsic/icon/proto:generic_part_config_cc_proto",
        "//intrinsic/icon/proto:part_status_cc_proto",
        "//intrinsic/kinematics/types:joint_limits",
        "//intrinsic/util/grpc:channel_interface",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
    ],
)

cc_binary(
    name = "joint_then_cart_move",
    srcs = ["joint_then_cart_move.cc"],
    deps = [
        ":joint_then_cart_move_lib",
        "//intrinsic/icon/release/portable:init_xfa_absl",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
    ],
)

cc_binary(
    name = "introspection",
    srcs = ["introspection.cc"],
    deps = [
        "//intrinsic/icon/cc_client:client",
        "//intrinsic/icon/proto:cart_space_cc_proto",
        "//intrinsic/icon/proto:io_block_cc_proto",
        "//intrinsic/icon/proto:part_status_cc_proto",
        "//intrinsic/icon/proto/v1:service_cc_proto",
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "//intrinsic/icon/release/portable:init_xfa_absl",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "joint_move_loop_lib",
    srcs = ["joint_move_loop_lib.cc"],
    hdrs = ["joint_move_loop_lib.h"],
    visibility = ["//visibility:private"],
    deps = [
        ":joint_move_positions_cc_proto",
        "//intrinsic/eigenmath",
        "//intrinsic/icon/actions:point_to_point_move_info",
        "//intrinsic/icon/cc_client:client",
        "//intrinsic/icon/cc_client:client_utils",
        "//intrinsic/icon/cc_client:condition",
        "//intrinsic/icon/cc_client:session",
        "//intrinsic/icon/common:id_types",
        "//intrinsic/kinematics/types:joint_limits",
        "//intrinsic/util:eigen",
        "//intrinsic/util/grpc:channel_interface",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
    ],
)

cc_binary(
    name = "joint_move_loop",
    srcs = ["joint_move_loop.cc"],
    deps = [
        ":joint_move_loop_lib",
        ":joint_move_positions_cc_proto",
        "//intrinsic/icon/release/portable:init_xfa_absl",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/proto:get_text_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "simple_gripper_lib",
    srcs = ["simple_gripper_lib.cc"],
    hdrs = ["simple_gripper_lib.h"],
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/icon/actions:simple_gripper_cc_proto",
        "//intrinsic/icon/actions:simple_gripper_info",
        "//intrinsic/icon/cc_client:client",
        "//intrinsic/icon/cc_client:condition",
        "//intrinsic/icon/cc_client:session",
        "//intrinsic/icon/common:id_types",
        "//intrinsic/icon/proto:part_status_cc_proto",
        "//intrinsic/util/grpc:channel_interface",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_binary(
    name = "simple_gripper",
    srcs = ["simple_gripper.cc"],
    deps = [
        ":simple_gripper_lib",
        "//intrinsic/icon/release/portable:init_xfa_absl",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
    ],
)
