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

package(default_visibility = [
    "//visibility:public",
])

cc_library(
    name = "robot_payload",
    srcs = ["robot_payload.cc"],
    hdrs = ["robot_payload.h"],
    deps = [
        ":robot_payload_base",
        "//intrinsic/eigenmath",
        "//intrinsic/kinematics:validate_link_parameters",
        "//intrinsic/math:almost_equals",
        "//intrinsic/math:inertia_utils",
        "//intrinsic/math:pose3",
        "//intrinsic/math:proto_conversion",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/world/proto:robot_payload_cc_proto",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "robot_payload_base",
    srcs = ["robot_payload_base.cc"],
    hdrs = ["robot_payload_base.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/math:almost_equals",
        "//intrinsic/math:pose3",
        "//intrinsic/world/proto:robot_payload_cc_proto",
    ],
)
