# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("//bazel:go_macros.bzl", "go_proto_library")
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

go_proto_library(
    name = "cluster_go_proto",
    importpath = "intrinsic/kubernetes/intrinsic_os/protos/v1/cluster_go_proto",
    proto = ":cluster_proto",
    deps = [
        "//intrinsic/util/status:extended_status_go_proto",
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

go_proto_library(
    name = "shutdown_go_proto",
    importpath = "intrinsic/kubernetes/intrinsic_os/protos/v1/shutdown_go_proto",
    proto = ":shutdown_proto",
)

proto_library(
    name = "cluster_service_proto",
    srcs = [
        "cluster_service.proto",
    ],
    deps = [
        ":cluster_proto",
        ":dns_leases_proto",
        ":shutdown_proto",
        "@com_google_googleapis//google/api:annotations_proto",
    ],
)

go_proto_library(
    name = "cluster_service_go_proto",
    compilers = [
        "//bazel:go_gen_grpc_gateway",
        "@io_bazel_rules_go//proto:go_grpc_v2",
        "@io_bazel_rules_go//proto:go_proto",
    ],
    importpath = "intrinsic/kubernetes/intrinsic_os/protos/v1/cluster_service_go_proto",
    proto = ":cluster_service_proto",
    deps = [
        ":cluster_go_proto",
        ":dns_leases_go_proto",
        ":shutdown_go_proto",
        "@org_golang_google_genproto_googleapis_api//annotations",
    ],
)

proto_source_code_info_transitive_descriptor_set(
    name = "cluster_service_proto_descriptor_set",
    deps = [":cluster_service_proto"],
)

py_proto_library(
    name = "cluster_py_pb2",
    deps = [":cluster_proto"],
)

py_proto_library(
    name = "cluster_service_py_pb2",
    deps = [":cluster_service_proto"],
)

py_proto_library(
    name = "shutdown_py_pb2",
    deps = [":shutdown_proto"],
)

proto_library(
    name = "dns_leases_proto",
    srcs = ["dns_leases.proto"],
    deps = ["@com_google_googleapis//google/api:field_behavior_proto"],
)

go_proto_library(
    name = "dns_leases_go_proto",
    importpath = "intrinsic/kubernetes/intrinsic_os/protos/v1/dns_leases_go_proto",
    proto = ":dns_leases_proto",
    deps = ["@org_golang_google_genproto_googleapis_api//annotations"],
)
