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
    name = "follow_joint_trajectory_proto",
    srcs = ["follow_joint_trajectory.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_proto",
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
        "//third_party/ros2/ros_interfaces/jazzy/control_msgs/msg:joint_component_tolerance_proto",
        "//third_party/ros2/ros_interfaces/jazzy/control_msgs/msg:joint_tolerance_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:joint_trajectory_point_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:joint_trajectory_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:multi_dof_joint_trajectory_point_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:multi_dof_joint_trajectory_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_proto",
    ],
)

py_proto_library(
    name = "follow_joint_trajectory_py_pb2",
    deps = ["follow_joint_trajectory_proto"],
)

cc_proto_library(
    name = "follow_joint_trajectory_cc_proto",
    deps = ["follow_joint_trajectory_proto"],
)

go_proto_library(
    name = "follow_joint_trajectory_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/action/follow_joint_trajectory_go_proto",
    protos = ["follow_joint_trajectory_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/control_msgs/msg:joint_component_tolerance_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/control_msgs/msg:joint_tolerance_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:joint_trajectory_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:joint_trajectory_point_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:multi_dof_joint_trajectory_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:multi_dof_joint_trajectory_point_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_go_proto",
    ],
)

proto_library(
    name = "gripper_command_proto",
    srcs = ["gripper_command.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
        "//third_party/ros2/ros_interfaces/jazzy/control_msgs/msg:gripper_command_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_proto",
    ],
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
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/action/gripper_command_go_proto",
    protos = ["gripper_command_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/control_msgs/msg:gripper_command_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_go_proto",
    ],
)

proto_library(
    name = "joint_trajectory_proto",
    srcs = ["joint_trajectory.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:joint_trajectory_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_proto",
    ],
)

py_proto_library(
    name = "joint_trajectory_py_pb2",
    deps = ["joint_trajectory_proto"],
)

cc_proto_library(
    name = "joint_trajectory_cc_proto",
    deps = ["joint_trajectory_proto"],
)

go_proto_library(
    name = "joint_trajectory_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/action/joint_trajectory_go_proto",
    protos = ["joint_trajectory_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:joint_trajectory_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_go_proto",
    ],
)

proto_library(
    name = "parallel_gripper_command_proto",
    srcs = ["parallel_gripper_command.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:joint_state_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_proto",
    ],
)

py_proto_library(
    name = "parallel_gripper_command_py_pb2",
    deps = ["parallel_gripper_command_proto"],
)

cc_proto_library(
    name = "parallel_gripper_command_cc_proto",
    deps = ["parallel_gripper_command_proto"],
)

go_proto_library(
    name = "parallel_gripper_command_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/action/parallel_gripper_command_go_proto",
    protos = ["parallel_gripper_command_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:joint_state_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_go_proto",
    ],
)

proto_library(
    name = "point_head_proto",
    srcs = ["point_head.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_proto",
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_stamped_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_proto",
    ],
)

py_proto_library(
    name = "point_head_py_pb2",
    deps = ["point_head_proto"],
)

cc_proto_library(
    name = "point_head_cc_proto",
    deps = ["point_head_proto"],
)

go_proto_library(
    name = "point_head_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/action/point_head_go_proto",
    protos = ["point_head_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:point_stamped_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:vector3_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_go_proto",
    ],
)

proto_library(
    name = "single_joint_position_proto",
    srcs = ["single_joint_position.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_proto",
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_proto",
    ],
)

py_proto_library(
    name = "single_joint_position_py_pb2",
    deps = ["single_joint_position_proto"],
)

cc_proto_library(
    name = "single_joint_position_cc_proto",
    deps = ["single_joint_position_proto"],
)

go_proto_library(
    name = "single_joint_position_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/control_msgs/action/single_joint_position_go_proto",
    protos = ["single_joint_position_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_go_proto",
    ],
)
