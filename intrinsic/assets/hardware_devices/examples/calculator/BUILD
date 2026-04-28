# Copyright 2023 Intrinsic Innovation LLC

load("//intrinsic/assets/data/build_defs:data.bzl", "intrinsic_data")
load("//intrinsic/assets/hardware_devices/build_defs:hardware_device.bzl", "intrinsic_hardware_device")

intrinsic_data(
    name = "calculator_config_data",
    manifest = "calculator_config_data_manifest.textproto",
    visibility = ["//visibility:private"],
    deps = [
        "//intrinsic/assets/services/examples/calcserver:calc_server_proto",
    ],
)

intrinsic_hardware_device(
    name = "calculator",
    assets = [
        ":calculator_config_data",
        "//intrinsic/assets/scene_objects/examples:calculator_scene_object_reference",
        "//intrinsic/assets/services/examples/calcserver:calculator_service",
        "//intrinsic/assets/services/examples/calcserver:units_conversion_data",
    ],
    manifest = "calculator_manifest.textproto",
)
