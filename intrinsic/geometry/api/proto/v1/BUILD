# Copyright 2023 Intrinsic Innovation LLC

""" Intrinsic package geometry/api/proto/v1. """

load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "material_proto",
    srcs = ["material.proto"],
    deps = ["@com_google_googleapis//google/type:color_proto"],
)

cc_proto_library(
    name = "material_cc_proto",
    deps = [":material_proto"],
)

py_proto_library(
    name = "material_py_pb2",
    deps = [":material_proto"],
)

go_proto_library(
    name = "material_go_proto",
    go_deps = [
        "@org_golang_google_genproto//googleapis/type/color",
    ],
    deps = [":material_proto"],
)
