# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

proto_library(
    name = "field_metadata_proto",
    srcs = ["field_metadata.proto"],
    visibility = ["//intrinsic:__subpackages__"],
    deps = ["@com_google_protobuf//:descriptor_proto"],
)

cc_proto_library(
    name = "field_metadata_cc_proto",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":field_metadata_proto"],
)

go_proto_library(
    name = "field_metadata_go_proto",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":field_metadata_proto"],
)

py_proto_library(
    name = "field_metadata_py_pb2",
    deps = [":field_metadata_proto"],
)

proto_library(
    name = "asset_tag_proto",
    srcs = ["asset_tag.proto"],
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":asset_type_proto"],
)

cc_proto_library(
    name = "asset_tag_cc_proto",
    visibility = ["//visibility:private"],
    deps = [":asset_tag_proto"],
)

py_proto_library(
    name = "asset_tag_py_pb2",
    visibility = ["//visibility:private"],
    deps = [":asset_tag_proto"],
)

go_proto_library(
    name = "asset_tag_go_proto",
    go_deps = [
        ":asset_type_go_proto",
    ],
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":asset_tag_proto"],
)

filegroup(
    name = "asset_tags",
    srcs = [
        "asset_tags.textproto",
    ],
    visibility = [
        "//intrinsic/assets:__subpackages__",
    ],
)

proto_library(
    name = "asset_type_proto",
    srcs = ["asset_type.proto"],
    visibility = ["//intrinsic:__subpackages__"],
)

cc_proto_library(
    name = "asset_type_cc_proto",
    visibility = [
        "//intrinsic/world/tools:__pkg__",
    ],
    deps = [":asset_type_proto"],
)

py_proto_library(
    name = "asset_type_py_pb2",
    visibility = [
        "//intrinsic:__subpackages__",
        "//intrinsic:public_api_users",
    ],
    deps = [":asset_type_proto"],
)

go_proto_library(
    name = "asset_type_go_proto",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":asset_type_proto"],
)

proto_library(
    name = "documentation_proto",
    srcs = ["documentation.proto"],
    visibility = ["//intrinsic:__subpackages__"],
)

cc_proto_library(
    name = "documentation_cc_proto",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":documentation_proto"],
)

py_proto_library(
    name = "documentation_py_pb2",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":documentation_proto"],
)

go_proto_library(
    name = "documentation_go_proto",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":documentation_proto"],
)

proto_library(
    name = "status_spec_proto",
    srcs = ["status_spec.proto"],
    visibility = ["//intrinsic:__subpackages__"],
)

cc_proto_library(
    name = "status_spec_cc_proto",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":status_spec_proto"],
)

py_proto_library(
    name = "status_spec_py_pb2",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":status_spec_proto"],
)

go_proto_library(
    name = "status_spec_go_proto",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":status_spec_proto"],
)

proto_library(
    name = "id_proto",
    srcs = ["id.proto"],
    visibility = ["//intrinsic:__subpackages__"],
)

py_proto_library(
    name = "id_py_pb2",
    visibility = [
        "//intrinsic:__subpackages__",
        "//intrinsic:public_api_users",
    ],
    deps = [":id_proto"],
)

cc_proto_library(
    name = "id_cc_proto",
    visibility = [
        "//intrinsic:__subpackages__",
        "//visibility:public",
    ],
    deps = [":id_proto"],
)

go_proto_library(
    name = "id_go_proto",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":id_proto"],
)

proto_library(
    name = "metadata_proto",
    srcs = ["metadata.proto"],
    visibility = ["//intrinsic:__subpackages__"],
    deps = [
        ":asset_tag_proto",
        ":asset_type_proto",
        ":documentation_proto",
        ":id_proto",
        ":release_tag_proto",
        ":vendor_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
        "@com_google_protobuf//:descriptor_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

cc_proto_library(
    name = "metadata_cc_proto",
    visibility = ["//visibility:private"],
    deps = [":metadata_proto"],
)

py_proto_library(
    name = "metadata_py_pb2",
    visibility = [
        "//intrinsic/solutions:__pkg__",
        "//visibility:private",
    ],
    deps = [":metadata_proto"],
)

go_proto_library(
    name = "metadata_go_proto",
    go_deps = [
        ":asset_tag_go_proto",
        ":asset_type_go_proto",
        ":documentation_go_proto",
        ":id_go_proto",
        ":release_tag_go_proto",
        ":vendor_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":metadata_proto"],
)

proto_library(
    name = "release_tag_proto",
    srcs = ["release_tag.proto"],
    visibility = ["//intrinsic:__subpackages__"],
)

go_proto_library(
    name = "release_tag_go_proto",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":release_tag_proto"],
)

proto_library(
    name = "vendor_proto",
    srcs = ["vendor.proto"],
    visibility = ["//intrinsic:__subpackages__"],
)

py_proto_library(
    name = "vendor_py_pb2",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":vendor_proto"],
)

cc_proto_library(
    name = "vendor_cc_proto",
    visibility = ["//intrinsic/perception/service:__pkg__"],
    deps = [":vendor_proto"],
)

go_proto_library(
    name = "vendor_go_proto",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":vendor_proto"],
)

