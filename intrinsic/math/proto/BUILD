# Copyright 2026 Intrinsic Innovation LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "vector3_proto",
    srcs = ["vector3.proto"],
    visibility = ["//visibility:public"],
)

proto_library(
    name = "accel_proto",
    srcs = ["accel.proto"],
    visibility = ["//visibility:public"],
    deps = [":vector3_proto"],
)

cc_proto_library(
    name = "vector3_cc_proto",
    deps = [":vector3_proto"],
)

go_proto_library(
    name = "vector3_go_proto",
    importpath = "intrinsic/math/proto/vector3_go_proto",
    protos = [":vector3_proto"],
)

py_proto_library(
    name = "vector3_py_pb2",
    deps = [":vector3_proto"],
)

cc_proto_library(
    name = "accel_cc_proto",
    deps = [":accel_proto"],
)

go_proto_library(
    name = "accel_go_proto",
    importpath = "intrinsic/math/proto/accel_go_proto",
    protos = [":accel_proto"],
    deps = [
        ":vector3_go_proto",
    ],
)

py_proto_library(
    name = "accel_py_pb2",
    deps = [":accel_proto"],
)

proto_library(
    name = "math_protos",
    visibility = ["//visibility:public"],
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
    name = "matrix_proto",
    srcs = ["matrix.proto"],
    visibility = ["//visibility:public"],
)

proto_library(
    name = "point_proto",
    srcs = ["point.proto"],
    visibility = ["//visibility:public"],
)

proto_library(
    name = "affine_proto",
    srcs = ["affine.proto"],
    visibility = ["//visibility:public"],
    deps = [
        ":matrix_proto",
        ":point_proto",
    ],
)

proto_library(
    name = "array_proto",
    srcs = ["array.proto"],
    visibility = ["//visibility:public"],
)

proto_library(
    name = "header_proto",
    srcs = ["header.proto"],
    visibility = ["//visibility:public"],
    deps = ["@com_google_protobuf//:timestamp_proto"],
)

proto_library(
    name = "quaternion_proto",
    srcs = ["quaternion.proto"],
    visibility = ["//visibility:public"],
)

proto_library(
    name = "pose_proto",
    srcs = ["pose.proto"],
    visibility = ["//visibility:public"],
    deps = [
        ":point_proto",
        ":quaternion_proto",
    ],
)

proto_library(
    name = "transform_proto",
    srcs = ["transform.proto"],
    visibility = ["//visibility:public"],
    deps = [
        ":quaternion_proto",
        ":vector3_proto",
    ],
)

proto_library(
    name = "transform_stamped_proto",
    srcs = ["transform_stamped.proto"],
    visibility = ["//visibility:public"],
    deps = [
        ":header_proto",
        ":transform_proto",
    ],
)

proto_library(
    name = "tf_message_proto",
    srcs = ["tf_message.proto"],
    visibility = ["//visibility:public"],
    deps = [":transform_stamped_proto"],
)

proto_library(
    name = "twist_proto",
    srcs = ["twist.proto"],
    visibility = ["//visibility:public"],
    deps = [":vector3_proto"],
)

proto_library(
    name = "vector2_proto",
    srcs = ["vector2.proto"],
    visibility = ["//visibility:public"],
)

cc_proto_library(
    name = "matrix_cc_proto",
    deps = [":matrix_proto"],
)

go_proto_library(
    name = "matrix_go_proto",
    importpath = "intrinsic/math/proto/matrix_go_proto",
    protos = [":matrix_proto"],
)

py_proto_library(
    name = "matrix_py_pb2",
    deps = [":matrix_proto"],
)

cc_proto_library(
    name = "point_cc_proto",
    deps = [":point_proto"],
)

go_proto_library(
    name = "point_go_proto",
    importpath = "intrinsic/math/proto/point_go_proto",
    protos = [":point_proto"],
)

py_proto_library(
    name = "point_py_pb2",
    deps = [":point_proto"],
)

cc_proto_library(
    name = "affine_cc_proto",
    deps = [":affine_proto"],
)

