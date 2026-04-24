# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("//intrinsic/platform:flatbuffers.bzl", "cc_flatbuffers_library", "flatbuffers_library")

package(default_visibility = [
    "//sdk:__subpackages__",
    "//visibility:public",
])

flatbuffers_library(
    name = "adio_fbs",
    srcs = ["adio.fbs"],
)

cc_flatbuffers_library(
    name = "adio_fbs_cc",
    deps = [":adio_fbs"],
)

cc_library(
    name = "adio_fbs_utils",
    srcs = ["adio_utils.cc"],
    hdrs = [
        "adio_utils.h",
    ],
    deps = [
        ":adio_fbs_cc",
        "@com_github_google_flatbuffers//:flatbuffers",
    ],
)

flatbuffers_library(
    name = "control_mode_fbs",
    srcs = [
        "control_mode.fbs",
    ],
)

cc_flatbuffers_library(
    name = "control_mode_fbs_cc",
    deps = [":control_mode_fbs"],
)

cc_library(
    name = "control_mode_fbs_utils",
    srcs = [
        "control_mode_utils.cc",
    ],
    hdrs = [
        "control_mode_utils.h",
    ],
    deps = [
        ":control_mode_fbs_cc",
        "@com_github_google_flatbuffers//:runtime_cc",
    ],
)

flatbuffers_library(
    name = "electrical_motor_fbs",
    srcs = [
        "electrical_motor.fbs",
    ],
)

cc_flatbuffers_library(
    name = "electrical_motor_fbs_cc",
    deps = [":electrical_motor_fbs"],
)

cc_library(
    name = "electrical_motor_fbs_utils",
    srcs = ["electrical_motor_utils.cc"],
    hdrs = ["electrical_motor_utils.h"],
    deps = [
        ":electrical_motor_fbs_cc",
        "//intrinsic/icon/flatbuffers:fixed_string",
        "//intrinsic/icon/utils:log",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/strings:string_view",
        "@com_github_google_flatbuffers//:flatbuffers",
    ],
)

cc_library(
    name = "electrical_motor_hardware_interfaces",
    hdrs = ["electrical_motor_hardware_interfaces.h"],
    deps = [
        ":electrical_motor_fbs_cc",
        ":electrical_motor_fbs_utils",
        "//intrinsic/icon/hal:hardware_interface_traits",
    ],
)

flatbuffers_library(
    name = "force_sensor_fbs",
    srcs = [
        "force_sensor.fbs",
    ],
    deps = [
        "//intrinsic/icon/flatbuffers:transform_types_fbs",
    ],
)

cc_flatbuffers_library(
    name = "force_sensor_fbs_cc",
    deps = [":force_sensor_fbs"],
)

cc_library(
    name = "force_sensor_fbs_utils",
    srcs = [
        "force_sensor_utils.cc",
    ],
    hdrs = [
        "force_sensor_utils.h",
    ],
    deps = [
        ":force_sensor_fbs_cc",
        "//intrinsic/icon/flatbuffers:transform_types_fbs_cc",
        "//intrinsic/icon/utils:fixed_string",
        "@com_github_google_flatbuffers//:flatbuffers",
    ],
)

flatbuffers_library(
    name = "force_torque_fbs",
    srcs = [
        "force_torque.fbs",
    ],
    deps = [
        ":force_sensor_fbs",
        "//intrinsic/icon/flatbuffers:transform_types_fbs",
    ],
)

cc_flatbuffers_library(
    name = "force_torque_fbs_cc",
    deps = [":force_torque_fbs"],
)

cc_library(
    name = "force_torque_fbs_utils",
    srcs = [
        "force_torque_utils.cc",
    ],
    hdrs = [
        "force_torque_utils.h",
    ],
    deps = [
        ":force_sensor_fbs_cc",
        ":force_torque_fbs_cc",
        "//intrinsic/icon/flatbuffers:transform_types",  # fixdeps: keep
        "//intrinsic/icon/flatbuffers:transform_types_fbs_cc",
        "//intrinsic/icon/utils:fixed_str_cat",
        "//intrinsic/icon/utils:fixed_string",
        "@com_github_google_flatbuffers//:flatbuffers",
    ],
)

flatbuffers_library(
    name = "hardware_module_state_fbs",
    srcs = [
        "hardware_module_state.fbs",
    ],
)

cc_flatbuffers_library(
    name = "hardware_module_state_fbs_cc",
    deps = [
        ":hardware_module_state_fbs",
    ],
)

cc_library(
    name = "hardware_module_state_fbs_utils",
    srcs = ["hardware_module_state_utils.cc"],
    hdrs = [
        "hardware_module_state_utils.h",
    ],
    deps = [
        ":hardware_module_state_fbs_cc",
        "@com_github_google_flatbuffers//:flatbuffers",
    ],
)

flatbuffers_library(
    name = "icon_state_fbs",
    srcs = [
        "icon_state.fbs",
    ],
)

cc_flatbuffers_library(
    name = "icon_state_fbs_cc",
    deps = [":icon_state_fbs"],
)

cc_library(
    name = "icon_state_fbs_utils",
    srcs = ["icon_state_utils.cc"],
    hdrs = [
        "icon_state_utils.h",
    ],
    deps = [
        ":icon_state_fbs_cc",
        "@com_github_google_flatbuffers//:flatbuffers",
    ],
)

flatbuffers_library(
    name = "io_controller_fbs",
    srcs = ["io_controller.fbs"],
    deps = [
        ":adio_fbs",
    ],
)

