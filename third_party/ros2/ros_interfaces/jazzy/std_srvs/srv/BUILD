load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/std_srvs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
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
    protos = ["empty_proto"],
)

proto_library(
    name = "set_bool_proto",
    srcs = ["set_bool.proto"],
)

py_proto_library(
    name = "set_bool_py_pb2",
    deps = ["set_bool_proto"],
)

cc_proto_library(
    name = "set_bool_cc_proto",
    deps = ["set_bool_proto"],
)

go_proto_library(
    name = "set_bool_go_proto",
    protos = ["set_bool_proto"],
)

proto_library(
    name = "trigger_proto",
    srcs = ["trigger.proto"],
)

py_proto_library(
    name = "trigger_py_pb2",
    deps = ["trigger_proto"],
)

cc_proto_library(
    name = "trigger_cc_proto",
    deps = ["trigger_proto"],
)

go_proto_library(
    name = "trigger_go_proto",
    protos = ["trigger_proto"],
)
