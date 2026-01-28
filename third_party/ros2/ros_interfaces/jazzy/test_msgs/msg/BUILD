load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/test_msgs:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "arrays_proto",
    srcs = ["arrays.proto"],
    deps = [
        ":basic_types_proto",
        ":constants_proto",
        ":defaults_proto",
    ],
)

py_proto_library(
    name = "arrays_py_pb2",
    deps = ["arrays_proto"],
)

cc_proto_library(
    name = "arrays_cc_proto",
    deps = ["arrays_proto"],
)

go_proto_library(
    name = "arrays_go_proto",
    protos = ["arrays_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:basic_types_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:constants_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:defaults_go_proto",
    ],
)

proto_library(
    name = "basic_types_proto",
    srcs = ["basic_types.proto"],
)

py_proto_library(
    name = "basic_types_py_pb2",
    deps = ["basic_types_proto"],
)

cc_proto_library(
    name = "basic_types_cc_proto",
    deps = ["basic_types_proto"],
)

go_proto_library(
    name = "basic_types_go_proto",
    protos = ["basic_types_proto"],
)

proto_library(
    name = "bounded_plain_sequences_proto",
    srcs = ["bounded_plain_sequences.proto"],
    deps = [
        ":basic_types_proto",
        ":constants_proto",
        ":defaults_proto",
    ],
)

py_proto_library(
    name = "bounded_plain_sequences_py_pb2",
    deps = ["bounded_plain_sequences_proto"],
)

cc_proto_library(
    name = "bounded_plain_sequences_cc_proto",
    deps = ["bounded_plain_sequences_proto"],
)

go_proto_library(
    name = "bounded_plain_sequences_go_proto",
    protos = ["bounded_plain_sequences_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:basic_types_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:constants_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:defaults_go_proto",
    ],
)

proto_library(
    name = "bounded_sequences_proto",
    srcs = ["bounded_sequences.proto"],
    deps = [
        ":basic_types_proto",
        ":constants_proto",
        ":defaults_proto",
    ],
)

py_proto_library(
    name = "bounded_sequences_py_pb2",
    deps = ["bounded_sequences_proto"],
)

cc_proto_library(
    name = "bounded_sequences_cc_proto",
    deps = ["bounded_sequences_proto"],
)

go_proto_library(
    name = "bounded_sequences_go_proto",
    protos = ["bounded_sequences_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:basic_types_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:constants_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:defaults_go_proto",
    ],
)

proto_library(
    name = "builtins_proto",
    srcs = ["builtins.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_proto",
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
    ],
)

py_proto_library(
    name = "builtins_py_pb2",
    deps = ["builtins_proto"],
)

cc_proto_library(
    name = "builtins_cc_proto",
    deps = ["builtins_proto"],
)

go_proto_library(
    name = "builtins_go_proto",
    protos = ["builtins_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:duration_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
    ],
)

proto_library(
    name = "constants_proto",
    srcs = ["constants.proto"],
)

py_proto_library(
    name = "constants_py_pb2",
    deps = ["constants_proto"],
)

cc_proto_library(
    name = "constants_cc_proto",
    deps = ["constants_proto"],
)

go_proto_library(
    name = "constants_go_proto",
    protos = ["constants_proto"],
)

proto_library(
    name = "defaults_proto",
    srcs = ["defaults.proto"],
)

py_proto_library(
    name = "defaults_py_pb2",
    deps = ["defaults_proto"],
)

cc_proto_library(
    name = "defaults_cc_proto",
    deps = ["defaults_proto"],
)

go_proto_library(
    name = "defaults_go_proto",
    protos = ["defaults_proto"],
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
    name = "multi_nested_proto",
    srcs = ["multi_nested.proto"],
    deps = [
        ":arrays_proto",
        ":bounded_sequences_proto",
        ":unbounded_sequences_proto",
    ],
)

py_proto_library(
    name = "multi_nested_py_pb2",
    deps = ["multi_nested_proto"],
)

cc_proto_library(
    name = "multi_nested_cc_proto",
    deps = ["multi_nested_proto"],
)

go_proto_library(
    name = "multi_nested_go_proto",
    protos = ["multi_nested_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:arrays_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:bounded_sequences_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:unbounded_sequences_go_proto",
    ],
)

proto_library(
    name = "nested_proto",
    srcs = ["nested.proto"],
    deps = [":basic_types_proto"],
)

py_proto_library(
    name = "nested_py_pb2",
    deps = ["nested_proto"],
)

cc_proto_library(
    name = "nested_cc_proto",
    deps = ["nested_proto"],
)

go_proto_library(
    name = "nested_go_proto",
    protos = ["nested_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:basic_types_go_proto",
    ],
)

proto_library(
    name = "strings_proto",
    srcs = ["strings.proto"],
)

py_proto_library(
    name = "strings_py_pb2",
    deps = ["strings_proto"],
)

cc_proto_library(
    name = "strings_cc_proto",
    deps = ["strings_proto"],
)

go_proto_library(
    name = "strings_go_proto",
    protos = ["strings_proto"],
)

proto_library(
    name = "unbounded_sequences_proto",
    srcs = ["unbounded_sequences.proto"],
    deps = [
        ":basic_types_proto",
        ":constants_proto",
        ":defaults_proto",
    ],
)

py_proto_library(
    name = "unbounded_sequences_py_pb2",
    deps = ["unbounded_sequences_proto"],
)

cc_proto_library(
    name = "unbounded_sequences_cc_proto",
    deps = ["unbounded_sequences_proto"],
)

go_proto_library(
    name = "unbounded_sequences_go_proto",
    protos = ["unbounded_sequences_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:basic_types_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:constants_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:defaults_go_proto",
    ],
)

proto_library(
    name = "w_strings_proto",
    srcs = ["w_strings.proto"],
)

py_proto_library(
    name = "w_strings_py_pb2",
    deps = ["w_strings_proto"],
)

cc_proto_library(
    name = "w_strings_cc_proto",
    deps = ["w_strings_proto"],
)

go_proto_library(
    name = "w_strings_go_proto",
    protos = ["w_strings_proto"],
)
