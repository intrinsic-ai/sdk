# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = [
    "//intrinsic/icon/hardware_modules:__subpackages__",
    "//sdk:__subpackages__",
])

cc_library(
    name = "ds402_driver",
    srcs = ["ds402_driver.cc"],
    hdrs = ["ds402_driver.h"],
    deps = [
        "//intrinsic/icon/utils:bitset",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
        "@abseil-cpp//absl/strings:string_view",
    ],
)

cc_library(
    name = "ds402_bus_component",
    srcs = ["ds402_bus_component.cc"],
    hdrs = ["ds402_bus_component.h"],
    deps = [
        ":ds402_bus_component_utils",
        ":ds402_driver",
        ":homing_task",
        "//intrinsic/icon/hal/interfaces:electrical_motor_hardware_interfaces",
        "//intrinsic/icon/hal/lib/adio/v1:adio_bus_component_config_cc_proto",
        "//intrinsic/icon/hal/lib/ds402/v1:ds402_bus_component_config_cc_proto",
        "//intrinsic/icon/hal/lib/fieldbus:async_device_request",
        "//intrinsic/icon/hal/lib/fieldbus:bus_component",
        "//intrinsic/icon/hal/lib/fieldbus:bus_component_factory",
        "//intrinsic/icon/hal/lib/fieldbus:device_init_context",
        "//intrinsic/icon/hal/lib/fieldbus:process_variable",
        "//intrinsic/icon/hal/lib/fieldbus:service_variable",
        "//intrinsic/icon/hal/lib/fieldbus:variable_registry",
        "//intrinsic/icon/hal/lib/fieldbus/v1:value_parsing_cc_proto",
        "//intrinsic/icon/testing:realtime_annotations",
        "//intrinsic/icon/utils:bitset",
        "//intrinsic/icon/utils:fixed_str_cat",
        "//intrinsic/icon/utils:fixed_string",
        "//intrinsic/icon/utils:log",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/platform/common/buffers:realtime_write_queue",
        "//intrinsic/platform/common/buffers:rt_queue",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/thread",
        "//intrinsic/util/thread:stop_token",
        "//intrinsic/util/thread:thread_options",
        "//intrinsic/util/thread:thread_utils",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "ds402_fake_device",
    srcs = ["ds402_fake_device.cc"],
    hdrs = ["ds402_fake_device.h"],
    deps = [
        ":ds402_driver",
        ":homing_task",
        "//intrinsic/icon/hal/lib/fieldbus:fake_variable_registry",
        "@abseil-cpp//absl/strings:string_view",
    ],
)

cc_library(
    name = "ds402_bus_component_utils",
    srcs = ["ds402_bus_component_utils.cc"],
    hdrs = ["ds402_bus_component_utils.h"],
    deps = [
        "//intrinsic/icon/hal/lib/adio/v1:adio_bus_component_config_cc_proto",
        "//intrinsic/icon/hal/lib/ds402/v1:ds402_bus_component_config_cc_proto",
        "//intrinsic/icon/hal/lib/fieldbus:process_variable",
        "//intrinsic/icon/hal/lib/fieldbus:value_pattern_utils",
        "//intrinsic/icon/hal/lib/fieldbus:variable_registry",
        "//intrinsic/icon/hal/lib/fieldbus/v1:value_parsing_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/functional:any_invocable",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "homing_task",
    srcs = ["homing_task.cc"],
    hdrs = ["homing_task.h"],
    deps = [
        ":ds402_driver",
        "//intrinsic/icon/flatbuffers:fixed_string",
        "//intrinsic/icon/hal:command_validator",
        "//intrinsic/icon/hal:hardware_interface_handle",
        "//intrinsic/icon/hal:hardware_interface_registry",
        "//intrinsic/icon/hal/interfaces:electrical_motor_fbs_cc",
        "//intrinsic/icon/hal/interfaces:electrical_motor_hardware_interfaces",
        "//intrinsic/icon/hal/lib/ds402/v1:ds402_bus_component_config_cc_proto",
        "//intrinsic/icon/hal/lib/fieldbus:process_variable",
        "//intrinsic/icon/hal/lib/fieldbus:service_variable",
        "//intrinsic/icon/hal/lib/fieldbus:variable_registry",
        "//intrinsic/icon/utils:async_function",
        "//intrinsic/icon/utils:bitset",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/platform/common/buffers:rt_promise",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/functional:any_invocable",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/time",
    ],
)
