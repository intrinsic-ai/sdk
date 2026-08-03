# Copyright 2023 Intrinsic Innovation LLC

""" Intrinsic package geometry/legacy/mesh/io. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")

cc_library(
    name = "ai_scene_to_mesh",
    srcs = ["ai_scene_to_mesh.cc"],
    hdrs = ["ai_scene_to_mesh.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    visibility = ["//intrinsic/geometry:__subpackages__"],
    deps = [
        "//intrinsic/geometry/internal/legacy/mesh",
        "//intrinsic/geometry/internal/legacy/mesh:concatenate_meshes",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@assimp",
    ],
)

cc_library(
    name = "load_ai_scene_from_buffer",
    srcs = ["load_ai_scene_from_buffer.cc"],
    hdrs = ["load_ai_scene_from_buffer.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    visibility = [
        "//intrinsic/geometry/api:__pkg__",
        "//intrinsic/geometry/internal/legacy/mesh:__pkg__",
        "//intrinsic/geometry/internal/legacy/point_cloud:__pkg__",
    ],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/internal/legacy/point_cloud:pts_to_ai_scene",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@assimp",
    ],
)

cc_library(
    name = "load_ai_scene_from_file",
    srcs = ["load_ai_scene_from_file.cc"],
    hdrs = ["load_ai_scene_from_file.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    visibility = [
        "//intrinsic/geometry/api:__pkg__",
        "//intrinsic/geometry/compatibility:__pkg__",
        "//intrinsic/geometry/internal/legacy/mesh:__pkg__",
        "//intrinsic/geometry/internal/legacy/point_cloud:__pkg__",
    ],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/internal/legacy/point_cloud:pts_to_ai_scene",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@assimp",
        "@or_tools//ortools/base:file",
        "@or_tools//ortools/base:path",
    ],
)

cc_library(
    name = "load_mesh_from_file",
    srcs = ["load_mesh_from_file.cc"],
    hdrs = ["load_mesh_from_file.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    visibility = [
        "//intrinsic/geometry/benchmarks:__pkg__",
        "//intrinsic/geometry/internal/distance:__pkg__",
        "//intrinsic/geometry/internal/distance_3:__pkg__",
        "//intrinsic/geometry/internal/legacy/mesh:__pkg__",
        "//intrinsic/geometry/internal/testing:__pkg__",
        "//intrinsic/scene/sdf:__pkg__",
    ],
    deps = [
        ":ai_scene_to_mesh",
        ":load_ai_scene_from_file",
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/internal/legacy/mesh",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status:statusor",
        "@assimp",
    ],
)

cc_library(
    name = "mesh_to_ai_scene",
    srcs = ["mesh_to_ai_scene.cc"],
    hdrs = ["mesh_to_ai_scene.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    visibility = ["//intrinsic/geometry/internal/legacy/utils:__pkg__"],
    deps = [
        "//intrinsic/geometry/api:material",
        "//intrinsic/geometry/internal/legacy/mesh",
        "@assimp",
        "@com_gitlab_libeigen_eigen//:eigen",
    ],
)
