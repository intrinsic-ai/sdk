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

package(default_visibility = [
    "//visibility:public",
])

cc_library(
    name = "init_intrinsic",
    srcs = ["init_intrinsic.cc"],
    hdrs = ["init_intrinsic.h"],
    deps = [
        "//intrinsic/icon/utils:log",
        "@abseil-cpp//absl/debugging:failure_signal_handler",
        "@abseil-cpp//absl/debugging:symbolize",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/flags:parse",
        "@abseil-cpp//absl/flags:usage",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:flags",
        "@abseil-cpp//absl/log:globals",
        "@abseil-cpp//absl/log:initialize",
        "@abseil-cpp//absl/time",
    ],
)

cc_test(
    name = "init_intrinsic_test",
    srcs = ["init_intrinsic_test.cc"],
    deps = [
        ":init_intrinsic",
        "@abseil-cpp//absl/base:log_severity",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log:globals",
        "@abseil-cpp//absl/strings",
    ],
)
