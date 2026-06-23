# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

proto_library(
    name = "webpubsub_proto",
    srcs = ["webpubsub.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_protobuf//:empty_proto",
    ],
)

go_proto_library(
    name = "webpubsub_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
        "//bazel:go_gen_grpc_gateway",
    ],
    importpath = "intrinsic/httpjson/proto/v1alpha/webpubsub_go_proto",
    protos = [
        ":webpubsub_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)
