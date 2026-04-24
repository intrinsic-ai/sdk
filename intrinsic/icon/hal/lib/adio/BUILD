# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = [
    "//intrinsic/icon/hardware_modules:__subpackages__",
    "//sdk:__subpackages__",
])

cc_library(
    name = "adio_bus_component",
    srcs = ["adio_bus_component.cc"],
    hdrs = ["adio_bus_component.h"],
    deps = [
        "//intrinsic/icon/hal:default_hardware_interfaces",
        "//intrinsic/icon/hal:hardware_interface_handle",
        "//intrinsic/icon/hal:hardware_interface_registry",
        "//intrinsic/icon/hal/interfaces:adio_fbs_cc",
        "//intrinsic/icon/hal/interfaces:io_controller_fbs_cc",
        "//intrinsic/icon/hal/lib/adio/v1:adio_bus_component_config_cc_proto",
        "//intrinsic/icon/hal/lib/fieldbus:async_device_request",
        "//intrinsic/icon/hal/lib/fieldbus:bus_component",
        "//intrinsic/icon/hal/lib/fieldbus:bus_component_factory",
        "//intrinsic/icon/hal/lib/fieldbus:device_init_context",
        "//intrinsic/icon/hal/lib/fieldbus:process_variable",
        "//intrinsic/icon/hal/lib/fieldbus:process_variable_utils",
        "//intrinsic/icon/hal/lib/fieldbus:variable_registry",
        "//intrinsic/icon/utils:fixed_str_cat",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/functional:any_invocable",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)
