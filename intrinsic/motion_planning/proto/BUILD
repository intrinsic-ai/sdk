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

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@intrinsic_js_rules//:defs.bzl", "js_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

js_library(
    name = "motion_target_jspb_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = ["//intrinsic/motion_planning/proto:motion_target_proto"],
)

proto_library(
    name = "motion_target_proto",
    srcs = ["motion_target.proto"],
    deps = [
        "//intrinsic/icon/proto:joint_space_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

cc_proto_library(
    name = "motion_target_cc_proto",
    deps = [":motion_target_proto"],
)

go_proto_library(
    name = "motion_target_go_proto",
    importpath = "intrinsic/motion_planning/proto/motion_target_go_proto",
    proto = ":motion_target_proto",
    deps = [
        "//intrinsic/icon/proto:joint_space_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
    ],
)

py_proto_library(
    name = "motion_target_py_pb2",
    deps = [":motion_target_proto"],
)
