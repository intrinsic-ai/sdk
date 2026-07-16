# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = [
    "//intrinsic/icon/skills:__subpackages__",
])

cc_library(
    name = "analog_digital_io",
    srcs = ["analog_digital_io.cc"],
    hdrs = ["analog_digital_io.h"],
    deps = [
        ":adio_util",
        "//intrinsic/connect/cc/grpc:channel",
        "//intrinsic/hardware/gpio:gpio_service_equipment_cc_proto",
        "//intrinsic/hardware/gpio/v1:gpio_service_cc_grpc_proto",
        "//intrinsic/hardware/gpio/v1:gpio_service_cc_proto",
        "//intrinsic/hardware/gpio/v1:signal_cc_proto",
        "//intrinsic/icon/actions:adio_cc_proto",
        "//intrinsic/icon/actions:adio_info",
        "//intrinsic/icon/cc_client:client",
        "//intrinsic/icon/cc_client:condition",
        "//intrinsic/icon/cc_client:session",
        "//intrinsic/icon/cc_client:state_variable_path",
        "//intrinsic/icon/common:id_types",
        "//intrinsic/icon/control/parts:io_block",
        "//intrinsic/icon/equipment:channel_factory",
        "//intrinsic/icon/equipment:icon_equipment_cc_proto",
        "//intrinsic/icon/proto:generic_part_config_cc_proto",
        "//intrinsic/icon/proto:io_block_cc_proto",
        "//intrinsic/icon/release:grpc_time_support",
        "//intrinsic/resources/proto:resource_handle_cc_proto",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/util/grpc:channel_interface",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_conversion_rpc",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/container:btree",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/time",
        "@abseil-cpp//absl/types:span",
        "@com_github_grpc_grpc//:grpc++",
    ],
)

cc_library(
    name = "adio_util",
    srcs = ["adio_util.cc"],
    hdrs = ["adio_util.h"],
    deps = [
        "//intrinsic/hardware/gpio:gpio_service_equipment_cc_proto",
        "//intrinsic/hardware/gpio/v1:gpio_service_cc_proto",
        "//intrinsic/icon/actions:adio_cc_proto",
        "//intrinsic/icon/cc_client:robot_config",
        "//intrinsic/icon/equipment:icon_equipment_cc_proto",
        "//intrinsic/icon/proto:generic_part_config_cc_proto",
        "//intrinsic/icon/proto:io_block_cc_proto",
        "//intrinsic/icon/release:grpc_time_support",
        "//intrinsic/resources/proto:resource_handle_cc_proto",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:str_format",
    ],
)

cc_library(
    name = "analog_digital_io_factory",
    srcs = ["analog_digital_io_factory.cc"],
    hdrs = ["analog_digital_io_factory.h"],
    deps = [
        ":analog_digital_io",
        "//intrinsic/icon/equipment:channel_factory",
        "//intrinsic/icon/equipment:equipment_utils",
        "//intrinsic/icon/equipment:icon_equipment_cc_proto",
        "//intrinsic/resources/proto:resource_handle_cc_proto",
        "//intrinsic/skills/cc:equipment_pack",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
    ],
)

cc_library(
    name = "position_part_util",
    srcs = ["position_part_util.cc"],
    hdrs = ["position_part_util.h"],
    visibility = [
        "//intrinsic/icon/skills:__subpackages__",
        "//intrinsic/motion_planning/skills:__pkg__",
        "//intrinsic/motion_planning/skills:__subpackages__",
        "//intrinsic/perception/skills:__subpackages__",
    ],
    deps = [
        "//intrinsic/icon/equipment:icon_equipment_cc_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/objects:object_world_client",
        "//intrinsic/world/objects:world_object",
        "//intrinsic/world/proto:object_world_refs_cc_proto",
        "//intrinsic/world/util:object_reference_utils",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
    ],
)
