# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(
    default_visibility = [
        "//intrinsic/icon/hal:__subpackages__",
        "//sdk:__subpackages__",
        "//visibility:public",
    ],
)

cc_library(
    name = "process_variable",
    srcs = ["process_variable.cc"],
    hdrs = ["process_variable.h"],
    deps = [
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
    ],
)

cc_library(
    name = "process_variable_utils",
    srcs = ["process_variable_utils.cc"],
    hdrs = ["process_variable_utils.h"],
    deps = [
        ":process_variable",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
    ],
)

cc_library(
    name = "value_pattern_utils",
    srcs = ["value_pattern_utils.cc"],
    hdrs = ["value_pattern_utils.h"],
    deps = [
        ":process_variable",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "service_variable",
    srcs = ["service_variable.cc"],
    hdrs = ["service_variable.h"],
    visibility = ["//intrinsic/icon/hal:__subpackages__"],
    deps = [
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "variable_registry",
    hdrs = ["variable_registry.h"],
    deps = [
        ":process_variable",
        ":service_variable",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "async_device_request",
    srcs = ["async_device_request.cc"],
    hdrs = ["async_device_request.h"],
    deps = [
        "//intrinsic/icon/utils:async_request",
        "//intrinsic/icon/utils:realtime_status",
    ],
)

cc_library(
    name = "bus_component",
    hdrs = ["bus_component.h"],
    deps = [
        ":async_device_request",
        "//intrinsic/icon/utils:realtime_status_or",
    ],
)

cc_library(
    name = "bus_component_factory",
    hdrs = ["bus_component_factory.h"],
    deps = [
        ":bus_component",
        ":device_init_context",
        ":variable_registry",
        "//intrinsic/icon/hal:hardware_module_init_context",
        "@abseil-cpp//absl/log",
        "@com_github_grpc_grpc//:grpc++",
    ],
)

cc_library(
    name = "fake_variable_registry",
    srcs = ["fake_variable_registry.cc"],
    hdrs = ["fake_variable_registry.h"],
    visibility = ["//intrinsic/icon/hal:__subpackages__"],
    deps = [
        ":process_variable",
        ":service_variable",
        ":variable_registry",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "device_init_context",
    hdrs = ["device_init_context.h"],
    deps = [
        ":variable_registry",
        "//intrinsic/icon/hal:hardware_interface_registry",
        "@com_github_grpc_grpc//:grpc++",
    ],
)
