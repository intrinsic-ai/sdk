# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

proto_library(
    name = "webpubsub_proto",
    srcs = ["webpubsub.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":controlplane_proto",
        "@com_google_googleapis//google/api:annotations_proto",
        "@com_google_googleapis//google/api:field_behavior_proto",
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
        ":controlplane_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)

proto_library(
    name = "signalling_proto",
    srcs = ["signalling.proto"],
    visibility = ["//intrinsic:public_api_users"],
)

go_proto_library(
    name = "signalling_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_proto"],
    importpath = "intrinsic/httpjson/proto/v1alpha/signalling_go_proto",
    protos = [
        ":signalling_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
)

proto_library(
    name = "controlplane_proto",
    srcs = ["controlplane.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":signalling_proto",
        "@com_google_googleapis//google/rpc:status_proto",
    ],
)

go_proto_library(
    name = "controlplane_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_proto"],
    importpath = "intrinsic/httpjson/proto/v1alpha/controlplane_go_proto",
    protos = [
        ":controlplane_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        ":signalling_go_proto",
        "@org_golang_google_genproto_googleapis_rpc//status",
    ],
)

proto_library(
    name = "dataplane_proto",
    srcs = ["dataplane.proto"],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/platform/pubsub/adapters:pubsub_proto",
        "@com_google_googleapis//google/rpc:status_proto",
    ],
)

go_proto_library(
    name = "dataplane_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_proto"],
    importpath = "intrinsic/httpjson/proto/v1alpha/dataplane_go_proto",
    protos = [
        ":dataplane_proto",
    ],
    visibility = ["//intrinsic:public_api_users"],
    deps = [
        "//intrinsic/platform/pubsub/adapters:pubsub_go_proto",
        "@org_golang_google_genproto_googleapis_rpc//status",
    ],
)
