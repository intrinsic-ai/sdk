# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = [
    "//intrinsic/icon/control/parts/proto:__subpackages__",
    "//intrinsic/icon/testing:__subpackages__",
    "//visibility:public",
])

proto_library(
    name = "realtime_part_config_proto",
    srcs = ["realtime_part_config.proto"],
    visibility = [
        "//intrinsic:__subpackages__",
        "//intrinsic/icon/control/parts/proto:__subpackages__",
    ],
    deps = ["@com_google_protobuf//:any_proto"],
)

cc_proto_library(
    name = "realtime_part_config_cc_proto",
    visibility = [
        "//intrinsic:__subpackages__",
        "//intrinsic/icon/control/parts/proto:__subpackages__",
    ],
    deps = [":realtime_part_config_proto"],
)

go_proto_library(
    name = "realtime_part_config_go_proto",
    visibility = [
        "//intrinsic:__subpackages__",
        "//intrinsic/icon/control/parts/proto:__subpackages__",
    ],
    deps = [":realtime_part_config_proto"],
)

proto_library(
    name = "mode_of_safe_operation_limits_config_proto",
    srcs = ["mode_of_safe_operation_limits_config.proto"],
    visibility = [
        "//intrinsic:__subpackages__",
        "//intrinsic/icon/control/parts/proto:__subpackages__",
    ],
)

cc_proto_library(
    name = "mode_of_safe_operation_limits_config_cc_proto",
    deps = [":mode_of_safe_operation_limits_config_proto"],
)

go_proto_library(
    name = "mode_of_safe_operation_limits_config_go_proto",
    visibility = [
        "//intrinsic:__subpackages__",
        "//intrinsic/icon/control/parts/proto:__subpackages__",
    ],
    deps = [":mode_of_safe_operation_limits_config_proto"],
)
