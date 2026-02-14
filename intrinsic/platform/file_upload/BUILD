# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_binary", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "file_upload_service_proto",
    srcs = ["file_upload_service.proto"],
)

go_proto_library(
    name = "file_upload_service_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/platform/file_upload/file_upload_service_go_proto",
    protos = [":file_upload_service_proto"],
)

go_binary(
    name = "file_upload_cli",
    srcs = [
        "file_upload_cli.go",
    ],
    importpath = "intrinsic/platform/file_upload/file_upload_cli",
    deps = [
        ":file_upload_service_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_google_subcommands//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials/insecure:go_default_library",
    ],
)
