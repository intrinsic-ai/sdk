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
    name = "query_calibration_state_proto",
    srcs = ["query_calibration_state.proto"],
)

py_proto_library(
    name = "query_calibration_state_py_pb2",
    deps = ["query_calibration_state_proto"],
)

cc_proto_library(
    name = "query_calibration_state_cc_proto",
    deps = ["query_calibration_state_proto"],
)

go_proto_library(
    name = "query_calibration_state_go_proto",
    deps = ["query_calibration_state_proto"],
)

proto_library(
    name = "query_trajectory_state_proto",
    srcs = ["query_trajectory_state.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto"],
)

py_proto_library(
    name = "query_trajectory_state_py_pb2",
    deps = ["query_trajectory_state_proto"],
)

cc_proto_library(
    name = "query_trajectory_state_cc_proto",
    deps = ["query_trajectory_state_proto"],
)

go_proto_library(
    name = "query_trajectory_state_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
    ],
    deps = ["query_trajectory_state_proto"],
)
