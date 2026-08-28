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

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "forwarding_service_proto",
    srcs = ["forwarding_service.proto"],
    deps = ["//intrinsic/platform/pubsub/connect/onprem/relay_router_service/proto:endpoint_spec_proto"],
)

go_proto_library(
    name = "forwarding_service_go_proto",
    importpath = "intrinsic/platform/pubsub/connect/onprem/forwarding_service/forwarding_service_go_proto",
    protos = [":forwarding_service_proto"],
    deps = [
        "//intrinsic/platform/pubsub/connect/onprem/relay_router_service/proto:endpoint_spec_go_proto",
    ],
)

py_proto_library(
    name = "forwarding_service_proto_py_pb2",
    deps = [":forwarding_service_proto"],
)
