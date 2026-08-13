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

# This package contains public protos that are exposed to clients of Intrinsic perception.

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

proto_library(
    name = "vector_proto",
    srcs = ["vector.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

py_proto_library(
    name = "vector_py_pb2",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":vector_proto"],
)

cc_proto_library(
    name = "vector_cc_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = [":vector_proto"],
)

go_proto_library(
    name = "vector_go_proto",
    importpath = "intrinsic/perception/proto/vector_go_proto",
    protos = [":vector_proto"],
    visibility = ["//intrinsic:public_api_users"],
)
