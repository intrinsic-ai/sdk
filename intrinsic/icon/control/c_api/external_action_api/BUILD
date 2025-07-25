# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@rules_cc//cc:cc_binary.bzl", "cc_binary")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

package(default_visibility = [
    "//visibility:public",
])

cc_library(
    name = "icon_action_interface",
    hdrs = ["icon_action_interface.h"],
    deps = [
        ":icon_action_factory_context",
        ":icon_realtime_signal_access",
        ":icon_realtime_slot_map",
        ":icon_streaming_io_access",
        "//intrinsic/icon/cc_client:condition",
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "//intrinsic/icon/utils:realtime_status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "icon_feature_interfaces",
    srcs = ["icon_feature_interfaces.cc"],
    hdrs = ["icon_feature_interfaces.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/icon/control:joint_position_command",
        "//intrinsic/icon/control/c_api:convert_c_realtime_status",
        "//intrinsic/icon/control/c_api:convert_c_types",
        "//intrinsic/icon/control/c_api:icon_c_api",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/kinematics/types:joint_limits",
        "//intrinsic/kinematics/types:joint_state",
        "//intrinsic/math:pose3",
        "//intrinsic/math:twist",
    ],
)

cc_library(
    name = "icon_realtime_slot_map",
    srcs = ["icon_realtime_slot_map.cc"],
    hdrs = ["icon_realtime_slot_map.h"],
    deps = [
        ":icon_feature_interfaces",
        "//intrinsic/icon/control:slot_types",
        "//intrinsic/icon/control/c_api:icon_c_api",
    ],
)

cc_library(
    name = "icon_realtime_signal_access",
    srcs = ["icon_realtime_signal_access.cc"],
    hdrs = ["icon_realtime_signal_access.h"],
    deps = [
        "//intrinsic/icon/control:realtime_signal_types",
        "//intrinsic/icon/control/c_api:convert_c_realtime_status",
        "//intrinsic/icon/control/c_api:convert_c_types",
        "//intrinsic/icon/control/c_api:icon_c_api",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/icon/utils:realtime_status_or",
    ],
)

cc_library(
    name = "icon_action_factory_context",
    srcs = ["icon_action_factory_context.cc"],
    hdrs = ["icon_action_factory_context.h"],
    deps = [
        "//intrinsic/icon/control:realtime_signal_types",
        "//intrinsic/icon/control:slot_types",
        "//intrinsic/icon/control:streaming_io_types",
        "//intrinsic/icon/control/c_api:convert_c_realtime_status",
        "//intrinsic/icon/control/c_api:icon_c_api",
        "//intrinsic/icon/control/c_api/wrappers:streaming_io_wrapper",
        "//intrinsic/icon/control/c_api/wrappers:string_wrapper",
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "icon_streaming_io_access",
    hdrs = ["icon_streaming_io_access.h"],
    deps = [
        "//intrinsic/icon/control:streaming_io_types",
        "//intrinsic/icon/control/c_api:convert_c_realtime_status",
        "//intrinsic/icon/control/c_api:icon_c_api",
        "//intrinsic/icon/control/c_api/wrappers:streaming_io_wrapper",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/util/status:status_macros",
    ],
)

cc_library(
    name = "icon_plugin_register_macro",
    hdrs = ["icon_plugin_register_macro.h"],
    deps = [
        ":make_icon_action_vtable",
        "//intrinsic/icon/control/c_api:icon_c_api",
    ],
)

cc_library(
    name = "make_icon_action_vtable",
    hdrs = ["make_icon_action_vtable.h"],
    deps = [
        ":icon_action_factory_context",
        ":icon_action_interface",
        ":icon_realtime_signal_access",
        ":icon_realtime_slot_map",
        ":icon_streaming_io_access",
        "//intrinsic/icon/control/c_api:convert_c_realtime_status",
        "//intrinsic/icon/control/c_api:icon_c_api",
        "//intrinsic/icon/control/c_api/wrappers:string_wrapper",
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "//intrinsic/icon/utils:realtime_status",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

proto_library(
    name = "sine_wave_action_proto",
    srcs = ["sine_wave_action.proto"],
)

cc_proto_library(
    name = "sine_wave_action_cc_proto",
    deps = [":sine_wave_action_proto"],
)

cc_library(
    name = "sine_wave_plugin_action",
    srcs = ["sine_wave_plugin_action.cc"],
    hdrs = ["sine_wave_plugin_action.h"],
    deps = [
        ":icon_action_factory_context",
        ":icon_action_interface",
        ":icon_feature_interfaces",
        ":icon_realtime_signal_access",
        ":icon_realtime_slot_map",
        ":icon_streaming_io_access",
        ":sine_wave_action_cc_proto",
        "//intrinsic/eigenmath",
        "//intrinsic/icon/actions:action_utils",
        "//intrinsic/icon/cc_client:condition",
        "//intrinsic/icon/control:joint_position_command",
        "//intrinsic/icon/control:slot_types",
        "//intrinsic/icon/control:streaming_io_types",
        "//intrinsic/icon/proto:generic_part_config_cc_proto",
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/kinematics/types:joint_limits",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@com_google_protobuf//:protobuf",
    ],
    alwayslink = True,
)

cc_test(
    name = "sine_wave_plugin_action_test",
    srcs = ["sine_wave_plugin_action_test.cc"],
    deps = [
        ":sine_wave_plugin_action",
        "//intrinsic/icon/cc_client:condition",
        "//intrinsic/icon/control/c_api/external_action_api/testing:action_test_helper",
        "//intrinsic/icon/control/c_api/external_action_api/testing:loopback_fake_arm",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_binary(
    name = "sine_wave_plugin.so",
    srcs = ["sine_wave_plugin.cc"],
    linkshared = 1,
    linkstatic = 1,
    deps = [
        ":icon_plugin_register_macro",
        ":sine_wave_plugin_action",
    ],
)

cc_binary(
    name = "move_sine_wave",
    srcs = ["move_sine_wave.cc"],
    deps = [
        ":sine_wave_action_cc_proto",
        ":sine_wave_plugin_action",
        "//intrinsic/icon/cc_client:condition",
        "//intrinsic/icon/cc_client:session",
        "//intrinsic/icon/common:id_types",
        "//intrinsic/icon/release/portable:init_xfa_absl",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
    ],
)
