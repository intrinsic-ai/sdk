# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = [
    "//intrinsic/icon/hardware_modules:__subpackages__",
    "//sdk:__subpackages__",
])

cc_library(
    name = "process_variable_device",
    srcs = ["process_variable_device.cc"],
    hdrs = ["process_variable_device.h"],
    deps = [
        "//intrinsic/icon/hal:default_hardware_interfaces",
        "//intrinsic/icon/hal/lib/fieldbus:async_device_request",
        "//intrinsic/icon/hal/lib/fieldbus:bus_component",
        "//intrinsic/icon/hal/lib/fieldbus:bus_component_factory",
        "//intrinsic/icon/hal/lib/fieldbus:device_init_context",
        "//intrinsic/icon/hal/lib/fieldbus:process_variable",
        "//intrinsic/icon/hal/lib/fieldbus:variable_registry",
        "//intrinsic/icon/hal/lib/process_variable_value/v1:process_variable_config_cc_proto",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)
