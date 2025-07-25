# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("@rules_python//python:defs.bzl", "py_library", "py_test")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

proto_library(
    name = "generic_part_config_proto",
    srcs = ["generic_part_config.proto"],
    deps = [
        ":cart_space_proto",
        "//intrinsic/kinematics/proto:skeleton_proto",
        "//intrinsic/kinematics/types:joint_limits_proto",
        "//intrinsic/math/proto:pose_proto",
    ],
)

cc_proto_library(
    name = "generic_part_config_cc_proto",
    deps = [":generic_part_config_proto"],
)

go_proto_library(
    name = "generic_part_config_go_proto",
    go_deps = [
        ":cart_space_go_proto",
        "//intrinsic/kinematics/proto:skeleton_go_proto",
        "//intrinsic/kinematics/types:joint_limits_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
    deps = [":generic_part_config_proto"],
)

py_proto_library(
    name = "generic_part_config_py_pb2",
    deps = [":generic_part_config_proto"],
)

proto_library(
    name = "part_status_proto",
    srcs = ["part_status.proto"],
    deps = [
        ":cart_space_proto",
        ":io_block_proto",
        ":safety_status_proto",
        "//intrinsic/icon/proto/v1:types_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/math/proto:quaternion_proto",
        "//intrinsic/math/proto:vector3_proto",
    ],
)

cc_proto_library(
    name = "part_status_cc_proto",
    deps = [":part_status_proto"],
)

py_proto_library(
    name = "part_status_py_pb2",
    deps = [":part_status_proto"],
)

go_proto_library(
    name = "part_status_go_proto",
    go_deps = [
        "//intrinsic/icon/proto/v1:types_go_proto",
        ":cart_space_go_proto",
        "//intrinsic/math/proto:quaternion_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
        ":io_block_go_proto",
        ":safety_status_go_proto",
    ],
    deps = [":part_status_proto"],
)

proto_library(
    name = "streaming_output_proto",
    srcs = ["streaming_output.proto"],
    deps = [
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "streaming_output_cc_proto",
    deps = [":streaming_output_proto"],
)

go_proto_library(
    name = "streaming_output_go_proto",
    deps = [":streaming_output_proto"],
)

py_proto_library(
    name = "streaming_output_py_pb2",
    deps = [":streaming_output_proto"],
)

cc_library(
    name = "eigen_conversion",
    srcs = ["eigen_conversion.cc"],
    hdrs = ["eigen_conversion.h"],
    deps = [
        ":cart_space_cc_proto",
        ":joint_space_cc_proto",
        ":matrix_cc_proto",
        "//intrinsic/eigenmath",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_test(
    name = "eigen_conversion_test",
    srcs = ["eigen_conversion_test.cc"],
    deps = [
        ":eigen_conversion",
        ":matrix_cc_proto",
        "//intrinsic/eigenmath",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@com_google_protobuf//:protobuf_lite",
    ],
)

proto_library(
    name = "matrix_proto",
    srcs = ["matrix.proto"],
)

cc_proto_library(
    name = "matrix_cc_proto",
    deps = [":matrix_proto"],
)

py_proto_library(
    name = "matrix_py_pb2",
    deps = [":matrix_proto"],
)

go_proto_library(
    name = "matrix_go_proto",
    deps = [":matrix_proto"],
)

proto_library(
    name = "cart_space_proto",
    srcs = ["cart_space.proto"],
)

cc_proto_library(
    name = "cart_space_cc_proto",
    deps = [":cart_space_proto"],
)

py_proto_library(
    name = "cart_space_py_pb2",
    deps = [":cart_space_proto"],
)

go_proto_library(
    name = "cart_space_go_proto",
    deps = [":cart_space_proto"],
)

cc_library(
    name = "cart_space_conversion",
    srcs = ["cart_space_conversion.cc"],
    hdrs = ["cart_space_conversion.h"],
    deps = [
        ":cart_space_cc_proto",
        ":eigen_conversion",
        "//intrinsic/eigenmath",
        "//intrinsic/kinematics/types:cartesian_limits",
        "//intrinsic/math:pose3",
        "//intrinsic/math:twist",
        "//intrinsic/util/status:status_builder",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_gitlab_libeigen_eigen//:eigen",
    ],
)

cc_library(
    name = "concatenate_trajectory_protos",
    srcs = ["concatenate_trajectory_protos.cc"],
    hdrs = ["concatenate_trajectory_protos.h"],
    deps = [
        ":joint_space_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

proto_library(
    name = "joint_space_proto",
    srcs = ["joint_space.proto"],
    deps = [
        "//intrinsic/kinematics/types:dynamic_limits_check_mode_proto",
        "//intrinsic/skills/proto:skill_parameter_metadata_proto",
        "@com_google_protobuf//:duration_proto",
    ],
)

cc_proto_library(
    name = "joint_space_cc_proto",
    deps = [":joint_space_proto"],
)

py_proto_library(
    name = "joint_space_py_pb2",
    deps = [":joint_space_proto"],
)

go_proto_library(
    name = "joint_space_go_proto",
    go_deps = [
        "//intrinsic/kinematics/types:dynamic_limits_check_mode_go_proto",
        "//intrinsic/skills/proto:skill_parameter_metadata_go_proto",
    ],
    deps = [":joint_space_proto"],
)

proto_library(
    name = "ik_options_proto",
    srcs = ["ik_options.proto"],
)

cc_proto_library(
    name = "ik_options_cc_proto",
    deps = [":ik_options_proto"],
)

py_proto_library(
    name = "ik_options_py_pb2",
    deps = [":ik_options_proto"],
)

proto_library(
    name = "io_block_proto",
    srcs = ["io_block.proto"],
)

py_proto_library(
    name = "io_block_py_pb2",
    deps = [":io_block_proto"],
)

go_proto_library(
    name = "io_block_go_proto",
    deps = [":io_block_proto"],
)

cc_proto_library(
    name = "io_block_cc_proto",
    deps = [":io_block_proto"],
)

py_library(
    name = "matrix_conversions",
    srcs = ["matrix_conversions.py"],
    deps = [
        ":matrix_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
    ],
)

py_test(
    name = "matrix_conversions_test",
    srcs = ["matrix_conversions_test.py"],
    deps = [
        ":matrix_conversions",
        ":matrix_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

proto_library(
    name = "linear_joint_acceleration_filter_config_proto",
    srcs = [
        "linear_joint_acceleration_filter_config.proto",
    ],
)

cc_proto_library(
    name = "linear_joint_acceleration_filter_config_cc_proto",
    deps = [":linear_joint_acceleration_filter_config_proto"],
)

go_proto_library(
    name = "linear_joint_acceleration_filter_config_go_proto",
    deps = [":linear_joint_acceleration_filter_config_proto"],
)

proto_library(
    name = "joint_position_pid_torque_controller_config_proto",
    srcs = [
        "joint_position_pid_torque_controller_config.proto",
    ],
)

go_proto_library(
    name = "joint_position_pid_torque_controller_config_go_proto",
    deps = [":joint_position_pid_torque_controller_config_proto"],
)

proto_library(
    name = "safety_status_proto",
    srcs = ["safety_status.proto"],
)

cc_proto_library(
    name = "safety_status_cc_proto",
    deps = [":safety_status_proto"],
)

py_proto_library(
    name = "safety_status_py_pb2",
    deps = [":safety_status_proto"],
)

go_proto_library(
    name = "safety_status_go_proto",
    deps = [":safety_status_proto"],
)

proto_library(
    name = "logging_mode_proto",
    srcs = ["logging_mode.proto"],
)

go_proto_library(
    name = "logging_mode_go_proto",
    deps = [":logging_mode_proto"],
)

cc_proto_library(
    name = "logging_mode_cc_proto",
    deps = [":logging_mode_proto"],
)

py_proto_library(
    name = "logging_mode_py_pb2",
    deps = [":logging_mode_proto"],
)

# Convenience proto library to add icon protos to icon python wheel
proto_library(
    name = "icon_protos",
    deps = [
        ":cart_space_proto",
        ":generic_part_config_proto",
        ":ik_options_proto",
        ":io_block_proto",
        ":joint_space_proto",
        ":matrix_proto",
        ":part_status_proto",
        ":safety_status_proto",
        ":streaming_output_proto",
        "//intrinsic/icon/proto/v1:service_proto",
        "//intrinsic/icon/proto/v1:types_proto",
    ],
)
