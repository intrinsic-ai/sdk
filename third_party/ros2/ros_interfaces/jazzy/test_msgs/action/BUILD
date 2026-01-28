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
    name = "fibonacci_proto",
    srcs = ["fibonacci.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_proto",
    ],
)

py_proto_library(
    name = "fibonacci_py_pb2",
    deps = ["fibonacci_proto"],
)

cc_proto_library(
    name = "fibonacci_cc_proto",
    deps = ["fibonacci_proto"],
)

go_proto_library(
    name = "fibonacci_go_proto",
    protos = ["fibonacci_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_go_proto",
    ],
)

proto_library(
    name = "nested_message_proto",
    srcs = ["nested_message.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:basic_types_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:builtins_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_proto",
    ],
)

py_proto_library(
    name = "nested_message_py_pb2",
    deps = ["nested_message_proto"],
)

cc_proto_library(
    name = "nested_message_cc_proto",
    deps = ["nested_message_proto"],
)

go_proto_library(
    name = "nested_message_go_proto",
    protos = ["nested_message_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/builtin_interfaces/msg:time_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:basic_types_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/test_msgs/msg:builtins_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/unique_identifier_msgs/msg:uuid_go_proto",
    ],
)
