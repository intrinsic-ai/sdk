# Copyright 2023 Intrinsic Innovation LLC

""" Intrinsic package geometry/api. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "shapes",
    srcs = ["shapes.cc"],
    hdrs = ["shapes.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/api:affine_transform_of_geometry",
        "//intrinsic/geometry/api:apply_material_properties",
        "//intrinsic/geometry/api:exact_geometry",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/api:geometry_options",
        "//intrinsic/geometry/api:material",
        "//intrinsic/geometry/api:renderable",
        "//intrinsic/geometry/internal/legacy/mesh:load_mesh_file_with_scale",
        "//intrinsic/geometry/internal/legacy/mesh/io:load_ai_scene_from_file",
        "//intrinsic/geometry/internal/legacy/point_cloud:pts_to_point_cloud",
        "//intrinsic/geometry/internal/legacy/utils:export_as_gltf",
        "//intrinsic/geometry/proto/v1:material_cc_proto",
        "//intrinsic/geometry/shapes",
        "//intrinsic/geometry/shapes:geometry_shape",
        "//intrinsic/math:pose3",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@assimp",
        "@or_tools//ortools/base:file",
        "@or_tools//ortools/base:path",
    ],
)

cc_test(
    name = "shapes_test",
    srcs = ["shapes_test.cc"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    data = [
        "//intrinsic/geometry/internal/testdata:cuboid_2_3_4_partially_transparent.mtl",
        "//intrinsic/geometry/internal/testdata:cuboid_2_3_4_partially_transparent.obj",
        "//intrinsic/geometry/internal/testdata:cuboid_2_3_4_transparent.mtl",
        "//intrinsic/geometry/internal/testdata:cuboid_2_3_4_transparent.obj",
    ],
    deps = [
        ":shapes",
        "//intrinsic/eigenmath",
        "//intrinsic/eigenmath:matchers",
        "//intrinsic/geometry/api:affine_transform_of_geometry",
        "//intrinsic/geometry/api:exact_geometry",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/api:material",
        "//intrinsic/geometry/internal/legacy/mesh",
        "//intrinsic/geometry/internal/legacy/mesh:compute_two_sided_hausdorff_distance",
        "//intrinsic/geometry/shapes",
        "//intrinsic/math:pose3",
        "//intrinsic/util/object_store:object_ref",
        "//intrinsic/util/path_resolver",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
    ],
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
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/api:exact_geometry",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/api:geometry_options",
        "//intrinsic/geometry/api:renderable",
        "//intrinsic/geometry/api:renderable_generation",
        "//intrinsic/geometry/internal/legacy/mesh",
        "//intrinsic/geometry/proto:geometry_cc_proto",
        "//intrinsic/geometry/proto:geometry_storage_refs_cc_proto",
        "//intrinsic/geometry/proto:lazy_exact_geometry_cc_proto",
        "//intrinsic/geometry/proto:mesh_octree_wrapping_cc_proto",
        "//intrinsic/geometry/proto:primitives_cc_proto",
        "//intrinsic/geometry/proto:renderable_cc_proto",
        "//intrinsic/geometry/proto:triangle_mesh_cc_proto",
        "//intrinsic/geometry/proto/v1:material_cc_proto",
        "//intrinsic/geometry/shapes",
        "//intrinsic/geometry/shapes:geometry_shape",
        "//intrinsic/math/proto:vector3_cc_proto",
        "//intrinsic/util/object_store",
        "//intrinsic/util/object_store:object_ref",
        "//intrinsic/util/status:status_builder",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@assimp",
        "@com_google_protobuf//:protobuf_lite",
    ],
)

cc_test(
    name = "io_test",
    srcs = ["io_test.cc"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    data = [
        "//intrinsic/geometry/testdata:RGB.pts",
        "//intrinsic/geometry/testdata:Triangle.gltf",
        "//intrinsic/geometry/testdata:TriangleWithoutIndices.gltf",
        "//intrinsic/geometry/testdata:stl_with_material.stl",
    ],
    deps = [
        ":io",
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/api:exact_geometry",
        "//intrinsic/geometry/api:file_io",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/proto:lazy_exact_geometry_cc_proto",
        "//intrinsic/geometry/shapes",
        "//intrinsic/util/path_resolver",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@tinygltf",
    ],
)
