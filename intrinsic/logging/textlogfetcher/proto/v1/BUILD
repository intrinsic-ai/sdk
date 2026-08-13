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
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "textlogfetcher_proto",
    srcs = ["textlogfetcher.proto"],
    deps = [
        "@com_google_googleapis//google/longrunning:operations_proto",
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "textlogfetcher_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/logging/textlogfetcher/proto/v1/textlogfetcher_go_proto",
    protos = [":textlogfetcher_proto"],
    deps = [
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_protobuf//types/known/timestamppb",
    ],
)
