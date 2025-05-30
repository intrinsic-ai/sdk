# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_binary", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_binary(
    name = "clear_faults",
    srcs = ["clear_faults.cc"],
    deps = [
        "//intrinsic/icon/cc_client:client",
        "//intrinsic/icon/cc_client:operational_status",
        "//intrinsic/icon/release/portable:init_xfa_absl",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_macros",
        "@com_google_absl//absl/flags:flag",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/log:check",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/strings",
    ],
)

cc_library(
    name = "generate_documentation",
    srcs = ["generate_documentation.cc"],
    hdrs = ["generate_documentation.h"],
    deps = [
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/types:span",
    ],
)

cc_binary(
    name = "list_actions",
    srcs = ["list_actions.cc"],
    deps = [
        ":generate_documentation",
        "//intrinsic/icon/cc_client:client",
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "//intrinsic/icon/release/portable:init_xfa_absl",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_macros",
        "@com_google_absl//absl/flags:flag",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
    ],
)

cc_binary(
    name = "list_parts",
    srcs = ["list_parts.cc"],
    deps = [
        "//intrinsic/icon/cc_client:client",
        "//intrinsic/icon/release/portable:init_xfa_absl",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_macros",
        "@com_google_absl//absl/flags:flag",
        "@com_google_absl//absl/log:check",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/strings",
    ],
)

cc_binary(
    name = "restart_server",
    srcs = ["restart_server.cc"],
    deps = [
        "//intrinsic/icon/cc_client:client",
        "//intrinsic/icon/release/portable:init_xfa_absl",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_macros",
        "@com_google_absl//absl/flags:flag",
        "@com_google_absl//absl/log:check",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/strings",
    ],
)

cc_binary(
    name = "show_joint_angles",
    srcs = ["show_joint_angles.cc"],
    deps = [
        "//intrinsic/icon/cc_client:client",
        "//intrinsic/icon/cc_client:operational_status",
        "//intrinsic/icon/proto:part_status_cc_proto",
        "//intrinsic/icon/release/portable:init_xfa_absl",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_macros",
        "@com_google_absl//absl/flags:flag",
        "@com_google_absl//absl/log:check",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/strings:str_format",
        "@com_google_absl//absl/time",
    ],
)
