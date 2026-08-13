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

""" Intrinsic package geometry/api. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")

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
