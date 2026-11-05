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
load("@rules_cc//cc:cc_test.bzl", "cc_test")

cc_library(
    name = "scene_object_view",
    srcs = ["scene_object_view.cc"],
    hdrs = ["scene_object_view.h"],
    visibility = [
        "//intrinsic:intrinsic_internal",
    ],
    deps = [
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/scene/validate:scene_object_validation",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/log:absl_check",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_test(
    name = "scene_object_view_test",
    srcs = ["scene_object_view_test.cc"],
    deps = [
        ":scene_object_view",
        "//intrinsic/scene/proto/v1:entity_cc_proto",
        "//intrinsic/scene/proto/v1:scene_object_cc_proto",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/world/proto:kinematics_component_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:status_matchers",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/types:span",
        "@com_github_inazarenko_protobuf_matchers//protobuf-matchers",
        "@com_google_cel_cpp//internal:testing_no_main",
        "@com_google_googletest//:gtest_main",
    ],
)
