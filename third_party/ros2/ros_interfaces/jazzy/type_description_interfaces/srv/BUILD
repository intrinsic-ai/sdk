load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(
    default_applicable_licenses = ["//third_party/ros2/ros_interfaces/jazzy/type_description_interfaces:license"],
    default_visibility = ["//visibility:public"],
    licenses = ["notice"],
)

proto_library(
    name = "get_type_description_proto",
    srcs = ["get_type_description.proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/type_description_interfaces/msg:key_value_proto",
        "//third_party/ros2/ros_interfaces/jazzy/type_description_interfaces/msg:type_description_proto",
        "//third_party/ros2/ros_interfaces/jazzy/type_description_interfaces/msg:type_source_proto",
    ],
)

py_proto_library(
    name = "get_type_description_py_pb2",
    deps = ["get_type_description_proto"],
)

cc_proto_library(
    name = "get_type_description_cc_proto",
    deps = ["get_type_description_proto"],
)

go_proto_library(
    name = "get_type_description_go_proto",
    importpath = "third_party/ros2/ros_interfaces/jazzy/type_description_interfaces/srv/get_type_description_go_proto",
    protos = ["get_type_description_proto"],
    deps = [
        "//third_party/ros2/ros_interfaces/jazzy/type_description_interfaces/msg:key_value_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/type_description_interfaces/msg:type_description_go_proto",
        "//third_party/ros2/ros_interfaces/jazzy/type_description_interfaces/msg:type_source_go_proto",
    ],
)
