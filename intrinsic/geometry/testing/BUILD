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

""" Intrinsic package geometry/testing. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "gltf_matchers",
    testonly = True,
    hdrs = ["gltf_matchers.h"],
    visibility = [
        "//intrinsic/geometry:__subpackages__",
    ],
    deps = [
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/api:material_properties",
        "//intrinsic/geometry/api:renderable",
        "//intrinsic/geometry/internal/util:gltf_helpers",
        "//intrinsic/geometry/proto/v1:material_cc_proto",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:str_format",
        "@com_google_googletest//:gtest",
        "@tinygltf",
    ],
)

cc_test(
    name = "gltf_matchers_test",
    srcs = ["gltf_matchers_test.cc"],
    deps = [
        ":gltf_matchers",
        "//intrinsic/eigenmath",
        "//intrinsic/geometry/api:exact_geometry",
        "//intrinsic/geometry/api:geometry",
        "//intrinsic/geometry/api:material_properties",
        "//intrinsic/geometry/api:renderable",
        "//intrinsic/geometry/internal/util:gltf_helpers",
        "//intrinsic/geometry/proto/v1:material_cc_proto",
        "//intrinsic/geometry/shapes",
        "@abseil-cpp//absl/status:statusor",
        "@com_google_googletest//:gtest_main",
        "@tinygltf",
    ],
)
