# Copyright 2023 Intrinsic Innovation LLC

# Grasp annotator service definition.

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

cc_proto_library(
    name = "grasp_annotator_service_cc_proto",
    deps = [":grasp_annotator_service_proto"],
)

proto_library(
    name = "grasp_annotator_service_proto",
    srcs = ["grasp_annotator_service.proto"],
    deps = [
        "//intrinsic/manipulation/grasping:grasp_annotations_proto",
        "//intrinsic/manipulation/grasping:grasp_annotator_proto",
        "@com_google_googleapis//google/rpc:status_proto",
    ],
)

py_proto_library(
    name = "grasp_annotator_service_py_pb2",
    deps = [":grasp_annotator_service_proto"],
)

py_grpc_library(
    name = "grasp_annotator_service_py_pb2_grpc",
    srcs = [":grasp_annotator_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":grasp_annotator_service_py_pb2"],
)
