# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:defs.bzl", "cc_library")

package(
    default_visibility = ["//visibility:public"],
)

cc_library(
    name = "gpio_service_proto_utils",
    srcs = ["gpio_service_proto_utils.cc"],
    hdrs = ["gpio_service_proto_utils.h"],
    visibility = ["//intrinsic:__subpackages__"],
    deps = [
        "//intrinsic/hardware/gpio/v1:gpio_service_cc_proto",
        "//intrinsic/hardware/gpio/v1:signal_cc_proto",
        "@com_google_absl//absl/container:flat_hash_map",
    ],
)

cc_library(
    name = "gpio_client",
    srcs = ["gpio_client.cc"],
    hdrs = ["gpio_client.h"],
    visibility = ["//intrinsic:__subpackages__"],
    deps = [
        ":gpio_service_proto_utils",
        "//intrinsic/hardware/gpio/v1:gpio_service_cc_grpc_proto",
        "//intrinsic/hardware/gpio/v1:gpio_service_cc_proto",
        "//intrinsic/icon/release:grpc_time_support",
        "//intrinsic/util/grpc",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:channel_interface",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_conversion_rpc",
        "//intrinsic/util/status:status_macros",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_absl//absl/container:flat_hash_set",
        "@com_google_absl//absl/log",
        "@com_google_absl//absl/status",
        "@com_google_absl//absl/status:statusor",
        "@com_google_absl//absl/strings",
        "@com_google_absl//absl/synchronization",
        "@com_google_absl//absl/time",
    ],
)
