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
load("//bazel:go_macros.bzl", "go_proto_library")

proto_library(
    name = "test_service_proto",
    testonly = True,
    srcs = ["test_service.proto"],
    visibility = ["//intrinsic/assets:__subpackages__"],
)

go_proto_library(
    name = "test_service_go_proto",
    testonly = True,
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/assets/instances/connect/testing/test_service_go_proto",
    protos = [":test_service_proto"],
    visibility = ["//intrinsic/assets:__subpackages__"],
)

py_proto_library(
    name = "test_service_py_pb2",
    testonly = True,
    visibility = ["//intrinsic/assets:__subpackages__"],
    deps = [":test_service_proto"],
)

py_grpc_library(
    name = "test_service_py_grpc_proto",
    testonly = True,
    srcs = [":test_service_proto"],
    grpc_library = requirement("grpcio"),
    visibility = ["//intrinsic/assets:__subpackages__"],
    deps = [":test_service_py_pb2"],
)

cc_proto_library(
    name = "test_service_cc_proto",
    testonly = True,
    visibility = ["//intrinsic/assets:__subpackages__"],
    deps = [":test_service_proto"],
)

cc_grpc_library(
    name = "test_service_cc_grpc_proto",
    testonly = True,
    srcs = [":test_service_proto"],
    grpc_only = True,
    visibility = ["//intrinsic/assets:__subpackages__"],
    deps = [":test_service_cc_proto"],
)
