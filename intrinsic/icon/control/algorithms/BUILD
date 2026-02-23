# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "joint_position_pid_velocity_controller",
    srcs = [
        "joint_position_pid_velocity_controller.cc",
    ],
    hdrs = [
        "joint_position_pid_velocity_controller.h",
    ],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/eigenmath:clamp",
        "//intrinsic/icon/proto:joint_position_pid_velocity_controller_config_cc_proto",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/math/signals:butter_filter2",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status:statusor",
    ],
)
