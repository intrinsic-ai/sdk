# Copyright 2023 Intrinsic Innovation LLC

""" Intrinsic package geometry/legacy/mesh. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//intrinsic/geometry:__subpackages__"])

cc_library(
    name = "concatenate_meshes",
    srcs = ["concatenate_meshes.cc"],
    hdrs = ["concatenate_meshes.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":mesh",
        ":remove_duplicate_vertices_from_mesh",
        "@abseil-cpp//absl/log:check",
    ],
)

cc_library(
    name = "mesh",
    srcs = ["mesh.cc"],
    hdrs = ["mesh.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    visibility = [
        "//intrinsic/geometry:__subpackages__",
        "//intrinsic/scene/service:__pkg__",
        "//intrinsic/scene/usd:__subpackages__",
    ],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/api:triangle",
        "//intrinsic/geometry/internal/util:eigen_serialization",
        "//intrinsic/geometry/proto:mesh_cc_proto",
        "//intrinsic/geometry/proto:triangle_mesh_cc_proto",
        "//intrinsic/geometry/proto/v1:exact_geometry_cc_proto",
        "//intrinsic/geometry/proto/v1:triangle_mesh_cc_proto",
        "//intrinsic/marshal:riegeli_proto_coder",
        "//intrinsic/math:pose3",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "mesh_primitives",
    srcs = ["mesh_primitives.cc"],
    hdrs = ["mesh_primitives.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    visibility = [
        "//intrinsic/geometry:__subpackages__",
        "//intrinsic/scene:__subpackages__",
        "//intrinsic/world/conversion:__subpackages__",
    ],
    deps = [
        ":mesh",
        "//intrinsic/eigenmath",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/numeric:bits",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@com_gitlab_libeigen_eigen//:eigen",
    ],
)

cc_library(
    name = "remove_duplicate_vertices_from_mesh",
    srcs = ["remove_duplicate_vertices_from_mesh.cc"],
    hdrs = ["remove_duplicate_vertices_from_mesh.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":mesh",
        "//intrinsic/eigenmath:eigen_matrix_hash",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/meta:type_traits",
        "@com_gitlab_libeigen_eigen//:eigen",
    ],
)
