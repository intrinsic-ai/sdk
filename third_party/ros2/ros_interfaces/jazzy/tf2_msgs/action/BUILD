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
    name = "lookup_transform_proto",
    srcs = ["lookup_transform.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_proto",
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:transform_stamped_proto",
        "//third_party/ros2/ros_interfaces/jazzy/tf2_msgs/msg:tf2_error_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_proto",
    ],
)

py_proto_library(
    name = "lookup_transform_py_pb2",
    deps = ["lookup_transform_proto"],
)

cc_proto_library(
    name = "lookup_transform_cc_proto",
    deps = ["lookup_transform_proto"],
)

go_proto_library(
    name = "lookup_transform_go_proto",
    protos = ["lookup_transform_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/geometry_msgs/msg:transform_stamped_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/tf2_msgs/msg:tf2_error_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_go_proto",
    ],
)
