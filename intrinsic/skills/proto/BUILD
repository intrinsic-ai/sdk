# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "equipment_proto",
    srcs = ["equipment.proto"],
    deps = [":footprint_proto"],
)

cc_proto_library(
    name = "equipment_cc_proto",
    deps = [":equipment_proto"],
)

py_proto_library(
    name = "equipment_py_pb2",
    deps = [":equipment_proto"],
)

go_proto_library(
    name = "equipment_go_proto",
    go_deps = [":footprint_go_proto"],
    deps = [":equipment_proto"],
)

proto_library(
    name = "skills_proto",
    srcs = ["skills.proto"],
    deps = [
        ":equipment_proto",
        "//intrinsic/resources/proto:resource_handle_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:descriptor_proto",
    ],
)

cc_proto_library(
    name = "skills_cc_proto",
    deps = [":skills_proto"],
)

py_proto_library(
    name = "skills_py_pb2",
    deps = [":skills_proto"],
)

go_proto_library(
    name = "skills_go_proto",
    go_deps = [
        ":equipment_go_proto",
        "//intrinsic/resources/proto:resource_handle_go_proto",
    ],
    deps = [":skills_proto"],
)

proto_library(
    name = "skill_registry_proto",
    srcs = ["skill_registry.proto"],
    deps = [
        ":skills_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

cc_proto_library(
    name = "skill_registry_cc_proto",
    deps = [":skill_registry_proto"],
)

cc_grpc_library(
    name = "skill_registry_cc_grpc_proto",
    srcs = [":skill_registry_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":skill_registry_cc_proto"],
)

py_proto_library(
    name = "skill_registry_py_pb2",
    deps = [":skill_registry_proto"],
)

py_grpc_library(
    name = "skill_registry_py_pb2_grpc",
    srcs = [":skill_registry_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":skill_registry_py_pb2"],
)

go_proto_library(
    name = "skill_registry_go_proto",
    go_deps = [":skills_go_proto"],
    deps = [":skill_registry_proto"],
)

go_grpc_library(
    name = "skill_registry_go_grpc_proto",
    srcs = [":skill_registry_proto"],
    deps = [":skills_go_proto"],
)

proto_library(
    name = "behavior_tree_registry_proto",
    srcs = ["behavior_tree_registry.proto"],
    deps = [
        ":skill_registry_config_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

cc_proto_library(
    name = "behavior_tree_registry_cc_proto",
    deps = [":behavior_tree_registry_proto"],
)

cc_grpc_library(
    name = "behavior_tree_registry_cc_grpc_proto",
    srcs = [":behavior_tree_registry_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":behavior_tree_registry_cc_proto"],
)

py_proto_library(
    name = "behavior_tree_registry_py_pb2",
    deps = [":behavior_tree_registry_proto"],
)

py_grpc_library(
    name = "behavior_tree_registry_py_pb2_grpc",
    srcs = [":behavior_tree_registry_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":behavior_tree_registry_py_pb2"],
)

go_proto_library(
    name = "behavior_tree_registry_go_proto",
    go_deps = [
        ":skills_go_proto",
        ":skill_registry_config_go_proto",
    ],
    deps = [":behavior_tree_registry_proto"],
)

go_grpc_library(
    name = "behavior_tree_registry_go_grpc_proto",
    srcs = [":behavior_tree_registry_proto"],
    deps = [
        ":skill_registry_config_go_proto",
        ":skills_go_proto",
    ],
)

proto_library(
    name = "skill_service_proto",
    srcs = ["skill_service.proto"],
    deps = [
        ":error_proto",
        ":footprint_proto",
        ":prediction_proto",
        ":skills_proto",
        "//intrinsic/logging/proto:context_proto",
        "//intrinsic/util/status:extended_status_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:duration_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

cc_proto_library(
    name = "skill_service_cc_proto",
    deps = [":skill_service_proto"],
)

cc_grpc_library(
    name = "skill_service_cc_grpc_proto",
    srcs = [":skill_service_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":skill_service_cc_proto"],
)

py_proto_library(
    name = "skill_service_py_pb2",
    deps = [":skill_service_proto"],
)

py_grpc_library(
    name = "skill_service_py_pb2_grpc",
    srcs = [":skill_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":skill_service_py_pb2"],
)

go_proto_library(
    name = "skill_service_go_proto",
    go_deps = [
        ":error_go_proto",
        ":footprint_go_proto",
        ":prediction_go_proto",
        ":skills_go_proto",
        "//intrinsic/logging/proto:context_go_proto",
        "//intrinsic/util/status:extended_status_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
    ],
    deps = [":skill_service_proto"],
)

go_grpc_library(
    name = "skill_service_go_grpc_proto",
    srcs = [":skill_service_proto"],
    deps = [
        ":error_go_proto",
        ":footprint_go_proto",
        ":prediction_go_proto",
        ":skills_go_proto",
        "//intrinsic/logging/proto:context_go_proto",
        "//intrinsic/util/status:extended_status_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
    ],
)

proto_library(
    name = "skill_service_config_proto",
    srcs = ["skill_service_config.proto"],
    deps = [
        ":skills_proto",
        "//intrinsic/assets/proto:status_spec_proto",
        "@com_google_protobuf//:duration_proto",
    ],
)

cc_proto_library(
    name = "skill_service_config_cc_proto",
    deps = [":skill_service_config_proto"],
)

py_proto_library(
    name = "skill_service_config_py_pb2",
    deps = [":skill_service_config_proto"],
)

go_proto_library(
    name = "skill_service_config_go_proto",
    go_deps = [
        ":skills_go_proto",
        "//intrinsic/assets/proto:status_spec_go_proto",
    ],
    deps = [":skill_service_config_proto"],
)

proto_library(
    name = "skill_registry_config_proto",
    srcs = ["skill_registry_config.proto"],
    deps = [
        ":skills_proto",
        "//intrinsic/executive/proto:behavior_tree_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "skill_registry_config_cc_proto",
    deps = [":skill_registry_config_proto"],
)

py_proto_library(
    name = "skill_registry_config_py_pb2",
    deps = [":skill_registry_config_proto"],
)

go_proto_library(
    name = "skill_registry_config_go_proto",
    go_deps = [
        ":skills_go_proto",
        "//intrinsic/executive/proto:behavior_tree_go_proto",
    ],
    deps = [":skill_registry_config_proto"],
)

proto_library(
    name = "footprint_proto",
    srcs = ["footprint.proto"],
    deps = [
        "//intrinsic/geometry/proto:transformed_geometry_storage_refs_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

py_proto_library(
    name = "footprint_py_pb2",
    deps = [":footprint_proto"],
)

cc_proto_library(
    name = "footprint_cc_proto",
    deps = [":footprint_proto"],
)

go_proto_library(
    name = "footprint_go_proto",
    go_deps = [
        "//intrinsic/geometry/proto:transformed_geometry_storage_refs_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
    ],
    deps = [":footprint_proto"],
)

proto_library(
    name = "prediction_proto",
    srcs = ["prediction.proto"],
    deps = [
        "//intrinsic/world/proto:object_world_updates_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:duration_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

py_proto_library(
    name = "prediction_py_pb2",
    deps = [":prediction_proto"],
)

go_proto_library(
    name = "prediction_go_proto",
    go_deps = [
        "//intrinsic/world/proto:object_world_updates_go_proto",
    ],
    deps = [":prediction_proto"],
)

cc_proto_library(
    name = "prediction_cc_proto",
    deps = [":prediction_proto"],
)

proto_library(
    name = "motion_targets_proto",
    srcs = ["motion_targets.proto"],
    deps = [
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/world/proto:entity_search_proto",
    ],
)

proto_library(
    name = "error_proto",
    srcs = ["error.proto"],
)

py_proto_library(
    name = "error_py_pb2",
    deps = [":error_proto"],
)

cc_proto_library(
    name = "error_cc_proto",
    deps = [":error_proto"],
)

go_proto_library(
    name = "error_go_proto",
    deps = [":error_proto"],
)

py_proto_library(
    name = "motion_targets_py_pb2",
    deps = [":motion_targets_proto"],
)

proto_library(
    name = "skill_manifest_proto",
    srcs = ["skill_manifest.proto"],
    deps = [
        ":equipment_proto",
        "//intrinsic/assets/proto:documentation_proto",
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/assets/proto:status_spec_proto",
        "//intrinsic/assets/proto:vendor_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:duration_proto",
    ],
)

cc_proto_library(
    name = "skill_manifest_cc_proto",
    deps = [":skill_manifest_proto"],
)

go_proto_library(
    name = "skill_manifest_go_proto",
    go_deps = [
        ":equipment_go_proto",
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:status_spec_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
    ],
    deps = [":skill_manifest_proto"],
)

py_proto_library(
    name = "skill_manifest_py_pb2",
    deps = [":skill_manifest_proto"],
)

proto_library(
    name = "processed_skill_manifest_proto",
    srcs = ["processed_skill_manifest.proto"],
    deps = [
        ":skill_manifest_proto",
        "//intrinsic/assets/proto:documentation_proto",
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/assets/proto:status_spec_proto",
        "//intrinsic/assets/proto:vendor_proto",
        "//intrinsic/kubernetes/workcell_spec/proto:image_proto",
        "@com_google_protobuf//:descriptor_proto",
    ],
)

go_proto_library(
    name = "processed_skill_manifest_go_proto",
    go_deps = [
        ":skill_manifest_go_proto",
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:status_spec_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
        "//intrinsic/kubernetes/workcell_spec/proto:image_go_proto",
    ],
    deps = [":processed_skill_manifest_proto"],
)

proto_library(
    name = "skill_parameter_metadata_proto",
    srcs = ["skill_parameter_metadata.proto"],
    deps = ["@com_google_protobuf//:descriptor_proto"],
)

cc_proto_library(
    name = "skill_parameter_metadata_cc_proto",
    deps = [":skill_parameter_metadata_proto"],
)

go_proto_library(
    name = "skill_parameter_metadata_go_proto",
    deps = [":skill_parameter_metadata_proto"],
)

py_proto_library(
    name = "skill_parameter_metadata_py_pb2",
    deps = [":skill_parameter_metadata_proto"],
)

proto_library(
    name = "skills_all_protos",
    deps = [
        ":equipment_proto",
        ":error_proto",
        ":footprint_proto",
        ":motion_targets_proto",
        ":skill_registry_config_proto",
        ":skill_registry_proto",
        ":skill_service_config_proto",
        ":skill_service_proto",
        ":skills_proto",
        "//intrinsic/resources/proto:resource_handle_proto",
    ],
)
