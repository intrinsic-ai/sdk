# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "lease_api_proto",
    srcs = ["lease_api.proto"],
    deps = [
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "lease_api_go_proto",
    deps = [":lease_api_proto"],
)

go_grpc_library(
    name = "lease_api_go_grpc_proto",
    srcs = [":lease_api_proto"],
    deps = [
        ":lease_api_go_proto",
        "@org_golang_google_protobuf//types/known/emptypb",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)
