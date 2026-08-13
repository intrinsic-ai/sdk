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
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")

package(default_visibility = ["//visibility:public"])

exports_files(["opcua_equipment_service.proto"])

proto_library(
    name = "opcua_equipment_service_proto",
    srcs = ["opcua_equipment_service.proto"],
    deps = [
        "//intrinsic/assets/proto:field_metadata_proto",
        "//intrinsic/assets/proto/v1:resolved_dependency_proto",
        "//intrinsic/hardware/gpio/v1:gpio_service_proto",
        "//intrinsic/hardware/gpio/v1:signal_proto",
        "//intrinsic/resources/proto:resource_handle_proto",
        "//intrinsic/world/proto:object_world_updates_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:duration_proto",
    ],
)

cc_proto_library(
    name = "opcua_equipment_service_cc_proto",
    deps = [":opcua_equipment_service_proto"],
)

cc_grpc_library(
    name = "opcua_equipment_service_cc_grpc_proto",
    srcs = [":opcua_equipment_service_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":opcua_equipment_service_cc_proto"],
)

py_proto_library(
    name = "opcua_equipment_service_py_pb2",
    deps = [":opcua_equipment_service_proto"],
)

py_grpc_library(
    name = "opcua_equipment_service_py_pb2_grpc",
    srcs = [":opcua_equipment_service_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":opcua_equipment_service_py_pb2"],
)
