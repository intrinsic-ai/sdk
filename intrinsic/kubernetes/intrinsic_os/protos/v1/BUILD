# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//intrinsic/util/proto/build_defs:descriptor_set.bzl", "proto_source_code_info_transitive_descriptor_set")

package(default_visibility = ["//visibility:public"])

proto_library(
    name = "cluster_proto",
    srcs = [
        "cluster.proto",
    ],
    deps = [
        "//intrinsic/util/status:extended_status_proto",
    ],
)

proto_library(
    name = "shutdown_proto",
    srcs = [
        "shutdown.proto",
    ],
    deps = [
        "@com_google_protobuf//:duration_proto",
    ],
)

proto_library(
    name = "cluster_service_proto",
    srcs = [
        "cluster_service.proto",
    ],
    deps = [
        ":cluster_proto",
        ":shutdown_proto",
    ],
)

proto_source_code_info_transitive_descriptor_set(
    name = "cluster_service_proto_descriptor_set",
    deps = [":cluster_service_proto"],
)
