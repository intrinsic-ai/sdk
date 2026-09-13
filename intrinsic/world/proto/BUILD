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
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

js_library(
    name = "collision_action_jspb_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = ["//intrinsic/world/proto:collision_action_proto"],
)

js_library(
    name = "collision_settings_jspb_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = ["//intrinsic/world/proto:collision_settings_proto"],
)

js_library(
    name = "geometry_component_jspb_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = ["//intrinsic/world/proto:geometry_component_proto"],
)

js_library(
    name = "kinematics_component_jspb_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = ["//intrinsic/world/proto:kinematics_component_proto"],
)

js_library(
    name = "object_world_refs_jspb_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = ["//intrinsic/world/proto:object_world_refs_proto"],
)

proto_library(
    name = "object_world_service_proto",
    srcs = ["object_world_service.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":object_world_updates_proto",
        "//intrinsic/icon/proto:cart_space_proto",
        "//intrinsic/kinematics/types:joint_limits_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/skills/proto:footprint_proto",
        "//intrinsic/world/proto:collision_settings_proto",
        "//intrinsic/world/proto:geometry_component_proto",
        "//intrinsic/world/proto:kinematics_component_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
        "//intrinsic/world/proto:physics_component_proto",
        "//intrinsic/world/proto:robot_payload_proto",
        "//intrinsic/world/proto:sensor_component_proto",
        "//intrinsic/world/proto:simulation_component_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/rpc:status_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

cc_proto_library(
    name = "object_world_service_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":object_world_service_proto"],
)

cc_grpc_library(
    name = "object_world_service_cc_grpc_proto",
    srcs = [":object_world_service_proto"],
    generate_mocks = True,
    grpc_only = True,
    visibility = ["//intrinsic:public_api_users"],
    deps = [":object_world_service_cc_proto"],
)

py_proto_library(
    name = "object_world_service_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":object_world_service_proto"],
)

py_grpc_library(
    name = "object_world_service_py_pb2_grpc",
    srcs = [":object_world_service_proto"],
    grpc_library = requirement("grpcio"),
    visibility = ["//intrinsic:public_api_users"],
    deps = [":object_world_service_py_pb2"],
)

js_library(
    name = "object_world_service_jspb_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [":object_world_service_proto"],
)

proto_library(
    name = "object_world_updates_proto",
    srcs = ["object_world_updates.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/icon/proto:cart_space_proto",
        "//intrinsic/kinematics/types:joint_limits_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/scene/proto/v1:scene_object_proto",
        "//intrinsic/world/proto:collision_action_proto",
        "//intrinsic/world/proto:collision_settings_proto",
        "//intrinsic/world/proto:geometry_component_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
        "//intrinsic/world/proto:physics_component_proto",
        "//intrinsic/world/proto:robot_payload_proto",
        "//intrinsic/world/proto:simulation_component_proto",
        "//intrinsic/world/proto:user_data_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:wrappers_proto",
    ],
)

cc_proto_library(
    name = "object_world_updates_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":object_world_updates_proto"],
)

py_proto_library(
    name = "object_world_updates_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":object_world_updates_proto"],
)

js_library(
    name = "object_world_updates_jspb_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [":object_world_updates_proto"],
)

go_proto_library(
    name = "object_world_updates_go_proto",
    importpath = "intrinsic/world/proto/object_world_updates_go_proto",
    protos = [":object_world_updates_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/icon/proto:cart_space_go_proto",
        "//intrinsic/kinematics/types:joint_limits_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "//intrinsic/world/proto:collision_action_go_proto",
        "//intrinsic/world/proto:collision_settings_go_proto",
        "//intrinsic/world/proto:geometry_component_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
        "//intrinsic/world/proto:physics_component_go_proto",
        "//intrinsic/world/proto:robot_payload_go_proto",
        "//intrinsic/world/proto:simulation_component_go_proto",
        "//intrinsic/world/proto:user_data_go_proto",
    ],
)

js_library(
    name = "physics_component_jspb_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = ["//intrinsic/world/proto:physics_component_proto"],
)

js_library(
    name = "simulation_component_jspb_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = ["//intrinsic/world/proto:simulation_component_proto"],
)

js_library(
    name = "sensor_component_jspb_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = ["//intrinsic/world/proto:sensor_component_proto"],
)

js_library(
    name = "robot_payload_jspb_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = ["//intrinsic/world/proto:robot_payload_proto"],
)

js_library(
    name = "user_data_jspb_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = ["//intrinsic/world/proto:user_data_proto"],
)

proto_library(
    name = "world_protos",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":object_world_service_proto",
        "//intrinsic/world/proto:entity_search_proto",
        "//intrinsic/world/proto:geometry_component_proto",
        "//intrinsic/world/proto:kinematics_component_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
        "//intrinsic/world/proto:physics_component_proto",
        "//intrinsic/world/proto:simulation_component_proto",
    ],
)

js_library(
    name = "collision_checker_config_jspb_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = ["//intrinsic/world/proto:collision_checker_config_proto"],
)

proto_library(
    name = "collision_action_proto",
    srcs = ["collision_action.proto"],
)

cc_proto_library(
    name = "collision_action_cc_proto",
    deps = [":collision_action_proto"],
)

go_proto_library(
    name = "collision_action_go_proto",
    importpath = "intrinsic/world/proto/collision_action_go_proto",
    protos = [":collision_action_proto"],
)

py_proto_library(
    name = "collision_action_py_pb2",
    deps = [":collision_action_proto"],
)

proto_library(
    name = "collision_checker_config_proto",
    srcs = ["collision_checker_config.proto"],
)

cc_proto_library(
    name = "collision_checker_config_cc_proto",
    deps = [":collision_checker_config_proto"],
)

go_proto_library(
    name = "collision_checker_config_go_proto",
    importpath = "intrinsic/world/proto/collision_checker_config_go_proto",
    protos = [":collision_checker_config_proto"],
)

py_proto_library(
    name = "collision_checker_config_py_pb2",
    deps = [":collision_checker_config_proto"],
)

proto_library(
    name = "object_world_refs_proto",
    srcs = ["object_world_refs.proto"],
)

cc_proto_library(
    name = "object_world_refs_cc_proto",
    deps = [":object_world_refs_proto"],
)

go_proto_library(
    name = "object_world_refs_go_proto",
    importpath = "intrinsic/world/proto/object_world_refs_go_proto",
    protos = [":object_world_refs_proto"],
)

py_proto_library(
    name = "object_world_refs_py_pb2",
    deps = [":object_world_refs_proto"],
)

proto_library(
    name = "collision_settings_proto",
    srcs = ["collision_settings.proto"],
    deps = [
        ":collision_action_proto",
        ":object_world_refs_proto",
        "//intrinsic/skills/proto:skill_parameter_metadata_proto",
    ],
)

cc_proto_library(
    name = "collision_settings_cc_proto",
    deps = [":collision_settings_proto"],
)

go_proto_library(
    name = "collision_settings_go_proto",
    importpath = "intrinsic/world/proto/collision_settings_go_proto",
    protos = [":collision_settings_proto"],
    deps = [
        ":collision_action_go_proto",
        ":object_world_refs_go_proto",
        "//intrinsic/skills/proto:skill_parameter_metadata_go_proto",
    ],
)

py_proto_library(
    name = "collision_settings_py_pb2",
    deps = [":collision_settings_proto"],
)

proto_library(
    name = "entity_search_proto",
    srcs = ["entity_search.proto"],
)

cc_proto_library(
    name = "entity_search_cc_proto",
    deps = [":entity_search_proto"],
)

go_proto_library(
    name = "entity_search_go_proto",
    protos = [":entity_search_proto"],
)

py_proto_library(
    name = "entity_search_py_pb2",
    deps = [":entity_search_proto"],
)

proto_library(
    name = "geometry_component_proto",
    srcs = ["geometry_component.proto"],
    deps = [
        "//intrinsic/geometry/proto:geometry_storage_refs_proto",
        "//intrinsic/geometry/proto:lazy_exact_geometry_proto",
        "//intrinsic/geometry/proto/v1:transformed_geometry_proto",
        "//intrinsic/math/proto:matrix_proto",
    ],
)

cc_proto_library(
    name = "geometry_component_cc_proto",
    deps = [":geometry_component_proto"],
)

go_proto_library(
    name = "geometry_component_go_proto",
    importpath = "intrinsic/world/proto/geometry_component_go_proto",
    protos = [":geometry_component_proto"],
    deps = [
        "//intrinsic/geometry/proto:geometry_storage_refs_go_proto",
        "//intrinsic/geometry/proto/v1:transformed_geometry_go_proto",
        "//intrinsic/math/proto:matrix_go_proto",
    ],
)

py_proto_library(
    name = "geometry_component_py_pb2",
    deps = [":geometry_component_proto"],
)

proto_library(
    name = "kinematics_component_proto",
    srcs = ["kinematics_component.proto"],
    deps = [
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/math/proto:vector3_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

cc_proto_library(
    name = "kinematics_component_cc_proto",
    deps = [":kinematics_component_proto"],
)

go_proto_library(
    name = "kinematics_component_go_proto",
    importpath = "intrinsic/world/proto/kinematics_component_go_proto",
    protos = [":kinematics_component_proto"],
    deps = [
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
    ],
)

proto_library(
    name = "physics_component_proto",
    srcs = ["physics_component.proto"],
    deps = [
        "//intrinsic/math/proto:matrix_proto",
        "//intrinsic/math/proto:point_proto",
        "//intrinsic/math/proto:pose_proto",
    ],
)

cc_proto_library(
    name = "physics_component_cc_proto",
    deps = [":physics_component_proto"],
)

go_proto_library(
    name = "physics_component_go_proto",
    importpath = "intrinsic/world/proto/physics_component_go_proto",
    protos = [":physics_component_proto"],
    deps = [
        "//intrinsic/math/proto:matrix_go_proto",
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
)

py_proto_library(
    name = "physics_component_py_pb2",
    deps = [":physics_component_proto"],
)

proto_library(
    name = "robot_payload_proto",
    srcs = ["robot_payload.proto"],
    deps = [
        "//intrinsic/math/proto:matrix_proto",
        "//intrinsic/math/proto:pose_proto",
    ],
)

cc_proto_library(
    name = "robot_payload_cc_proto",
    deps = [":robot_payload_proto"],
)

go_proto_library(
    name = "robot_payload_go_proto",
    importpath = "intrinsic/world/proto/robot_payload_go_proto",
    protos = [":robot_payload_proto"],
    deps = [
        "//intrinsic/math/proto:matrix_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
)

py_proto_library(
    name = "robot_payload_py_pb2",
    deps = [":robot_payload_proto"],
)

proto_library(
    name = "sensor_component_proto",
    srcs = ["sensor_component.proto"],
)

cc_proto_library(
    name = "sensor_component_cc_proto",
    deps = [":sensor_component_proto"],
)

go_proto_library(
    name = "sensor_component_go_proto",
    importpath = "intrinsic/world/proto/sensor_component_go_proto",
    protos = ["sensor_component_proto"],
)

proto_library(
    name = "simulation_component_proto",
    srcs = ["simulation_component.proto"],
)

cc_proto_library(
    name = "simulation_component_cc_proto",
    deps = [":simulation_component_proto"],
)

go_proto_library(
    name = "simulation_component_go_proto",
    importpath = "intrinsic/world/proto/simulation_component_go_proto",
    protos = [":simulation_component_proto"],
)

py_proto_library(
    name = "simulation_component_py_pb2",
    deps = [":simulation_component_proto"],
)

proto_library(
    name = "user_data_proto",
    srcs = ["user_data.proto"],
    deps = ["@com_google_protobuf//:any_proto"],
)

cc_proto_library(
    name = "user_data_cc_proto",
    deps = [":user_data_proto"],
)

go_proto_library(
    name = "user_data_go_proto",
    importpath = "intrinsic/world/proto/user_data_go_proto",
    protos = [":user_data_proto"],
    deps = [
        "//intrinsic/math/proto:matrix_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
)

py_proto_library(
    name = "user_data_py_pb2",
    deps = [":user_data_proto"],
)
