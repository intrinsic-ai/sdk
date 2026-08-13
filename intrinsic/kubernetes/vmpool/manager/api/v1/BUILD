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
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@intrinsic_js_rules//:defs.bzl", "js_library")
load("//bazel:go_macros.bzl", "go_proto_library")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "lease_api_proto",
    srcs = ["lease_api.proto"],
    deps = [
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

proto_source_code_info_transitive_descriptor_set(
    name = "lease_api_proto_descriptor_set",
    deps = [":lease_api_proto"],
)

go_proto_library(
    name = "lease_api_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_proto",
    protos = [":lease_api_proto"],
    deps = [
        "@org_golang_google_protobuf//types/known/emptypb",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)

js_library(
    name = "lease_api_jspb_proto",
    deps = [":lease_api_proto"],
)

py_proto_library(
    name = "lease_api_py_pb2",
    deps = [":lease_api_proto"],
)

py_grpc_library(
    name = "lease_api_py_pb2_grpc",
    srcs = [":lease_api_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":lease_api_py_pb2"],
)
