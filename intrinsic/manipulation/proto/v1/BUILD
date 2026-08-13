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

# Commonly used protos for manipulation

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "gripper_assembly_metadata_proto",
    srcs = ["gripper_assembly_metadata.proto"],
    visibility = ["//intrinsic:internal_api_users"],
    deps = [
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/math/proto:pose_proto",
        "//intrinsic/math/proto:vector3_proto",
        "//intrinsic/scene/proto/v1:scene_object_proto",
    ],
)

cc_proto_library(
    name = "gripper_assembly_metadata_cc_proto",
    deps = [":gripper_assembly_metadata_proto"],
)

py_proto_library(
    name = "gripper_assembly_metadata_py_pb2",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [":gripper_assembly_metadata_proto"],
)

proto_library(
    name = "gripper_exchange_service_proto",
    srcs = ["gripper_exchange_service.proto"],
    deps = [
        ":gripper_assembly_metadata_proto",
        "//intrinsic/world/proto:object_world_refs_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

py_proto_library(
    name = "gripper_exchange_service_py_pb2",
    deps = [":gripper_exchange_service_proto"],
)

py_grpc_library(
    name = "gripper_exchange_service_py_pb2_grpc",
    srcs = [":gripper_exchange_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":gripper_exchange_service_py_pb2"],
)

cc_proto_library(
    name = "gripper_exchange_service_cc_proto",
    deps = [":gripper_exchange_service_proto"],
)
