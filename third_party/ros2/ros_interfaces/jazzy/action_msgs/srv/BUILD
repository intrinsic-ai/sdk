load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/action_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "cancel_goal_proto",
    srcs = ["cancel_goal.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/action_msgs/msg:goal_info_proto"],
)

py_proto_library(
    name = "cancel_goal_py_pb2",
    deps = ["cancel_goal_proto"],
)

cc_proto_library(
    name = "cancel_goal_cc_proto",
    deps = ["cancel_goal_proto"],
)

go_proto_library(
    name = "cancel_goal_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/action_msgs/srv/cancel_goal_go_proto",
    protos = ["cancel_goal_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/action_msgs/msg:goal_info_go_proto",
    ],
)
