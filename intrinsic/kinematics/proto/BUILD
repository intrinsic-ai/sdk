# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "kinematics_proto",
    srcs = ["kinematics.proto"],
    deps = [
        "//intrinsic/math/proto:matrix_proto",
        "//intrinsic/math/proto:point_proto",
        "//intrinsic/math/proto:pose_proto",
    ],
)

cc_proto_library(
    name = "kinematics_cc_proto",
    deps = [":kinematics_proto"],
)

py_proto_library(
    name = "kinematics_py_proto",
    deps = [":kinematics_proto"],
)

go_proto_library(
    name = "kinematics_go_proto",
    importpath = "intrinsic/kinematics/proto/kinematics_go_proto",
    protos = [":kinematics_proto"],
    deps = [
        "//intrinsic/math/proto:matrix_go_proto",
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
)

proto_library(
    name = "skeleton_proto",
    srcs = ["skeleton.proto"],
    deps = [
        ":kinematics_proto",
        "//intrinsic/math/proto:pose_proto",
    ],
)

cc_proto_library(
    name = "skeleton_cc_proto",
    deps = [":skeleton_proto"],
)

py_proto_library(
    name = "skeleton_py_proto",
    deps = [":skeleton_proto"],
)

go_proto_library(
    name = "skeleton_go_proto",
    importpath = "intrinsic/kinematics/proto/skeleton_go_proto",
    protos = [":skeleton_proto"],
    deps = [
        ":kinematics_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
)
