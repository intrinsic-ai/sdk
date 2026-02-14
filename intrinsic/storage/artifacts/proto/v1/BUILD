# Copyright 2023 Intrinsic Innovation LLC

# gRPC service definition for the workcell-spec-service.

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "artifact_proto",
    srcs = ["artifact.proto"],
    deps = ["@com_google_protobuf//:timestamp_proto"],
)

go_proto_library(
    name = "artifact_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/storage/artifacts/proto/v1/artifact_go_proto",
    protos = [":artifact_proto"],
)

proto_library(
    name = "articat_proto",
    srcs = ["articat.proto"],
)

go_proto_library(
    name = "articat_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/storage/artifacts/proto/v1/articat_go_proto",
    protos = [
        ":articat_proto",
    ],
)
