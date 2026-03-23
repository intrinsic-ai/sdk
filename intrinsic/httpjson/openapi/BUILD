# Copyright 2023 Intrinsic Innovation LLC

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("//bazel:go_macros.bzl", "go_library")

bzl_library(
    name = "protoc_gen_openapi_bzl",
    srcs = ["protoc_gen_openapi.bzl"],
    visibility = ["//intrinsic:public_api_users"],
    deps = ["@com_google_protobuf//bazel/common:proto_info_bzl"],
)

go_library(
    name = "handlers",
    srcs = ["handlers.go"],
    importpath = "intrinsic/httpjson/openapi/handlers",
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "@com_github_grpc_ecosystem_grpc_gateway_v2//runtime",
        "@io_bazel_rules_go//go/runfiles:go_default_library",
    ],
)