cc_flatbuffers_library(
    name = "io_controller_fbs_cc",
    deps = [":io_controller_fbs"],
)

cc_library(
    name = "io_controller_fbs_utils",
    srcs = ["io_controller_utils.cc"],
    hdrs = ["io_controller_utils.h"],
    deps = [
        ":adio_fbs_cc",
        ":io_controller_fbs_cc",
        "@com_github_google_flatbuffers//:runtime_cc",
    ],
)

flatbuffers_library(
    name = "joint_command_fbs",
    srcs = [
        "joint_command.fbs",
    ],
)

cc_flatbuffers_library(
    name = "joint_command_fbs_cc",
    deps = [":joint_command_fbs"],
)

cc_library(
    name = "joint_command_fbs_utils",
    srcs = ["joint_command_utils.cc"],
    hdrs = [
        "joint_command_utils.h",
    ],
    deps = [
        ":joint_command_fbs_cc",
        "//intrinsic/icon/flatbuffers:flatbuffer_utils",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "@com_github_google_flatbuffers//:runtime_cc",
    ],
)

flatbuffers_library(
    name = "joint_limits_fbs",
    srcs = [
        "joint_limits.fbs",
    ],
)

cc_flatbuffers_library(
    name = "joint_limits_fbs_cc",
    deps = [":joint_limits_fbs"],
)

cc_library(
    name = "joint_limits_fbs_utils",
    srcs = ["joint_limits_utils.cc"],
    hdrs = [
        "joint_limits_utils.h",
    ],
    deps = [
        ":joint_limits_fbs_cc",
        "//intrinsic/icon/hal:hardware_interface_handle",
        "//intrinsic/icon/utils:fixed_str_cat",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/kinematics/types:joint_limits",
        "//intrinsic/kinematics/types:joint_limits_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/types:span",
        "@com_github_google_flatbuffers//:flatbuffers",
    ],
)

flatbuffers_library(
    name = "joint_state_fbs",
    srcs = [
        "joint_state.fbs",
    ],
)

cc_flatbuffers_library(
    name = "joint_state_fbs_cc",
    deps = [":joint_state_fbs"],
)

cc_library(
    name = "joint_state_fbs_utils",
    srcs = ["joint_state_utils.cc"],
    hdrs = [
        "joint_state_utils.h",
    ],
    deps = [
        ":joint_state_fbs_cc",
        "@com_github_google_flatbuffers//:flatbuffers",
    ],
)

flatbuffers_library(
    name = "robot_controller_fbs",
    srcs = [
        "robot_controller.fbs",
    ],
)

cc_flatbuffers_library(
    name = "robot_controller_fbs_cc",
    deps = [":robot_controller_fbs"],
)

cc_library(
    name = "robot_controller_fbs_utils",
    srcs = ["robot_controller_utils.cc"],
    hdrs = [
        "robot_controller_utils.h",
    ],
    deps = [
        ":robot_controller_fbs_cc",
        "@com_github_google_flatbuffers//:flatbuffers",
    ],
)

flatbuffers_library(
    name = "robot_payload_fbs",
    srcs = [
        "robot_payload.fbs",
    ],
    deps = [
        "//intrinsic/icon/flatbuffers:matrix_types_fbs",
        "//intrinsic/icon/flatbuffers:transform_types_fbs",
    ],
)

cc_flatbuffers_library(
    name = "robot_payload_fbs_cc",
    deps = [":robot_payload_fbs"],
)

cc_library(
    name = "robot_payload_utils",
    srcs = ["robot_payload_utils.cc"],
    hdrs = ["robot_payload_utils.h"],
    deps = [
        ":robot_payload_fbs_cc",
        "//intrinsic/eigenmath",
        "//intrinsic/icon/control/parts:realtime_robot_payload",
        "//intrinsic/icon/flatbuffers:control_types_utils",
        "//intrinsic/icon/flatbuffers:matrix_types_fbs_cc",
        "//intrinsic/icon/flatbuffers:matrix_types_utils",
        "//intrinsic/icon/flatbuffers:transform_types_fbs_cc",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/math:pose3",
        "//intrinsic/world/robot_payload",
        "//intrinsic/world/robot_payload:robot_payload_base",
        "@com_github_google_flatbuffers//:runtime_cc",
    ],
)

flatbuffers_library(
    name = "payload_command_fbs",
    srcs = [
        "payload_command.fbs",
    ],
    deps = [
        ":robot_payload_fbs",
    ],
)

cc_flatbuffers_library(
    name = "payload_command_fbs_cc",
    deps = [":payload_command_fbs"],
)

cc_library(
    name = "payload_command_utils",
    srcs = ["payload_command_utils.cc"],
    hdrs = ["payload_command_utils.h"],
    deps = [
        ":payload_command_fbs_cc",
        ":robot_payload_fbs_cc",
        ":robot_payload_utils",
        "@com_github_google_flatbuffers//:runtime_cc",
    ],
)

flatbuffers_library(
    name = "payload_state_fbs",
    srcs = [
        "payload_state.fbs",
    ],
    deps = [
        ":robot_payload_fbs",
    ],
)

cc_flatbuffers_library(
    name = "payload_state_fbs_cc",
    deps = [":payload_state_fbs"],
)

cc_library(
    name = "payload_state_utils",
    srcs = ["payload_state_utils.cc"],
    hdrs = ["payload_state_utils.h"],
    deps = [
        ":payload_state_fbs_cc",
        ":robot_payload_fbs_cc",
        ":robot_payload_utils",
        "@com_github_google_flatbuffers//:runtime_cc",
    ],
)
