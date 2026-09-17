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

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "connectivity_diagnostics_proto",
    srcs = ["connectivity_diagnostics.proto"],
    deps = [
        "@com_google_googleapis//google/api:field_behavior_proto",
        "@com_google_googleapis//google/api:field_info_proto",
        "@com_google_protobuf//:duration_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "connectivity_diagnostics_go_proto",
    importpath = "intrinsic/kubernetes/intrinsic_os/connectivity_diagnostics/proto/v1/connectivity_diagnostics_go_proto",
    proto = ":connectivity_diagnostics_proto",
    deps = [
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)

proto_library(
    name = "connectivity_diagnostics_service_proto",
    srcs = ["connectivity_diagnostics_service.proto"],
    deps = [
        ":connectivity_diagnostics_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "connectivity_diagnostics_service_go_proto",
    compilers = [
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/kubernetes/intrinsic_os/connectivity_diagnostics/proto/v1/connectivity_diagnostics_service_go_proto",
    proto = ":connectivity_diagnostics_service_proto",
    deps = [
        ":connectivity_diagnostics_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)
