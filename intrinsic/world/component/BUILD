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
    name = "kinematics_component",
    srcs = ["kinematics_component.cc"],
    hdrs = ["kinematics_component.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/math:pose3",
        "//intrinsic/math:proto_conversion",
        "//intrinsic/math/proto:vector3_cc_proto",
        "//intrinsic/scene:constants",
        "//intrinsic/util/status:status_builder",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world:labels",
        "//intrinsic/world/proto:kinematics_component_cc_proto",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "physics_component",
    srcs = ["physics_component.cc"],
    hdrs = ["physics_component.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/math:pose3",
        "//intrinsic/math:proto_conversion",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/proto:physics_component_cc_proto",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "sensor_component",
    srcs = ["sensor_component.cc"],
    hdrs = ["sensor_component.h"],
    deps = [
        "//intrinsic/util/status:status_builder",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/proto:sensor_component_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "geometry_component",
    srcs = ["geometry_component.cc"],
    hdrs = ["geometry_component.h"],
    visibility = ["//intrinsic:intrinsic_internal"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/api:affine_transform_of",
        "//intrinsic/geometry/api:affine_transform_of_geometry",
        "//intrinsic/geometry/api:exact_geometry",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/api:geometry_options",
        "//intrinsic/geometry/api:io",
        "//intrinsic/geometry/proto/v1:geometry_cc_proto",
        "//intrinsic/geometry/proto/v1:transformed_geometry_cc_proto",
        "//intrinsic/geometry/storage:geometry_deserializer",
        "//intrinsic/geometry/storage:geometry_serializer",
        "//intrinsic/math:proto_conversion",
        "//intrinsic/math/proto:matrix_cc_proto",
        "//intrinsic/util:eigen",
        "//intrinsic/util/status:status_builder",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/hashing",
        "//intrinsic/world/proto:geometry_component_cc_proto",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_google_protobuf//:protobuf_lite",
    ],
)
