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
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "admin_set_proto",
    srcs = ["admin_set.proto"],
    deps = ["@com_google_protobuf//:any_proto"],
)

cc_proto_library(
    name = "admin_set_cc_proto",
    deps = [":admin_set_proto"],
)

cc_grpc_library(
    name = "admin_set_cc_grpc",
    srcs = [":admin_set_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":admin_set_cc_proto"],
)

go_proto_library(
    name = "admin_set_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/platform/pubsub/admin_set_grpc/v1/admin_set_go_proto",
    protos = [":admin_set_proto"],
)

py_proto_library(
    name = "admin_set_py_pb2",
    deps = [":admin_set_proto"],
)

py_grpc_library(
    name = "admin_set_py_pb2_grpc",
    srcs = [":admin_set_proto"],
    grpc_library = requirement("grpcio"),
    deps = [":admin_set_py_pb2"],
)
