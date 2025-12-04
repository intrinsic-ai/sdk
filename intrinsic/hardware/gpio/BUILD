# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(
    default_visibility = ["//visibility:public"],
)

cc_library(
    name = "gpio_service_proto_utils",
    srcs = ["gpio_service_proto_utils.cc"],
    hdrs = ["gpio_service_proto_utils.h"],
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [
        "//intrinsic/hardware/gpio/v1:gpio_service_cc_proto",
        "//intrinsic/hardware/gpio/v1:signal_cc_proto",
        "@abseil-cpp//absl/container:flat_hash_map",
    ],
)

cc_library(
    name = "gpio_client",
    srcs = ["gpio_client.cc"],
    hdrs = ["gpio_client.h"],
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [
        ":gpio_service_proto_utils",
        "//intrinsic/connect/cc/grpc:channel",
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
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
    ],
)