proto_library(
    name = "view_proto",
    srcs = ["view.proto"],
    visibility = ["//intrinsic:__subpackages__"],
)

cc_proto_library(
    name = "view_cc_proto",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":view_proto"],
)

go_proto_library(
    name = "view_go_proto",
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":view_proto"],
)

py_proto_library(
    name = "view_py_pb2",
    visibility = [
        "//intrinsic:__subpackages__",
        "//intrinsic:public_api_users",
    ],
    deps = [":view_proto"],
)

proto_library(
    name = "asset_deployment_proto",
    srcs = ["asset_deployment.proto"],
    visibility = [
        "//intrinsic/kubernetes/api_relay:__subpackages__",
        "//intrinsic/scene:__subpackages__",
        "//intrinsic/world/service:__subpackages__",
    ],
    deps = [
        ":asset_type_proto",
        ":id_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/scene/proto/v1:scene_object_config_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

py_proto_library(
    name = "asset_deployment_py_pb2",
    visibility = [
        "//intrinsic:__subpackages__",
        "//intrinsic:public_api_users",
    ],
    deps = [":asset_deployment_proto"],
)

py_grpc_library(
    name = "asset_deployment_pb2_grpc",
    srcs = [":asset_deployment_proto"],
    grpc_library = requirement("grpcio"),
    visibility = [
        "//intrinsic:__subpackages__",
        "//intrinsic:public_api_users",
    ],
    deps = [":asset_deployment_py_pb2"],
)

go_proto_library(
    name = "asset_deployment_go_proto",
    go_deps = [
        ":asset_type_go_proto",
        ":id_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_config_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
    ],
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":asset_deployment_proto"],
)

go_grpc_library(
    name = "asset_deployment_go_grpc_proto",
    srcs = [":asset_deployment_proto"],
    visibility = ["//intrinsic:__subpackages__"],
    deps = [
        ":asset_type_go_proto",
        ":id_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_config_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_protobuf//types/known/anypb",
    ],
)

proto_library(
    name = "installed_assets_proto",
    srcs = ["installed_assets.proto"],
    visibility = [
        "//intrinsic/kubernetes/api_relay:__subpackages__",
        "//intrinsic/world/service:__subpackages__",
    ],
    deps = [
        ":asset_tag_proto",
        ":asset_type_proto",
        ":id_proto",
        ":metadata_proto",
        ":view_proto",
        "//intrinsic/assets/data/proto/v1:data_asset_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_proto",
        "//intrinsic/assets/processes/proto:process_manifest_proto",
        "//intrinsic/assets/proto/v1:search_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_proto",
        "//intrinsic/assets/services/proto:service_manifest_proto",
        "//intrinsic/skills/proto:processed_skill_manifest_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
    ],
)

go_proto_library(
    name = "installed_assets_go_proto",
    go_deps = [
        ":asset_tag_go_proto",
        ":asset_type_go_proto",
        ":id_go_proto",
        ":metadata_go_proto",
        ":view_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/processes/proto:process_manifest_go_proto",
        "//intrinsic/assets/proto/v1:search_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/skills/proto:processed_skill_manifest_go_proto",
    ],
    visibility = ["//intrinsic:__subpackages__"],
    deps = [":installed_assets_proto"],
)

py_proto_library(
    name = "installed_assets_py_pb2",
    visibility = [
        "//intrinsic:public_api_users",
        "//intrinsic/solutions:__pkg__",
    ],
    deps = [":installed_assets_proto"],
)

py_grpc_library(
    name = "installed_assets_pb2_grpc",
    srcs = [":installed_assets_proto"],
    grpc_library = requirement("grpcio"),
    visibility = [
        "//intrinsic:public_api_users",
        "//intrinsic/solutions:__pkg__",
    ],
    deps = [":installed_assets_py_pb2"],
)

go_grpc_library(
    name = "installed_assets_go_grpc_proto",
    srcs = [":installed_assets_proto"],
    visibility = ["//intrinsic:__subpackages__"],
    deps = [
        ":asset_tag_go_proto",
        ":asset_type_go_proto",
        ":id_go_proto",
        ":metadata_go_proto",
        ":view_go_proto",
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/processes/proto:process_manifest_go_proto",
        "//intrinsic/assets/proto/v1:search_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/skills/proto:processed_skill_manifest_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
    ],
)
