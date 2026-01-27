# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

proto_library(
    name = "asset_proto",
    srcs = ["asset.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":processed_asset_proto",
        "//intrinsic/assets/proto:id_proto",
    ],
)

cc_proto_library(
    name = "asset_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":asset_proto"],
)

py_proto_library(
    name = "asset_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":asset_proto"],
)

go_proto_library(
    name = "asset_go_proto",
    go_deps = [
        ":processed_asset_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":asset_proto"],
)

proto_library(
    name = "asset_graph_proto",
    srcs = ["asset_graph.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "asset_graph_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":asset_graph_proto"],
)

py_proto_library(
    name = "asset_graph_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":asset_graph_proto"],
)

go_proto_library(
    name = "asset_graph_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":asset_graph_proto"],
)

proto_library(
    name = "grpc_connection_proto",
    srcs = ["grpc_connection.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

go_proto_library(
    name = "grpc_connection_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":grpc_connection_proto"],
)

cc_proto_library(
    name = "grpc_connection_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":grpc_connection_proto"],
)

py_proto_library(
    name = "grpc_connection_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":grpc_connection_proto"],
)

proto_library(
    name = "processed_asset_proto",
    srcs = ["processed_asset.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/data/proto/v1:data_asset_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_proto",
        "//intrinsic/assets/processes/proto:process_asset_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_proto",
        "//intrinsic/assets/services/proto:service_manifest_proto",
        "//intrinsic/skills/proto:processed_skill_manifest_proto",
    ],
)

cc_proto_library(
    name = "processed_asset_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":processed_asset_proto"],
)

py_proto_library(
    name = "processed_asset_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":processed_asset_proto"],
)

go_proto_library(
    name = "processed_asset_go_proto",
    go_deps = [
        "//intrinsic/assets/data/proto/v1:data_asset_go_proto",
        "//intrinsic/assets/hardware_devices/proto/v1:hardware_device_manifest_go_proto",
        "//intrinsic/assets/processes/proto:process_asset_go_proto",
        "//intrinsic/assets/scene_objects/proto:scene_object_manifest_go_proto",
        "//intrinsic/assets/services/proto:service_manifest_go_proto",
        "//intrinsic/skills/proto:processed_skill_manifest_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":processed_asset_proto"],
)

proto_library(
    name = "reference_proto",
    srcs = ["reference.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto:asset_type_proto",
        "//intrinsic/assets/proto:id_proto",
    ],
)

cc_proto_library(
    name = "reference_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":reference_proto"],
)

go_proto_library(
    name = "reference_go_proto",
    go_deps = [
        "//intrinsic/assets/proto:asset_type_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":reference_proto"],
)

py_proto_library(
    name = "reference_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":reference_proto"],
)

proto_library(
    name = "search_proto",
    srcs = ["search.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

py_proto_library(
    name = "search_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":search_proto"],
)

go_proto_library(
    name = "search_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":search_proto"],
)

cc_proto_library(
    name = "search_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":search_proto"],
)

proto_library(
    name = "resolved_dependency_proto",
    srcs = ["resolved_dependency.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/assets/proto/v1:grpc_connection_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
    ],
)

go_proto_library(
    name = "resolved_dependency_go_proto",
    go_deps = [
        "@org_golang_google_genproto_googleapis_api//annotations",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto/v1:grpc_connection_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":resolved_dependency_proto"],
)

cc_proto_library(
    name = "resolved_dependency_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":resolved_dependency_proto"],
)

py_proto_library(
    name = "resolved_dependency_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":resolved_dependency_proto"],
)

proto_library(
    name = "solution_deployment_proto",
    srcs = ["solution_deployment.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/config/proto:operation_mode_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
    ],
)

go_proto_library(
    name = "solution_deployment_go_proto",
    go_deps = [
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "//intrinsic/config/proto:operation_mode_go_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":solution_deployment_proto"],
)

go_grpc_library(
    name = "solution_deployment_go_grpc_proto",
    protos = [":solution_deployment_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/config/proto:operation_mode_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
    ],
)

cc_proto_library(
    name = "solution_deployment_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":solution_deployment_proto"],
)

py_proto_library(
    name = "solution_deployment_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":solution_deployment_proto"],
)

proto_library(
    name = "dependency_proto",
    srcs = ["dependency.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

go_proto_library(
    name = "dependency_go_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dependency_proto"],
)

cc_proto_library(
    name = "dependency_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dependency_proto"],
)

py_proto_library(
    name = "dependency_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":dependency_proto"],
)
