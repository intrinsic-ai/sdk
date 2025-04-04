load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "floating_point_range_proto",
    srcs = ["floating_point_range.proto"],
)

py_proto_library(
    name = "floating_point_range_py_pb2",
    deps = ["floating_point_range_proto"],
)

cc_proto_library(
    name = "floating_point_range_cc_proto",
    deps = ["floating_point_range_proto"],
)

go_proto_library(
    name = "floating_point_range_go_proto",
    deps = ["floating_point_range_proto"],
)

proto_library(
    name = "integer_range_proto",
    srcs = ["integer_range.proto"],
)

py_proto_library(
    name = "integer_range_py_pb2",
    deps = ["integer_range_proto"],
)

cc_proto_library(
    name = "integer_range_cc_proto",
    deps = ["integer_range_proto"],
)

go_proto_library(
    name = "integer_range_go_proto",
    deps = ["integer_range_proto"],
)

proto_library(
    name = "list_parameters_result_proto",
    srcs = ["list_parameters_result.proto"],
)

py_proto_library(
    name = "list_parameters_result_py_pb2",
    deps = ["list_parameters_result_proto"],
)

cc_proto_library(
    name = "list_parameters_result_cc_proto",
    deps = ["list_parameters_result_proto"],
)

go_proto_library(
    name = "list_parameters_result_go_proto",
    deps = ["list_parameters_result_proto"],
)

proto_library(
    name = "log_proto",
    srcs = ["log.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto"],
)

py_proto_library(
    name = "log_py_pb2",
    deps = ["log_proto"],
)

cc_proto_library(
    name = "log_cc_proto",
    deps = ["log_proto"],
)

go_proto_library(
    name = "log_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
    ],
    deps = ["log_proto"],
)

proto_library(
    name = "logger_level_proto",
    srcs = ["logger_level.proto"],
)

py_proto_library(
    name = "logger_level_py_pb2",
    deps = ["logger_level_proto"],
)

cc_proto_library(
    name = "logger_level_cc_proto",
    deps = ["logger_level_proto"],
)

go_proto_library(
    name = "logger_level_go_proto",
    deps = ["logger_level_proto"],
)

proto_library(
    name = "parameter_proto",
    srcs = ["parameter.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_value_proto"],
)

py_proto_library(
    name = "parameter_py_pb2",
    deps = ["parameter_proto"],
)

cc_proto_library(
    name = "parameter_cc_proto",
    deps = ["parameter_proto"],
)

go_proto_library(
    name = "parameter_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_value_go_proto",
    ],
    deps = ["parameter_proto"],
)

proto_library(
    name = "parameter_descriptor_proto",
    srcs = ["parameter_descriptor.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:floating_point_range_proto",
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:integer_range_proto",
    ],
)

py_proto_library(
    name = "parameter_descriptor_py_pb2",
    deps = ["parameter_descriptor_proto"],
)

cc_proto_library(
    name = "parameter_descriptor_cc_proto",
    deps = ["parameter_descriptor_proto"],
)

go_proto_library(
    name = "parameter_descriptor_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:floating_point_range_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:integer_range_go_proto",
    ],
    deps = ["parameter_descriptor_proto"],
)

proto_library(
    name = "parameter_event_proto",
    srcs = ["parameter_event.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_proto",
    ],
)

py_proto_library(
    name = "parameter_event_py_pb2",
    deps = ["parameter_event_proto"],
)

cc_proto_library(
    name = "parameter_event_cc_proto",
    deps = ["parameter_event_proto"],
)

go_proto_library(
    name = "parameter_event_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_go_proto",
    ],
    deps = ["parameter_event_proto"],
)

proto_library(
    name = "parameter_event_descriptors_proto",
    srcs = ["parameter_event_descriptors.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_descriptor_proto"],
)

py_proto_library(
    name = "parameter_event_descriptors_py_pb2",
    deps = ["parameter_event_descriptors_proto"],
)

cc_proto_library(
    name = "parameter_event_descriptors_cc_proto",
    deps = ["parameter_event_descriptors_proto"],
)

go_proto_library(
    name = "parameter_event_descriptors_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/rcl_interfaces/msg:parameter_descriptor_go_proto",
    ],
    deps = ["parameter_event_descriptors_proto"],
)

proto_library(
    name = "parameter_type_proto",
    srcs = ["parameter_type.proto"],
)

py_proto_library(
    name = "parameter_type_py_pb2",
    deps = ["parameter_type_proto"],
)

cc_proto_library(
    name = "parameter_type_cc_proto",
    deps = ["parameter_type_proto"],
)

go_proto_library(
    name = "parameter_type_go_proto",
    deps = ["parameter_type_proto"],
)

proto_library(
    name = "parameter_value_proto",
    srcs = ["parameter_value.proto"],
)

py_proto_library(
    name = "parameter_value_py_pb2",
    deps = ["parameter_value_proto"],
)

cc_proto_library(
    name = "parameter_value_cc_proto",
    deps = ["parameter_value_proto"],
)

go_proto_library(
    name = "parameter_value_go_proto",
    deps = ["parameter_value_proto"],
)

proto_library(
    name = "set_logger_levels_result_proto",
    srcs = ["set_logger_levels_result.proto"],
)

py_proto_library(
    name = "set_logger_levels_result_py_pb2",
    deps = ["set_logger_levels_result_proto"],
)

cc_proto_library(
    name = "set_logger_levels_result_cc_proto",
    deps = ["set_logger_levels_result_proto"],
)

go_proto_library(
    name = "set_logger_levels_result_go_proto",
    deps = ["set_logger_levels_result_proto"],
)

proto_library(
    name = "set_parameters_result_proto",
    srcs = ["set_parameters_result.proto"],
)

py_proto_library(
    name = "set_parameters_result_py_pb2",
    deps = ["set_parameters_result_proto"],
)

cc_proto_library(
    name = "set_parameters_result_cc_proto",
    deps = ["set_parameters_result_proto"],
)

go_proto_library(
    name = "set_parameters_result_go_proto",
    deps = ["set_parameters_result_proto"],
)
