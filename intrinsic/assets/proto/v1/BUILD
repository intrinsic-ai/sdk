# Copyright 2023 Intrinsic Innovation LLC

load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

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
