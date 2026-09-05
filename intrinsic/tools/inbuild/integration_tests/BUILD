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

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "inbuild_service_proto",
    srcs = ["inbuild_service.proto"],
)

cc_proto_library(
    name = "inbuild_service_cc_proto",
    deps = [":inbuild_service_proto"],
)

py_proto_library(
    name = "inbuild_service_py_pb2",
    deps = [":inbuild_service_proto"],
)

proto_library(
    name = "inbuild_skill_proto",
    srcs = ["inbuild_skill.proto"],
)

cc_proto_library(
    name = "inbuild_skill_cc_proto",
    deps = [":inbuild_skill_proto"],
)

py_proto_library(
    name = "inbuild_skill_py_pb2",
    deps = [":inbuild_skill_proto"],
)
