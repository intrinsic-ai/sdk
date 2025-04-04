load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "uuid_proto",
    srcs = ["uuid.proto"],
)

py_proto_library(
    name = "uuid_py_pb2",
    deps = ["uuid_proto"],
)

cc_proto_library(
    name = "uuid_cc_proto",
    deps = ["uuid_proto"],
)

go_proto_library(
    name = "uuid_go_proto",
    deps = ["uuid_proto"],
)
