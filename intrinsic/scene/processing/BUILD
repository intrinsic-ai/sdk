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

load("@rules_cc//cc:cc_library.bzl", "cc_library")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "merge_scene_object_links",
    srcs = ["merge_scene_object_links.cc"],
    hdrs = ["merge_scene_object_links.h"],
    deps = [
        "//intrinsic/scene/processing/util:transform_util",
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/scene_object:scene_object_view",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/proto:geometry_component_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_library(
    name = "transform_scene_object",
    srcs = ["transform_scene_object.cc"],
    hdrs = ["transform_scene_object.h"],
    deps = [
        "//intrinsic/math:pose3",
        "//intrinsic/scene/processing/util:transform_util",
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/scene_object:root_entity",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "scale_scene_object",
    srcs = ["scale_scene_object.cc"],
    hdrs = ["scale_scene_object.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/api:io",
        "//intrinsic/math:pose3",
        "//intrinsic/math:proto_conversion",
        "//intrinsic/math/proto:vector3_cc_proto",
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/proto:geometry_component_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)
