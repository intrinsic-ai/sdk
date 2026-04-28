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
    name = "pubsub_connect_proto",
    srcs = ["pubsub_connect.proto"],
    deps = [
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

go_proto_library(
    name = "pubsub_connect_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/platform/pubsub/connect/cloud/proto/v1alpha1/pubsub_connect_go_proto",
    protos = [":pubsub_connect_proto"],
    deps = [
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)

cc_proto_library(
    name = "pubsub_connect_cc_proto",
    deps = [":pubsub_connect_proto"],
)

cc_grpc_library(
    name = "pubsub_connect_cc_grpc_proto",
    srcs = [":pubsub_connect_proto"],
    grpc_only = True,
    deps = [":pubsub_connect_cc_proto"],
)

py_proto_library(
    name = "pubsub_connect_py_pb2",
    deps = [":pubsub_connect_proto"],
)

py_grpc_library(
    name = "pubsub_connect_py_pb2_grpc",
    srcs = [":pubsub_connect_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":pubsub_connect_py_pb2"],
)
