# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "accesscontrol_proto",
    srcs = ["accesscontrol.proto"],
    deps = [
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-openapiv2/options:options_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
        "@com_google_googleapis//google/longrunning:operations_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

proto_source_code_info_transitive_descriptor_set(
    name = "accesscontrol_proto_descriptor_set",
    deps = [":accesscontrol_proto"],
)

go_proto_library(
    name = "accesscontrol_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/kubernetes/accounts/service/api/accesscontrol/v1/accesscontrol_go_proto",
    protos = [":accesscontrol_proto"],
    deps = [
        "@com_github_grpc_ecosystem_grpc_gateway_v2//protoc-gen-openapiv2/options",
        "@com_google_cloud_go_longrunning//autogen/longrunningpb",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)
