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
    name = "triangle",
    hdrs = ["triangle.h"],
    deps = [
        "//intrinsic/eigenmath",
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
    name = "geometry_options",
    hdrs = ["geometry_options.h"],
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
