# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "tokens_proto",
    srcs = ["tokens.proto"],
    deps = [
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

py_proto_library(
    name = "tokens_py_pb2",
    deps = [
        ":tokens_proto",
    ],
)

py_grpc_library(
    name = "tokens_py_pb2_grpc",
    srcs = [":tokens_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":tokens_py_pb2"],
)

go_proto_library(
    name = "tokens_go_proto",
    protos = [":tokens_proto"],
    deps = [
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)

go_grpc_library(
    name = "tokens_go_grpc_proto",
    protos = [":tokens_proto"],
    deps = [
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)
