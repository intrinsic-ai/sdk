# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

# AssetCatalog protos.

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "asset_catalog_proto",
    srcs = ["asset_catalog.proto"],
    deps = [
        ":release_metadata_proto",
        "//intrinsic/assets/data/proto/v1:data_asset_proto",
        "//intrinsic/assets/data/proto/v1:referenced_data_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_proto",
        "//intrinsic/assets/processes/proto:process_asset_proto",
        "//intrinsic/assets/proto:asset_tag_proto",
        "//intrinsic/assets/proto:asset_type_proto",
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/assets/proto:metadata_proto",
        "//intrinsic/assets/proto:view_proto",
        "//intrinsic/assets/proto/v1:search_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_proto",
        "//intrinsic/assets/services/proto:service_manifest_proto",
        "//intrinsic/skills/proto:processed_skill_manifest_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
        "@com_google_protobuf//:descriptor_proto",
        "@com_google_protobuf//:field_mask_proto",
    ],
)

cc_proto_library(
    name = "asset_catalog_cc_proto",
    deps = [":asset_catalog_proto"],
)

cc_grpc_library(
    name = "asset_catalog_cc_grpc_proto",
    srcs = [":asset_catalog_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":asset_catalog_cc_proto"],
)

go_proto_library(
    name = "asset_catalog_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    protos = [":asset_catalog_proto"],
    deps = [
        ":release_metadata_go_proto",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/data/proto/v1:referenced_data_go_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/processes/proto:process_asset_go_proto",
        "//intrinsic/assets/proto:asset_tag_go_proto",
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
        "//intrinsic/assets/proto:view_go_proto",
        "//intrinsic/assets/proto/v1:search_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/skills/proto:processed_skill_manifest_go_proto",
        "@io_bazel_rules_go//proto/wkt:descriptor_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)

py_proto_library(
    name = "asset_catalog_py_pb2",
    deps = [":asset_catalog_proto"],
)

py_grpc_library(
    name = "asset_catalog_py_pb2_grpc",
    srcs = [":asset_catalog_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":asset_catalog_py_pb2"],
)

proto_library(
    name = "release_metadata_proto",
    srcs = ["release_metadata.proto"],
)

py_proto_library(
    name = "release_metadata_py_pb2",
    deps = [":release_metadata_proto"],
)

cc_proto_library(
    name = "release_metadata_cc_proto",
    deps = [":release_metadata_proto"],
)

go_proto_library(
    name = "release_metadata_go_proto",
    protos = [":release_metadata_proto"],
)
