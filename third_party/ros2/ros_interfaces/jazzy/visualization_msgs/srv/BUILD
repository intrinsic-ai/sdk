load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/visualization_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "get_interactive_markers_proto",
    srcs = ["get_interactive_markers.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/visualization_msgs/msg:interactive_marker_proto"],
)

py_proto_library(
    name = "get_interactive_markers_py_pb2",
    deps = ["get_interactive_markers_proto"],
)

cc_proto_library(
    name = "get_interactive_markers_cc_proto",
    deps = ["get_interactive_markers_proto"],
)

go_proto_library(
    name = "get_interactive_markers_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/visualization_msgs/msg:interactive_marker_go_proto",
    ],
    deps = ["get_interactive_markers_proto"],
)
