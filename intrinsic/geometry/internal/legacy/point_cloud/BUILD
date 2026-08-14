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

""" Intrinsic package geometry/legacy/point_cloud. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = [
    "//intrinsic/geometry:__subpackages__",
])

exports_files(
    srcs = [
        "test_data/Raw_IPS.ply",
        "test_data/example_point_cloud.ply",
    ],
    visibility = [
        "//intrinsic/geometry:__subpackages__",
        "//intrinsic/world/service/benchmarks:__subpackages__",
    ],
)

cc_library(
    name = "get_bounding_box_from_point_cloud",
    srcs = ["get_bounding_box_from_point_cloud.cc"],
    hdrs = ["get_bounding_box_from_point_cloud.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        "//intrinsic/geometry/api:axis_aligned_bounding_box_3d",
        "//intrinsic/geometry/shapes",
    ],
)

cc_library(
    name = "pts_to_ai_scene",
    srcs = ["pts_to_ai_scene.cc"],
    hdrs = ["pts_to_ai_scene.h"],
    visibility = [
        "//intrinsic/geometry/internal/legacy:__subpackages__",
    ],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/shapes",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@assimp",
    ],
)

cc_library(
    name = "pts_to_point_cloud",
    srcs = ["pts_to_point_cloud.cc"],
    hdrs = ["pts_to_point_cloud.h"],
    visibility = [
        "//intrinsic/geometry:__subpackages__",
    ],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/shapes",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@assimp",
    ],
)

cc_library(
    name = "point_cloud_riegeli_coder",
    srcs = ["point_cloud_riegeli_coder.cc"],
    hdrs = ["point_cloud_riegeli_coder.h"],
    copts = ["-fexceptions"] + select({
        "@platforms//cpu:x86_32": ["-frounding-math"],
        "@platforms//cpu:x86_64": ["-frounding-math"],
        "//conditions:default": [],
    }),
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/proto/v1:point_cloud_cc_proto",
        "//intrinsic/geometry/shapes",
        "//intrinsic/marshal:riegeli_proto_coder",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)
