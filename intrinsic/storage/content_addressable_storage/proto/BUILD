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

"""Build definitions for content addressable storage service protos."""

load("@ai_intrinsic_sdks_pip_deps//:requirements.bzl", "requirement")
load("@com_github_grpc_grpc//bazel:cc_grpc_library.bzl", "cc_grpc_library")
load("@com_github_grpc_grpc//bazel:python_rules.bzl", "py_grpc_library")
load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

proto_library(
    name = "cas_service_proto",
    srcs = ["cas_service.proto"],
)

go_proto_library(
    name = "cas_service_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/storage/content_addressable_storage/proto/cas_service_go_proto",
    protos = [":cas_service_proto"],
    visibility = ["//intrinsic:public_api_users"],
)

cc_proto_library(
    name = "cas_service_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":cas_service_proto"],
)

cc_grpc_library(
    name = "cas_service_cc_grpc_proto",
    srcs = [":cas_service_proto"],
    generate_mocks = False,  # Please use FakeCAS instead.
    grpc_only = True,
    visibility = ["//intrinsic:public_api_users"],
    deps = [":cas_service_cc_proto"],
)

py_proto_library(
    name = "cas_service_py_pb2",
    visibility = [
        "//intrinsic:public_api_users",
        "//intrinsic/perception/assets:__subpackages__",
    ],
    deps = [":cas_service_proto"],
)

py_grpc_library(
    name = "cas_service_py_pb2_grpc",
    srcs = [":cas_service_proto"],
    grpc_library = requirement("grpcio"),
    visibility = ["//intrinsic:public_api_users"],
    deps = [":cas_service_py_pb2"],
)
