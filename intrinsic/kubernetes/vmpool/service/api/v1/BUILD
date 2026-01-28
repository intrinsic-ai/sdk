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
    name = "vmpool_api_proto",
    srcs = ["vmpool_api.proto"],
    deps = [
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:field_mask_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "vmpool_api_go_proto",
    protos = [":vmpool_api_proto"],
    deps = ["@org_golang_google_genproto_googleapis_api//annotations"],
)

go_grpc_library(
    name = "vmpool_api_go_grpc_proto",
    protos = [":vmpool_api_proto"],
    deps = [
        ":vmpool_api_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
        "@org_golang_google_protobuf//types/known/emptypb",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

proto_source_code_info_transitive_descriptor_set(
    name = "vmpool_api_proto_descriptor_set",
    deps = [":vmpool_api_proto"],
)

py_proto_library(
    name = "vmpool_api_py_pb2",
    deps = [":vmpool_api_proto"],
)

py_grpc_library(
    name = "vmpool_api_py_pb2_grpc",
    srcs = [":vmpool_api_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":vmpool_api_py_pb2"],
)
