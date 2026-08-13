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

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@rules_python//python:defs.bzl", "py_library")

package(default_visibility = [
    "//visibility:public",
])

py_library(
    name = "grasp_planner_client",
    srcs = ["grasp_planner_client.py"],
    deps = [
        "//intrinsic/assets:interface_utils_py",
        "//intrinsic/assets/dependencies:utils_py",
        "//intrinsic/assets/proto/v1:resolved_dependency_py_pb2",
        "//intrinsic/manipulation/grasping:grasp_execution_planner_params_py_pb2",
        "//intrinsic/manipulation/grasping:grasp_planner_params_py_pb2",
        "//intrinsic/manipulation/grasping:grasp_planning_interfaces",
        "//intrinsic/manipulation/grasping:grasp_py_pb2",
        "//intrinsic/manipulation/grasping:grasp_ranker_params_py_pb2",
        "//intrinsic/manipulation/service/grasp_planner_service/v1:grasp_planner_service_py_pb2",
        "//intrinsic/manipulation/service/grasp_planner_service/v1:grasp_planner_service_py_pb2_grpc",
        "//intrinsic/world/proto:object_world_refs_py_pb2",
        "@ai_intrinsic_sdks_pip_deps//numpy",
        requirement("grpcio"),
        "@com_google_absl_py//absl/logging",
        "@com_google_protobuf//:protobuf_python",
    ],
)
