# Copyright 2023 Intrinsic Innovation LLC

""" Intrinsic package grasping.

This includes grasping related protos, grasp planning service and grasp skills."""

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_python//python:defs.bzl", "py_library")

cc_proto_library(
    name = "grasp_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":grasp_proto"],
)

proto_library(
    name = "grasp_proto",
    srcs = ["grasp.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto:field_metadata_proto",
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/hardware/gripper:gripper_proto",
        "//intrinsic/hardware/gripper/eoat:eoat_service_proto",
        "//intrinsic/icon/proto:joint_space_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/motion_planning/proto/v1:geometric_constraints_proto",
        "//intrinsic/perception/proto/v1:pose_estimate_in_root_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

py_proto_library(
    name = "grasp_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":grasp_proto"],
)

py_library(
    name = "grasp_planning_interfaces",
    srcs = ["grasp_planning_interfaces.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":grasp_py_pb2",
        "//intrinsic/icon/proto:joint_space_py_pb2",
        "//intrinsic/math/proto:pose_py_pb2",
        "//intrinsic/motion_planning:motion_planner_client_py",
        "//intrinsic/motion_planning/proto/v1:geometric_constraints_py_pb2",
        "//intrinsic/world/proto:object_world_refs_py_pb2",
        "//intrinsic/world/python:object_world_client",
    ],
)

proto_library(
    name = "grasp_execution_planner_params_proto",
    srcs = ["grasp_execution_planner_params.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/motion_planning/proto/v1:motion_planning_limits_proto",
        "//intrinsic/motion_planning/proto/v1:motion_specification_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

py_proto_library(
    name = "grasp_execution_planner_params_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":grasp_execution_planner_params_proto"],
)

cc_proto_library(
    name = "grasp_execution_planner_params_cc_proto",
    visibility = ["//visibility:private"],
    deps = [":grasp_execution_planner_params_proto"],
)

proto_library(
    name = "grasp_planner_params_proto",
    srcs = ["grasp_planner_params.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":grasp_adapter_params_proto",
        ":grasp_proto",
        ":grasp_ranker_params_proto",
        "//intrinsic/manipulation/proto:workspace_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

py_proto_library(
    name = "grasp_planner_params_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":grasp_planner_params_proto"],
)

cc_proto_library(
    name = "grasp_planner_params_cc_proto",
    visibility = ["//visibility:private"],
    deps = [":grasp_planner_params_proto"],
)

proto_library(
    name = "grasp_placement_constraint_proto",
    srcs = ["grasp_placement_constraint.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/math/proto:point_proto",
        "//intrinsic/world/proto:collision_settings_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

py_proto_library(
    name = "grasp_placement_constraint_py_pb2",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [":grasp_placement_constraint_proto"],
)

cc_proto_library(
    name = "grasp_placement_constraint_cc_proto",
    visibility = ["//visibility:private"],
    deps = [":grasp_placement_constraint_proto"],
)

proto_library(
    name = "grasp_ranker_params_proto",
    srcs = ["grasp_ranker_params.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":grasp_placement_constraint_proto",
        "//intrinsic/kinematics/types:joint_limits_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

py_proto_library(
    name = "grasp_ranker_params_py_pb2",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [":grasp_ranker_params_proto"],
)

cc_proto_library(
    name = "grasp_ranker_params_cc_proto",
    visibility = ["//visibility:private"],
    deps = [":grasp_ranker_params_proto"],
)

proto_library(
    name = "grasp_adapter_params_proto",
    srcs = ["grasp_adapter_params.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

py_proto_library(
    name = "grasp_adapter_params_py_pb2",
    visibility = ["//visibility:private"],
    deps = [":grasp_adapter_params_proto"],
)

cc_proto_library(
    name = "grasp_adapter_params_cc_proto",
    visibility = ["//visibility:private"],
    deps = [":grasp_adapter_params_proto"],
)

proto_library(
    name = "grasp_annotations_proto",
    srcs = ["grasp_annotations.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/math/proto:vector3_proto",
    ],
)

py_proto_library(
    name = "grasp_annotations_py_pb2",
    visibility = [
        "//intrinsic:internal_api_users",
        "//intrinsic/apps:__subpackages__",
    ],
    deps = [":grasp_annotations_proto"],
)

proto_library(
    name = "grasp_annotator_proto",
    srcs = ["grasp_annotator.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/geometry/proto:oriented_bounding_box_proto",
        "//intrinsic/geometry/proto:triangle_mesh_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/math/proto:vector2_proto",
        "//intrinsic/math/proto:vector3_proto",
    ],
)

py_proto_library(
    name = "grasp_annotator_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":grasp_annotator_proto"],
)

cc_proto_library(
    name = "grasp_annotator_cc_proto",
    visibility = ["//visibility:private"],
    deps = [":grasp_annotator_proto"],
)

py_library(
    name = "grasp_annotator_client",
    srcs = ["grasp_annotator_client.py"],
    deps = [
        ":grasp_annotations_py_pb2",
        ":grasp_annotator_py_pb2",
        "//intrinsic/geometry/proto:triangle_mesh_py_pb2",
        "//intrinsic/manipulation/service/grasp_annotator_service/v1:grasp_annotator_service_py_pb2",
        "//intrinsic/manipulation/service/grasp_annotator_service/v1:grasp_annotator_service_py_pb2_grpc",
        "//intrinsic/world/python:object_world_resources",
        requirement("grpcio"),
        "@com_google_absl_py//absl/logging",
        "@com_google_googleapis//google/rpc:code_py_proto",
    ],
)

proto_library(
    name = "schmalz_multi_axis_grasp_proto",
    srcs = ["schmalz_multi_axis_grasp.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/math/proto:pose_proto"],
)

py_proto_library(
    name = "schmalz_multi_axis_grasp_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":schmalz_multi_axis_grasp_proto"],
)

py_library(
    name = "multi_axis_auto_annotator_client",
    srcs = ["multi_axis_auto_annotator_client.py"],
    deps = [
        ":schmalz_multi_axis_grasp_py_pb2",
        "//intrinsic/manipulation/service:multi_axis_auto_annotator_service_py_pb2",
        "//intrinsic/manipulation/service:multi_axis_auto_annotator_service_py_pb2_grpc",
        requirement("grpcio"),
        "@com_google_absl_py//absl/logging",
    ],
)
