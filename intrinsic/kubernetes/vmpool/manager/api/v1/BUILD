# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

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

proto_source_code_info_transitive_descriptor_set(
    name = "lease_api_proto_descriptor_set",
    deps = [":lease_api_proto"],
)

go_proto_library(
    name = "lease_api_go_proto",
    deps = [":lease_api_proto"],
)

go_grpc_library(
    name = "lease_api_go_grpc_proto",
    protos = [":lease_api_proto"],
    deps = [
        ":lease_api_go_proto",
        "@org_golang_google_protobuf//types/known/emptypb",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

py_proto_library(
    name = "lease_api_py_pb2",
    deps = [":lease_api_proto"],
)

py_grpc_library(
    name = "lease_api_py_pb2_grpc",
    srcs = [":lease_api_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":lease_api_py_pb2"],
)
