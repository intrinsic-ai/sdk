# Copyright 2023 Intrinsic Innovation LLC

# Basic math protos (e.g. matrices, vectors, and a pose)

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "array_proto",
    srcs = ["array.proto"],
)

cc_proto_library(
    name = "array_cc_proto",
    deps = [":array_proto"],
)

py_proto_library(
    name = "array_py_pb2",
    deps = [":array_proto"],
)

go_proto_library(
    name = "array_go_proto",
    deps = [":array_proto"],
)

proto_library(
    name = "point_proto",
    srcs = ["point.proto"],
)

cc_proto_library(
    name = "point_cc_proto",
    deps = [":point_proto"],
)

py_proto_library(
    name = "point_py_pb2",
    deps = [":point_proto"],
)

go_proto_library(
    name = "point_go_proto",
    deps = [":point_proto"],
)

proto_library(
    name = "quaternion_proto",
    srcs = ["quaternion.proto"],
)

cc_proto_library(
    name = "quaternion_cc_proto",
    deps = [":quaternion_proto"],
)

py_proto_library(
    name = "quaternion_py_pb2",
    deps = [":quaternion_proto"],
)

go_proto_library(
    name = "quaternion_go_proto",
    deps = [":quaternion_proto"],
)

proto_library(
    name = "pose_proto",
    srcs = ["pose.proto"],
    deps = [
        ":point_proto",
        ":quaternion_proto",
    ],
)

cc_proto_library(
    name = "pose_cc_proto",
    deps = [":pose_proto"],
)

py_proto_library(
    name = "pose_py_pb2",
    deps = [":pose_proto"],
)

go_proto_library(
    name = "pose_go_proto",
    go_deps = [
        ":point_go_proto",
        ":quaternion_go_proto",
    ],
    deps = [":pose_proto"],
)

proto_library(
    name = "affine_proto",
    srcs = ["affine.proto"],
    deps = [
        ":matrix_proto",
        ":point_proto",
    ],
)

cc_proto_library(
    name = "affine_cc_proto",
    deps = [":affine_proto"],
)

py_proto_library(
    name = "affine_py_pb2",
    deps = [":affine_proto"],
)

go_proto_library(
    name = "affine_go_proto",
    go_deps = [
        ":matrix_go_proto",
        ":point_go_proto",
    ],
    deps = [":affine_proto"],
)

proto_library(
    name = "matrix_proto",
    srcs = ["matrix.proto"],
)

cc_proto_library(
    name = "matrix_cc_proto",
    deps = [":matrix_proto"],
)

py_proto_library(
    name = "matrix_py_pb2",
    deps = [":matrix_proto"],
)

go_proto_library(
    name = "matrix_go_proto",
    deps = [":matrix_proto"],
)

proto_library(
    name = "accel_proto",
    srcs = ["accel.proto"],
    deps = [
        ":vector3_proto",
    ],
)

py_proto_library(
    name = "accel_py_pb2",
    deps = [":accel_proto"],
)

go_proto_library(
    name = "accel_go_proto",
    go_deps = [
        ":vector3_go_proto",
    ],
    deps = [":accel_proto"],
)

cc_proto_library(
    name = "accel_cc_proto",
    deps = [":accel_proto"],
)

proto_library(
    name = "twist_proto",
    srcs = ["twist.proto"],
    deps = [
        ":vector3_proto",
    ],
)

py_proto_library(
    name = "twist_py_pb2",
    deps = [":twist_proto"],
)

go_proto_library(
    name = "twist_go_proto",
    go_deps = [
        ":vector3_go_proto",
    ],
    deps = [":twist_proto"],
)

cc_proto_library(
    name = "twist_cc_proto",
    deps = [":twist_proto"],
)

proto_library(
    name = "vector2_proto",
    srcs = ["vector2.proto"],
)

cc_proto_library(
    name = "vector2_cc_proto",
    deps = [":vector2_proto"],
)

py_proto_library(
    name = "vector2_py_pb2",
    deps = [":vector2_proto"],
)

go_proto_library(
    name = "vector2_go_proto",
    deps = [":vector2_proto"],
)

proto_library(
    name = "vector3_proto",
    srcs = ["vector3.proto"],
)

cc_proto_library(
    name = "vector3_cc_proto",
    deps = [":vector3_proto"],
)

go_proto_library(
    name = "vector3_go_proto",
    deps = [":vector3_proto"],
)

py_proto_library(
    name = "vector3_py_pb2",
    deps = [":vector3_proto"],
)

proto_library(
    name = "math_protos",
    deps = [
        ":array_proto",
        ":header_proto",
        ":matrix_proto",
        ":point_proto",
        ":pose_proto",
        ":quaternion_proto",
        ":tf_message_proto",
        ":transform_proto",
        ":transform_stamped_proto",
        ":vector2_proto",
        ":vector3_proto",
    ],
)

proto_library(
    name = "header_proto",
    srcs = ["header.proto"],
    deps = [
        "@com_google_protobuf//:timestamp_proto",
    ],
)

cc_proto_library(
    name = "header_cc_proto",
    deps = [":header_proto"],
)

py_proto_library(
    name = "header_py_pb2",
    deps = [":header_proto"],
)

go_proto_library(
    name = "header_go_proto",
    deps = [":header_proto"],
)

proto_library(
    name = "transform_proto",
    srcs = ["transform.proto"],
    deps = [
        ":quaternion_proto",
        ":vector3_proto",
    ],
)

cc_proto_library(
    name = "transform_cc_proto",
    deps = [":transform_proto"],
)

py_proto_library(
    name = "transform_py_pb2",
    deps = [":transform_proto"],
)

go_proto_library(
    name = "transform_go_proto",
    go_deps = [
        ":vector3_go_proto",
        ":quaternion_go_proto",
    ],
    deps = [":transform_proto"],
)

proto_library(
    name = "transform_stamped_proto",
    srcs = ["transform_stamped.proto"],
    deps = [
        ":header_proto",
        ":transform_proto",
    ],
)

cc_proto_library(
    name = "transform_stamped_cc_proto",
    deps = [":transform_stamped_proto"],
)

py_proto_library(
    name = "transform_stamped_py_pb2",
    deps = [":transform_stamped_proto"],
)

go_proto_library(
    name = "transform_stamped_go_proto",
    go_deps = [
        ":header_go_proto",
        ":transform_go_proto",
    ],
    deps = [":transform_stamped_proto"],
)

proto_library(
    name = "tf_message_proto",
    srcs = ["tf_message.proto"],
    deps = [
        ":transform_stamped_proto",
    ],
)

cc_proto_library(
    name = "tf_message_cc_proto",
    deps = [":tf_message_proto"],
)

py_proto_library(
    name = "tf_message_py_pb2",
    deps = [":tf_message_proto"],
)

go_proto_library(
    name = "tf_message_go_proto",
    go_deps = [
        ":transform_stamped_go_proto",
    ],
    deps = [":tf_message_proto"],
)
