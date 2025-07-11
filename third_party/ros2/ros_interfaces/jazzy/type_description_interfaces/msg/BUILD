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
    name = "field_proto",
    srcs = ["field.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/type_description_interfaces/msg:field_type_proto"],
)

py_proto_library(
    name = "field_py_pb2",
    deps = ["field_proto"],
)

cc_proto_library(
    name = "field_cc_proto",
    deps = ["field_proto"],
)

go_proto_library(
    name = "field_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/type_description_interfaces/msg:field_type_go_proto",
    ],
    deps = ["field_proto"],
)

proto_library(
    name = "field_type_proto",
    srcs = ["field_type.proto"],
)

py_proto_library(
    name = "field_type_py_pb2",
    deps = ["field_type_proto"],
)

cc_proto_library(
    name = "field_type_cc_proto",
    deps = ["field_type_proto"],
)

go_proto_library(
    name = "field_type_go_proto",
    deps = ["field_type_proto"],
)

proto_library(
    name = "individual_type_description_proto",
    srcs = ["individual_type_description.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/type_description_interfaces/msg:field_proto"],
)

py_proto_library(
    name = "individual_type_description_py_pb2",
    deps = ["individual_type_description_proto"],
)

cc_proto_library(
    name = "individual_type_description_cc_proto",
    deps = ["individual_type_description_proto"],
)

go_proto_library(
    name = "individual_type_description_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/type_description_interfaces/msg:field_go_proto",
    ],
    deps = ["individual_type_description_proto"],
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

proto_library(
    name = "type_description_proto",
    srcs = ["type_description.proto"],
    deps = ["//third_party/ros2/ros_interfaces/jazzy/type_description_interfaces/msg:individual_type_description_proto"],
)

py_proto_library(
    name = "type_description_py_pb2",
    deps = ["type_description_proto"],
)

cc_proto_library(
    name = "type_description_cc_proto",
    deps = ["type_description_proto"],
)

go_proto_library(
    name = "type_description_go_proto",
    go_deps = [
        "//third_party/ros2/ros_interfaces/jazzy/type_description_interfaces/msg:individual_type_description_go_proto",
    ],
    deps = ["type_description_proto"],
)

proto_library(
    name = "type_source_proto",
    srcs = ["type_source.proto"],
)

py_proto_library(
    name = "type_source_py_pb2",
    deps = ["type_source_proto"],
)

cc_proto_library(
    name = "type_source_cc_proto",
    deps = ["type_source_proto"],
)

go_proto_library(
    name = "type_source_go_proto",
    deps = ["type_source_proto"],
)
