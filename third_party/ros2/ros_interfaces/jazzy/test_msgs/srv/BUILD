load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/test_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "arrays_proto",
    srcs = ["arrays.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:basic_types_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:constants_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:defaults_proto",
    ],
)

py_proto_library(
    name = "arrays_py_pb2",
    deps = ["arrays_proto"],
)

cc_proto_library(
    name = "arrays_cc_proto",
    deps = ["arrays_proto"],
)

go_proto_library(
    name = "arrays_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:basic_types_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:constants_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:defaults_go_proto",
    ],
    deps = ["arrays_proto"],
)

proto_library(
    name = "basic_types_proto",
    srcs = ["basic_types.proto"],
)

py_proto_library(
    name = "basic_types_py_pb2",
    deps = ["basic_types_proto"],
)

cc_proto_library(
    name = "basic_types_cc_proto",
    deps = ["basic_types_proto"],
)

go_proto_library(
    name = "basic_types_go_proto",
    deps = ["basic_types_proto"],
)

proto_library(
    name = "empty_proto",
    srcs = ["empty.proto"],
)

py_proto_library(
    name = "empty_py_pb2",
    deps = ["empty_proto"],
)

cc_proto_library(
    name = "empty_cc_proto",
    deps = ["empty_proto"],
)

go_proto_library(
    name = "empty_go_proto",
    deps = ["empty_proto"],
)
