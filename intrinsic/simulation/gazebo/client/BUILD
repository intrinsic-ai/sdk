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

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "gazebo_client",
    srcs = ["gazebo_client.cc"],
    hdrs = ["gazebo_client.h"],
    deps = [
        "//intrinsic/assets/dependencies:utils_cc",
        "//intrinsic/assets/proto/v1:resolved_dependency_cc_proto",
        "//intrinsic/connect/cc/grpc:channel",
        "//intrinsic/simulation/gazebo/proto/v1:gazebo_service_cc_grpc_proto",
        "//intrinsic/simulation/gazebo/proto/v1:gazebo_service_cc_proto",
        "//intrinsic/util/grpc:channel",
        "//intrinsic/util/grpc:connection_params",
        "//intrinsic/util/status:status_conversion_grpc",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/time",
        "@gz-transport",
    ],
)

cc_library(
    name = "joint_control",
    srcs = ["joint_control.cc"],
    hdrs = ["joint_control.h"],
    deps = [
        ":gazebo_client",
        "//intrinsic/util/status:status_builder",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@gz-msgs",
        "@gz-transport",
    ],
)

cc_library(
    name = "joint_state",
    srcs = ["joint_state.cc"],
    hdrs = ["joint_state.h"],
    deps = [
        ":gazebo_client",
        "//intrinsic/util/status:ret_check",
        "//intrinsic/util/status:status_builder",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:nullability",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@gz-msgs",
        "@gz-transport",
    ],
)

cc_test(
    name = "gazebo_client_test",
    srcs = ["gazebo_client_test.cc"],
    deps = [
        ":gazebo_client",
        ":gazebo_client_test_env",
        "//intrinsic/assets/proto/v1:resolved_dependency_cc_proto",
        "//intrinsic/simulation/gazebo/proto/v1:gazebo_service_cc_grpc_proto",
        "//intrinsic/simulation/gazebo/proto/v1:gazebo_service_cc_proto",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/time",
        "@com_github_grpc_grpc//:grpc++",
        "@gz-transport",
    ],
)

cc_test(
    name = "joint_control_test",
    srcs = ["joint_control_test.cc"],
    deps = [
        ":gazebo_client",
        ":gazebo_client_test_env",
        ":joint_control",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@gz-msgs",
        "@gz-transport",
    ],
)

cc_test(
    name = "joint_state_test",
    srcs = ["joint_state_test.cc"],
    deps = [
        ":gazebo_client",
        ":gazebo_client_test_env",
        ":joint_state",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//intrinsic/util/thread",
        "//intrinsic/util/thread:stop_token",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/synchronization",
        "@abseil-cpp//absl/time",
        "@gz-msgs",
        "@gz-transport",
    ],
)

cc_library(
    name = "topic_utils",
    srcs = ["topic_utils.cc"],
    hdrs = ["topic_utils.h"],
    deps = [
        ":gazebo_client",
        ":joint_control",
        ":joint_state",
        "//intrinsic/simulation/gazebo:gz_topic_constants",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/base:nullability",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings:str_format",
    ],
)

cc_test(
    name = "topic_utils_test",
    srcs = ["topic_utils_test.cc"],
    deps = [
        ":gazebo_client",
        ":gazebo_client_test_env",
        ":joint_control",
        ":joint_state",
        ":topic_utils",
        "//intrinsic/util/proto:parse_text_proto",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
    ],
)

cc_library(
    name = "gazebo_client_test_env",
    testonly = True,
    hdrs = ["gazebo_client_test_env.h"],
    visibility = [
        ":__subpackages__",
    ],
    deps = [
        ":gazebo_client",
        "//intrinsic/simulation/gazebo/proto/v1:gazebo_service_cc_grpc_proto",
        "//intrinsic/simulation/gazebo/proto/v1:gazebo_service_cc_proto",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@com_github_grpc_grpc//:grpc++",
        "@com_google_googletest//:gtest",
    ],
)
