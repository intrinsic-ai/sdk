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
    name = "binary_futex",
    srcs = ["binary_futex.cc"],
    hdrs = [
        "binary_futex.h",
    ],
    deps = [
        "//intrinsic/icon/testing:realtime_annotations",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:strerror",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "lockable_binary_futex",
    hdrs = [
        "lockable_binary_futex.h",
    ],
    deps = [
        ":binary_futex",
        "//intrinsic/icon/testing:realtime_annotations",
        "//intrinsic/icon/utils:realtime_status",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "binary_futex_condition_variable",
    srcs = ["binary_futex_condition_variable.cc"],
    hdrs = [
        "binary_futex_condition_variable.h",
    ],
    deps = [
        ":binary_futex",
        ":lockable_binary_futex",
        "//intrinsic/icon/testing:realtime_annotations",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "@abseil-cpp//absl/base:core_headers",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
    ],
)
