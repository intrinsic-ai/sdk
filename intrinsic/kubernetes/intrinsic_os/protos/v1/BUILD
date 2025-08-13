# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "cluster_service_proto",
    srcs = [
        "cluster.proto",
        "cluster_service.proto",
        "shutdown.proto",
    ],
    deps = [
        "//intrinsic/util/status:extended_status_proto",
        "@com_google_protobuf//:duration_proto",
    ],
)
