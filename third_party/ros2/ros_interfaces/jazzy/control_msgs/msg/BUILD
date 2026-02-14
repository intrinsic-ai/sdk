load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/control_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "admittance_controller_state_proto",
    srcs = ["admittance_controller_state.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:quaternion_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:transform_stamped_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_stamped_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:wrench_stamped_proto",
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:joint_state_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:float64_multi_array_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:int8_multi_array_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:string_proto",
    ],
)

py_proto_library(
    name = "admittance_controller_state_py_pb2",
    deps = ["admittance_controller_state_proto"],
)

cc_proto_library(
    name = "admittance_controller_state_cc_proto",
    deps = ["admittance_controller_state_proto"],
)

go_proto_library(
    name = "admittance_controller_state_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/admittance_controller_state_go_proto",
    protos = ["admittance_controller_state_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:quaternion_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:transform_stamped_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_stamped_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:wrench_stamped_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:joint_state_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:float64_multi_array_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:int8_multi_array_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:string_go_proto",
    ],
)

proto_library(
    name = "dynamic_interface_group_values_proto",
    srcs = ["dynamic_interface_group_values.proto"],
    deps = [
        ":interface_value_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "dynamic_interface_group_values_py_pb2",
    deps = ["dynamic_interface_group_values_proto"],
)

cc_proto_library(
    name = "dynamic_interface_group_values_cc_proto",
    deps = ["dynamic_interface_group_values_proto"],
)

go_proto_library(
    name = "dynamic_interface_group_values_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/dynamic_interface_group_values_go_proto",
    protos = ["dynamic_interface_group_values_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/control_msgs/msg:interface_value_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
)

proto_library(
    name = "dynamic_interface_values_proto",
    srcs = ["dynamic_interface_values.proto"],
    deps = [
        ":interface_value_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "dynamic_interface_values_py_pb2",
    deps = ["dynamic_interface_values_proto"],
)

cc_proto_library(
    name = "dynamic_interface_values_cc_proto",
    deps = ["dynamic_interface_values_proto"],
)

go_proto_library(
    name = "dynamic_interface_values_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/dynamic_interface_values_go_proto",
    protos = ["dynamic_interface_values_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/control_msgs/msg:interface_value_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
)

proto_library(
    name = "dynamic_joint_state_proto",
    srcs = ["dynamic_joint_state.proto"],
    deps = [
        ":interface_value_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "dynamic_joint_state_py_pb2",
    deps = ["dynamic_joint_state_proto"],
)

cc_proto_library(
    name = "dynamic_joint_state_cc_proto",
    deps = ["dynamic_joint_state_proto"],
)

go_proto_library(
    name = "dynamic_joint_state_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/dynamic_joint_state_go_proto",
    protos = ["dynamic_joint_state_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/control_msgs/msg:interface_value_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
)

proto_library(
    name = "gripper_command_proto",
    srcs = ["gripper_command.proto"],
)

py_proto_library(
    name = "gripper_command_py_pb2",
    deps = ["gripper_command_proto"],
)

cc_proto_library(
    name = "gripper_command_cc_proto",
    deps = ["gripper_command_proto"],
)

go_proto_library(
    name = "gripper_command_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/gripper_command_go_proto",
    protos = ["gripper_command_proto"],
)

proto_library(
    name = "interface_value_proto",
    srcs = ["interface_value.proto"],
)

py_proto_library(
    name = "interface_value_py_pb2",
    deps = ["interface_value_proto"],
)

cc_proto_library(
    name = "interface_value_cc_proto",
    deps = ["interface_value_proto"],
)

go_proto_library(
    name = "interface_value_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/interface_value_go_proto",
    protos = ["interface_value_proto"],
)

proto_library(
    name = "joint_component_tolerance_proto",
    srcs = ["joint_component_tolerance.proto"],
)

py_proto_library(
    name = "joint_component_tolerance_py_pb2",
    deps = ["joint_component_tolerance_proto"],
)

cc_proto_library(
    name = "joint_component_tolerance_cc_proto",
    deps = ["joint_component_tolerance_proto"],
)

go_proto_library(
    name = "joint_component_tolerance_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/joint_component_tolerance_go_proto",
    protos = ["joint_component_tolerance_proto"],
)

proto_library(
    name = "joint_controller_state_proto",
    srcs = ["joint_controller_state.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "joint_controller_state_py_pb2",
    deps = ["joint_controller_state_proto"],
)

cc_proto_library(
    name = "joint_controller_state_cc_proto",
    deps = ["joint_controller_state_proto"],
)

go_proto_library(
    name = "joint_controller_state_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/joint_controller_state_go_proto",
    protos = ["joint_controller_state_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
)

proto_library(
    name = "joint_jog_proto",
    srcs = ["joint_jog.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "joint_jog_py_pb2",
    deps = ["joint_jog_proto"],
)

cc_proto_library(
    name = "joint_jog_cc_proto",
    deps = ["joint_jog_proto"],
)

go_proto_library(
    name = "joint_jog_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/joint_jog_go_proto",
    protos = ["joint_jog_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
)

proto_library(
    name = "joint_tolerance_proto",
    srcs = ["joint_tolerance.proto"],
)

py_proto_library(
    name = "joint_tolerance_py_pb2",
    deps = ["joint_tolerance_proto"],
)

cc_proto_library(
    name = "joint_tolerance_cc_proto",
    deps = ["joint_tolerance_proto"],
)

go_proto_library(
    name = "joint_tolerance_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/joint_tolerance_go_proto",
    protos = ["joint_tolerance_proto"],
)

proto_library(
    name = "joint_trajectory_controller_state_proto",
    srcs = ["joint_trajectory_controller_state.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:joint_trajectory_point_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:multi_dof_joint_trajectory_point_proto",
    ],
)

py_proto_library(
    name = "joint_trajectory_controller_state_py_pb2",
    deps = ["joint_trajectory_controller_state_proto"],
)

cc_proto_library(
    name = "joint_trajectory_controller_state_cc_proto",
    deps = ["joint_trajectory_controller_state_proto"],
)

go_proto_library(
    name = "joint_trajectory_controller_state_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/joint_trajectory_controller_state_go_proto",
    protos = ["joint_trajectory_controller_state_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:joint_trajectory_point_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:multi_dof_joint_trajectory_point_go_proto",
    ],
)

proto_library(
    name = "mecanum_drive_controller_state_proto",
    srcs = ["mecanum_drive_controller_state.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "mecanum_drive_controller_state_py_pb2",
    deps = ["mecanum_drive_controller_state_proto"],
)

cc_proto_library(
    name = "mecanum_drive_controller_state_cc_proto",
    deps = ["mecanum_drive_controller_state_proto"],
)

go_proto_library(
    name = "mecanum_drive_controller_state_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/mecanum_drive_controller_state_go_proto",
    protos = ["mecanum_drive_controller_state_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
)

proto_library(
    name = "multi_dof_command_proto",
    srcs = ["multi_dof_command.proto"],
)

py_proto_library(
    name = "multi_dof_command_py_pb2",
    deps = ["multi_dof_command_proto"],
)

cc_proto_library(
    name = "multi_dof_command_cc_proto",
    deps = ["multi_dof_command_proto"],
)

go_proto_library(
    name = "multi_dof_command_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/multi_dof_command_go_proto",
    protos = ["multi_dof_command_proto"],
)

proto_library(
    name = "multi_dof_state_stamped_proto",
    srcs = ["multi_dof_state_stamped.proto"],
    deps = [
        ":single_dof_state_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "multi_dof_state_stamped_py_pb2",
    deps = ["multi_dof_state_stamped_proto"],
)

cc_proto_library(
    name = "multi_dof_state_stamped_cc_proto",
    deps = ["multi_dof_state_stamped_proto"],
)

go_proto_library(
    name = "multi_dof_state_stamped_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/multi_dof_state_stamped_go_proto",
    protos = ["multi_dof_state_stamped_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/control_msgs/msg:single_dof_state_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
)

proto_library(
    name = "pid_state_proto",
    srcs = ["pid_state.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "pid_state_py_pb2",
    deps = ["pid_state_proto"],
)

cc_proto_library(
    name = "pid_state_cc_proto",
    deps = ["pid_state_proto"],
)

go_proto_library(
    name = "pid_state_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/pid_state_go_proto",
    protos = ["pid_state_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
)

proto_library(
    name = "single_dof_state_proto",
    srcs = ["single_dof_state.proto"],
)

py_proto_library(
    name = "single_dof_state_py_pb2",
    deps = ["single_dof_state_proto"],
)

cc_proto_library(
    name = "single_dof_state_cc_proto",
    deps = ["single_dof_state_proto"],
)

go_proto_library(
    name = "single_dof_state_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/single_dof_state_go_proto",
    protos = ["single_dof_state_proto"],
)

proto_library(
    name = "single_dof_state_stamped_proto",
    srcs = ["single_dof_state_stamped.proto"],
    deps = [
        ":single_dof_state_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "single_dof_state_stamped_py_pb2",
    deps = ["single_dof_state_stamped_proto"],
)

cc_proto_library(
    name = "single_dof_state_stamped_cc_proto",
    deps = ["single_dof_state_stamped_proto"],
)

go_proto_library(
    name = "single_dof_state_stamped_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/single_dof_state_stamped_go_proto",
    protos = ["single_dof_state_stamped_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/control_msgs/msg:single_dof_state_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
)

proto_library(
    name = "speed_scaling_factor_proto",
    srcs = ["speed_scaling_factor.proto"],
)

py_proto_library(
    name = "speed_scaling_factor_py_pb2",
    deps = ["speed_scaling_factor_proto"],
)

cc_proto_library(
    name = "speed_scaling_factor_cc_proto",
    deps = ["speed_scaling_factor_proto"],
)

go_proto_library(
    name = "speed_scaling_factor_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/speed_scaling_factor_go_proto",
    protos = ["speed_scaling_factor_proto"],
)

proto_library(
    name = "steering_controller_status_proto",
    srcs = ["steering_controller_status.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto"],
)

py_proto_library(
    name = "steering_controller_status_py_pb2",
    deps = ["steering_controller_status_proto"],
)

cc_proto_library(
    name = "steering_controller_status_cc_proto",
    deps = ["steering_controller_status_proto"],
)

go_proto_library(
    name = "steering_controller_status_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/msg/steering_controller_status_go_proto",
    protos = ["steering_controller_status_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
)
