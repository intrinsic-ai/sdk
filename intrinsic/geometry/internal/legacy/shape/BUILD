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
