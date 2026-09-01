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
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

js_library(
    name = "webpubsub_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/httpjson/proto/v1alpha:webpubsub_proto"],
)

js_library(
    name = "signalling_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/httpjson/proto/v1alpha:signalling_proto"],
)

js_library(
    name = "controlplane_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/httpjson/proto/v1alpha:controlplane_proto"],
)

js_library(
    name = "dataplane_jspb_proto",
    visibility = ["//intrinsic:public_api_users"],
    deps = ["//intrinsic/httpjson/proto/v1alpha:dataplane_proto"],
)

proto_library(
    name = "controlplane_proto",
    srcs = ["controlplane.proto"],
    deps = [
        ":signalling_proto",
        "@com_google_googleapis//google/rpc:status_proto",
    ],
)

proto_library(
    name = "dataplane_proto",
    srcs = ["dataplane.proto"],
    deps = [
        "//intrinsic/platform/pubsub/adapters:pubsub_proto",
        "@com_google_googleapis//google/rpc:status_proto",
    ],
)

proto_library(
    name = "signalling_proto",
    srcs = ["signalling.proto"],
)

proto_library(
    name = "webpubsub_proto",
    srcs = ["webpubsub.proto"],
    deps = [
        ":controlplane_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

go_proto_library(
    name = "controlplane_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_proto"],
    importpath = "intrinsic/httpjson/proto/v1alpha/controlplane_go_proto",
    protos = [
        ":controlplane_proto",
    ],
    deps = [
        ":signalling_go_proto",
        "@org_golang_google_genproto_googleapis_rpc//status",
    ],
)

go_proto_library(
    name = "dataplane_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_proto"],
    importpath = "intrinsic/httpjson/proto/v1alpha/dataplane_go_proto",
    protos = [
        ":dataplane_proto",
    ],
    deps = [
        "//intrinsic/platform/pubsub/adapters:pubsub_go_proto",
        "@org_golang_google_genproto_googleapis_rpc//status",
    ],
)

go_proto_library(
    name = "signalling_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_proto"],
    importpath = "intrinsic/httpjson/proto/v1alpha/signalling_go_proto",
    protos = [
        ":signalling_proto",
    ],
)

go_proto_library(
    name = "webpubsub_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
    ],
    importpath = "intrinsic/httpjson/proto/v1alpha/webpubsub_go_proto",
    protos = [
        ":webpubsub_proto",
    ],
    deps = [
        ":controlplane_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)
