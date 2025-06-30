load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/tf2_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "tf2_error_proto",
    srcs = ["tf2_error.proto"],
)

py_proto_library(
    name = "tf2_error_py_pb2",
    deps = ["tf2_error_proto"],
)

cc_proto_library(
    name = "tf2_error_cc_proto",
    deps = ["tf2_error_proto"],
)

go_proto_library(
    name = "tf2_error_go_proto",
    deps = ["tf2_error_proto"],
)

proto_library(
    name = "tf_message_proto",
    srcs = ["tf_message.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:transform_stamped_proto"],
)

py_proto_library(
    name = "tf_message_py_pb2",
    deps = ["tf_message_proto"],
)

cc_proto_library(
    name = "tf_message_cc_proto",
    deps = ["tf_message_proto"],
)

go_proto_library(
    name = "tf_message_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:transform_stamped_go_proto",
    ],
    deps = ["tf_message_proto"],
)
