# Copyright 2023 Intrinsic Innovation LLC

# Shared Bazel macros for use in external workspaces and general Bazel helpers.

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")
load("@io_bazel_rules_go//proto:compiler.bzl", "go_proto_compiler")
load("//bazel:go_macros.bzl", "go_binary")
load(":exec_properties.bzl", "DEFAULT")

package(default_visibility = ["//visibility:public"])

exports_files(glob(["content_mirror/*"]))

bzl_library(
    name = "container_bzl",
    srcs = ["container.bzl"],
    deps = [
        "@rules_oci//oci:defs",
        "@rules_pkg//pkg:bzl_srcs",
    ],
)

bzl_library(
    name = "cc_oci_image_bzl",
    srcs = ["cc_oci_image.bzl"],
    deps = [
        ":container_bzl",
        "@bazel_skylib//lib:paths",
    ],
)

bzl_library(
    name = "python_oci_image_bzl",
    srcs = ["python_oci_image.bzl"],
    deps = [
        ":container_bzl",
        "@aspect_bazel_lib//lib:tar",
        "@bazel_skylib//lib:paths",
    ],
)

bzl_library(
    name = "python_paths_aspect_bzl",
    srcs = ["python_paths_aspect.bzl"],
)

bzl_library(
    name = "go_macros",
    srcs = ["go_macros.bzl"],
    deps = [
        "@bazel_skylib//lib:paths",
        "@io_bazel_rules_go//go:def",
        "@io_bazel_rules_go//proto:def",
    ],
)

platform(
    name = "linux_x86_64",
    constraint_values = [
        "@platforms//os:linux",
        "@platforms//cpu:x86_64",
        "@bazel_tools//tools/cpp:clang",
    ],
    exec_properties = DEFAULT,
    parents = ["@platforms//host"],
)

go_proto_compiler(
    name = "go_gen_grpc_gateway",
    options = [
        "paths=source_relative",
        "logtostderr=true",
    ],
    plugin = "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-grpc-gateway",
    suffix = ".pb.gw.go",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_grpc_ecosystem_grpc_gateway_v2//runtime:go_default_library",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//utilities:go_default_library",
        "@org_golang_google_grpc//grpclog:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

go_binary(
    name = "checkimportpathmatchesgopackage",
    testonly = True,
    srcs = ["checkimportpathmatchesgopackage.go"],
    visibility = ["//intrinsic:internal_api_users"],
)
