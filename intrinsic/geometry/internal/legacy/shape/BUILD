# Copyright 2023 Intrinsic Innovation LLC

""" Intrinsic package geometry/legacy/shape. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//intrinsic/geometry:__subpackages__"])

cc_library(
    name = "mesh_from_primitives",
    srcs = ["mesh_from_primitives.cc"],
    hdrs = ["mesh_from_primitives.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        "//intrinsic/geometry/api:affine_transform_of",
        "//intrinsic/geometry/internal/legacy/mesh",
        "//intrinsic/geometry/internal/legacy/mesh:concatenate_meshes",
        "//intrinsic/geometry/internal/legacy/mesh:mesh_primitives",
        "//intrinsic/geometry/shapes",
        "//intrinsic/geometry/shapes:geometry_shape",
        "//intrinsic/math:pose3",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)
