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

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@intrinsic_js_rules//:defs.bzl", "js_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "motion_planner_service_proto",
    srcs = ["motion_planner_service.proto"],
    deps = [
        ":compute_ik_proto",
        ":geometric_constraints_proto",
        ":motion_planner_config_proto",
        ":motion_planning_proto",
        ":motion_specification_proto",
        ":robot_specification_proto",
        "//intrinsic/assets/proto:field_metadata_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_proto",
        "//intrinsic/geometry/proto:transformed_geometry_storage_refs_proto",
        "//intrinsic/icon/proto:joint_space_proto",
        "//intrinsic/kinematics/types:joint_limits_proto",
        "//intrinsic/logging/proto:context_proto",
        "//intrinsic/math/proto:point_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/world/proto:collision_checker_config_proto",
        "//intrinsic/world/proto:collision_settings_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_protobuf//:duration_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

cc_proto_library(
    name = "motion_planner_service_cc_proto",
    deps = [":motion_planner_service_proto"],
)

cc_grpc_library(
    name = "motion_planner_service_cc_grpc_proto",
    srcs = [":motion_planner_service_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":motion_planner_service_cc_proto"],
)

proto_library(
    name = "motion_blending_parameter_proto",
    srcs = ["motion_blending_parameter.proto"],
    deps = [
        "//intrinsic/skills/proto:skill_parameter_metadata_proto",
    ],
)

cc_proto_library(
    name = "motion_blending_parameter_cc_proto",
    deps = [":motion_blending_parameter_proto"],
)

js_library(
    name = "motion_blending_parameter_jspb_proto",
    deps = [":motion_blending_parameter_proto"],
)

go_proto_library(
    name = "motion_blending_parameter_go_proto",
    importpath = "intrinsic/motion_planning/proto/v1/motion_blending_parameter_go_proto",
    protos = [":motion_blending_parameter_proto"],
    deps = [
        "//intrinsic/skills/proto:skill_parameter_metadata_go_proto",
    ],
)

py_proto_library(
    name = "motion_blending_parameter_py_pb2",
    deps = [":motion_blending_parameter_proto"],
)

proto_library(
    name = "motion_planning_limits_proto",
    srcs = ["motion_planning_limits.proto"],
    deps = [
        "//intrinsic/kinematics/types:joint_limits_proto",
        "//intrinsic/skills/proto:skill_parameter_metadata_proto",
    ],
)

cc_proto_library(
    name = "motion_planning_limits_cc_proto",
    deps = [":motion_planning_limits_proto"],
)

js_library(
    name = "motion_planning_limits_jspb_proto",
    deps = [":motion_planning_limits_proto"],
)

go_proto_library(
    name = "motion_planning_limits_go_proto",
    importpath = "intrinsic/motion_planning/proto/v1/motion_planning_limits_go_proto",
    protos = [":motion_planning_limits_proto"],
    deps = [
        "//intrinsic/kinematics/types:joint_limits_go_proto",
        "//intrinsic/skills/proto:skill_parameter_metadata_go_proto",
    ],
)

py_proto_library(
    name = "motion_planning_limits_py_pb2",
    deps = [":motion_planning_limits_proto"],
)

proto_library(
    name = "motion_planning_proto",
    srcs = ["motion_planning.proto"],
    deps = ["//intrinsic/icon/proto:joint_space_proto"],
)

py_proto_library(
    name = "motion_planning_py_pb2",
    deps = [":motion_planning_proto"],
)

cc_proto_library(
    name = "motion_planning_cc_proto",
    deps = [":motion_planning_proto"],
)

js_library(
    name = "motion_planning_jspb_proto",
    deps = [":motion_planning_proto"],
)

go_proto_library(
    name = "robot_specification_go_proto",
    importpath = "intrinsic/motion_planning/proto/v1/robot_specification_go_proto",
    protos = [":robot_specification_proto"],
    deps = [
        "//intrinsic/geometry/proto:geometry_storage_refs_go_proto",
        "//intrinsic/icon/proto:cart_space_go_proto",
        "//intrinsic/icon/proto:joint_space_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
    ],
)

go_proto_library(
    name = "motion_planning_go_proto",
    importpath = "intrinsic/motion_planning/proto/v1/motion_planning_go_proto",
    protos = [":motion_planning_proto"],
    deps = [
        "//intrinsic/icon/proto:joint_space_go_proto",
    ],
)

go_proto_library(
    name = "motion_planner_service_go_proto",
    compilers = [
        "//bazel:go_gen_grpc_gateway",
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/motion_planning/proto/v1/motion_planner_service_go_proto",
    protos = [":motion_planner_service_proto"],
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [
        ":compute_ik_go_proto",
        ":motion_planner_config_go_proto",
        ":motion_planning_go_proto",
        ":robot_specification_go_proto",
        "//intrinsic/assets/proto:field_metadata_go_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_go_proto",
        "//intrinsic/geometry/proto:geometry_storage_refs_go_proto",
        "//intrinsic/geometry/proto:transformed_geometry_storage_refs_go_proto",
        "//intrinsic/icon/proto:cart_space_go_proto",
        "//intrinsic/icon/proto:joint_space_go_proto",
        "//intrinsic/kinematics/types:joint_limits_go_proto",
        "//intrinsic/logging/proto:context_go_proto",
        "//intrinsic/math/proto:matrix_go_proto",
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
        "//intrinsic/motion_planning/proto/v1:geometric_constraints_go_proto",
        "//intrinsic/motion_planning/proto/v1:motion_planning_error_go_proto",
        "//intrinsic/motion_planning/proto/v1:motion_specification_go_proto",
        "//intrinsic/world/proto:collision_action_go_proto",
        "//intrinsic/world/proto:collision_checker_config_go_proto",
        "//intrinsic/world/proto:collision_settings_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)

js_library(
    name = "motion_planner_service_jspb_proto",
    deps = [":motion_planner_service_proto"],
)

py_proto_library(
    name = "motion_planner_service_py_pb2",
    deps = [":motion_planner_service_proto"],
)

py_grpc_library(
    name = "motion_planner_service_py_pb2_grpc",
    srcs = [":motion_planner_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":motion_planner_service_py_pb2"],
)

proto_library(
    name = "motion_specification_proto",
    srcs = ["motion_specification.proto"],
    deps = [
        ":geometric_constraints_proto",
        ":motion_blending_parameter_proto",
        ":motion_planning_limits_proto",
        "//intrinsic/world/proto:collision_settings_proto",
    ],
)

proto_library(
    name = "motion_planning_error_proto",
    srcs = ["motion_planning_error.proto"],
    deps = [
        "//intrinsic/icon/proto:joint_space_proto",
        "//intrinsic/kinematics/types:joint_limits_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/motion_planning/proto/v1:geometric_constraints_proto",
        "//intrinsic/world/proto:collision_settings_proto",
    ],
)

py_proto_library(
    name = "motion_planning_error_py_pb2",
    deps = [":motion_planning_error_proto"],
)

cc_proto_library(
    name = "motion_planning_error_cc_proto",
    deps = [":motion_planning_error_proto"],
)

js_library(
    name = "motion_planning_error_jspb_proto",
    deps = [":motion_planning_error_proto"],
)

proto_library(
    name = "motion_planner_config_proto",
    srcs = ["motion_planner_config.proto"],
    deps = [
        "//intrinsic/world/proto:collision_checker_config_proto",
        "@com_google_protobuf//:duration_proto",
    ],
)

proto_library(
    name = "robot_specification_proto",
    srcs = ["robot_specification.proto"],
    deps = [
        "//intrinsic/icon/proto:joint_space_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

proto_library(
    name = "compute_ik_proto",
    srcs = ["compute_ik.proto"],
    deps = [
        "//intrinsic/icon/proto:joint_space_proto",
        "//intrinsic/world/proto:collision_settings_proto",
    ],
)

cc_proto_library(
    name = "compute_ik_cc_proto",
    deps = [":compute_ik_proto"],
)

go_proto_library(
    name = "compute_ik_go_proto",
    importpath = "intrinsic/motion_planning/proto/v1/compute_ik_go_proto",
    protos = [":compute_ik_proto"],
    deps = [
        "//intrinsic/icon/proto:joint_space_go_proto",
        "//intrinsic/world/proto:collision_settings_go_proto",
    ],
)

js_library(
    name = "compute_ik_jspb_proto",
    deps = [":compute_ik_proto"],
)

py_proto_library(
    name = "compute_ik_py_pb2",
    deps = [":compute_ik_proto"],
)

cc_proto_library(
    name = "robot_specification_cc_proto",
    deps = [":robot_specification_proto"],
)

py_proto_library(
    name = "robot_specification_py_pb2",
    deps = [":robot_specification_proto"],
)

js_library(
    name = "robot_specification_jspb_proto",
    deps = [":robot_specification_proto"],
)

cc_proto_library(
    name = "motion_planner_config_cc_proto",
    deps = [":motion_planner_config_proto"],
)

py_proto_library(
    name = "motion_planner_config_py_pb2",
    deps = [":motion_planner_config_proto"],
)

go_proto_library(
    name = "motion_planner_config_go_proto",
    importpath = "intrinsic/motion_planning/proto/v1/motion_planner_config_go_proto",
    protos = [":motion_planner_config_proto"],
    deps = [
        "//intrinsic/world/proto:collision_checker_config_go_proto",
    ],
)

cc_proto_library(
    name = "motion_specification_cc_proto",
    deps = [":motion_specification_proto"],
)

js_library(
    name = "motion_specification_jspb_proto",
    deps = [":motion_specification_proto"],
)

py_proto_library(
    name = "motion_specification_py_pb2",
    deps = [":motion_specification_proto"],
)

go_proto_library(
    name = "motion_specification_go_proto",
    importpath = "intrinsic/motion_planning/proto/v1/motion_specification_go_proto",
    protos = [":motion_specification_proto"],
    deps = [
        ":motion_blending_parameter_go_proto",
        ":motion_planning_limits_go_proto",
        "//intrinsic/motion_planning/proto/v1:geometric_constraints_go_proto",
        "//intrinsic/world/proto:collision_settings_go_proto",
    ],
)

go_proto_library(
    name = "motion_planning_error_go_proto",
    importpath = "intrinsic/motion_planning/proto/v1/motion_planning_error_go_proto",
    protos = [":motion_planning_error_proto"],
    deps = [
        "//intrinsic/icon/proto:joint_space_go_proto",
        "//intrinsic/kinematics/types:joint_limits_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/motion_planning/proto/v1:geometric_constraints_go_proto",
        "//intrinsic/world/proto:collision_settings_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
    ],
)

proto_library(
    name = "geometric_constraints_proto",
    srcs = ["geometric_constraints.proto"],
    deps = [
        "//intrinsic/icon/proto:joint_space_proto",
        "//intrinsic/math/proto:point_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/math/proto:quaternion_proto",
        "//intrinsic/math/proto:vector3_proto",
        "//intrinsic/skills/proto:skill_parameter_metadata_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

cc_proto_library(
    name = "geometric_constraints_cc_proto",
    deps = [":geometric_constraints_proto"],
)

py_proto_library(
    name = "geometric_constraints_py_pb2",
    deps = [":geometric_constraints_proto"],
)

js_library(
    name = "geometric_constraints_jspb_proto",
    deps = [":geometric_constraints_proto"],
)

go_proto_library(
    name = "geometric_constraints_go_proto",
    importpath = "intrinsic/motion_planning/proto/v1/geometric_constraints_go_proto",
    protos = [":geometric_constraints_proto"],
    deps = [
        "//intrinsic/icon/proto:joint_space_go_proto",
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/math/proto:quaternion_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
        "//intrinsic/skills/proto:skill_parameter_metadata_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
    ],
)

proto_library(
    name = "motion_planning_v1_protos",
    deps = [
        ":geometric_constraints_proto",
    ],
)
