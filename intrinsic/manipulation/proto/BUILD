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
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "workspace_proto",
    srcs = ["workspace.proto"],
    deps = [
        "//intrinsic/geometry/proto:oriented_bounding_box_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
    ],
)

cc_proto_library(
    name = "workspace_cc_proto",
    deps = [":workspace_proto"],
)

go_proto_library(
    name = "workspace_go_proto",
    importpath = "intrinsic/manipulation/proto/workspace_go_proto",
    protos = [":workspace_proto"],
    deps = [
        "//intrinsic/geometry/proto:axis_aligned_bounding_box_go_proto",
        "//intrinsic/geometry/proto:oriented_bounding_box_go_proto",
        "//intrinsic/math/proto:point_go_proto",
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/math/proto:quaternion_go_proto",
        "//intrinsic/math/proto:vector2_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
    ],
)

py_proto_library(
    name = "workspace_py_pb2",
    deps = [":workspace_proto"],
)
