# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "product_reader_proto",
    srcs = ["product_reader.proto"],
    deps = [
        "//intrinsic/scene/product/proto:product_proto",
        "//intrinsic/scene/proto/v1:scene_object_proto",
        "@com_google_protobuf//:struct_proto",
    ],
)

go_proto_library(
    name = "product_reader_go_proto",
    go_deps = [
        "//intrinsic/scene/product/proto:product_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
    ],
    deps = [":product_reader_proto"],
)

go_grpc_library(
    name = "product_reader_go_grpc_proto",
    srcs = [":product_reader_proto"],
    deps = [
        "//intrinsic/scene/product/proto:product_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "@org_golang_google_protobuf//types/known/structpb",
    ],
)

py_proto_library(
    name = "product_reader_py_pb2",
    deps = [":product_reader_proto"],
)

py_grpc_library(
    name = "product_reader_py_pb2_grpc",
    srcs = [":product_reader_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":product_reader_py_pb2"],
)

cc_proto_library(
    name = "product_reader_cc_proto",
    deps = [":product_reader_proto"],
)

cc_grpc_library(
    name = "product_reader_cc_grpc",
    srcs = [":product_reader_proto"],
    grpc_only = True,
    deps = [":product_reader_cc_proto"],
)
