# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "simulation_service_proto",
    srcs = [
        "simulation_service.proto",
    ],
    deps = [
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:wrappers_proto",
    ],
)

cc_proto_library(
    name = "simulation_service_cc_proto",
    deps = [":simulation_service_proto"],
)

cc_grpc_library(
    name = "simulation_service_cc_grpc_proto",
    srcs = [":simulation_service_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":simulation_service_cc_proto"],
)

go_proto_library(
    name = "simulation_service_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/simulation/service/proto/v1/simulation_service_go_proto",
    protos = [":simulation_service_proto"],
    deps = [
        "@org_golang_google_protobuf//types/known/emptypb",
        "@org_golang_google_protobuf//types/known/wrapperspb",
    ],
)

py_proto_library(
    name = "simulation_service_py_pb2",
    deps = [":simulation_service_proto"],
)

py_grpc_library(
    name = "simulation_service_py_pb2_grpc",
    srcs = [":simulation_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":simulation_service_py_pb2"],
)
