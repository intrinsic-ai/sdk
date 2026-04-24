# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = [
    "//intrinsic/icon/hardware_modules:__subpackages__",
    "//sdk:__subpackages__",
])

cc_library(
    name = "service_variable_device",
    srcs = ["service_variable_device.cc"],
    hdrs = ["service_variable_device.h"],
    deps = [
        "//intrinsic/icon/hal:default_hardware_interfaces",
        "//intrinsic/icon/hal:hardware_interface_handle",
        "//intrinsic/icon/hal:hardware_interface_registry",
        "//intrinsic/icon/hal/interfaces:io_controller_fbs_cc",
        "//intrinsic/icon/hal/lib/fieldbus:async_device_request",
        "//intrinsic/icon/hal/lib/fieldbus:bus_component",
        "//intrinsic/icon/hal/lib/fieldbus:bus_component_factory",
        "//intrinsic/icon/hal/lib/fieldbus:device_init_context",
        "//intrinsic/icon/hal/lib/fieldbus:service_variable",
        "//intrinsic/icon/hal/lib/fieldbus:variable_registry",
        "//intrinsic/icon/hal/lib/fieldbus/v1:value_parsing_cc_proto",
        "//intrinsic/icon/hal/lib/service_variable_value/v1:service_variable_config_cc_proto",
        "//intrinsic/icon/hal/lib/service_variable_value/v1:service_variable_export_cc_proto",
        "//intrinsic/icon/utils:async_buffer",
        "//intrinsic/icon/utils:log",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/logging/proto:log_item_cc_proto",
        "//intrinsic/platform/pubsub",
        "//intrinsic/platform/pubsub:publisher",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/thread",
        "//intrinsic/util/thread:stop_token",
        "//intrinsic/util/thread:thread_utils",
        "@abseil-cpp//absl/functional:any_invocable",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@com_google_protobuf//:protobuf",
    ],
)
