# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "discoveryapi_proto",
    srcs = ["discovery_api.proto"],
    deps = [
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-openapiv2/options:options_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

go_proto_library(
    name = "discoveryapi_go_proto",
    go_deps = [
        "@org_golang_google_genproto_googleapis_api//annotations",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-openapiv2/options:options",
    ],
    deps = [
        ":discoveryapi_proto",
    ],
)

go_grpc_library(
    name = "discoveryapi_go_grpc_proto",
    srcs = [":discoveryapi_proto"],
    deps = [
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-openapiv2/options",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)
