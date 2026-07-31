# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

licenses(["unencumbered"])

cc_library(
    name = "geometry_shape",
    hdrs = [
        "shape_base.h",
        "shape_ptrs.h",
    ],
    deps = [
        "//intrinsic/geometry/api:affine_transform_of",
        "//intrinsic/math:pose3",
        "@abseil-cpp//absl/log:check",
    ],
)

cc_library(
    name = "shapes",
    hdrs = [
        "axes.h",
        "box.h",
        "capsule.h",
        "convex_hull.h",
        "cylinder.h",
        "discs.h",
        "ellipsoid.h",
        "frustum.h",
        "line_strip.h",
        "mesh_file.h",
        "point_cloud.h",
        "shapes.h",
        "sphere.h",
        "spheres.h",
        "triangle_mesh.h",
    ],
    deps = [
        ":geometry_shape",
        "//intrinsic/eigenmath",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)
