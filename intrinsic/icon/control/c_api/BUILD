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
    default_visibility = ["//visibility:public"],
)

cc_library(
    name = "icon_c_api",
    hdrs = [
        "c_action_factory_context.h",
        "c_feature_interfaces.h",
        "c_plugin_api.h",
        "c_realtime_signal_access.h",
        "c_realtime_slot_map.h",
        "c_realtime_status.h",
        "c_rtcl_action.h",
        "c_streaming_io_realtime_access.h",
        "c_types.h",
    ],
)

cc_library(
    name = "convert_c_realtime_status",
    srcs = ["convert_c_realtime_status.cc"],
    hdrs = ["convert_c_realtime_status.h"],
    deps = [
        ":icon_c_api",
        "//intrinsic/icon/utils:realtime_status",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "convert_c_types",
    srcs = ["convert_c_types.cc"],
    hdrs = ["convert_c_types.h"],
    deps = [
        ":icon_c_api",
        "//intrinsic/eigenmath",
        "//intrinsic/icon/control:joint_position_command",
        "//intrinsic/icon/control:realtime_signal_types",
        "//intrinsic/kinematics/types:joint_limits",
        "//intrinsic/kinematics/types:joint_state",
        "//intrinsic/math:pose3",
        "//intrinsic/math:twist",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/types:optional",
        "@com_gitlab_libeigen_eigen//:eigen",
    ],
)
