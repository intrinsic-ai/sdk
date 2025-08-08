# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "resourcemanager_proto",
    srcs = ["resourcemanager.proto"],
    deps = [
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-openapiv2/options:options_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

proto_source_code_info_transitive_descriptor_set(
    name = "resourcemanager_proto_descriptor_set",
    deps = [":resourcemanager_proto"],
)

go_proto_library(
    name = "resourcemanager_go_proto",
    go_deps = [
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-openapiv2/options:options",
        "@org_golang_google_genproto_googleapis_api//annotations",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
    ],
    deps = [":resourcemanager_proto"],
)

go_grpc_library(
    name = "resourcemanager_go_grpc_proto",
    srcs = [":resourcemanager_proto"],
    deps = [
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-openapiv2/options",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)
