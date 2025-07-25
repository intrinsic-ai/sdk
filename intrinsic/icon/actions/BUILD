# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_python//python:defs.bzl", "py_library", "py_test")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "action_utils",
    srcs = ["action_utils.cc"],
    hdrs = ["action_utils.h"],
    deps = [
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "//intrinsic/icon/release:source_location",
        "//intrinsic/util/proto:descriptors",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

proto_library(
    name = "point_to_point_move_proto",
    srcs = ["point_to_point_move.proto"],
    deps = [
        "//intrinsic/icon/proto:joint_space_proto",
        "//intrinsic/kinematics/types:joint_limits_proto",
    ],
)

cc_proto_library(
    name = "point_to_point_move_cc_proto",
    deps = [":point_to_point_move_proto"],
)

py_proto_library(
    name = "point_to_point_move_py_pb2",
    deps = [":point_to_point_move_proto"],
)

cc_library(
    name = "point_to_point_move_info",
    srcs = ["point_to_point_move_info.cc"],
    hdrs = ["point_to_point_move_info.h"],
    deps = [
        ":point_to_point_move_cc_proto",
        "//intrinsic/icon/proto:joint_space_cc_proto",
        "//intrinsic/kinematics/types:joint_limits",
        "@abseil-cpp//absl/types:optional",
        "@abseil-cpp//absl/types:span",
    ],
)

proto_library(
    name = "tare_force_torque_sensor_proto",
    srcs = ["tare_force_torque_sensor.proto"],
)

cc_proto_library(
    name = "tare_force_torque_sensor_cc_proto",
    deps = [":tare_force_torque_sensor_proto"],
)

py_proto_library(
    name = "tare_force_torque_sensor_py_pb2",
    deps = [":tare_force_torque_sensor_proto"],
)

cc_library(
    name = "tare_force_torque_sensor_info",
    hdrs = ["tare_force_torque_sensor_info.h"],
    deps = [":tare_force_torque_sensor_cc_proto"],
)

proto_library(
    name = "cartesian_jogging_proto",
    srcs = ["cartesian_jogging.proto"],
    deps = [
        "//intrinsic/icon/proto:cart_space_proto",
        "//intrinsic/kinematics/types:joint_limits_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/math/proto:quaternion_proto",
    ],
)

cc_proto_library(
    name = "cartesian_jogging_cc_proto",
    deps = [":cartesian_jogging_proto"],
)

py_proto_library(
    name = "cartesian_jogging_py_pb2",
    deps = [":cartesian_jogging_proto"],
)

cc_library(
    name = "cartesian_jogging_info",
    hdrs = ["cartesian_jogging_info.h"],
    deps = [
        ":cartesian_jogging_cc_proto",
    ],
)

proto_library(
    name = "adio_proto",
    srcs = ["adio.proto"],
)

py_proto_library(
    name = "adio_py_pb2",
    deps = [":adio_proto"],
)

cc_proto_library(
    name = "adio_cc_proto",
    deps = [":adio_proto"],
)

go_proto_library(
    name = "adio_go_proto",
    deps = [":adio_proto"],
)

cc_library(
    name = "adio_info",
    hdrs = ["adio_info.h"],
    deps = [":adio_cc_proto"],
)

proto_library(
    name = "simple_gripper_proto",
    srcs = ["simple_gripper.proto"],
)

cc_proto_library(
    name = "simple_gripper_cc_proto",
    deps = [":simple_gripper_proto"],
)

py_proto_library(
    name = "simple_gripper_py_pb2",
    deps = [":simple_gripper_proto"],
)

cc_library(
    name = "simple_gripper_info",
    hdrs = ["simple_gripper_info.h"],
    deps = [":simple_gripper_cc_proto"],
)

proto_library(
    name = "trajectory_tracking_action_proto",
    srcs = ["trajectory_tracking_action.proto"],
    deps = ["//intrinsic/icon/proto:joint_space_proto"],
)

cc_proto_library(
    name = "trajectory_tracking_action_cc_proto",
    deps = [":trajectory_tracking_action_proto"],
)

py_proto_library(
    name = "trajectory_tracking_action_py_pb2",
    deps = [":trajectory_tracking_action_proto"],
)

cc_library(
    name = "trajectory_tracking_action_info",
    hdrs = ["trajectory_tracking_action_info.h"],
    deps = [":trajectory_tracking_action_cc_proto"],
)

cc_library(
    name = "stop_info",
    srcs = [],
    hdrs = ["stop_info.h"],
)

py_library(
    name = "trajectory_tracking_action_utils",
    srcs = ["trajectory_tracking_action_utils.py"],
    visibility = ["//intrinsic/icon/python:__subpackages__"],
    deps = [
        ":trajectory_tracking_action_py_pb2",
        "//intrinsic/icon/proto:joint_space_py_pb2",
        "//intrinsic/icon/python:icon",
    ],
)

py_library(
    name = "tare_force_torque_sensor_utils",
    srcs = ["tare_force_torque_sensor_utils.py"],
    visibility = ["//intrinsic/icon/python:__subpackages__"],
    deps = [
        ":tare_force_torque_sensor_py_pb2",
        "//intrinsic/icon/python:icon",
    ],
)

py_library(
    name = "point_to_point_move_utils",
    srcs = ["point_to_point_move_utils.py"],
    visibility = ["//intrinsic/icon/python:__subpackages__"],
    deps = [
        ":point_to_point_move_py_pb2",
        "//intrinsic/icon/proto:joint_space_py_pb2",
        "//intrinsic/icon/python:icon",
        "//intrinsic/kinematics/types:joint_limits_py_pb2",
    ],
)

py_library(
    name = "stop_utils",
    srcs = ["stop_utils.py"],
    visibility = ["//intrinsic:__subpackages__"],
    deps = ["//intrinsic/icon/python:icon"],
)

py_library(
    name = "wait_for_settling_utils",
    srcs = ["wait_for_settling_utils.py"],
    visibility = ["//intrinsic/icon/python:__subpackages__"],
    deps = [
        ":wait_for_settling_action_py_pb2",
        "//intrinsic/icon/python:icon",
    ],
)

py_library(
    name = "adio_utils",
    srcs = ["adio_utils.py"],
    deps = [
        ":adio_py_pb2",
        "//intrinsic/icon/python:icon",
    ],
)

py_test(
    name = "adio_utils_test",
    srcs = ["adio_utils_test.py"],
    deps = [
        ":adio_py_pb2",
        ":adio_utils",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

py_test(
    name = "wait_for_settling_utils_test",
    srcs = ["wait_for_settling_utils_test.py"],
    deps = [
        ":wait_for_settling_action_py_pb2",
        ":wait_for_settling_utils",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

py_test(
    name = "point_to_point_move_utils_test",
    srcs = ["point_to_point_move_utils_test.py"],
    deps = [
        ":point_to_point_move_py_pb2",
        ":point_to_point_move_utils",
        "//intrinsic/kinematics/types:joint_limits_py_pb2",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

py_test(
    name = "trajectory_tracking_action_utils_test",
    srcs = ["trajectory_tracking_action_utils_test.py"],
    deps = [
        ":trajectory_tracking_action_py_pb2",
        ":trajectory_tracking_action_utils",
        "//intrinsic/icon/proto:joint_space_py_pb2",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)

py_test(
    name = "tare_force_torque_sensor_utils_test",
    srcs = ["tare_force_torque_sensor_utils_test.py"],
    deps = [
        ":tare_force_torque_sensor_py_pb2",
        ":tare_force_torque_sensor_utils",
        "@com_google_absl_py//absl/testing:absltest",
    ],
)

py_test(
    name = "stop_utils_test",
    srcs = ["stop_utils_test.py"],
    deps = [
        ":stop_utils",
        "@com_google_absl_py//absl/testing:absltest",
        "@com_google_protobuf//:protobuf_python",
    ],
)

proto_library(
    name = "wait_for_settling_action_proto",
    srcs = ["wait_for_settling_action.proto"],
)

cc_proto_library(
    name = "wait_for_settling_action_cc_proto",
    deps = [":wait_for_settling_action_proto"],
)

py_proto_library(
    name = "wait_for_settling_action_py_pb2",
    deps = [":wait_for_settling_action_proto"],
)

cc_library(
    name = "wait_for_settling_action_info",
    srcs = [],
    hdrs = ["wait_for_settling_action_info.h"],
    deps = [
        ":wait_for_settling_action_cc_proto",
    ],
)

cc_library(
    name = "force_stop_action_info",
    hdrs = ["force_stop_action_info.h"],
    visibility = ["//intrinsic/icon/control/actions:__pkg__"],
)
