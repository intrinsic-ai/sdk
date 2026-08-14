# Copyright 2026 Intrinsic Innovation LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

""" Intrinsic package geometry/legacy/mesh. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = [
    "//intrinsic/geometry:__subpackages__",
])

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
    name = "load_mesh_file_with_scale",
    srcs = ["load_mesh_file_with_scale.cc"],
    hdrs = ["load_mesh_file_with_scale.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":mesh",
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/internal/legacy/mesh/io:load_mesh_from_file",
        "//intrinsic/scene/sdf:sdf_path_resolver",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
    ],
)

cc_library(
    name = "triangle_mesh_riegeli_coder",
    srcs = ["triangle_mesh_riegeli_coder.cc"],
    hdrs = ["triangle_mesh_riegeli_coder.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/proto/v1:triangle_mesh_cc_proto",
        "//intrinsic/geometry/shapes",
        "//intrinsic/marshal:riegeli_proto_coder",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "memoized_mesh_from_triangle_mesh",
    srcs = ["memoized_mesh_from_triangle_mesh.cc"],
    hdrs = ["memoized_mesh_from_triangle_mesh.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":mesh",
        ":mesh_from_triangle_mesh",
        ":triangle_mesh_riegeli_coder",
        "//intrinsic/geometry/shapes",
        "//intrinsic/util/object_store:memoize",
        "//intrinsic/util/object_store:object_ref",
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
    name = "mesh_from_triangle_mesh",
    srcs = ["mesh_from_triangle_mesh.cc"],
    hdrs = ["mesh_from_triangle_mesh.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        ":mesh",
        "//intrinsic/geometry/shapes",
        "@com_gitlab_libeigen_eigen//:eigen",
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
