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

""" Intrinsic package geometry/gltf. """

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "gltf_io",
    srcs = ["gltf_io.cc"],
    hdrs = ["gltf_io.h"],
    visibility = [
        "//intrinsic/geometry:__subpackages__",
    ],
    deps = [
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
        "@tinygltf",
    ],
)

cc_test(
    name = "gltf_io_test",
    srcs = ["gltf_io_test.cc"],
    deps = [
        ":gltf_io",
        "@abseil-cpp//absl/base:log_severity",
        "@abseil-cpp//absl/log:scoped_mock_log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:status_matchers",
        "@abseil-cpp//absl/strings",
        "@com_google_cel_cpp//internal:testing_no_main",
        "@com_google_googletest//:gtest_main",
        "@tinygltf",
    ],
)

cc_library(
    name = "gltf_materials",
    srcs = ["gltf_materials.cc"],
    hdrs = ["gltf_materials.h"],
    visibility = [
        "//intrinsic/geometry:__subpackages__",
    ],
    deps = [
        "@abseil-cpp//absl/strings:string_view",
        "@tinygltf",
    ],
)

cc_test(
    name = "gltf_materials_test",
    srcs = ["gltf_materials_test.cc"],
    deps = [
        ":gltf_io",
        ":gltf_materials",
        "@abseil-cpp//absl/strings:string_view",
        "@com_google_cel_cpp//internal:testing_no_main",
        "@com_google_googletest//:gtest_main",
        "@tinygltf",
    ],
)

cc_library(
    name = "gltf_merge",
    srcs = ["gltf_merge.cc"],
    hdrs = ["gltf_merge.h"],
    visibility = [
        "//intrinsic/geometry:__subpackages__",
    ],
    deps = [
        "@abseil-cpp//absl/strings",
        "@tinygltf",
    ],
)
