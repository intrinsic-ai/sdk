# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "motion_planner_service_proto",
    srcs = ["motion_planner_service.proto"],
    deps = [
        ":compute_ik_proto",
        ":geometric_constraints_proto",
        ":motion_planner_config_proto",
        ":motion_specification_proto",
        ":robot_specification_proto",
        "//intrinsic/geometry/proto:transformed_geometry_storage_refs_proto",
        "//intrinsic/icon/proto:joint_space_proto",
        "//intrinsic/kinematics/types:joint_limits_proto",
        "//intrinsic/logging/proto:context_proto",
        "//intrinsic/math/proto:point_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/world/proto:collision_settings_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
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

py_proto_library(
    name = "motion_planning_limits_py_pb2",
    deps = [":motion_planning_limits_proto"],
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

cc_proto_library(
    name = "motion_planning_error_cc_proto",
    deps = [":motion_planning_error_proto"],
)

proto_library(
    name = "motion_planner_config_proto",
    srcs = ["motion_planner_config.proto"],
    deps = ["@com_google_protobuf//:duration_proto"],
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
    go_deps = [
        "//intrinsic/icon/proto:joint_space_go_proto",
        "//intrinsic/world/proto:collision_settings_go_proto",
    ],
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
    deps = [":motion_planner_config_proto"],
)

cc_proto_library(
    name = "motion_specification_cc_proto",
    deps = [":motion_specification_proto"],
)

py_proto_library(
    name = "motion_specification_py_pb2",
    deps = [":motion_specification_proto"],
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

proto_library(
    name = "motion_planning_v1_protos",
    deps = [
        ":geometric_constraints_proto",
    ],
)
