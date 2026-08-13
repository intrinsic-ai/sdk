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

""" Intrinsic package geometry/legacy/utils. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = [
    "//intrinsic/geometry:__subpackages__",
    "//intrinsic/perception:__subpackages__",
])

cc_library(
    name = "export_as_gltf",
    srcs = ["export_as_gltf.cc"],
    hdrs = ["export_as_gltf.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/api:exact_geometry",
        "//intrinsic/geometry/api:material",
        "//intrinsic/geometry/internal/legacy/mesh",
        "//intrinsic/geometry/internal/legacy/mesh/io:mesh_to_ai_scene",
        "//intrinsic/geometry/internal/legacy/point_cloud:pts_to_ai_scene",
        "//intrinsic/geometry/internal/legacy/renderable_info",
        "//intrinsic/util/object_store:object_ref",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@assimp",
        "@com_gitlab_libeigen_eigen//:eigen",
        "@or_tools//ortools/base:path",
        "@tinygltf",
    ],
)

cc_library(
    name = "scale_shape",
    srcs = ["scale_shape.cc"],
    hdrs = ["scale_shape.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    visibility = [
        "//intrinsic/geometry:__subpackages__",
        "//intrinsic/perception:__subpackages__",
        "//intrinsic/scene/sdf:__pkg__",
        "//intrinsic/simulation/world:__pkg__",
        "//intrinsic/world:__subpackages__",
    ],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/api:exact_geometry",
        "//intrinsic/geometry/shapes",
        "//intrinsic/geometry/shapes:geometry_shape",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)
