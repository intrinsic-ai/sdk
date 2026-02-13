# Copyright 2023 Intrinsic Innovation LLC

load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "device_service_config_proto",
    srcs = ["device_service_config.proto"],
    deps = [
        "//intrinsic/assets/proto:field_metadata_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_proto",
    ],
)

cc_proto_library(
    name = "device_service_config_cc_proto",
    deps = [":device_service_config_proto"],
)

go_proto_library(
    name = "device_service_config_go_proto",
    protos = [":device_service_config_proto"],
    deps = [
        "//intrinsic/assets/proto:field_metadata_go_proto",
        "//intrinsic/assets/proto/v1:dependency_go_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_go_proto",
    ],
)

proto_library(
    name = "device_service_proto",
    srcs = ["device_service.proto"],
    deps = [
        ":device_service_config_proto",
        ":esi_proto",
    ],
)

proto_source_code_info_transitive_descriptor_set(
    name = "device_service_proto_descriptor_set",
    deps = [":device_service_proto"],
)

proto_library(
    name = "esi_proto",
    srcs = ["esi.proto"],
)

cc_proto_library(
    name = "esi_cc_proto",
    deps = [":esi_proto"],
)

go_proto_library(
    name = "esi_go_proto",
    protos = [":esi_proto"],
)

cc_proto_library(
    name = "device_service_cc_proto",
    deps = [":device_service_proto"],
)

cc_grpc_library(
    name = "device_service_cc_grpc_proto",
    srcs = [":device_service_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":device_service_cc_proto"],
)

go_proto_library(
    name = "device_service_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    protos = [
        ":device_service_proto",
    ],
    deps = [
        ":device_service_config_go_proto",
        ":esi_go_proto",
        "//intrinsic/assets/proto:id_go_proto",
    ],
)

proto_library(
    name = "eni_proto",
    srcs = ["eni.proto"],
)

cc_proto_library(
    name = "eni_cc_proto",
    deps = [":eni_proto"],
)

go_proto_library(
    name = "eni_go_proto",
    protos = [":eni_proto"],
)
