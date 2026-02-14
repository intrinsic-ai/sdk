# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

cc_library(
    name = "state_rn",
    hdrs = ["state_rn.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/util:aggregate_type",
        "@abseil-cpp//absl/log:check",
        "@com_gitlab_libeigen_eigen//:eigen",
    ],
)

cc_library(
    name = "joint_state",
    hdrs = ["joint_state.h"],
    deps = [":state_rn"],
)

cc_library(
    name = "to_fixed_string",
    srcs = ["to_fixed_string.cc"],
    hdrs = ["to_fixed_string.h"],
    deps = [
        ":cartesian_limits",
        ":joint_limits",
        ":state_rn",
        "//intrinsic/eigenmath",
        "//intrinsic/icon/utils:fixed_str_cat",
        "//intrinsic/icon/utils:fixed_string",
        "//intrinsic/math:pose3",
        "@abseil-cpp//absl/strings",
    ],
)

proto_library(
    name = "joint_limits_proto",
    srcs = ["joint_limits.proto"],
)

cc_proto_library(
    name = "joint_limits_cc_proto",
    deps = [":joint_limits_proto"],
)

py_proto_library(
    name = "joint_limits_py_pb2",
    deps = [":joint_limits_proto"],
)

go_proto_library(
    name = "joint_limits_go_proto",
    importpath = "intrinsic/kinematics/types/joint_limits_go_proto",
    protos = [":joint_limits_proto"],
)

proto_library(
    name = "dynamic_limits_check_mode_proto",
    srcs = ["dynamic_limits_check_mode.proto"],
)

proto_library(
    name = "rotational_limits_proto",
    srcs = ["rotational_limits.proto"],
    deps = [
        "//intrinsic/math/proto:quaternion_proto",
        "//intrinsic/math/proto:vector3_proto",
    ],
)

cc_proto_library(
    name = "rotational_limits_cc_proto",
    deps = [":rotational_limits_proto"],
)

py_proto_library(
    name = "rotational_limits_py_pb2",
    deps = [":rotational_limits_proto"],
)

cc_proto_library(
    name = "dynamic_limits_check_mode_cc_proto",
    deps = [":dynamic_limits_check_mode_proto"],
)

py_proto_library(
    name = "dynamic_limits_check_mode_py_pb2",
    deps = [":dynamic_limits_check_mode_proto"],
)

go_proto_library(
    name = "dynamic_limits_check_mode_go_proto",
    importpath = "intrinsic/kinematics/types/dynamic_limits_check_mode_go_proto",
    protos = [":dynamic_limits_check_mode_proto"],
)

cc_library(
    name = "joint_limits",
    srcs = ["joint_limits.cc"],
    hdrs = ["joint_limits.h"],
    deps = [
        ":joint_limits_cc_proto",
        "//intrinsic/eigenmath",
        "//intrinsic/icon/proto:eigen_conversion",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "joint_limits_xd",
    srcs = ["joint_limits_xd.cc"],
    hdrs = ["joint_limits_xd.h"],
    deps = [
        ":joint_limits",
        ":joint_limits_cc_proto",
        "//intrinsic/eigenmath",
        "//intrinsic/icon/proto:eigen_conversion",
        "//intrinsic/icon/testing:realtime_annotations",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "dynamic_limits_check_mode",
    srcs = ["dynamic_limits_check_mode.cc"],
    hdrs = ["dynamic_limits_check_mode.h"],
    deps = [
        ":dynamic_limits_check_mode_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "cartesian_limits",
    srcs = ["cartesian_limits.cc"],
    hdrs = ["cartesian_limits.h"],
    deps = [
        "//intrinsic/eigenmath",
    ],
)

cc_library(
    name = "to_string",
    srcs = ["to_string.cc"],
    hdrs = ["to_string.h"],
    deps = ["//intrinsic/eigenmath"],
)