go_proto_library(
    name = "affine_go_proto",
    importpath = "intrinsic/math/proto/affine_go_proto",
    protos = [":affine_proto"],
    deps = [
        ":matrix_go_proto",
        ":point_go_proto",
    ],
)

py_proto_library(
    name = "affine_py_pb2",
    deps = [":affine_proto"],
)

cc_proto_library(
    name = "array_cc_proto",
    deps = [":array_proto"],
)

go_proto_library(
    name = "array_go_proto",
    importpath = "intrinsic/math/proto/array_go_proto",
    protos = [":array_proto"],
)

py_proto_library(
    name = "array_py_pb2",
    deps = [":array_proto"],
)

cc_proto_library(
    name = "header_cc_proto",
    deps = [":header_proto"],
)

go_proto_library(
    name = "header_go_proto",
    importpath = "intrinsic/math/proto/header_go_proto",
    protos = [":header_proto"],
)

py_proto_library(
    name = "header_py_pb2",
    deps = [":header_proto"],
)

cc_proto_library(
    name = "quaternion_cc_proto",
    deps = [":quaternion_proto"],
)

go_proto_library(
    name = "quaternion_go_proto",
    importpath = "intrinsic/math/proto/quaternion_go_proto",
    protos = [":quaternion_proto"],
)

py_proto_library(
    name = "quaternion_py_pb2",
    deps = [":quaternion_proto"],
)

cc_proto_library(
    name = "pose_cc_proto",
    deps = [":pose_proto"],
)

go_proto_library(
    name = "pose_go_proto",
    importpath = "intrinsic/math/proto/pose_go_proto",
    protos = [":pose_proto"],
    deps = [
        ":point_go_proto",
        ":quaternion_go_proto",
    ],
)

py_proto_library(
    name = "pose_py_pb2",
    deps = [":pose_proto"],
)

cc_proto_library(
    name = "transform_cc_proto",
    deps = [":transform_proto"],
)

go_proto_library(
    name = "transform_go_proto",
    importpath = "intrinsic/math/proto/transform_go_proto",
    protos = [":transform_proto"],
    deps = [
        ":quaternion_go_proto",
        ":vector3_go_proto",
    ],
)

py_proto_library(
    name = "transform_py_pb2",
    deps = [":transform_proto"],
)

cc_proto_library(
    name = "transform_stamped_cc_proto",
    deps = [":transform_stamped_proto"],
)

go_proto_library(
    name = "transform_stamped_go_proto",
    importpath = "intrinsic/math/proto/transform_stamped_go_proto",
    protos = [":transform_stamped_proto"],
    deps = [
        ":header_go_proto",
        ":transform_go_proto",
    ],
)

py_proto_library(
    name = "transform_stamped_py_pb2",
    deps = [":transform_stamped_proto"],
)

cc_proto_library(
    name = "tf_message_cc_proto",
    deps = [":tf_message_proto"],
)

go_proto_library(
    name = "tf_message_go_proto",
    importpath = "intrinsic/math/proto/tf_message_go_proto",
    protos = [":tf_message_proto"],
    deps = [
        ":transform_stamped_go_proto",
    ],
)

py_proto_library(
    name = "tf_message_py_pb2",
    deps = [":tf_message_proto"],
)

cc_proto_library(
    name = "twist_cc_proto",
    deps = [":twist_proto"],
)

go_proto_library(
    name = "twist_go_proto",
    importpath = "intrinsic/math/proto/twist_go_proto",
    protos = [":twist_proto"],
    deps = [
        ":vector3_go_proto",
    ],
)

py_proto_library(
    name = "twist_py_pb2",
    deps = [":twist_proto"],
)

cc_proto_library(
    name = "vector2_cc_proto",
    deps = [":vector2_proto"],
)

go_proto_library(
    name = "vector2_go_proto",
    importpath = "intrinsic/math/proto/vector2_go_proto",
    protos = [":vector2_proto"],
)

py_proto_library(
    name = "vector2_py_pb2",
    deps = [":vector2_proto"],
)
