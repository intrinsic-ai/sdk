load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "describe_parameters_proto",
    srcs = ["describe_parameters.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_descriptor_proto"],
)

py_proto_library(
    name = "describe_parameters_py_pb2",
    deps = ["describe_parameters_proto"],
)

cc_proto_library(
    name = "describe_parameters_cc_proto",
    deps = ["describe_parameters_proto"],
)

go_proto_library(
    name = "describe_parameters_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_descriptor_go_proto",
    ],
    deps = ["describe_parameters_proto"],
)

proto_library(
    name = "get_logger_levels_proto",
    srcs = ["get_logger_levels.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:logger_level_proto"],
)

py_proto_library(
    name = "get_logger_levels_py_pb2",
    deps = ["get_logger_levels_proto"],
)

cc_proto_library(
    name = "get_logger_levels_cc_proto",
    deps = ["get_logger_levels_proto"],
)

go_proto_library(
    name = "get_logger_levels_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:logger_level_go_proto",
    ],
    deps = ["get_logger_levels_proto"],
)

proto_library(
    name = "get_parameter_types_proto",
    srcs = ["get_parameter_types.proto"],
)

py_proto_library(
    name = "get_parameter_types_py_pb2",
    deps = ["get_parameter_types_proto"],
)

cc_proto_library(
    name = "get_parameter_types_cc_proto",
    deps = ["get_parameter_types_proto"],
)

go_proto_library(
    name = "get_parameter_types_go_proto",
    deps = ["get_parameter_types_proto"],
)

proto_library(
    name = "get_parameters_proto",
    srcs = ["get_parameters.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_value_proto"],
)

py_proto_library(
    name = "get_parameters_py_pb2",
    deps = ["get_parameters_proto"],
)

cc_proto_library(
    name = "get_parameters_cc_proto",
    deps = ["get_parameters_proto"],
)

go_proto_library(
    name = "get_parameters_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_value_go_proto",
    ],
    deps = ["get_parameters_proto"],
)

proto_library(
    name = "list_parameters_proto",
    srcs = ["list_parameters.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:list_parameters_result_proto"],
)

py_proto_library(
    name = "list_parameters_py_pb2",
    deps = ["list_parameters_proto"],
)

cc_proto_library(
    name = "list_parameters_cc_proto",
    deps = ["list_parameters_proto"],
)

go_proto_library(
    name = "list_parameters_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:list_parameters_result_go_proto",
    ],
    deps = ["list_parameters_proto"],
)

proto_library(
    name = "set_logger_levels_proto",
    srcs = ["set_logger_levels.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:logger_level_proto",
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:set_logger_levels_result_proto",
    ],
)

py_proto_library(
    name = "set_logger_levels_py_pb2",
    deps = ["set_logger_levels_proto"],
)

cc_proto_library(
    name = "set_logger_levels_cc_proto",
    deps = ["set_logger_levels_proto"],
)

go_proto_library(
    name = "set_logger_levels_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:logger_level_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:set_logger_levels_result_go_proto",
    ],
    deps = ["set_logger_levels_proto"],
)

proto_library(
    name = "set_parameters_proto",
    srcs = ["set_parameters.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_proto",
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:set_parameters_result_proto",
    ],
)

py_proto_library(
    name = "set_parameters_py_pb2",
    deps = ["set_parameters_proto"],
)

cc_proto_library(
    name = "set_parameters_cc_proto",
    deps = ["set_parameters_proto"],
)

go_proto_library(
    name = "set_parameters_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:set_parameters_result_go_proto",
    ],
    deps = ["set_parameters_proto"],
)

proto_library(
    name = "set_parameters_atomically_proto",
    srcs = ["set_parameters_atomically.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_proto",
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:set_parameters_result_proto",
    ],
)

py_proto_library(
    name = "set_parameters_atomically_py_pb2",
    deps = ["set_parameters_atomically_proto"],
)

cc_proto_library(
    name = "set_parameters_atomically_cc_proto",
    deps = ["set_parameters_atomically_proto"],
)

go_proto_library(
    name = "set_parameters_atomically_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:set_parameters_result_go_proto",
    ],
    deps = ["set_parameters_atomically_proto"],
)
