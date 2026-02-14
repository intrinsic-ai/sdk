# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = [
    "//visibility:public",
])

proto_library(
    name = "discoveryapi_proto",
    srcs = ["discovery_api.proto"],
    deps = [
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

go_proto_library(
    name = "discoveryapi_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/kubernetes/accounts/service/api/v1/discoveryapi_go_proto",
    protos = [
        ":discoveryapi_proto",
    ],
    deps = ["@org_golang_google_genproto_googleapis_api//annotations"],
)
