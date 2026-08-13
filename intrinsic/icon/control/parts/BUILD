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
    name = "feature_interfaces",
    hdrs = ["feature_interfaces.h"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/icon/control:joint_acceleration_command",
        "//intrinsic/icon/control:joint_position_command",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/kinematics/types:cartesian_limits",
        "//intrinsic/kinematics/types:joint_limits",
        "//intrinsic/kinematics/types:joint_state",
        "//intrinsic/math:pose3",
        "//intrinsic/math:twist",
        "//intrinsic/util:fixed_vector",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/types:optional",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_library(
    name = "io_block",
    srcs = ["io_block.cc"],
    hdrs = ["io_block.h"],
    deps = [
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/util:fixed_vector",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_library(
    name = "realtime_robot_payload",
    srcs = ["realtime_robot_payload.cc"],
    hdrs = ["realtime_robot_payload.h"],
    visibility = ["//intrinsic/icon:__subpackages__"],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/icon/utils:fixed_str_cat",
        "//intrinsic/icon/utils:fixed_string",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/kinematics:validate_link_parameters",
        "//intrinsic/kinematics/types:to_fixed_string",
        "//intrinsic/math:almost_equals",
        "//intrinsic/math:pose3",
        "//intrinsic/world/proto:robot_payload_cc_proto",
        "//intrinsic/world/robot_payload:robot_payload_base",
    ],
)
