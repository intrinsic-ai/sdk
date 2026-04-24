# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(
    default_visibility = [
        "//sdk:__subpackages__",
        "//visibility:public",
    ],
)

cc_library(
    name = "ati_ecat_objects",
    hdrs = ["ati_constants.h"],
)

cc_library(
    name = "ati_force_torque_bus_component",
    srcs = ["ati_force_torque_bus_component.cc"],
    hdrs = ["ati_force_torque_bus_component.h"],
    deps = [
        ":ati_ecat_objects",
        "//intrinsic/icon/flatbuffers:transform_types_fbs_cc",
        "//intrinsic/icon/hal:default_hardware_interfaces",
        "//intrinsic/icon/hal:hardware_interface_handle",
        "//intrinsic/icon/hal:hardware_interface_registry",
        "//intrinsic/icon/hal:module_config",
        "//intrinsic/icon/hal/interfaces:force_sensor_fbs_cc",
        "//intrinsic/icon/hal/interfaces:force_torque_fbs_cc",
        "//intrinsic/icon/hal/lib/ati/v1:ati_force_torque_bus_component_config_cc_proto",
        "//intrinsic/icon/hal/lib/fieldbus:async_device_request",
        "//intrinsic/icon/hal/lib/fieldbus:bus_component",
        "//intrinsic/icon/hal/lib/fieldbus:bus_component_factory",
        "//intrinsic/icon/hal/lib/fieldbus:device_init_context",
        "//intrinsic/icon/hal/lib/fieldbus:process_variable",
        "//intrinsic/icon/hal/lib/fieldbus:variable_registry",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/icon/utils:sensor_utils",
        "//intrinsic/util/status:annotate",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
    ],
    alwayslink = 1,
)
