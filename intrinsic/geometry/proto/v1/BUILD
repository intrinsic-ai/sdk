# Copyright 2023 Intrinsic Innovation LLC

""" Intrinsic package geometry/proto/v1. """

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
    protos = [":axis_aligned_bounding_box_proto"],
    deps = [
        "//intrinsic/math/proto:vector2_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
    ],
)

proto_library(
    name = "geometric_transform_proto",
    srcs = ["geometric_transform.proto"],
    deps = [
        "//intrinsic/math/proto:matrix_proto",
        "//intrinsic/math/proto:vector3_proto",
    ],
)

cc_proto_library(
    name = "geometric_transform_cc_proto",
    deps = [":geometric_transform_proto"],
)

py_proto_library(
    name = "geometric_transform_py_pb2",
    deps = [":geometric_transform_proto"],
)

go_proto_library(
    name = "geometric_transform_go_proto",
    protos = [":geometric_transform_proto"],
    deps = [
        "//intrinsic/math/proto:matrix_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
    ],
)

proto_library(
    name = "geometry_proto",
    srcs = ["geometry.proto"],
    deps = [
        ":geometry_storage_refs_proto",
        ":inline_geometry_proto",
        ":material_proto",
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

go_proto_library(
    name = "geometry_go_proto",
    protos = [":geometry_proto"],
    deps = [
        ":geometry_storage_refs_go_proto",
        ":inline_geometry_go_proto",
        ":material_go_proto",
    ],
)

proto_library(
    name = "geometry_options_proto",
    srcs = ["geometry_options.proto"],
)

cc_proto_library(
    name = "geometry_options_cc_proto",
    deps = [":geometry_options_proto"],
)

py_proto_library(
    name = "geometry_options_py_pb2",
    deps = [":geometry_options_proto"],
)

go_proto_library(
    name = "geometry_options_go_proto",
    protos = [":geometry_options_proto"],
)

proto_library(
    name = "geometry_storage_refs_proto",
    srcs = ["geometry_storage_refs.proto"],
)

cc_proto_library(
    name = "geometry_storage_refs_cc_proto",
    deps = [":geometry_storage_refs_proto"],
)

py_proto_library(
    name = "geometry_storage_refs_py_pb2",
    deps = [":geometry_storage_refs_proto"],
)

go_proto_library(
    name = "geometry_storage_refs_go_proto",
    protos = [":geometry_storage_refs_proto"],
)

proto_library(
    name = "inline_geometry_proto",
    srcs = ["inline_geometry.proto"],
    deps = [
        ":exact_geometry_proto",
        ":renderable_proto",
    ],
)

cc_proto_library(
    name = "inline_geometry_cc_proto",
    deps = [":inline_geometry_proto"],
)

py_proto_library(
    name = "inline_geometry_py_pb2",
    deps = [":inline_geometry_proto"],
)

go_proto_library(
    name = "inline_geometry_go_proto",
    protos = [":inline_geometry_proto"],
    deps = [
        ":exact_geometry_go_proto",
        ":renderable_go_proto",
    ],
)

proto_library(
    name = "exact_geometry_proto",
    srcs = ["exact_geometry.proto"],
    deps = [
        ":geometry_options_proto",
        ":octree_wrapping_proto",
        ":transformed_primitive_shape_set_proto",
        ":triangle_mesh_proto",
    ],
)

cc_proto_library(
    name = "exact_geometry_cc_proto",
    deps = [":exact_geometry_proto"],
)

py_proto_library(
    name = "exact_geometry_py_pb2",
    deps = [":exact_geometry_proto"],
)

go_proto_library(
    name = "exact_geometry_go_proto",
    protos = [":exact_geometry_proto"],
    deps = [
        ":geometry_options_go_proto",
        ":octree_wrapping_go_proto",
        ":transformed_primitive_shape_set_go_proto",
        ":triangle_mesh_go_proto",
    ],
)

proto_library(
    name = "material_proto",
    srcs = ["material.proto"],
    deps = ["@com_google_googleapis//google/type:color_proto"],
)

cc_proto_library(
    name = "material_cc_proto",
    deps = [":material_proto"],
)

py_proto_library(
    name = "material_py_pb2",
    deps = [":material_proto"],
)

go_proto_library(
    name = "material_go_proto",
    protos = [":material_proto"],
    deps = [
        "@org_golang_google_genproto//googleapis/type/color",
    ],
)

proto_library(
    name = "octree_proto",
    srcs = ["octree.proto"],
    deps = [
        ":axis_aligned_bounding_box_proto",
    ],
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
    protos = [":octree_proto"],
    deps = [
        ":axis_aligned_bounding_box_go_proto",
    ],
)

proto_library(
    name = "octree_wrapping_proto",
    srcs = ["octree_wrapping.proto"],
    deps = [
        ":octree_proto",
    ],
)

cc_proto_library(
    name = "octree_wrapping_cc_proto",
    deps = [":octree_wrapping_proto"],
)

py_proto_library(
    name = "octree_wrapping_py_pb2",
    deps = [":octree_wrapping_proto"],
)

go_proto_library(
    name = "octree_wrapping_go_proto",
    protos = [":octree_wrapping_proto"],
    deps = [
        ":octree_go_proto",
    ],
)

proto_library(
    name = "point_cloud_proto",
    srcs = ["point_cloud.proto"],
)

cc_proto_library(
    name = "point_cloud_cc_proto",
    deps = [":point_cloud_proto"],
)

py_proto_library(
    name = "point_cloud_py_pb2",
    deps = [":point_cloud_proto"],
)

go_proto_library(
    name = "point_cloud_go_proto",
    protos = [":point_cloud_proto"],
)

proto_library(
    name = "primitive_shape_proto",
    srcs = ["primitive_shape.proto"],
    deps = [
        ":primitives_proto",
    ],
)

cc_proto_library(
    name = "primitive_shape_cc_proto",
    deps = [":primitive_shape_proto"],
)

py_proto_library(
    name = "primitive_shape_py_pb2",
    deps = [":primitive_shape_proto"],
)

go_proto_library(
    name = "primitive_shape_go_proto",
    protos = [":primitive_shape_proto"],
    deps = [
        ":primitives_go_proto",
    ],
)

proto_library(
    name = "primitives_proto",
    srcs = ["primitives.proto"],
    deps = [
        "//intrinsic/math/proto:vector3_proto",
    ],
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

go_proto_library(
    name = "renderable_go_proto",
    protos = [":renderable_proto"],
)

proto_library(
    name = "transformed_geometry_proto",
    srcs = ["transformed_geometry.proto"],
    deps = [
        ":geometric_transform_proto",
        ":geometry_proto",
    ],
)

cc_proto_library(
    name = "transformed_geometry_cc_proto",
    deps = [":transformed_geometry_proto"],
)

py_proto_library(
    name = "transformed_geometry_py_pb2",
    deps = [":transformed_geometry_proto"],
)

go_proto_library(
    name = "transformed_geometry_go_proto",
    protos = [":transformed_geometry_proto"],
    deps = [
        ":geometric_transform_go_proto",
        ":geometry_go_proto",
    ],
)

proto_library(
    name = "transformed_primitive_shape_proto",
    srcs = ["transformed_primitive_shape.proto"],
    deps = [
        ":geometric_transform_proto",
        ":primitive_shape_proto",
    ],
)

cc_proto_library(
    name = "transformed_primitive_shape_cc_proto",
    deps = [":transformed_primitive_shape_proto"],
)

py_proto_library(
    name = "transformed_primitive_shape_py_pb2",
    deps = [":transformed_primitive_shape_proto"],
)

go_proto_library(
    name = "transformed_primitive_shape_go_proto",
    protos = [":transformed_primitive_shape_proto"],
    deps = [
        ":geometric_transform_go_proto",
        ":primitive_shape_go_proto",
    ],
)

proto_library(
    name = "transformed_primitive_shape_set_proto",
    srcs = ["transformed_primitive_shape_set.proto"],
    deps = [
        ":transformed_primitive_shape_proto",
    ],
)

cc_proto_library(
    name = "transformed_primitive_shape_set_cc_proto",
    deps = [":transformed_primitive_shape_set_proto"],
)

py_proto_library(
    name = "transformed_primitive_shape_set_py_pb2",
    deps = [":transformed_primitive_shape_set_proto"],
)

go_proto_library(
    name = "transformed_primitive_shape_set_go_proto",
    protos = [":transformed_primitive_shape_set_proto"],
    deps = [
        ":transformed_primitive_shape_go_proto",
    ],
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
    protos = [":triangle_mesh_proto"],
)
