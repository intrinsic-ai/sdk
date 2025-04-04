load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/pcl_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "update_filename_proto",
    srcs = ["update_filename.proto"],
)

py_proto_library(
    name = "update_filename_py_pb2",
    deps = ["update_filename_proto"],
)

cc_proto_library(
    name = "update_filename_cc_proto",
    deps = ["update_filename_proto"],
)

go_proto_library(
    name = "update_filename_go_proto",
    deps = ["update_filename_proto"],
)
