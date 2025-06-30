load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/statistics_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "metrics_message_proto",
    srcs = ["metrics_message.proto"],
    deps = [
        ":statistic_data_point_proto",
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
    ],
)

py_proto_library(
    name = "metrics_message_py_pb2",
    deps = ["metrics_message_proto"],
)

cc_proto_library(
    name = "metrics_message_cc_proto",
    deps = ["metrics_message_proto"],
)

go_proto_library(
    name = "metrics_message_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/statistics_msgs/msg:statistic_data_point_go_proto",
    ],
    deps = ["metrics_message_proto"],
)

proto_library(
    name = "statistic_data_point_proto",
    srcs = ["statistic_data_point.proto"],
)

py_proto_library(
    name = "statistic_data_point_py_pb2",
    deps = ["statistic_data_point_proto"],
)

cc_proto_library(
    name = "statistic_data_point_cc_proto",
    deps = ["statistic_data_point_proto"],
)

go_proto_library(
    name = "statistic_data_point_go_proto",
    deps = ["statistic_data_point_proto"],
)

proto_library(
    name = "statistic_data_type_proto",
    srcs = ["statistic_data_type.proto"],
)

py_proto_library(
    name = "statistic_data_type_py_pb2",
    deps = ["statistic_data_type_proto"],
)

cc_proto_library(
    name = "statistic_data_type_cc_proto",
    deps = ["statistic_data_type_proto"],
)

go_proto_library(
    name = "statistic_data_type_go_proto",
    deps = ["statistic_data_type_proto"],
)
