# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")

package(
    default_visibility = [
        "//visibility:public",
    ],
)

proto_library(
    name = "service_proto",
    srcs = ["service.proto"],
    deps = [
        "//intrinsic/icon/proto:joint_space_proto",
        "//intrinsic/icon/proto:logging_mode_proto",
        "//intrinsic/icon/proto:part_status_proto",
        "//intrinsic/icon/proto:safety_status_proto",
        "//intrinsic/icon/proto:streaming_output_proto",
        "//intrinsic/icon/proto:types_proto",
        "//intrinsic/logging/proto:context_proto",
        "//intrinsic/world/proto:robot_payload_proto",
        "@com_google_googleapis//google/rpc:status_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:duration_proto",
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

cc_proto_library(
    name = "service_cc_proto",
    deps = [":service_proto"],
)

py_proto_library(
    name = "service_py_pb2",
    deps = [":service_proto"],
)

cc_grpc_library(
    name = "service_cc_grpc_proto",
    srcs = [":service_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":service_cc_proto"],
)

py_grpc_library(
    name = "service_py_pb2_grpc",
    srcs = [":service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":service_py_pb2"],
)
