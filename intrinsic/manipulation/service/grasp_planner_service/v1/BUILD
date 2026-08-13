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

# Service definition for the grasp planner service.

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")

package(
    default_visibility = [
        "//intrinsic/logging:__subpackages__",
        "//intrinsic/manipulation:__subpackages__",
    ],
)

cc_proto_library(
    name = "grasp_planner_service_cc_proto",
    deps = [":grasp_planner_service_proto"],
)

proto_library(
    name = "grasp_planner_service_proto",
    srcs = ["grasp_planner_service.proto"],
    deps = [
        "//intrinsic/icon/proto:joint_space_proto",
        "//intrinsic/manipulation/grasping:grasp_execution_planner_params_proto",
        "//intrinsic/manipulation/grasping:grasp_planner_params_proto",
        "//intrinsic/manipulation/grasping:grasp_proto",
        "//intrinsic/manipulation/grasping:grasp_ranker_params_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

py_proto_library(
    name = "grasp_planner_service_py_pb2",
    deps = [":grasp_planner_service_proto"],
)

py_grpc_library(
    name = "grasp_planner_service_py_pb2_grpc",
    srcs = [":grasp_planner_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":grasp_planner_service_py_pb2"],
)
