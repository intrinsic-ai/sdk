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

# Provides the bindings and other utility with externally built zenoh pubsub

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")

exports_files(
    [
        "zenoh_handle.h",
    ],
    visibility = [
        "//intrinsic:intrinsic_internal",
    ],
)

cc_library(
    name = "zenoh_handle",
    srcs = ["zenoh_handle.cc"],
    hdrs = ["zenoh_handle.h"],
    visibility = [
        "//intrinsic/platform/pubsub:__subpackages__",
    ],
    deps = [
        ":liveliness_get_context",
        ":zenoh_helpers",
        "//intrinsic/middleware/zenoh:imw_zenoh",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_test(
    name = "zenoh_handle_test",
    size = "small",
    srcs = ["zenoh_handle_test.cc"],
    deps = [
        ":zenoh_handle",
        "//intrinsic/util/testing:gtest_wrapper_main",
    ],
)

cc_library(
    name = "liveliness_get_context",
    hdrs = ["liveliness_get_context.h"],
    visibility = ["//intrinsic/platform/pubsub:__subpackages__"],
    deps = [
        "@abseil-cpp//absl/strings:string_view",
    ],
)

cc_library(
    name = "zenoh_helpers",
    srcs = ["zenoh_helpers.cc"],
    hdrs = ["zenoh_helpers.h"],
    visibility = [
        "//intrinsic/platform/pubsub:__subpackages__",
    ],
    deps = [
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/strings:string_view",
    ],
)

cc_library(
    name = "zenoh_config",
    srcs = ["zenoh_config.cc"],
    hdrs = ["zenoh_config.h"],
    visibility = [
        "//intrinsic/platform/pubsub:__subpackages__",
    ],
    deps = [
        ":zenoh_helpers",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/strings:string_view",
        "@com_github_nlohmann_json//:json",
    ],
)

cc_test(
    name = "zenoh_config_test",
    size = "small",
    srcs = ["zenoh_config_test.cc"],
    deps = [
        ":zenoh_config",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/flags:flag",
        "@com_github_nlohmann_json//:json",
    ],
)
