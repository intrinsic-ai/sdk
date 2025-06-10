# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "collision_rules_proto",
    srcs = ["collision_rules.proto"],
)

cc_proto_library(
    name = "collision_rules_cc_proto",
    deps = [":collision_rules_proto"],
)

go_proto_library(
    name = "collision_rules_go_proto",
    deps = [":collision_rules_proto"],
)

py_proto_library(
    name = "collision_rules_py_pb2",
    visibility = [
        "//visibility:public",
    ],
    deps = [":collision_rules_proto"],
)

proto_library(
    name = "entity_proto",
    srcs = ["entity.proto"],
    deps = [
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/world/proto:geometry_component_proto",
        "//intrinsic/world/proto:kinematics_component_proto",
        "//intrinsic/world/proto:physics_component_proto",
        "//intrinsic/world/proto:sensor_component_proto",
    ],
)

cc_proto_library(
    name = "entity_cc_proto",
    deps = [
        ":entity_proto",
    ],
)

go_proto_library(
    name = "entity_go_proto",
    go_deps = [
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/world/proto:geometry_component_go_proto",
        "//intrinsic/world/proto:kinematics_component_go_proto",
        "//intrinsic/world/proto:physics_component_go_proto",
        "//intrinsic/world/proto:sensor_component_go_proto",
    ],
    deps = [":entity_proto"],
)

py_proto_library(
    name = "entity_py_pb2",
    visibility = [
        "//visibility:public",
    ],
    deps = [":entity_proto"],
)

proto_library(
    name = "object_properties_proto",
    srcs = ["object_properties.proto"],
)

cc_proto_library(
    name = "object_properties_cc_proto",
    deps = [":object_properties_proto"],
)

go_proto_library(
    name = "object_properties_go_proto",
    deps = [":object_properties_proto"],
)

py_proto_library(
    name = "object_properties_py_pb2",
    visibility = [
        "//visibility:public",
    ],
    deps = [":object_properties_proto"],
)

proto_library(
    name = "scene_object_proto",
    srcs = ["scene_object.proto"],
    deps = [
        ":collision_rules_proto",
        ":entity_proto",
        ":object_properties_proto",
        ":simulation_spec_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "scene_object_cc_proto",
    deps = [":scene_object_proto"],
)

go_proto_library(
    name = "scene_object_go_proto",
    go_deps = [
        ":collision_rules_go_proto",
        ":entity_go_proto",
        ":object_properties_go_proto",
        ":simulation_spec_go_proto",
    ],
    deps = [":scene_object_proto"],
)

py_proto_library(
    name = "scene_object_py_pb2",
    visibility = [
        "//visibility:public",
    ],
    deps = [":scene_object_proto"],
)

proto_library(
    name = "scene_object_import_proto",
    srcs = ["scene_object_import.proto"],
    deps = [
        ":scene_object_proto",
        "//intrinsic/geometry/proto/v1:material_proto",
        "//intrinsic/math/proto:quaternion_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "scene_object_import_cc_proto",
    deps = [":scene_object_import_proto"],
)

go_proto_library(
    name = "scene_object_import_go_proto",
    go_deps = [
        ":scene_object_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "//intrinsic/geometry/proto/v1:material_go_proto",
        "//intrinsic/math/proto:quaternion_go_proto",
    ],
    deps = [":scene_object_import_proto"],
)

py_proto_library(
    name = "scene_object_import_py_pb2",
    deps = [":scene_object_import_proto"],
)

cc_grpc_library(
    name = "scene_object_import_cc_grpc_proto",
    srcs = [":scene_object_import_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":scene_object_import_cc_proto"],
)

go_grpc_library(
    name = "scene_object_import_go_grpc_proto",
    srcs = [":scene_object_import_proto"],
    deps = [
        ":scene_object_go_proto",
        "//intrinsic/geometry/proto/v1:material_go_proto",
        "//intrinsic/math/proto:quaternion_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
    ],
)

py_grpc_library(
    name = "scene_object_import_py_pb2_grpc",
    srcs = [":scene_object_import_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":scene_object_import_py_pb2"],
)

proto_library(
    name = "scene_object_updates_proto",
    srcs = ["scene_object_updates.proto"],
    deps = [
        ":collision_rules_proto",
        ":object_properties_proto",
        "//intrinsic/kinematics/types:joint_limits_proto",
        "//intrinsic/math/proto:pose_proto",
    ],
)

cc_proto_library(
    name = "scene_object_updates_cc_proto",
    deps = [":scene_object_updates_proto"],
)

go_proto_library(
    name = "scene_object_updates_go_proto",
    go_deps = [
        ":collision_rules_go_proto",
        ":object_properties_go_proto",
        "//intrinsic/kinematics/types:joint_limits_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
    deps = [":scene_object_updates_proto"],
)

py_proto_library(
    name = "scene_object_updates_py_pb2",
    visibility = [
        "//visibility:public",
    ],
    deps = [":scene_object_updates_proto"],
)

proto_library(
    name = "simulation_spec_proto",
    srcs = ["simulation_spec.proto"],
)

cc_proto_library(
    name = "simulation_spec_cc_proto",
    deps = [":simulation_spec_proto"],
)

go_proto_library(
    name = "simulation_spec_go_proto",
    deps = [":simulation_spec_proto"],
)

proto_library(
    name = "scene_object_config_proto",
    srcs = ["scene_object_config.proto"],
    deps = [
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "scene_object_config_cc_proto",
    deps = [":scene_object_config_proto"],
)

go_proto_library(
    name = "scene_object_config_go_proto",
    deps = [":scene_object_config_proto"],
)

py_proto_library(
    name = "scene_object_config_py_pb2",
    visibility = [
        "//visibility:public",
    ],
    deps = [":scene_object_config_proto"],
)
