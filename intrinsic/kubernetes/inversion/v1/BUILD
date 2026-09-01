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
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

package(default_visibility = ["//intrinsic:public_api_users"])

proto_source_code_info_transitive_descriptor_set(
    name = "inversion_proto_descriptor_set",
    deps = ["//intrinsic/kubernetes/inversion/v1:inversion_proto"],
)

proto_library(
    name = "inversion_proto",
    srcs = ["inversion.proto"],
    deps = ["@com_google_googleapis//google/api:annotations_proto"],
)

cc_proto_library(
    name = "inversion_cc_proto",
    deps = [":inversion_proto"],
)

go_proto_library(
    name = "inversion_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
    ],
    importpath = "intrinsic/kubernetes/inversion/v1/inversion_go_proto",
    protos = [":inversion_proto"],
    deps = ["@org_golang_google_genproto_googleapis_api//annotations"],
)

cc_grpc_library(
    name = "inversion_cc_grpc_proto",
    srcs = ["//intrinsic/kubernetes/inversion/v1:inversion_proto"],
    grpc_only = True,
    deps = [":inversion_cc_proto"],
)
