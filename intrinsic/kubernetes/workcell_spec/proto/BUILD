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

# gRPC service definition for the workcell-spec-service.

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@intrinsic_js_rules//:defs.bzl", "js_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "image_proto",
    srcs = ["image.proto"],
)

go_proto_library(
    name = "image_go_proto",
    importpath = "intrinsic/kubernetes/workcell_spec/proto/image_go_proto",
    protos = [":image_proto"],
)

py_proto_library(
    name = "image_py_pb2",
    deps = [":image_proto"],
)

js_library(
    name = "image_jspb_proto",
    deps = [":image_proto"],
)

proto_library(
    name = "deploy_proto",
    srcs = ["deploy.proto"],
    deps = ["//intrinsic/config/proto:operation_mode_proto"],
)

go_proto_library(
    name = "deploy_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/kubernetes/workcell_spec/proto/deploy_go_proto",
    protos = [":deploy_proto"],
    deps = ["//intrinsic/config/proto:operation_mode_go_proto"],
)
