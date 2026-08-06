# Copyright 2023 Intrinsic Innovation LLC

""" Intrinsic package geometry/api. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "affine_transform_of",
    hdrs = ["affine_transform_of.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/math:pose3",
    ],
)

cc_library(
    name = "affine_transform_of_geometry",
    hdrs = ["affine_transform_of_geometry.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":affine_transform_of",
        ":geometry",
    ],
)

cc_library(
    name = "apply_material_properties",
    srcs = ["apply_material_properties.cc"],
    hdrs = ["apply_material_properties.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":material",
        ":renderable",
        "//intrinsic/geometry/proto/v1:material_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@assimp",
        "@com_google_googleapis//google/type:color_cc_proto",
        "@com_google_protobuf//:protobuf",
        "@tinygltf",
    ],
)

cc_library(
    name = "material_properties",
    srcs = ["material_properties.cc"],
    hdrs = ["material_properties.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/proto/v1:material_cc_proto",
    ],
)

cc_library(
    name = "apply_transform",
    srcs = ["apply_transform.cc"],
    hdrs = ["apply_transform.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":exact_geometry",
        ":geometry",
        ":geometry_options",
        ":renderable",
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/internal/legacy/mesh",
        "//intrinsic/geometry/internal/legacy/utils:export_as_gltf",
        "//intrinsic/math:pose3",
        "//intrinsic/util/object_store",
        "//intrinsic/util/object_store:object_ref",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "triangle",
    hdrs = ["triangle.h"],
    deps = [
        "//intrinsic/eigenmath",
    ],
)

cc_library(
    name = "axis_aligned_bounding_box_3d",
    srcs = ["axis_aligned_bounding_box_3d.cc"],
    hdrs = ["axis_aligned_bounding_box_3d.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":triangle",
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/proto:axis_aligned_bounding_box_cc_proto",
        "//intrinsic/geometry/proto/v1:axis_aligned_bounding_box_cc_proto",
        "//intrinsic/marshal:riegeli_proto_coder",
        "//intrinsic/math/proto:vector3_cc_proto",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "compute_axis_aligned_bounding_box_3d",
    srcs = ["compute_axis_aligned_bounding_box_3d.cc"],
    hdrs = ["compute_axis_aligned_bounding_box_3d.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    visibility = [
        "//intrinsic/scene/conversion:__pkg__",
        "//intrinsic/scene/mesh:__pkg__",
        "//intrinsic/scene/sdf:__pkg__",
        "//intrinsic/scene/service:__pkg__",
        "//intrinsic/scene/util:__pkg__",
        "//intrinsic/scene/validate:__pkg__",
        "//intrinsic/world/conversion/sdf:__pkg__",
        "//intrinsic/world/util:__pkg__",
    ],
    deps = [
        ":apply_transform",
        ":axis_aligned_bounding_box_3d",
        ":exact_geometry",
        ":geometry",
        "//intrinsic/geometry/internal/legacy/mesh",
        "//intrinsic/geometry/internal/legacy/point_cloud:get_bounding_box_from_point_cloud",
        "//intrinsic/geometry/shapes",
        "//intrinsic/math:pose3",
        "//intrinsic/util/object_store:object_ref",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "exact_geometry",
    srcs = ["exact_geometry.cc"],
    hdrs = ["exact_geometry.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":affine_transform_of",
        ":geometry_options",
        "//intrinsic/geometry/internal/legacy/mesh",
        "//intrinsic/geometry/internal/legacy/mesh:remove_duplicate_vertices_from_mesh",
        "//intrinsic/geometry/internal/legacy/point_cloud:point_cloud_riegeli_coder",
        "//intrinsic/geometry/internal/legacy/shape:mesh_from_primitives",
        "//intrinsic/geometry/proto:mesh_cc_proto",
        "//intrinsic/geometry/shapes",
        "//intrinsic/geometry/shapes:geometry_shape",
        "//intrinsic/util:macros",
        "//intrinsic/util/object_store",
        "//intrinsic/util/object_store:object_ref",
        "//intrinsic/util/status:status_builder",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "geometry",
    srcs = ["geometry.cc"],
    hdrs = ["geometry.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":exact_geometry",
        ":renderable",
        "//intrinsic/geometry/proto/v1:material_cc_proto",
        "//intrinsic/util/object_store:object_ref",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log:die_if_null",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
    ],
)

cc_library(
    name = "geometry_fingerprint",
    srcs = ["geometry_fingerprint.cc"],
    hdrs = ["geometry_fingerprint.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":exact_geometry",
        ":geometry",
        ":io",
        ":renderable",
        "//intrinsic/util:hash",
        "//intrinsic/util:macros",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/synchronization",
    ],
)

cc_library(
    name = "geometry_options",
    hdrs = ["geometry_options.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
)

cc_library(
    name = "io",
    srcs = ["io.cc"],
    hdrs = ["io.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":affine_transform_of_geometry",
        ":exact_geometry",
        ":geometry",
        ":geometry_options",
        ":renderable",
        "//intrinsic/eigenmath",
        "//intrinsic/eigenmath:rotation_utils",
        "//intrinsic/geometry/internal/legacy/mesh",
        "//intrinsic/geometry/internal/legacy/point_cloud:point_cloud_riegeli_coder",
        "//intrinsic/geometry/proto/v1:exact_geometry_cc_proto",
        "//intrinsic/geometry/proto/v1:geometric_transform_cc_proto",
        "//intrinsic/geometry/proto/v1:geometry_cc_proto",
        "//intrinsic/geometry/proto/v1:geometry_options_cc_proto",
        "//intrinsic/geometry/proto/v1:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto/v1:inline_geometry_cc_proto",
        "//intrinsic/geometry/proto/v1:octree_cc_proto",
        "//intrinsic/geometry/proto/v1:octree_wrapping_cc_proto",
        "//intrinsic/geometry/proto/v1:point_cloud_cc_proto",
        "//intrinsic/geometry/proto/v1:primitive_shape_cc_proto",
        "//intrinsic/geometry/proto/v1:primitives_cc_proto",
        "//intrinsic/geometry/proto/v1:renderable_cc_proto",
        "//intrinsic/geometry/proto/v1:transformed_geometry_cc_proto",
        "//intrinsic/geometry/proto/v1:transformed_primitive_shape_cc_proto",
        "//intrinsic/geometry/proto/v1:transformed_primitive_shape_set_cc_proto",
        "//intrinsic/geometry/proto/v1:triangle_mesh_cc_proto",
        "//intrinsic/geometry/shapes",
        "//intrinsic/geometry/shapes:geometry_shape",
        "//intrinsic/geometry/storage:geometry_deserializer",
        "//intrinsic/geometry/storage:geometry_serializer",
        "//intrinsic/math:proto_conversion",
        "//intrinsic/math/proto:vector3_cc_proto",
        "//intrinsic/util/object_store",
        "//intrinsic/util/object_store:object_ref",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "material",
    hdrs = ["material.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
)

cc_library(
    name = "renderable",
    srcs = ["renderable.cc"],
    hdrs = ["renderable.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
)

cc_library(
    name = "renderable_generation",
    srcs = ["renderable_generation.cc"],
    hdrs = ["renderable_generation.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":geometry",
        ":material_properties",
        ":renderable",
        "//intrinsic/geometry/api:apply_material_properties",
        "//intrinsic/geometry/internal/legacy/utils:export_as_gltf",
        "//intrinsic/util/status:ret_check",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "validate_mesh",
    srcs = ["validate_mesh.cc"],
    hdrs = ["validate_mesh.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/internal/legacy/mesh/io:load_ai_scene_from_buffer",
        "//intrinsic/geometry/internal/legacy/mesh/io:load_ai_scene_from_file",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/container:flat_hash_set",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@assimp",
        "@or_tools//ortools/base:file",
        "@or_tools//ortools/base:path",
    ],
)
