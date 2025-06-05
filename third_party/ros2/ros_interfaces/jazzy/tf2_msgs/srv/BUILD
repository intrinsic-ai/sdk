load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/tf2_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "frame_graph_proto",
    srcs = ["frame_graph.proto"],
)

py_proto_library(
    name = "frame_graph_py_pb2",
    deps = ["frame_graph_proto"],
)

cc_proto_library(
    name = "frame_graph_cc_proto",
    deps = ["frame_graph_proto"],
)

go_proto_library(
    name = "frame_graph_go_proto",
    deps = ["frame_graph_proto"],
)
