load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/sensor_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "set_camera_info_proto",
    srcs = ["set_camera_info.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:camera_info_proto"],
)

py_proto_library(
    name = "set_camera_info_py_pb2",
    deps = ["set_camera_info_proto"],
)

cc_proto_library(
    name = "set_camera_info_cc_proto",
    deps = ["set_camera_info_proto"],
)

go_proto_library(
    name = "set_camera_info_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/sensor_msgs/msg:camera_info_go_proto",
    ],
    deps = ["set_camera_info_proto"],
)
