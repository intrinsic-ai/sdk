# Copyright 2023 Intrinsic Innovation LLC

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")

package(default_visibility = ["//visibility:public"])

exports_files(["opcua_equipment_service.proto"])

proto_library(
    name = "opcua_equipment_service_proto",
    srcs = ["opcua_equipment_service.proto"],
    deps = [
        "//intrinsic/hardware/gpio/v1:gpio_service_proto",
        "//intrinsic/hardware/gpio/v1:signal_proto",
        "//intrinsic/resources/proto:resource_handle_proto",
        "//intrinsic/world/proto:object_world_updates_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:duration_proto",
    ],
)

cc_proto_library(
    name = "opcua_equipment_service_cc_proto",
    deps = [":opcua_equipment_service_proto"],
)

cc_grpc_library(
    name = "opcua_equipment_service_cc_grpc_proto",
    srcs = [":opcua_equipment_service_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":opcua_equipment_service_cc_proto"],
)

py_proto_library(
    name = "opcua_equipment_service_py_pb2",
    deps = [":opcua_equipment_service_proto"],
)

py_grpc_library(
    name = "opcua_equipment_service_py_pb2_grpc",
    srcs = [":opcua_equipment_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":opcua_equipment_service_py_pb2"],
)
