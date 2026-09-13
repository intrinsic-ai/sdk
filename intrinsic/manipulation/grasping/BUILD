# Copyright 2026 Intrinsic Innovation LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

""" Intrinsic package grasping.

This includes grasping related protos, grasp planning service and grasp skills."""

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@rules_python//python:defs.bzl", "py_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

py_library(
    name = "grasp_planning_interfaces",
    srcs = ["grasp_planning_interfaces.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/icon/proto:joint_space_py_pb2",
        "//intrinsic/manipulation/grasping:grasp_py_pb2",
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

py_library(
    name = "grasp_annotator_client",
    srcs = ["grasp_annotator_client.py"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/geometry/proto:triangle_mesh_py_pb2",
        "//intrinsic/manipulation/grasping:grasp_annotations_py_pb2",
        "//intrinsic/manipulation/grasping:grasp_annotator_py_pb2",
        "//intrinsic/manipulation/service/grasp_annotator_service/v1:grasp_annotator_service_py_pb2",
        "//intrinsic/manipulation/service/grasp_annotator_service/v1:grasp_annotator_service_py_pb2_grpc",
        "//intrinsic/world/python:object_world_resources",
        requirement("grpcio"),
        "@com_google_absl_py//absl/logging",
        "@com_google_googleapis//google/rpc:code_py_proto",
    ],
)

py_library(
    name = "multi_axis_auto_annotator_client",
    srcs = ["multi_axis_auto_annotator_client.py"],
    deps = [
        "//intrinsic/manipulation/grasping:schmalz_multi_axis_grasp_py_pb2",
        "//intrinsic/manipulation/service:multi_axis_auto_annotator_service_py_pb2",
        "//intrinsic/manipulation/service:multi_axis_auto_annotator_service_py_pb2_grpc",
        requirement("grpcio"),
        "@com_google_absl_py//absl/logging",
    ],
)

proto_library(
    name = "grasp_proto",
    srcs = ["grasp.proto"],
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

cc_proto_library(
    name = "grasp_cc_proto",
    deps = [":grasp_proto"],
)

go_proto_library(
    name = "grasp_go_proto",
    importpath = "intrinsic/manipulation/grasping/grasp_go_proto",
    protos = [":grasp_proto"],
    deps = [
        "//intrinsic/assets/proto:field_metadata_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/hardware/gpio/v1:signal_go_proto",
        "//intrinsic/hardware/gripper:gripper_go_proto",
        "//intrinsic/hardware/gripper/eoat:eoat_service_go_proto",
        "//intrinsic/hardware/gripper/eoat:gripper_config_go_proto",
        "//intrinsic/icon/proto:joint_space_go_proto",
        "//intrinsic/kinematics/types:dynamic_limits_check_mode_go_proto",
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/math/proto:quaternion_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
        "//intrinsic/motion_planning/proto/v1:geometric_constraints_go_proto",
        "//intrinsic/perception/proto/v1:pose_estimate_in_root_go_proto",
        "//intrinsic/skills/proto:skill_parameter_metadata_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
    ],
)

py_proto_library(
    name = "grasp_py_pb2",
    deps = [":grasp_proto"],
)

proto_library(
    name = "grasp_adapter_params_proto",
    srcs = ["grasp_adapter_params.proto"],
)

cc_proto_library(
    name = "grasp_adapter_params_cc_proto",
    deps = [":grasp_adapter_params_proto"],
)

go_proto_library(
    name = "grasp_adapter_params_go_proto",
    importpath = "intrinsic/manipulation/grasping/grasp_adapter_params_go_proto",
    protos = [":grasp_adapter_params_proto"],
    deps = [
        ":grasp_go_proto",
        ":grasp_ranker_params_go_proto",
        "//intrinsic/geometry/proto:axis_aligned_bounding_box_go_proto",
        "//intrinsic/geometry/proto:oriented_bounding_box_go_proto",
        "//intrinsic/hardware/gpio/v1:signal_go_proto",
        "//intrinsic/hardware/gripper:gripper_go_proto",
        "//intrinsic/hardware/gripper/eoat:eoat_service_go_proto",
        "//intrinsic/hardware/gripper/eoat:gripper_config_go_proto",
        "//intrinsic/icon/proto:cart_space_go_proto",
        "//intrinsic/icon/proto:joint_space_go_proto",
        "//intrinsic/kinematics/types:dynamic_limits_check_mode_go_proto",
        "//intrinsic/kinematics/types:joint_limits_go_proto",
        "//intrinsic/manipulation/proto:workspace_go_proto",
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/math/proto:quaternion_go_proto",
        "//intrinsic/math/proto:vector2_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
        "//intrinsic/motion_planning/proto/v1:geometric_constraints_go_proto",
        "//intrinsic/perception/proto/v1:pose_estimate_in_root_go_proto",
        "//intrinsic/skills/proto:skill_parameter_metadata_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
    ],
)

py_proto_library(
    name = "grasp_adapter_params_py_pb2",
    deps = [":grasp_adapter_params_proto"],
)

proto_library(
    name = "grasp_annotations_proto",
    srcs = ["grasp_annotations.proto"],
    deps = [
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/math/proto:vector3_proto",
    ],
)

cc_proto_library(
    name = "grasp_annotations_cc_proto",
    deps = [":grasp_annotations_proto"],
)

go_proto_library(
    name = "grasp_annotations_go_proto",
    importpath = "intrinsic/manipulation/grasping/grasp_annotations_go_proto",
    protos = [":grasp_annotations_proto"],
    deps = [
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/math/proto:quaternion_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
    ],
)

py_proto_library(
    name = "grasp_annotations_py_pb2",
    deps = [":grasp_annotations_proto"],
)

proto_library(
    name = "grasp_annotator_proto",
    srcs = ["grasp_annotator.proto"],
    deps = [
        "//intrinsic/geometry/proto:oriented_bounding_box_proto",
        "//intrinsic/geometry/proto:triangle_mesh_proto",
        "//intrinsic/geometry/proto/v1:triangle_mesh_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/math/proto:vector2_proto",
        "//intrinsic/math/proto:vector3_proto",
    ],
)

cc_proto_library(
    name = "grasp_annotator_cc_proto",
    deps = [":grasp_annotator_proto"],
)

go_proto_library(
    name = "grasp_annotator_go_proto",
    importpath = "intrinsic/manipulation/grasping/grasp_annotator_go_proto",
    protos = [":grasp_annotator_proto"],
    deps = [
        "//intrinsic/geometry/proto:axis_aligned_bounding_box_go_proto",
        "//intrinsic/geometry/proto:oriented_bounding_box_go_proto",
        "//intrinsic/geometry/proto:triangle_mesh_go_proto",
        "//intrinsic/geometry/proto/v1:triangle_mesh_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/math/proto:vector2_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
    ],
)

py_proto_library(
    name = "grasp_annotator_py_pb2",
    deps = [":grasp_annotator_proto"],
)

proto_library(
    name = "grasp_placement_constraint_proto",
    srcs = ["grasp_placement_constraint.proto"],
    deps = [
        "//intrinsic/math/proto:point_proto",
        "//intrinsic/world/proto:collision_settings_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

cc_proto_library(
    name = "grasp_placement_constraint_cc_proto",
    deps = [":grasp_placement_constraint_proto"],
)

go_proto_library(
    name = "grasp_placement_constraint_go_proto",
    importpath = "intrinsic/manipulation/grasping/grasp_placement_constraint_go_proto",
    protos = [":grasp_placement_constraint_proto"],
    deps = [
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/skills/proto:skill_parameter_metadata_go_proto",
        "//intrinsic/world/proto:collision_settings_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
    ],
)

py_proto_library(
    name = "grasp_placement_constraint_py_pb2",
    deps = [":grasp_placement_constraint_proto"],
)

proto_library(
    name = "grasp_ranker_params_proto",
    srcs = ["grasp_ranker_params.proto"],
    deps = [
        ":grasp_placement_constraint_proto",
        "//intrinsic/kinematics/types:joint_limits_proto",
        "//intrinsic/world/proto:collision_settings_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

cc_proto_library(
    name = "grasp_ranker_params_cc_proto",
    deps = [":grasp_ranker_params_proto"],
)

go_proto_library(
    name = "grasp_ranker_params_go_proto",
    importpath = "intrinsic/manipulation/grasping/grasp_ranker_params_go_proto",
    protos = [":grasp_ranker_params_proto"],
    deps = [
        ":grasp_placement_constraint_go_proto",
        "//intrinsic/icon/proto:cart_space_go_proto",
        "//intrinsic/kinematics/types:joint_limits_go_proto",
        "//intrinsic/world/proto:collision_settings_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
    ],
)

py_proto_library(
    name = "grasp_ranker_params_py_pb2",
    deps = [":grasp_ranker_params_proto"],
)

proto_library(
    name = "grasp_planner_params_proto",
    srcs = ["grasp_planner_params.proto"],
    deps = [
        ":grasp_adapter_params_proto",
        ":grasp_proto",
        ":grasp_ranker_params_proto",
        "//intrinsic/manipulation/proto:workspace_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

cc_proto_library(
    name = "grasp_planner_params_cc_proto",
    deps = [":grasp_planner_params_proto"],
)

go_proto_library(
    name = "grasp_planner_params_go_proto",
    importpath = "intrinsic/manipulation/grasping/grasp_planner_params_go_proto",
    protos = [":grasp_planner_params_proto"],
    deps = [
        ":grasp_adapter_params_go_proto",
        ":grasp_go_proto",
        ":grasp_ranker_params_go_proto",
        "//intrinsic/hardware/gpio/v1:signal_go_proto",
        "//intrinsic/hardware/gripper:gripper_go_proto",
        "//intrinsic/hardware/gripper/eoat:eoat_service_go_proto",
        "//intrinsic/hardware/gripper/eoat:gripper_config_go_proto",
        "//intrinsic/icon/proto:joint_space_go_proto",
        "//intrinsic/kinematics/types:dynamic_limits_check_mode_go_proto",
        "//intrinsic/manipulation/proto:workspace_go_proto",
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/math/proto:quaternion_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
        "//intrinsic/motion_planning/proto/v1:geometric_constraints_go_proto",
        "//intrinsic/perception/proto/v1:pose_estimate_in_root_go_proto",
        "//intrinsic/skills/proto:skill_parameter_metadata_go_proto",
        "//intrinsic/world/proto:collision_settings_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
    ],
)

py_proto_library(
    name = "grasp_planner_params_py_pb2",
    deps = [":grasp_planner_params_proto"],
)

proto_library(
    name = "schmalz_multi_axis_grasp_proto",
    srcs = ["schmalz_multi_axis_grasp.proto"],
    deps = ["//intrinsic/math/proto:pose_proto"],
)

cc_proto_library(
    name = "schmalz_multi_axis_grasp_cc_proto",
    deps = [":schmalz_multi_axis_grasp_proto"],
)

go_proto_library(
    name = "schmalz_multi_axis_grasp_go_proto",
    importpath = "intrinsic/manipulation/grasping/schmalz_multi_axis_grasp_go_proto",
    protos = [":schmalz_multi_axis_grasp_proto"],
    deps = [
        "//intrinsic/math/proto:pose_go_proto",
    ],
)

py_proto_library(
    name = "schmalz_multi_axis_grasp_py_pb2",
    deps = [":schmalz_multi_axis_grasp_proto"],
)
