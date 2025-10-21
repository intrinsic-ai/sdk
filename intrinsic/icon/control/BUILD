# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

cc_library(
    name = "joint_position_command",
    srcs = ["joint_position_command.cc"],
    hdrs = ["joint_position_command.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/kinematics/types:dynamic_limits_check_mode",
    ],
)

cc_library(
    name = "streaming_io_types",
    hdrs = ["streaming_io_types.h"],
    deps = [
        "//intrinsic/production/external/intops:strong_int",
        "@abseil-cpp//absl/status:statusor",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "slot_types",
    hdrs = ["slot_types.h"],
    deps = [
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "//intrinsic/production/external/intops:strong_int",
    ],
)

cc_library(
    name = "realtime_clock_interface",
    srcs = ["realtime_clock_interface.cc"],
    hdrs = ["realtime_clock_interface.h"],
    visibility = [
        "//intrinsic:public_api_users",
    ],
    deps = [
        "//intrinsic/icon/utils:core_time",
        "//intrinsic/icon/utils:realtime_status",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "realtime_signal_types",
    hdrs = ["realtime_signal_types.h"],
    deps = ["//intrinsic/production/external/intops:strong_int"],
)

cc_library(
    name = "logging_mode",
    srcs = ["logging_mode.cc"],
    hdrs = ["logging_mode.h"],
    visibility = ["//intrinsic/icon:__subpackages__"],
    deps = ["//intrinsic/icon/proto:logging_mode_cc_proto"],
)
