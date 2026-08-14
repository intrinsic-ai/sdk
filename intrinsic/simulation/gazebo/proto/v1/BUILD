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

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "gazebo_service_proto",
    srcs = ["gazebo_service.proto"],
)

cc_proto_library(
    name = "gazebo_service_cc_proto",
    deps = [":gazebo_service_proto"],
)

cc_grpc_library(
    name = "gazebo_service_cc_grpc_proto",
    srcs = [":gazebo_service_proto"],
    generate_mocks = True,
    grpc_only = True,
    deps = [":gazebo_service_cc_proto"],
)
