# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "vmpool_api_proto",
    srcs = ["vmpool_api.proto"],
    deps = [
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:field_mask_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "vmpool_api_go_proto",
    deps = [":vmpool_api_proto"],
)

go_grpc_library(
    name = "vmpool_api_go_grpc_proto",
    srcs = [":vmpool_api_proto"],
    deps = [
        ":vmpool_api_go_proto",
        "@org_golang_google_protobuf//types/known/emptypb",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

proto_source_code_info_transitive_descriptor_set(
    name = "vmpool_api_proto_descriptor_set",
    deps = [":vmpool_api_proto"],
)
