# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library")

package(
    default_visibility = ["//visibility:public"],
)

cc_library(
    name = "icon_c_api",
    hdrs = [
        "c_action_factory_context.h",
        "c_feature_interfaces.h",
        "c_plugin_api.h",
        "c_realtime_signal_access.h",
        "c_realtime_slot_map.h",
        "c_realtime_status.h",
        "c_rtcl_action.h",
        "c_streaming_io_realtime_access.h",
        "c_types.h",
    ],
)

cc_library(
    name = "convert_c_realtime_status",
    srcs = ["convert_c_realtime_status.cc"],
    hdrs = ["convert_c_realtime_status.h"],
    deps = [
        ":icon_c_api",
        "//intrinsic/icon/utils:realtime_status",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/strings",
    ],
)

cc_library(
    name = "convert_c_types",
    srcs = ["convert_c_types.cc"],
    hdrs = ["convert_c_types.h"],
    deps = [
        ":icon_c_api",
        "//intrinsic/eigenmath",
        "//intrinsic/icon/control:joint_position_command",
        "//intrinsic/icon/control:realtime_signal_types",
        "//intrinsic/kinematics/types:joint_limits",
        "//intrinsic/kinematics/types:joint_state",
        "//intrinsic/math:pose3",
        "//intrinsic/math:twist",
        "@com_gitlab_libeigen_eigen//:eigen",
        "@com_google_absl//absl/log:check",
        "@com_google_absl//absl/types:optional",
    ],
)
