# Copyright 2023 Intrinsic Innovation LLC

# PubSub adapters implementations.

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "pubsub_proto",
    srcs = ["pubsub.proto"],
    deps = [
        "@com_google_googleapis//google/rpc:status_proto",
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

cc_proto_library(
    name = "pubsub_cc_proto",
    deps = [":pubsub_proto"],
)

go_proto_library(
    name = "pubsub_go_proto",
    importpath = "intrinsic/platform/pubsub/adapters/pubsub_go_proto",
    protos = [":pubsub_proto"],
    deps = [
        "//intrinsic/logging/proto:context_go_proto",
        "@org_golang_google_genproto_googleapis_rpc//status",
    ],
)

py_proto_library(
    name = "pubsub_py_pb2",
    deps = [":pubsub_proto"],
)
