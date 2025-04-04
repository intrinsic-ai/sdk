load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "joint_trajectory_proto",
    srcs = ["joint_trajectory.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:joint_trajectory_point_proto",
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
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:joint_trajectory_point_go_proto",
    ],
    deps = ["joint_trajectory_proto"],
)

proto_library(
    name = "joint_trajectory_point_proto",
    srcs = ["joint_trajectory_point.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_proto"],
)

py_proto_library(
    name = "joint_trajectory_point_py_pb2",
    deps = ["joint_trajectory_point_proto"],
)

cc_proto_library(
    name = "joint_trajectory_point_cc_proto",
    deps = ["joint_trajectory_point_proto"],
)

go_proto_library(
    name = "joint_trajectory_point_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_go_proto",
    ],
    deps = ["joint_trajectory_point_proto"],
)

proto_library(
    name = "multi_dof_joint_trajectory_proto",
    srcs = ["multi_dof_joint_trajectory.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:multi_dof_joint_trajectory_point_proto",
    ],
)

py_proto_library(
    name = "multi_dof_joint_trajectory_py_pb2",
    deps = ["multi_dof_joint_trajectory_proto"],
)

cc_proto_library(
    name = "multi_dof_joint_trajectory_cc_proto",
    deps = ["multi_dof_joint_trajectory_proto"],
)

go_proto_library(
    name = "multi_dof_joint_trajectory_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/trajectory_msgs/msg:multi_dof_joint_trajectory_point_go_proto",
    ],
    deps = ["multi_dof_joint_trajectory_proto"],
)

proto_library(
    name = "multi_dof_joint_trajectory_point_proto",
    srcs = ["multi_dof_joint_trajectory_point.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:transform_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_proto",
    ],
)

py_proto_library(
    name = "multi_dof_joint_trajectory_point_py_pb2",
    deps = ["multi_dof_joint_trajectory_point_proto"],
)

cc_proto_library(
    name = "multi_dof_joint_trajectory_point_cc_proto",
    deps = ["multi_dof_joint_trajectory_point_proto"],
)

go_proto_library(
    name = "multi_dof_joint_trajectory_point_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:transform_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:twist_go_proto",
    ],
    deps = ["multi_dof_joint_trajectory_point_proto"],
)
