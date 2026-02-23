# Copyright 2023 Intrinsic Innovation LLC

""" Intrinsic package geometry/proto. """

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

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
    importpath = "intrinsic/geometry/proto/axis_aligned_bounding_box_go_proto",
    protos = [":axis_aligned_bounding_box_proto"],
    deps = [
        "//intrinsic/math/proto:vector2_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
    ],
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
    importpath = "intrinsic/geometry/proto/oriented_bounding_box_go_proto",
    protos = [":oriented_bounding_box_proto"],
    deps = [
        ":axis_aligned_bounding_box_go_proto",
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/math/proto:quaternion_go_proto",
        "//intrinsic/math/proto:vector2_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
    ],
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
        "//intrinsic/geometry/proto/v1:geometry_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_proto",
        "@com_google_googleapis//google/api:annotations_proto",
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
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
        "//bazel:go_gen_grpc_gateway",
    ],
    importpath = "intrinsic/geometry/proto/geometry_service_go_proto",
    protos = [":geometry_service_proto"],
    deps = [
        ":geometry_service_types_go_proto",
        ":geometry_storage_refs_go_proto",
        "//intrinsic/geometry/proto/v1:geometry_go_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
        "@org_golang_google_protobuf//types/known/anypb",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

cc_grpc_library(
    name = "geometry_service_cc_grpc_proto",
    srcs = [":geometry_service_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":geometry_service_cc_proto"],
)

proto_library(
    name = "geometry_service_types_proto",
    srcs = ["geometry_service_types.proto"],
    deps = [
        ":geometry_proto",
        ":geometry_storage_refs_proto",
        ":lazy_exact_geometry_proto",
        ":renderable_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_proto",
        "//intrinsic/geometry/proto/v1:inline_geometry_proto",
        "//intrinsic/geometry/proto/v1:renderable_proto",
        "//intrinsic/geometry/proto/v1:transformed_primitive_shape_set_proto",
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
    importpath = "intrinsic/geometry/proto/geometry_service_types_go_proto",
    protos = [":geometry_service_types_proto"],
    deps = [
        ":geometry_go_proto",
        ":geometry_storage_refs_go_proto",
        ":lazy_exact_geometry_go_proto",
        ":renderable_go_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_go_proto",
        "//intrinsic/geometry/proto/v1:inline_geometry_go_proto",
        "//intrinsic/geometry/proto/v1:renderable_go_proto",
        "//intrinsic/geometry/proto/v1:transformed_primitive_shape_set_go_proto",
    ],
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
    importpath = "intrinsic/geometry/proto/geometry_storage_refs_go_proto",
    protos = [":geometry_storage_refs_proto"],
)

go_proto_library(
    name = "geometry_go_proto",
    importpath = "intrinsic/geometry/proto/geometry_go_proto",
    protos = [":geometry_proto"],
    deps = [
        ":lazy_exact_geometry_go_proto",
        ":renderable_go_proto",
    ],
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
    importpath = "intrinsic/geometry/proto/lazy_exact_geometry_go_proto",
    protos = [":lazy_exact_geometry_proto"],
    deps = [
        ":mesh_octree_wrapping_go_proto",
        ":primitives_go_proto",
        ":triangle_mesh_go_proto",
    ],
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
    importpath = "intrinsic/geometry/proto/mesh_octree_wrapping_go_proto",
    protos = [":mesh_octree_wrapping_proto"],
    deps = [
        ":octree_go_proto",
    ],
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
    importpath = "intrinsic/geometry/proto/octree_go_proto",
    protos = [":octree_proto"],
    deps = [
        ":axis_aligned_bounding_box_go_proto",
    ],
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
    importpath = "intrinsic/geometry/proto/primitives_go_proto",
    protos = [":primitives_proto"],
    deps = [
        "//intrinsic/math/proto:vector3_go_proto",
    ],
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
    importpath = "intrinsic/geometry/proto/transformed_geometry_storage_refs_go_proto",
    protos = [":transformed_geometry_storage_refs_proto"],
    deps = [
        ":geometry_storage_refs_go_proto",
        "//intrinsic/math/proto:matrix_go_proto",
    ],
)

go_proto_library(
    name = "renderable_go_proto",
    importpath = "intrinsic/geometry/proto/renderable_go_proto",
    protos = [":renderable_proto"],
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
    importpath = "intrinsic/geometry/proto/triangle_mesh_go_proto",
    protos = [":triangle_mesh_proto"],
)
