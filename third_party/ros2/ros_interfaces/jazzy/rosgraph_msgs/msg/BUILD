load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/rosgraph_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "clock_proto",
    srcs = ["clock.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto"],
)

py_proto_library(
    name = "clock_py_pb2",
    deps = ["clock_proto"],
)

cc_proto_library(
    name = "clock_cc_proto",
    deps = ["clock_proto"],
)

go_proto_library(
    name = "clock_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
    ],
    deps = ["clock_proto"],
)
