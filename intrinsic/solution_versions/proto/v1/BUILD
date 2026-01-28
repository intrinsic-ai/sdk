# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("//bazel:go_macros.bzl", "go_grpc_library", "go_proto_library")

proto_library(
    name = "snapshot_proto",
    srcs = ["snapshot.proto"],
    visibility = ["//intrinsic/modified_solution/proto/v1:__pkg__"],
    deps = [
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "snapshot_go_proto",
    protos = [":snapshot_proto"],
    visibility = [
        "//intrinsic/solution_versions:__subpackages__",
        "//intrinsic/tools/inctl/cmd:__subpackages__",
    ],
)

proto_library(
    name = "branch_proto",
    srcs = ["branch.proto"],
    visibility = ["//intrinsic/frontend/cloud/portal_bridge:__subpackages__"],
    deps = [
        "//intrinsic/config/proto:operation_mode_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

go_proto_library(
    name = "branch_go_proto",
    protos = [":branch_proto"],
    visibility = [
        "//intrinsic/frontend/cloud:__subpackages__",
        "//intrinsic/frontend/cloud/portal_bridge:__subpackages__",
        "//intrinsic/kubernetes/data_store/http_endpoints:__pkg__",
        "//intrinsic/performance/inctl:__subpackages__",
        "//intrinsic/solution_versions:__subpackages__",
        "//intrinsic/tools/inctl/cmd:__subpackages__",
    ],
    deps = ["//intrinsic/config/proto:operation_mode_go_proto"],
)

proto_library(
    name = "solution_version_service_proto",
    srcs = ["solution_version_service.proto"],
    deps = [
        ":branch_proto",
        ":snapshot_proto",
        "@com_google_protobuf//:field_mask_proto",
    ],
)

go_proto_library(
    name = "solution_version_service_go_proto",
    protos = [":solution_version_service_proto"],
    visibility = [
        "//intrinsic/assets/deploy:__pkg__",
        "//intrinsic/frontend/cloud:__subpackages__",
        "//intrinsic/kubernetes/data_store/http_endpoints:__pkg__",
        "//intrinsic/kubernetes/workcell_spec:__pkg__",
        "//intrinsic/performance/inctl:__subpackages__",
        "//intrinsic/solution_versions:__subpackages__",
        "//intrinsic/tools/inctl/cmd:__subpackages__",
    ],
    deps = [
        ":branch_go_proto",
        ":snapshot_go_proto",
    ],
)

go_grpc_library(
    name = "solution_version_service_go_grpc_proto",
    protos = [":solution_version_service_proto"],
    visibility = [
        "//intrinsic/assets/deploy:__pkg__",
        "//intrinsic/frontend/cloud:__subpackages__",
        "//intrinsic/kubernetes/data_store/http_endpoints:__pkg__",
        "//intrinsic/kubernetes/workcell_spec:__pkg__",
        "//intrinsic/modified_solution:__subpackages__",
        "//intrinsic/performance/inctl:__subpackages__",
        "//intrinsic/solution_versions:__subpackages__",
        "//intrinsic/tools/inctl/cmd:__subpackages__",
    ],
    deps = [
        ":branch_go_proto",
        ":snapshot_go_proto",
    ],
)
