# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

proto_library(
    name = "collision_action_proto",
    srcs = ["collision_action.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "collision_action_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":collision_action_proto"],
)

py_proto_library(
    name = "collision_action_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":collision_action_proto"],
)

go_proto_library(
    name = "collision_action_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":collision_action_proto"],
)

proto_library(
    name = "collision_settings_proto",
    srcs = ["collision_settings.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":collision_action_proto",
        ":object_world_refs_proto",
        "//intrinsic/skills/proto:skill_parameter_metadata_proto",
    ],
)

cc_proto_library(
    name = "collision_settings_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":collision_settings_proto"],
)

py_proto_library(
    name = "collision_settings_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":collision_settings_proto"],
)

go_proto_library(
    name = "collision_settings_go_proto",
    go_deps = [
        ":collision_action_go_proto",
        ":object_world_refs_go_proto",
        "//intrinsic/skills/proto:skill_parameter_metadata_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":collision_settings_proto"],
)

proto_library(
    name = "geometry_component_proto",
    srcs = ["geometry_component.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/geometry/proto:geometry_storage_refs_proto",
        "//intrinsic/geometry/proto/v1:transformed_geometry_proto",
        "//intrinsic/math/proto:matrix_proto",
    ],
)

cc_proto_library(
    name = "geometry_component_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":geometry_component_proto"],
)

py_proto_library(
    name = "geometry_component_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":geometry_component_proto"],
)

go_proto_library(
    name = "geometry_component_go_proto",
    go_deps = [
        "//intrinsic/math/proto:matrix_go_proto",
        "//intrinsic/geometry/proto/v1:transformed_geometry_go_proto",
        "//intrinsic/geometry/proto:geometry_storage_refs_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":geometry_component_proto"],
)

proto_library(
    name = "kinematics_component_proto",
    srcs = ["kinematics_component.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/math/proto:vector3_proto",
    ],
)

cc_proto_library(
    name = "kinematics_component_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":kinematics_component_proto"],
)

go_proto_library(
    name = "kinematics_component_go_proto",
    go_deps = [
        "//intrinsic/math/proto:vector3_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":kinematics_component_proto"],
)

proto_library(
    name = "object_world_refs_proto",
    srcs = ["object_world_refs.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "object_world_refs_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":object_world_refs_proto"],
)

py_proto_library(
    name = "object_world_refs_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":object_world_refs_proto"],
)

go_proto_library(
    name = "object_world_refs_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":object_world_refs_proto"],
)

proto_library(
    name = "object_world_service_proto",
    srcs = ["object_world_service.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":collision_settings_proto",
        ":geometry_component_proto",
        ":gripper_component_proto",
        ":kinematics_component_proto",
        ":object_world_refs_proto",
        ":object_world_updates_proto",
        ":physics_component_proto",
        ":robot_payload_proto",
        ":sensor_component_proto",
        ":simulation_component_proto",
        "//intrinsic/icon/proto:cart_space_proto",
        "//intrinsic/kinematics/types:joint_limits_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/skills/proto:footprint_proto",
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

proto_library(
    name = "object_world_updates_proto",
    srcs = ["object_world_updates.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":collision_action_proto",
        ":collision_settings_proto",
        ":geometry_component_proto",
        ":object_world_refs_proto",
        ":physics_component_proto",
        ":robot_payload_proto",
        ":simulation_component_proto",
        "//intrinsic/icon/proto:cart_space_proto",
        "//intrinsic/kinematics/types:joint_limits_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/scene/proto/v1:scene_object_proto",
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

go_proto_library(
    name = "object_world_updates_go_proto",
    go_deps = [
        ":collision_action_go_proto",
        ":collision_settings_go_proto",
        ":geometry_component_go_proto",
        ":gripper_component_go_proto",
        ":robot_payload_go_proto",
        ":physics_component_go_proto",
        ":simulation_component_go_proto",
        ":object_world_refs_go_proto",
        "//intrinsic/icon/proto:cart_space_go_proto",
        "//intrinsic/kinematics/types:joint_limits_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":object_world_updates_proto"],
)

proto_library(
    name = "physics_component_proto",
    srcs = ["physics_component.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/math/proto:matrix_proto",
        "//intrinsic/math/proto:point_proto",
        "//intrinsic/math/proto:pose_proto",
    ],
)

py_proto_library(
    name = "physics_component_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":physics_component_proto"],
)

cc_proto_library(
    name = "physics_component_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":physics_component_proto"],
)

go_proto_library(
    name = "physics_component_go_proto",
    go_deps = [
        "//intrinsic/math/proto:matrix_go_proto",
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":physics_component_proto"],
)

proto_library(
    name = "simulation_component_proto",
    srcs = ["simulation_component.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

py_proto_library(
    name = "simulation_component_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":simulation_component_proto"],
)

cc_proto_library(
    name = "simulation_component_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":simulation_component_proto"],
)

go_proto_library(
    name = "simulation_component_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":simulation_component_proto"],
)

proto_library(
    name = "sensor_component_proto",
    srcs = ["sensor_component.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "sensor_component_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":sensor_component_proto"],
)

go_proto_library(
    name = "sensor_component_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["sensor_component_proto"],
)

proto_library(
    name = "gripper_component_proto",
    srcs = ["gripper_component.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "gripper_component_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":gripper_component_proto"],
)

go_proto_library(
    name = "gripper_component_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["gripper_component_proto"],
)

proto_library(
    name = "robot_payload_proto",
    srcs = ["robot_payload.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/math/proto:matrix_proto",
        "//intrinsic/math/proto:pose_proto",
    ],
)

py_proto_library(
    name = "robot_payload_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":robot_payload_proto"],
)

py_proto_library(
    name = "entity_search_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":entity_search_proto"],
)

go_proto_library(
    name = "robot_payload_go_proto",
    go_deps = [
        "//intrinsic/math/proto:matrix_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":robot_payload_proto"],
)

go_proto_library(
    name = "entity_search_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":entity_search_proto"],
)

cc_proto_library(
    name = "robot_payload_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":robot_payload_proto"],
)

proto_library(
    name = "world_protos",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":geometry_component_proto",
        ":kinematics_component_proto",
        ":object_world_refs_proto",
        ":object_world_service_proto",
        ":physics_component_proto",
        ":simulation_component_proto",
        "//intrinsic/world/proto:entity_search_proto",
    ],
)

proto_library(
    name = "entity_search_proto",
    srcs = ["entity_search.proto"],
    visibility = ["//intrinsic:public_api_users"],
)
