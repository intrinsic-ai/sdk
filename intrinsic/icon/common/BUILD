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

package(
    default_visibility = [
        "//visibility:public",
    ],
)

cc_library(
    name = "builtins",
    hdrs = ["builtins.h"],
)

cc_library(
    name = "id_types",
    hdrs = ["id_types.h"],
    deps = ["@or_tools//ortools/base:strong_int"],
)

cc_library(
    name = "slot_part_map",
    srcs = ["slot_part_map.cc"],
    hdrs = ["slot_part_map.h"],
    deps = [
        "//intrinsic/icon/proto/v1:types_cc_proto",
        "@abseil-cpp//absl/container:btree",
        "@abseil-cpp//absl/container:flat_hash_map",
    ],
)

cc_library(
    name = "part_properties",
    srcs = ["part_properties.cc"],
    hdrs = ["part_properties.h"],
    deps = [
        "//intrinsic/icon/proto/v1:service_cc_proto",
        "@abseil-cpp//absl/container:flat_hash_map",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/time",
        "@or_tools//ortools/base:strong_int",
    ],
)

cc_library(
    name = "state_variable_path_constants",
    hdrs = [
        "state_variable_path_constants.h",
    ],
)

cc_library(
    name = "state_variable_path_util",
    srcs = [
        "state_variable_path_util.cc",
    ],
    hdrs = [
        "state_variable_path_util.h",
    ],
    deps = [
        ":state_variable_path_constants",
        "@abseil-cpp//absl/algorithm:container",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/types:span",
    ],
)
