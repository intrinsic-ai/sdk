load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/action_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "goal_info_proto",
    srcs = ["goal_info.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_proto",
    ],
)

py_proto_library(
    name = "goal_info_py_pb2",
    deps = ["goal_info_proto"],
)

cc_proto_library(
    name = "goal_info_cc_proto",
    deps = ["goal_info_proto"],
)

go_proto_library(
    name = "goal_info_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_go_proto",
    ],
    deps = ["goal_info_proto"],
)

proto_library(
    name = "goal_status_proto",
    srcs = ["goal_status.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/action_msgs/msg:goal_info_proto"],
)

py_proto_library(
    name = "goal_status_py_pb2",
    deps = ["goal_status_proto"],
)

cc_proto_library(
    name = "goal_status_cc_proto",
    deps = ["goal_status_proto"],
)

go_proto_library(
    name = "goal_status_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/action_msgs/msg:goal_info_go_proto",
    ],
    deps = ["goal_status_proto"],
)

proto_library(
    name = "goal_status_array_proto",
    srcs = ["goal_status_array.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/action_msgs/msg:goal_status_proto"],
)

py_proto_library(
    name = "goal_status_array_py_pb2",
    deps = ["goal_status_array_proto"],
)

cc_proto_library(
    name = "goal_status_array_cc_proto",
    deps = ["goal_status_array_proto"],
)

go_proto_library(
    name = "goal_status_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/action_msgs/msg:goal_status_go_proto",
    ],
    deps = ["goal_status_array_proto"],
)
