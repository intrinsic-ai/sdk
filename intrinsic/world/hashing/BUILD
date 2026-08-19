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
    name = "hashing",
    hdrs = ["hashing.h"],
    visibility = [
        "//intrinsic:intrinsic_internal",
        "//intrinsic/motion_planner/viewer:__pkg__",
    ],
    deps = [
        "//intrinsic/util:string_type",
        "//intrinsic/world:entity_id",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/container:flat_hash_set",
    ],
)

cc_test(
    name = "hashing_test",
    srcs = ["hashing_test.cc"],
    deps = [
        ":hashing",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/world:entity_id",
        "//intrinsic/world:labels",
    ],
)
