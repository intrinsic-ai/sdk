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

load("//bazel:go_macros.bzl", "go_library", "go_test")

go_library(
    name = "protoregistryclient",
    srcs = ["proto_registry_client.go"],
    importpath = "intrinsic/proto_tools/registry/protoregistryclient",
    visibility = [
        "//intrinsic:public_api_users",
        "//intrinsic/proto_tools/registry:__subpackages__",
    ],
    deps = [
        "//intrinsic/proto_tools/proto:proto_registry_go_proto",
        "//intrinsic/util/proto:registryutil",
        "//intrinsic/util/proto:typeurl",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//reflect/protoregistry",
    ],
)

go_test(
    name = "protoregistryclient_test",
    srcs = ["proto_registry_client_test.go"],
    embed = [":protoregistryclient"],
    importpath = "intrinsic/proto_tools/registry/protoregistryclient_test",
    deps = [
        "//intrinsic/proto_tools/proto:proto_registry_go_proto",
        "//intrinsic/testing:grpctest",
        "//intrinsic/util/proto:registryutil",
        "@com_github_google_go_cmp//cmp",
        "@com_github_google_go_cmp//cmp/cmpopts",
        "@org_golang_google_grpc//:grpc",
        "@org_golang_google_grpc//codes",
        "@org_golang_google_grpc//credentials/insecure",
        "@org_golang_google_grpc//status",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//reflect/protoregistry",
        "@org_golang_google_protobuf//types/descriptorpb",
    ],
)
