# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(
    default_visibility = ["//visibility:public"],
)

cc_library(
    name = "client_utils",
    srcs = ["client_utils.cc"],
    hdrs = ["client_utils.h"],
    deps = [
        ":client",
        ":condition",
        ":robot_config",
        "//intrinsic/icon/common:builtins",
        "//intrinsic/icon/common:part_properties",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "state_variable_path",
    srcs = ["state_variable_path.cc"],
    hdrs = [
        "state_variable_path.h",
    ],
    deps = [
        "//intrinsic/icon/common:state_variable_path_constants",
        "//intrinsic/icon/common:state_variable_path_util",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "condition",
    srcs = ["condition.cc"],
    hdrs = ["condition.h"],
    deps = [
        "//intrinsic/icon/proto/v1:condition_types_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/types:span",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "session",
    srcs = ["session.cc"],
    hdrs = ["session.h"],
    deps = [
        ":condition",
        ":stream",
        "//intrinsic/icon/common:id_types",
        "//intrinsic/icon/common:slot_part_map",
        "//intrinsic/icon/proto:concatenate_trajectory_protos",
        "//intrinsic/icon/proto:joint_space_cc_proto",
        "//intrinsic/icon/proto:streaming_output_cc_proto",
        "//intrinsic/icon/proto/v1:service_cc_grpc_proto",
        "//intrinsic/icon/proto/v1:service_cc_proto",
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "//intrinsic/icon/release:grpc_time_support",
        "//intrinsic/icon/release:source_location",
        "//intrinsic/logging/proto:context_cc_proto",
        "//intrinsic/util:atomic_sequence_num",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/grpc:channel_interface",
        "//intrinsic/util/status:status_builder",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_conversion_rpc",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/thread",
        "@abseil-cpp//absl/algorithm:container",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@abseil-cpp//absl/types:span",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_googleapis//google/rpc:status_cc_proto",
        "@or_tools//ortools/base:strong_int",
    ],
)

cc_library(
    name = "client",
    srcs = ["client.cc"],
    hdrs = ["client.h"],
    deps = [
        ":operational_status",
        ":robot_config",
        "//intrinsic/icon/common:part_properties",
        "//intrinsic/icon/common:slot_part_map",
        "//intrinsic/icon/control:logging_mode",
        "//intrinsic/icon/proto:part_status_cc_proto",
        "//intrinsic/icon/proto/v1:service_cc_grpc_proto",
        "//intrinsic/icon/proto/v1:service_cc_proto",
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "//intrinsic/icon/release:grpc_time_support",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/grpc:channel_interface",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/robot_payload",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@abseil-cpp//absl/types:span",
        "@com_github_grpc_grpc//:grpc++",
    ],
)

cc_library(
    name = "operational_status",
    srcs = ["operational_status.cc"],
    hdrs = ["operational_status.h"],
    deps = [
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "robot_config",
    srcs = ["robot_config.cc"],
    hdrs = ["robot_config.h"],
    deps = [
        "//intrinsic/icon/proto:generic_part_config_cc_proto",
        "//intrinsic/icon/proto/v1:service_cc_proto",
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/algorithm:container",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "stream",
    srcs = ["stream.cc"],
    hdrs = ["stream.h"],
    deps = [
        "//intrinsic/icon/common:id_types",
        "//intrinsic/icon/proto/v1:service_cc_grpc_proto",
        "//intrinsic/icon/proto/v1:service_cc_proto",
        "//intrinsic/icon/release:grpc_time_support",
        "//intrinsic/util/grpc:channel_interface",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_conversion_rpc",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
    ],
)
