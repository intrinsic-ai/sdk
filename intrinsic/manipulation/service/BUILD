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

# Services for the manipulation package.

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

package(
    default_visibility = [
        "//intrinsic/logging:__subpackages__",
        "//intrinsic/manipulation:__subpackages__",
    ],
)

proto_library(
    name = "multi_axis_auto_annotator_service_proto",
    srcs = ["multi_axis_auto_annotator_service.proto"],
    deps = [
        "//intrinsic/geometry/proto:triangle_mesh_proto",
        "//intrinsic/manipulation/grasping:schmalz_multi_axis_grasp_proto",
    ],
)

# Multi axis auto annotator service definition.
cc_proto_library(
    name = "multi_axis_auto_annotator_service_cc_proto",
    deps = [":multi_axis_auto_annotator_service_proto"],
)

go_proto_library(
    name = "multi_axis_auto_annotator_service_go_proto",
    importpath = "intrinsic/manipulation/service/multi_axis_auto_annotator_service_go_proto",
    protos = [":multi_axis_auto_annotator_service_proto"],
    deps = [
        "//intrinsic/geometry/proto:triangle_mesh_go_proto",
        "//intrinsic/manipulation/grasping:schmalz_multi_axis_grasp_go_proto",
    ],
)

py_proto_library(
    name = "multi_axis_auto_annotator_service_py_pb2",
    deps = [":multi_axis_auto_annotator_service_proto"],
)

py_grpc_library(
    name = "multi_axis_auto_annotator_service_py_pb2_grpc",
    srcs = [":multi_axis_auto_annotator_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":multi_axis_auto_annotator_service_py_pb2"],
)
