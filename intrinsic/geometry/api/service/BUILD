# Copyright 2023 Intrinsic Innovation LLC

""" Intrinsic package geometry/api/service. """

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library", "py_proto_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "geometry_storage_refs_proto",
    srcs = ["geometry_storage_refs.proto"],
)

py_proto_library(
    name = "geometry_storage_refs_py_pb2",
    deps = [":geometry_storage_refs_proto"],
)

cc_proto_library(
    name = "geometry_storage_refs_cc_proto",
    deps = [":geometry_storage_refs_proto"],
)

go_proto_library(
    name = "geometry_storage_refs_go_proto",
    deps = [":geometry_storage_refs_proto"],
)

proto_library(
    name = "geometry_service_proto",
    srcs = ["geometry_service.proto"],
    deps = [
        ":geometry_service_types_proto",
        ":geometry_storage_refs_proto",
    ],
)

cc_proto_library(
    name = "geometry_service_cc_proto",
    deps = [":geometry_service_proto"],
)

py_proto_library(
    name = "geometry_service_py_pb2",
    deps = [":geometry_service_proto"],
)

py_grpc_library(
    name = "geometry_service_py_pb2_grpc",
    srcs = [":geometry_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":geometry_service_py_pb2"],
)

proto_library(
    name = "geometry_service_types_proto",
    srcs = ["geometry_service_types.proto"],
    deps = [
        ":geometry_storage_refs_proto",
        "//intrinsic/geometry/api/proto:geometry_proto",
        "//intrinsic/geometry/api/proto:lazy_exact_geometry_proto",
        "//intrinsic/geometry/api/proto:renderable_proto",
        "//intrinsic/math/proto:matrix_proto",
    ],
)

py_proto_library(
    name = "geometry_service_types_py_pb2",
    deps = [":geometry_service_types_proto"],
)

cc_proto_library(
    name = "geometry_service_types_cc_proto",
    deps = [":geometry_service_types_proto"],
)

go_proto_library(
    name = "geometry_service_types_go_proto",
    go_deps = [
        ":geometry_storage_refs_go_proto",
        "//intrinsic/geometry/api/proto:geometry_go_proto",
        "//intrinsic/geometry/api/proto:lazy_exact_geometry_go_proto",
        "//intrinsic/geometry/api/proto:renderable_go_proto",
        "//intrinsic/math/proto:matrix_go_proto",
    ],
    deps = [":geometry_service_types_proto"],
)

go_proto_library(
    name = "geometry_service_go_proto",
    go_deps = [
        ":geometry_service_types_go_proto",
        ":geometry_storage_refs_go_proto",
    ],
    deps = [":geometry_service_proto"],
)

cc_grpc_library(
    name = "geometry_service_cc_grpc_proto",
    srcs = [":geometry_service_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":geometry_service_cc_proto"],
)

go_grpc_library(
    name = "geometry_service_go_grpc_proto",
    srcs = [":geometry_service_proto"],
    deps = [
        ":geometry_service_types_go_proto",
        ":geometry_storage_refs_go_proto",
    ],
)

proto_library(
    name = "transformed_geometry_storage_refs_proto",
    srcs = ["transformed_geometry_storage_refs.proto"],
    deps = [
        ":geometry_storage_refs_proto",
        "//intrinsic/math/proto:matrix_proto",
    ],
)

cc_proto_library(
    name = "transformed_geometry_storage_refs_cc_proto",
    deps = [":transformed_geometry_storage_refs_proto"],
)

py_proto_library(
    name = "transformed_geometry_storage_refs_py_pb2",
    deps = [":transformed_geometry_storage_refs_proto"],
)

go_proto_library(
    name = "transformed_geometry_storage_refs_go_proto",
    go_deps = [
        ":geometry_storage_refs_go_proto",
        "//intrinsic/math/proto:matrix_go_proto",
    ],
    deps = [":transformed_geometry_storage_refs_proto"],
)
