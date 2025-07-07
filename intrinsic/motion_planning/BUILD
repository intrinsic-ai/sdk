# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_cc//cc:defs.bzl", "cc_library")
load("@rules_python//python:defs.bzl", "py_library")

cc_library(
    name = "motion_planner_client",
    srcs = ["motion_planner_client.cc"],
    hdrs = ["motion_planner_client.h"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":conversions",
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/proto:transformed_geometry_storage_refs_cc_proto",
        "//intrinsic/icon/proto:joint_space_cc_proto",
        "//intrinsic/logging/proto:context_cc_proto",
        "//intrinsic/math:pose3",
        "//intrinsic/math:proto_conversion",
        "//intrinsic/motion_planning/proto:motion_target_cc_proto",
        "//intrinsic/motion_planning/proto/v1:motion_planner_config_cc_proto",
        "//intrinsic/motion_planning/proto/v1:motion_planner_service_cc_grpc_proto",
        "//intrinsic/motion_planning/proto/v1:motion_planner_service_cc_proto",
        "//intrinsic/motion_planning/proto/v1:motion_specification_cc_proto",
        "//intrinsic/util:eigen",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/objects:kinematic_object",
        "//intrinsic/world/objects:transform_node",
        "//intrinsic/world/proto:collision_settings_cc_proto",
        "//intrinsic/world/proto:object_world_refs_cc_proto",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_protobuf//:protobuf",
    ],
)

py_library(
    name = "motion_planner_client_py",
    srcs = ["motion_planner_client.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/icon/proto:joint_space_py_pb2",
        "//intrinsic/math/python:data_types",
        "//intrinsic/math/python:proto_conversion",
        "//intrinsic/motion_planning/proto:motion_target_py_pb2",
        "//intrinsic/motion_planning/proto/v1:geometric_constraints_py_pb2",
        "//intrinsic/motion_planning/proto/v1:motion_planner_config_py_pb2",
        "//intrinsic/motion_planning/proto/v1:motion_planner_service_py_pb2",
        "//intrinsic/motion_planning/proto/v1:motion_planner_service_py_pb2_grpc",
        "//intrinsic/motion_planning/proto/v1:motion_specification_py_pb2",
        "//intrinsic/motion_planning/proto/v1:robot_specification_py_pb2",
        "//intrinsic/world/proto:collision_settings_py_pb2",
        "//intrinsic/world/python:object_world_ids",
        requirement("grpcio"),
        "@com_google_protobuf//:protobuf_python",
    ],
)

cc_library(
    name = "conversions",
    srcs = ["conversions.cc"],
    hdrs = ["conversions.h"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/icon/proto:joint_space_cc_proto",
        "//intrinsic/util:eigen",
        "@com_google_protobuf//:protobuf_lite",
    ],
)
