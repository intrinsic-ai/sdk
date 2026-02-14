# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

# Protos for the Data asset type.

proto_library(
    name = "data_asset_proto",
    srcs = ["data_asset.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto:metadata_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:descriptor_proto",
    ],
)

cc_proto_library(
    name = "data_asset_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":data_asset_proto"],
)

go_proto_library(
    name = "data_asset_go_proto",
    importpath = "intrinsic/assets/data/proto/v1/data_asset_go_proto",
    protos = [":data_asset_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/assets/proto:metadata_go_proto"],
)

py_proto_library(
    name = "data_asset_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":data_asset_proto"],
)

proto_library(
    name = "referenced_data_proto",
    srcs = ["referenced_data.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "referenced_data_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":referenced_data_proto"],
)

go_proto_library(
    name = "referenced_data_go_proto",
    importpath = "intrinsic/assets/data/proto/v1/referenced_data_go_proto",
    protos = [":referenced_data_proto"],
    visibility = ["//intrinsic:public_api_users"],
)

py_proto_library(
    name = "referenced_data_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":referenced_data_proto"],
)

proto_library(
    name = "referenced_data_struct_proto",
    srcs = ["referenced_data_struct.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [":referenced_data_proto"],
)

cc_proto_library(
    name = "referenced_data_struct_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":referenced_data_struct_proto"],
)

go_proto_library(
    name = "referenced_data_struct_go_proto",
    importpath = "intrinsic/assets/data/proto/v1/referenced_data_struct_go_proto",
    protos = [":referenced_data_struct_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":referenced_data_go_proto",
    ],
)

py_proto_library(
    name = "referenced_data_struct_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":referenced_data_struct_proto"],
)

proto_library(
    name = "data_manifest_proto",
    srcs = ["data_manifest.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto:documentation_proto",
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/assets/proto:vendor_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "data_manifest_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":data_manifest_proto"],
)

go_proto_library(
    name = "data_manifest_go_proto",
    importpath = "intrinsic/assets/data/proto/v1/data_manifest_go_proto",
    protos = [":data_manifest_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/assets/proto:documentation_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:vendor_go_proto",
    ],
)

py_proto_library(
    name = "data_manifest_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":data_manifest_proto"],
)

proto_library(
    name = "data_assets_proto",
    srcs = ["data_assets.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":data_asset_proto",
        ":referenced_data_proto",
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/assets/proto:metadata_proto",
    ],
)

cc_proto_library(
    name = "data_assets_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":data_assets_proto"],
)

go_proto_library(
    name = "data_assets_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/assets/data/proto/v1/data_assets_go_proto",
    protos = [":data_assets_proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":data_asset_go_proto",
        ":referenced_data_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/assets/proto:metadata_go_proto",
    ],
)

cc_grpc_library(
    name = "data_assets_cc_grpc_proto",
    srcs = [":data_assets_proto"],
    generate_mocks = True,
    grpc_only = True,
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":data_assets_cc_proto",
    ],
)

py_proto_library(
    name = "data_assets_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":data_assets_proto"],
)

py_grpc_library(
    name = "data_assets_py_pb2_grpc",
    srcs = [":data_assets_proto"],
    grpc_library = requirement("grpcio"),
    visibility = ["//intrinsic:public_api_users"],
    deps = [":data_assets_py_pb2"],
)
