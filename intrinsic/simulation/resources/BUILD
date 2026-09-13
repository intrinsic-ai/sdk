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

load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "spawner_server_proto",
    srcs = ["spawner_server.proto"],
    deps = [
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/math/proto:vector3_proto",
        "//intrinsic/scene/proto/v1:scene_object_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

cc_proto_library(
    name = "spawner_server_cc_proto",
    deps = [":spawner_server_proto"],
)

cc_grpc_library(
    name = "spawner_server_cc_grpc_proto",
    srcs = [":spawner_server_proto"],
    grpc_only = True,
    deps = [":spawner_server_cc_proto"],
)

go_proto_library(
    name = "spawner_server_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/simulation/resources/spawner_server_go_proto",
    protos = [":spawner_server_proto"],
    deps = [
        "//intrinsic/math/proto:pose_go_proto",
        "//intrinsic/math/proto:vector3_go_proto",
        "//intrinsic/scene/proto/v1:scene_object_go_proto",
        "//intrinsic/world/proto:object_world_refs_go_proto",
    ],
)
