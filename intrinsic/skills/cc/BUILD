# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "client_common",
    hdrs = ["client_common.h"],
    deps = [
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "equipment_pack",
    srcs = ["equipment_pack.cc"],
    hdrs = ["equipment_pack.h"],
    deps = [
        "//intrinsic/resources/proto:resource_handle_cc_proto",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/skills/proto:skill_service_cc_proto",
        "//intrinsic/skills/proto:skills_cc_proto",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "skill_interface",
    hdrs = ["skill_interface.h"],
    deps = [
        ":execute_context",
        ":execute_request",
        ":get_footprint_context",
        ":get_footprint_request",
        ":preview_context",
        ":preview_request",
        "//intrinsic/logging/proto:context_cc_proto",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/skills/proto:footprint_cc_proto",
        "//intrinsic/skills/proto:skill_service_cc_proto",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "skill_interface_utils",
    srcs = ["skill_interface_utils.cc"],
    hdrs = ["skill_interface_utils.h"],
    deps = [
        ":equipment_pack",
        ":skill_interface",
        "//intrinsic/skills/internal:execute_context_view",
        "//intrinsic/skills/proto:prediction_cc_proto",
        "//intrinsic/skills/proto:skill_service_cc_proto",
        "//intrinsic/util:proto_time",
        "//intrinsic/util/proto:any",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/proto:object_world_updates_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:str_format",
        "@abseil-cpp//absl/time",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "get_footprint_context",
    hdrs = ["get_footprint_context.h"],
    deps = [
        "//intrinsic/motion_planning:motion_planner_client",
        "//intrinsic/skills/proto:footprint_cc_proto",
        "//intrinsic/world/objects:frame",
        "//intrinsic/world/objects:kinematic_object",
        "//intrinsic/world/objects:object_world_client",
        "//intrinsic/world/objects:world_object",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "get_footprint_request",
    hdrs = ["get_footprint_request.h"],
    deps = [
        "//intrinsic/util/proto:any",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "execute_context",
    hdrs = ["execute_context.h"],
    deps = [
        ":equipment_pack",
        ":skill_canceller",
        ":skill_logging_context",
        "//intrinsic/logging/proto:context_cc_proto",
        "//intrinsic/motion_planning:motion_planner_client",
        "//intrinsic/skills/internal:skill_registry_client_interface",
        "//intrinsic/skills/proto:skill_service_cc_proto",
        "//intrinsic/world/objects:object_world_client",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@com_github_grpc_grpc//:grpc++",
    ],
)

cc_library(
    name = "execute_request",
    hdrs = ["execute_request.h"],
    deps = [
        "//intrinsic/util/proto:any",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "preview_context",
    hdrs = ["preview_context.h"],
    deps = [
        ":equipment_pack",
        ":skill_canceller",
        ":skill_logging_context",
        "//intrinsic/logging/proto:context_cc_proto",
        "//intrinsic/motion_planning:motion_planner_client",
        "//intrinsic/world/objects:frame",
        "//intrinsic/world/objects:kinematic_object",
        "//intrinsic/world/objects:object_world_client",
        "//intrinsic/world/objects:world_object",
        "//intrinsic/world/proto:object_world_updates_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "preview_request",
    hdrs = ["preview_request.h"],
    deps = [
        "//intrinsic/util/proto:any",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
        "@com_google_protobuf//:protobuf",
    ],
)

cc_library(
    name = "skill_logging_context",
    hdrs = ["skill_logging_context.h"],
    deps = [
        "//intrinsic/logging/proto:context_cc_proto",
        "@abseil-cpp//absl/strings:string_view",
    ],
)

cc_library(
    name = "skill_canceller",
    srcs = ["skill_canceller.cc"],
    hdrs = ["skill_canceller.h"],
    deps = [
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/functional:any_invocable",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "skill_utils",
    srcs = ["skill_utils.cc"],
    hdrs = ["skill_utils.h"],
    deps = [
        "//intrinsic/resources/proto:resource_handle_cc_proto",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/skills/proto:skills_cc_proto",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:str_format",
    ],
)

cc_test(
    name = "skill_interface_utils_test",
    srcs = ["skill_interface_utils_test.cc"],
    deps = [
        ":skill_interface",
        ":skill_interface_utils",
        "//intrinsic/skills/internal:preview_context_impl",
        "//intrinsic/skills/proto:prediction_cc_proto",
        "//intrinsic/skills/proto:skill_service_cc_proto",
        "//intrinsic/skills/testing:echo_skill_cc_proto",
        "//intrinsic/skills/testing:skill_test_utils_cc",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
        "@com_google_protobuf//:protobuf",
    ],
)
