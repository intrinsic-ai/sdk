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

# Gripper hardware protos and headers.

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@intrinsic_js_rules//:defs.bzl", "js_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "gripper_proto",
    srcs = ["gripper.proto"],
    deps = ["//intrinsic/perception/proto/v1:settings_proto"],
)

py_proto_library(
    name = "gripper_py_pb2",
    deps = [":gripper_proto"],
)

cc_proto_library(
    name = "gripper_cc_proto",
    deps = [":gripper_proto"],
)

js_library(
    name = "gripper_jspb_proto",
    deps = [":gripper_proto"],
)
