# Copyright 2023 Intrinsic Innovation LLC

# Perception capabilities that are exprosed to clients of Intrinsic perception.

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = [
    "//visibility:public",
])

cc_library(
    name = "pose_estimator_id_utils",
    srcs = ["pose_estimator_id_utils.cc"],
    hdrs = ["pose_estimator_id_utils.h"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/perception/proto/v1:pose_estimator_id_cc_proto",
    ],
)

cc_library(
    name = "calibration_service_utils",
    srcs = ["calibration_service_utils.cc"],
    hdrs = ["calibration_service_utils.h"],
    deps = [
        "//intrinsic/connect/cc/grpc:channel",
        "//intrinsic/icon/release:grpc_time_support",
        "//intrinsic/perception/proto/v1:calibration_service_cc_grpc_proto",
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/cc:skill_utils",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_builder",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/time:deadline_timeout",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
    ],
)
