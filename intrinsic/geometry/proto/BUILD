# Copyright 2023 Intrinsic Innovation LLC

""" Intrinsic package geometry/proto. """

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "axis_aligned_bounding_box_proto",
    srcs = ["axis_aligned_bounding_box.proto"],
    deps = [
        "//intrinsic/math/proto:vector2_proto",
        "//intrinsic/math/proto:vector3_proto",
    ],
)

cc_proto_library(
    name = "axis_aligned_bounding_box_cc_proto",
    deps = [":axis_aligned_bounding_box_proto"],
)

py_proto_library(
    name = "axis_aligned_bounding_box_py_pb2",
    deps = [":axis_aligned_bounding_box_proto"],
)

go_proto_library(
    name = "axis_aligned_bounding_box_go_proto",
    go_deps = [
        "//intrinsic/math/proto:vector2_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
    ],
    deps = [":axis_aligned_bounding_box_proto"],
)

proto_library(
    name = "oriented_bounding_box_proto",
    srcs = ["oriented_bounding_box.proto"],
    deps = [
        ":axis_aligned_bounding_box_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/math/proto:vector3_proto",
    ],
)

py_proto_library(
    name = "oriented_bounding_box_py_pb2",
    deps = [":oriented_bounding_box_proto"],
)

cc_proto_library(
    name = "oriented_bounding_box_cc_proto",
    deps = [":oriented_bounding_box_proto"],
)

go_proto_library(
    name = "oriented_bounding_box_go_proto",
    go_deps = [
        "//intrinsic/math/proto:vector2_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
        ":axis_aligned_bounding_box_go_proto",
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/math/proto:quaternion_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
    ],
    deps = [":oriented_bounding_box_proto"],
)

proto_library(
    name = "geometry_proto",
    srcs = ["geometry.proto"],
    deps = [
        ":lazy_exact_geometry_proto",
        ":renderable_proto",
    ],
)

cc_proto_library(
    name = "geometry_cc_proto",
    deps = [":geometry_proto"],
)

py_proto_library(
    name = "geometry_py_pb2",
    deps = [":geometry_proto"],
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
    name = "geometry_service_types_proto",
    srcs = ["geometry_service_types.proto"],
    deps = [
        ":geometry_proto",
        ":geometry_storage_refs_proto",
        ":lazy_exact_geometry_proto",
        ":renderable_proto",
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
        ":geometry_go_proto",
        ":geometry_storage_refs_go_proto",
        ":lazy_exact_geometry_go_proto",
        ":renderable_go_proto",
        "//intrinsic/math/proto:matrix_go_proto",
    ],
    deps = [":geometry_service_types_proto"],
)

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

go_proto_library(
    name = "geometry_go_proto",
    go_deps = [
        ":lazy_exact_geometry_go_proto",
        ":renderable_go_proto",
    ],
    deps = [":geometry_proto"],
)

proto_library(
    name = "lazy_exact_geometry_proto",
    srcs = ["lazy_exact_geometry.proto"],
    deps = [
        ":mesh_octree_wrapping_proto",
        ":primitives_proto",
        ":triangle_mesh_proto",
    ],
)

cc_proto_library(
    name = "lazy_exact_geometry_cc_proto",
    deps = [":lazy_exact_geometry_proto"],
)

py_proto_library(
    name = "lazy_exact_geometry_py_pb2",
    deps = [":lazy_exact_geometry_proto"],
)

go_proto_library(
    name = "lazy_exact_geometry_go_proto",
    go_deps = [
        ":mesh_octree_wrapping_go_proto",
        ":primitives_go_proto",
        ":triangle_mesh_go_proto",
    ],
    deps = [":lazy_exact_geometry_proto"],
)

proto_library(
    name = "mesh_octree_wrapping_proto",
    srcs = ["mesh_octree_wrapping.proto"],
    deps = [":octree_proto"],
)

cc_proto_library(
    name = "mesh_octree_wrapping_cc_proto",
    deps = [":mesh_octree_wrapping_proto"],
)

py_proto_library(
    name = "mesh_octree_wrapping_py_pb2",
    deps = [":mesh_octree_wrapping_proto"],
)

go_proto_library(
    name = "mesh_octree_wrapping_go_proto",
    go_deps = [
        ":octree_go_proto",
    ],
    deps = [":mesh_octree_wrapping_proto"],
)

proto_library(
    name = "octree_proto",
    srcs = ["octree.proto"],
    deps = [":axis_aligned_bounding_box_proto"],
)

cc_proto_library(
    name = "octree_cc_proto",
    deps = [":octree_proto"],
)

py_proto_library(
    name = "octree_py_pb2",
    deps = [":octree_proto"],
)

go_proto_library(
    name = "octree_go_proto",
    go_deps = [
        ":axis_aligned_bounding_box_go_proto",
    ],
    deps = [":octree_proto"],
)

proto_library(
    name = "primitives_proto",
    srcs = ["primitives.proto"],
    deps = ["//intrinsic/math/proto:vector3_proto"],
)

cc_proto_library(
    name = "primitives_cc_proto",
    deps = [":primitives_proto"],
)

py_proto_library(
    name = "primitives_py_pb2",
    deps = [":primitives_proto"],
)

go_proto_library(
    name = "primitives_go_proto",
    go_deps = [
        "//intrinsic/math/proto:vector3_go_proto",
    ],
    deps = [":primitives_proto"],
)

proto_library(
    name = "renderable_proto",
    srcs = ["renderable.proto"],
)

cc_proto_library(
    name = "renderable_cc_proto",
    deps = [":renderable_proto"],
)

py_proto_library(
    name = "renderable_py_pb2",
    deps = [":renderable_proto"],
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

go_proto_library(
    name = "renderable_go_proto",
    deps = [":renderable_proto"],
)

proto_library(
    name = "triangle_mesh_proto",
    srcs = ["triangle_mesh.proto"],
)

cc_proto_library(
    name = "triangle_mesh_cc_proto",
    deps = [":triangle_mesh_proto"],
)

py_proto_library(
    name = "triangle_mesh_py_pb2",
    deps = [":triangle_mesh_proto"],
)

go_proto_library(
    name = "triangle_mesh_go_proto",
    deps = [":triangle_mesh_proto"],
)
