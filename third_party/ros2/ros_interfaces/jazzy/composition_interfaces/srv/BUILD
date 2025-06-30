load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/composition_interfaces:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "list_nodes_proto",
    srcs = ["list_nodes.proto"],
)

py_proto_library(
    name = "list_nodes_py_pb2",
    deps = ["list_nodes_proto"],
)

cc_proto_library(
    name = "list_nodes_cc_proto",
    deps = ["list_nodes_proto"],
)

go_proto_library(
    name = "list_nodes_go_proto",
    deps = ["list_nodes_proto"],
)

proto_library(
    name = "load_node_proto",
    srcs = ["load_node.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_proto"],
)

py_proto_library(
    name = "load_node_py_pb2",
    deps = ["load_node_proto"],
)

cc_proto_library(
    name = "load_node_cc_proto",
    deps = ["load_node_proto"],
)

go_proto_library(
    name = "load_node_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_go_proto",
    ],
    deps = ["load_node_proto"],
)

proto_library(
    name = "unload_node_proto",
    srcs = ["unload_node.proto"],
)

py_proto_library(
    name = "unload_node_py_pb2",
    deps = ["unload_node_proto"],
)

cc_proto_library(
    name = "unload_node_cc_proto",
    deps = ["unload_node_proto"],
)

go_proto_library(
    name = "unload_node_go_proto",
    deps = ["unload_node_proto"],
)
