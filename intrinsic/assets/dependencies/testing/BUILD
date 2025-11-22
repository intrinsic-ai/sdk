# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

proto_library(
    name = "test_configs_proto",
    srcs = ["test_configs.proto"],
    visibility = ["//intrinsic/assets/dependencies:__pkg__"],
    deps = [
        "//intrinsic/assets/proto:field_metadata_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_proto",
        "@com_google_protobuf//:any_proto",
    ],
)

cc_proto_library(
    name = "test_configs_cc_proto",
    visibility = ["//intrinsic/assets/dependencies:__pkg__"],
    deps = [":test_configs_proto"],
)

go_proto_library(
    name = "test_configs_go_proto",
    testonly = True,
    go_deps = [
        "//intrinsic/assets/proto:field_metadata_go_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_go_proto",
    ],
    visibility = [
        "//intrinsic/assets:__subpackages__",
    ],
    deps = [":test_configs_proto"],
)

proto_library(
    name = "test_service_proto",
    srcs = ["test_service.proto"],
    visibility = ["//intrinsic/assets/dependencies:__pkg__"],
)

go_proto_library(
    name = "test_service_go_proto",
    visibility = ["//intrinsic/assets/dependencies:__pkg__"],
    deps = [":test_service_proto"],
)

go_grpc_library(
    name = "test_service_go_grpc_proto",
    srcs = [":test_service_proto"],
    visibility = ["//intrinsic/assets/dependencies:__pkg__"],
    deps = [":test_service_go_proto"],
)

py_proto_library(
    name = "test_service_py_pb2",
    visibility = ["//intrinsic/assets/dependencies:__pkg__"],
    deps = [":test_service_proto"],
)

py_grpc_library(
    name = "test_service_py_grpc_proto",
    srcs = [":test_service_proto"],
    grpc_library = requirement("grpcio"),
    visibility = ["//intrinsic/assets/dependencies:__pkg__"],
    deps = [":test_service_py_pb2"],
)

cc_proto_library(
    name = "test_service_cc_proto",
    visibility = ["//intrinsic/assets/dependencies:__pkg__"],
    deps = [":test_service_proto"],
)

cc_grpc_library(
    name = "test_service_cc_grpc_proto",
    srcs = [":test_service_proto"],
    grpc_only = True,
    visibility = ["//intrinsic/assets/dependencies:__pkg__"],
    deps = [":test_service_cc_proto"],
)
