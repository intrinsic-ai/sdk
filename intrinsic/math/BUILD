# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "pose3",
    hdrs = [
        "pose3.h",
    ],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/eigenmath:so3",
        "@com_gitlab_libeigen_eigen//:eigen",
        "@com_google_absl//absl/base:core_headers",
    ],
)

cc_library(
    name = "almost_equals",
    hdrs = ["almost_equals.h"],
)

cc_library(
    name = "proto_conversion",
    srcs = ["proto_conversion.cc"],
    hdrs = ["proto_conversion.h"],
    deps = [
        ":almost_equals",
        ":pose3",
        "//intrinsic/eigenmath",
        "//intrinsic/math/proto:matrix_cc_proto",
        "//intrinsic/math/proto:point_cc_proto",
        "//intrinsic/math/proto:pose_cc_proto",
        "//intrinsic/math/proto:quaternion_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
    ],
)

cc_library(
    name = "gaussian_noise",
    srcs = ["gaussian_noise.cc"],
    hdrs = ["gaussian_noise.h"],
    deps = [
        "//intrinsic/eigenmath",
        "@com_google_absl//absl/random:bit_gen_ref",
        "@com_google_absl//absl/random:distributions",
    ],
)

cc_library(
    name = "twist",
    srcs = ["twist.cc"],
    hdrs = ["twist.h"],
    deps = ["//intrinsic/eigenmath"],
)

cc_library(
    name = "transform_utils",
    srcs = ["transform_utils.cc"],
    hdrs = ["transform_utils.h"],
    deps = [
        ":pose3",
        ":twist",
    ],
)

cc_library(
    name = "inertia_utils",
    srcs = ["inertia_utils.cc"],
    hdrs = ["inertia_utils.h"],
    deps = [
        ":almost_equals",
        "//intrinsic/eigenmath",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/kinematics/types:to_fixed_string",
        "//intrinsic/world/proto:robot_payload_cc_proto",
        "@com_gitlab_libeigen_eigen//:eigen",
        "@com_google_absl//absl/strings:string_view",
    ],
)
