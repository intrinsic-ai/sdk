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
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

js_library(
    name = "tokens_jspb_proto",
    deps = ["//intrinsic/kubernetes/accounts/service/api/tokens/v2:tokens_proto"],
)

proto_library(
    name = "tokens_proto",
    srcs = ["tokens.proto"],
    deps = [
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "tokens_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
    ],
    importpath = "intrinsic/kubernetes/accounts/service/api/tokens/v2/tokens_go_proto",
    protos = [":tokens_proto"],
    deps = [
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)

py_proto_library(
    name = "tokens_py_pb2",
    deps = [
        ":tokens_proto",
    ],
)

py_grpc_library(
    name = "tokens_py_pb2_grpc",
    srcs = ["//intrinsic/kubernetes/accounts/service/api/tokens/v2:tokens_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":tokens_py_pb2"],
)
