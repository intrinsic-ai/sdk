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
    name = "diagnostic_array_proto",
    srcs = ["diagnostic_array.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/diagnostic_msgs/msg:diagnostic_status_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_proto",
    ],
)

py_proto_library(
    name = "diagnostic_array_py_pb2",
    deps = ["diagnostic_array_proto"],
)

cc_proto_library(
    name = "diagnostic_array_cc_proto",
    deps = ["diagnostic_array_proto"],
)

go_proto_library(
    name = "diagnostic_array_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/diagnostic_msgs/msg:diagnostic_status_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/std_msgs/msg:header_go_proto",
    ],
    deps = ["diagnostic_array_proto"],
)

proto_library(
    name = "diagnostic_status_proto",
    srcs = ["diagnostic_status.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/diagnostic_msgs/msg:key_value_proto"],
)

py_proto_library(
    name = "diagnostic_status_py_pb2",
    deps = ["diagnostic_status_proto"],
)

cc_proto_library(
    name = "diagnostic_status_cc_proto",
    deps = ["diagnostic_status_proto"],
)

go_proto_library(
    name = "diagnostic_status_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/diagnostic_msgs/msg:key_value_go_proto",
    ],
    deps = ["diagnostic_status_proto"],
)

proto_library(
    name = "key_value_proto",
    srcs = ["key_value.proto"],
)

py_proto_library(
    name = "key_value_py_pb2",
    deps = ["key_value_proto"],
)

cc_proto_library(
    name = "key_value_cc_proto",
    deps = ["key_value_proto"],
)

go_proto_library(
    name = "key_value_go_proto",
    deps = ["key_value_proto"],
)
