load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/diagnostic_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "add_diagnostics_proto",
    srcs = ["add_diagnostics.proto"],
)

py_proto_library(
    name = "add_diagnostics_py_pb2",
    deps = ["add_diagnostics_proto"],
)

cc_proto_library(
    name = "add_diagnostics_cc_proto",
    deps = ["add_diagnostics_proto"],
)

go_proto_library(
    name = "add_diagnostics_go_proto",
    deps = ["add_diagnostics_proto"],
)

proto_library(
    name = "self_test_proto",
    srcs = ["self_test.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/diagnostic_msgs/msg:diagnostic_status_proto"],
)

py_proto_library(
    name = "self_test_py_pb2",
    deps = ["self_test_proto"],
)

cc_proto_library(
    name = "self_test_cc_proto",
    deps = ["self_test_proto"],
)

go_proto_library(
    name = "self_test_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/diagnostic_msgs/msg:diagnostic_status_go_proto",
    ],
    deps = ["self_test_proto"],
)
