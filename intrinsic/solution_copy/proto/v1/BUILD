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
load("@intrinsic_js_rules//:defs.bzl", "js_library")
load("//bazel:go_macros.bzl", "go_proto_library")

js_library(
    name = "solution_copy_service_jspb_proto",
    visibility = [
        "//intrinsic/portal/frontend:__subpackages__",
    ],
    deps = [":solution_copy_service_proto"],
)

proto_library(
    name = "solution_copy_service_proto",
    srcs = ["solution_copy_service.proto"],
    visibility = [
        "//visibility:public",
    ],
    deps = [
        "//intrinsic/assets/proto:id_proto",
        "//intrinsic/solution_versions/proto/v1:snapshot_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
        "@com_google_protobuf//:go_features_proto",
    ],
)

go_proto_library(
    name = "solution_copy_service_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/solution_copy/proto/v1/solution_copy_service_go_proto",
    protos = [":solution_copy_service_proto"],
    visibility = [
        "//visibility:public",
    ],
    deps = [
        "//intrinsic/assets/proto:id_go_proto",
        "//intrinsic/solution_versions/proto/v1:snapshot_go_proto",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
    ],
)
